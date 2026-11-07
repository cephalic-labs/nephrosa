package sarvam

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BaseURL, Org, Workspace, App, Connection, Phone, Webhook, Key string
	Version                                                       int
	Enabled                                                       bool
}

func Load() (Config, error) {
	c := Config{BaseURL: os.Getenv("SARVAM_API_BASE_URL"), Org: os.Getenv("SARVAM_ORG_ID"), Workspace: os.Getenv("SARVAM_WORKSPACE_ID"), App: os.Getenv("SARVAM_APP_ID"), Connection: os.Getenv("SARVAM_CONNECTION_ID"), Phone: os.Getenv("SARVAM_AGENT_PHONE_NUMBER"), Webhook: os.Getenv("SARVAM_WEBHOOK_URL"), Enabled: os.Getenv("SARVAM_WEBHOOK_ENABLED") == "true"}
	if !c.Enabled {
		return c, nil
	}
	if c.BaseURL == "" {
		c.BaseURL = "https://apps.sarvam.ai"
	}
	c.Version, _ = strconv.Atoi(os.Getenv("SARVAM_APP_VERSION"))
	for _, v := range []string{c.Org, c.Workspace, c.App, c.Connection, c.Phone, c.Webhook} {
		if strings.TrimSpace(v) == "" {
			return c, errors.New("Samvaad configuration is incomplete")
		}
	}
	for _, v := range []string{c.BaseURL, c.Webhook} {
		u, e := url.Parse(v)
		if e != nil || u.Scheme != "https" || u.Host == "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
			return c, errors.New("Samvaad URLs require HTTPS and no credentials, query, or fragment")
		}
	}
	if !regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`).MatchString(c.Phone) {
		return c, errors.New("Samvaad agent phone must use E.164 format")
	}
	if !strings.HasSuffix(c.Webhook, "/api/v1/sarvam/webhooks") {
		return c, errors.New("Samvaad callback base must end in /api/v1/sarvam/webhooks")
	}
	if c.BaseURL != "https://apps.sarvam.ai" || c.Version < 1 {
		return c, errors.New("use the Sarvam API origin and a committed agent version")
	}
	b, e := os.ReadFile(os.Getenv("SARVAM_API_KEY_FILE"))
	if e != nil {
		return c, errors.New("cannot read Samvaad API key file")
	}
	c.Key = strings.TrimSpace(string(b))
	if c.Key == "" {
		return c, errors.New("Samvaad API key file is empty")
	}
	return c, nil
}

type Client struct {
	Config Config
	HTTP   *http.Client
}

func New(c Config) *Client {
	return &Client{Config: c, HTTP: &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}
func (c *Client) Start(ctx context.Context, phone, callback string, variables map[string]string) (string, error) {
	cfg := c.Config
	body, _ := json.Marshal(map[string]any{"app_config": map[string]any{"app_id": cfg.App, "app_version": cfg.Version, "connection_config": map[string]string{"connection_id": cfg.Connection, "agent_phone_number": cfg.Phone}, "agent_variables": variables, "app_overrides": map[string]string{"initial_language_name": "English"}}, "user_config": map[string]string{"user_phone_number": phone}, "webhook_config": map[string]any{"url": callback}})
	path := fmt.Sprintf("%s/api/outbounds/v1/orgs/%s/workspaces/%s/outbounds", cfg.BaseURL, url.PathEscape(cfg.Org), url.PathEscape(cfg.Workspace))
	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewReader(body))
	if err != nil {
		return "", errors.New("cannot build call request")
	}
	req.Header.Set("X-API-Key", cfg.Key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		var requestError *url.Error
		if errors.As(err, &requestError) {
			return "", fmt.Errorf("provider transport failure: %w", requestError.Err)
		}
		return "", errors.New("call request outcome is uncertain; reconcile before redialing")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		responseBody, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
		fields := validationFields(bytes.NewReader(responseBody))
		var rejection struct {
			Detail string `json:"detail"`
			Error  struct {
				Message string `json:"message"`
				Data    struct {
					Details string `json:"details"`
				} `json:"data"`
			} `json:"error"`
		}
		if json.Unmarshal(responseBody, &rejection) == nil {
			if rejection.Detail == "" {
				rejection.Detail = rejection.Error.Data.Details
			}
			if rejection.Detail == "" {
				rejection.Detail = rejection.Error.Message
			}
		}
		if rejection.Detail != "" {
			private := []string{cfg.Key, phone, callback, cfg.Org, cfg.Workspace, cfg.App, cfg.Connection, cfg.Phone}
			for _, value := range variables {
				private = append(private, value)
			}
			fields += "; reason: " + redactProviderDetail(rejection.Detail, private)
		}
		return "", fmt.Errorf("provider did not confirm acceptance (HTTP %d%s); reconcile before redialing", res.StatusCode, fields)
	}
	var out struct {
		AttemptID string `json:"attempt_id"`
	}
	err = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&out)
	if err != nil || out.AttemptID == "" || len(out.AttemptID) > 200 {
		return "", errors.New("provider response is uncertain; reconcile before redialing")
	}
	return out.AttemptID, nil
}

type Result struct {
	AttemptID     string         `json:"attempt_id"`
	Status        string         `json:"status"`
	InteractionID string         `json:"interaction_id"`
	Variables     map[string]any `json:"final_agent_variables"`
	Channel       struct {
		Type     string `json:"channel_type"`
		Provider string `json:"channel_provider"`
		Phone    string `json:"agent_phone_number"`
	} `json:"channel_info"`
}

func (r Result) Outcome() string {
	s, _ := r.Variables["outcome"].(string)
	switch s {
	case "booked", "test_reported_complete", "barrier", "callback_requested", "opted_out", "wrong_contact", "unable_to_talk":
		return s
	}
	return "needs_review"
}
func (r Result) Note() string {
	s, _ := r.Variables["follow_up_note"].(string)
	s = strings.TrimSpace(s)
	if len(s) > 240 {
		return ""
	}
	return s
}
func (r Result) Valid() bool {
	return r.AttemptID != "" && len(r.AttemptID) <= 200 && len(r.InteractionID) <= 300 && (r.Status == "connected" || r.Status == "no_answer" || r.Status == "busy" || r.Status == "failed") && r.Channel.Type == "v2v" && r.Channel.Provider != "" && r.Channel.Phone != ""
}

// Provider validation payloads may echo credentials and contacts. Retain only known field names.
func validationFields(body io.Reader) string {
	var response struct {
		Detail []struct {
			Loc []any `json:"loc"`
		} `json:"detail"`
	}
	if json.NewDecoder(io.LimitReader(body, 1<<20)).Decode(&response) != nil {
		return ""
	}
	allowed := map[string]bool{"body": true, "app_config": true, "app_id": true, "app_version": true, "connection_config": true, "connection_id": true, "agent_phone_number": true, "agent_variables": true, "app_type": true, "app_overrides": true, "initial_language_name": true, "user_config": true, "user_phone_number": true, "webhook_config": true, "url": true, "metadata": true, "value": true, "type": true, "name": true}
	paths := []string{}
	for _, detail := range response.Detail {
		path := []string{}
		for _, part := range detail.Loc {
			if name, ok := part.(string); ok && allowed[name] {
				path = append(path, name)
			}
		}
		if len(path) > 0 {
			paths = append(paths, strings.Join(path, "."))
		}
		if len(paths) >= 8 {
			break
		}
	}
	if len(paths) == 0 {
		return ""
	}
	return "; validation fields: " + strings.Join(paths, ", ")
}

func redactProviderDetail(detail string, private []string) string {
	for _, value := range private {
		if value != "" {
			detail = strings.ReplaceAll(detail, value, "[redacted]")
		}
	}
	detail = regexp.MustCompile(`https?://[^\s]+`).ReplaceAllString(detail, "[redacted URL]")
	detail = strings.Join(strings.Fields(detail), " ")
	if len(detail) > 240 {
		detail = detail[:240]
	}
	return detail
}
