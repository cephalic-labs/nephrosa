package sarvam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInstantCallContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/outbounds/v1/orgs/org/workspaces/work/outbounds" || r.Header.Get("X-API-Key") != "test-key" {
			t.Error("incorrect outbound request")
		}
		var body map[string]any
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			t.Fatal("invalid JSON")
		}
		app := body["app_config"].(map[string]any)
		if app["app_version"] != float64(2) || body["user_config"].(map[string]any)["user_phone_number"] != "+919876543210" {
			t.Error("missing version or contact")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"attempt_id":"test-attempt"}`))
	}))
	defer server.Close()
	c := New(Config{BaseURL: server.URL, Org: "org", Workspace: "work", App: "app", Version: 2, Connection: "conn", Phone: "+918000000000", Key: "test-key"})
	id, err := c.Start(context.Background(), "+919876543210", "https://example.test/callback/secret", map[string]string{"episode_reference": "test"})
	if err != nil || id != "test-attempt" {
		t.Fatalf("unexpected result %q %v", id, err)
	}
}
func TestAmbiguousProviderResponseIsNeverRetried(t *testing.T) {
	for _, response := range []string{`{}`, `invalid`} {
		count := 0
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { count++; w.Write([]byte(response)) }))
		c := New(Config{BaseURL: s.URL})
		_, err := c.Start(context.Background(), "+919876543210", "https://example.test", nil)
		s.Close()
		if err == nil || count != 1 {
			t.Fatal("uncertain response must fail with one request")
		}
	}
}
func TestInstantWebhookUsesItsOwnContract(t *testing.T) {
	var r Result
	if err := json.Unmarshal([]byte(`{"attempt_id":"test","status":"connected","channel_info":{"channel_type":"v2v","channel_provider":"exotel","agent_phone_number":"+918000000000"},"final_agent_variables":{"outcome":"opted_out"}}`), &r); err != nil || !r.Valid() || r.Outcome() != "opted_out" {
		t.Fatal("valid instant webhook rejected")
	}
	r.Variables = map[string]any{"outcome": "close_episode"}
	if r.Outcome() != "needs_review" {
		t.Fatal("unsafe output accepted")
	}
	r.Status = "completed"
	if r.Valid() {
		t.Fatal("campaign status accepted as instant status")
	}
}

func TestProviderFailureReportsStatusWithoutResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"detail":{"input":"private provider response"}}`))
	}))
	defer server.Close()
	client := New(Config{BaseURL: server.URL})
	_, err := client.Start(context.Background(), "+919876543210", "https://example.test/callback", nil)
	if err == nil || !strings.Contains(err.Error(), "HTTP 422") || strings.Contains(err.Error(), "private provider response") {
		t.Fatalf("expected safe HTTP status diagnostic, got %v", err)
	}
}

func TestValidationDiagnosticExcludesSubmittedValues(t *testing.T) {
	fields := validationFields(strings.NewReader(`{"detail":[{"loc":["body","app_config","agent_variables","private-field","value"],"input":"private-input","msg":"private-message"}]}`))
	if fields != "; validation fields: body.app_config.agent_variables.value" {
		t.Fatalf("unsafe or missing diagnostic: %s", fields)
	}
	if validationFields(strings.NewReader(`{"detail":"private-message"}`)) != "" {
		t.Fatal("unexpected provider message disclosure")
	}
}

func TestProviderReasonRedactsCredentialsAndContacts(t *testing.T) {
	result := redactProviderDetail("Rejected test-key +919876543210 https://example.test/callback/secret", []string{"test-key", "+919876543210"})
	if strings.Contains(result, "test-key") || strings.Contains(result, "919876543210") || strings.Contains(result, "secret") {
		t.Fatal("provider diagnostic contains private values")
	}
}

func TestSarvamErrorEnvelopeReportsRedactedReason(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`{"error":{"message":"Invalid Parameter","data":{"details":"Agent private-agent was not found"}}}`))
	}))
	defer server.Close()
	client := New(Config{BaseURL: server.URL, App: "private-agent"})
	_, err := client.Start(context.Background(), "+919876543210", "https://example.test/callback", nil)
	if err == nil || !strings.Contains(err.Error(), "Agent [redacted] was not found") || strings.Contains(err.Error(), "private-agent") {
		t.Fatalf("missing safe provider reason: %v", err)
	}
}
