package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/cephalic-labs/nephrosa/backend/internal/sarvam"
	"github.com/cephalic-labs/nephrosa/db/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func testDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is required for disposable PostgreSQL integration tests")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal("cannot create test pool")
	}
	schema := "test_" + strings.ToLower(token())
	schema = strings.ReplaceAll(schema, "-", "")
	if _, err = admin.Exec(ctx, `CREATE SCHEMA `+pgx.Identifier{schema}.Sanitize()); err != nil {
		admin.Close()
		t.Fatal("test role must be able to create a disposable schema")
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid test database configuration")
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		pool.Close()
		_, _ = admin.Exec(ctx, `DROP SCHEMA `+pgx.Identifier{schema}.Sanitize()+` CASCADE`)
		admin.Close()
	})
	files, _ := fs.Glob(migrations.Files, "*.sql")
	sort.Strings(files)
	for _, f := range files {
		sql, _ := migrations.Files.ReadFile(f)
		if _, err = pool.Exec(ctx, string(sql)); err != nil {
			t.Fatalf("migration %s failed: %v", f, err)
		}
	}
	return pool
}
func TestDemoWorkflowAndVoiceIntegration(t *testing.T) {
	pool := testDatabase(t)
	ctx := context.Background()
	hash, _ := bcrypt.GenerateFromPassword([]byte("fictional-test-password"), bcrypt.MinCost)
	for _, role := range []string{"clinician", "coordinator"} {
		if _, err := pool.Exec(ctx, `INSERT INTO users(role,display_name,login,password_hash) VALUES($1,$1,$1,$2)`, role, string(hash)); err != nil {
			t.Fatal(err)
		}
	}
	var dispatched struct{ Callback string }
	providerCalls := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providerCalls++
		var body struct {
			Webhook struct {
				URL string `json:"url"`
			} `json:"webhook_config"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			t.Error("invalid call request")
		}
		dispatched.Callback = body.Webhook.URL
		w.Write([]byte(`{"attempt_id":"attempt-demo"}`))
	}))
	defer provider.Close()
	handler := New(pool, sarvam.New(sarvam.Config{Enabled: true, BaseURL: provider.URL, Org: "org", Workspace: "ws", App: "app", Version: 1, Connection: "conn", Phone: "+918000000000", Webhook: "https://example.test/api/v1/sarvam/webhooks", Key: "test"}))
	request := func(method, path, token string, body any, want int) map[string]any {
		t.Helper()
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != want {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		out := map[string]any{}
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		return out
	}
	login := func(role string) string {
		out := request("POST", "/api/v1/sessions", "", map[string]string{"login": role, "password": "fictional-test-password"}, 200)
		return out["token"].(string)
	}
	request("GET", "/api/v1/dashboard/overview", "", nil, 401)
	doctor := login("clinician")
	coordinator := login("coordinator")
	var owner string
	_ = pool.QueryRow(ctx, `SELECT id::text FROM users WHERE role='coordinator'`).Scan(&owner)
	create := map[string]any{"name": "Fictional demo", "reference": "DEMO-TEST", "phone": "+919876543210", "language": "en", "owner": owner, "due": "2026-12-01", "tests": []string{"uacr", "egfr"}, "consent": true, "consentSource": "consenting demonstrator"}
	request("POST", "/api/v1/episodes", coordinator, create, 403)
	id := request("POST", "/api/v1/episodes", doctor, create, 201)["id"].(string)
	path := "/api/v1/episodes/" + id
	action := func(who string, body map[string]string, want int) { request("POST", path+"/actions", who, body, want) }
	action(doctor, map[string]string{"action": "close"}, 409)
	action(coordinator, map[string]string{"action": "review", "nextAction": "Follow-up documented."}, 403)
	action(coordinator, map[string]string{"action": "report", "reference": "DEMO-REPORT"}, 409)
	request("POST", path+"/calls", coordinator, nil, 202)
	request("POST", path+"/calls", coordinator, nil, 409)
	if providerCalls != 1 {
		t.Fatal("duplicate provider call")
	}
	callback := strings.TrimPrefix(dispatched.Callback, "https://example.test")
	event := map[string]any{"attempt_id": "attempt-demo", "status": "connected", "channel_info": map[string]string{"channel_type": "v2v", "channel_provider": "test", "agent_phone_number": "+918000000000"}, "final_agent_variables": map[string]string{"outcome": "opted_out", "follow_up_note": "Please stop automated calls."}}
	request("POST", "/api/v1/sarvam/webhooks/invalid", "", event, 404)
	request("POST", callback, "", event, 204)
	request("POST", callback, "", event, 204)
	var followUpNote string
	if err := pool.QueryRow(ctx, `SELECT follow_up_note FROM call_attempts WHERE episode_id=$1`, id).Scan(&followUpNote); err != nil || followUpNote != "Please stop automated calls." {
		t.Fatalf("callback did not persist the patient response summary: note=%q err=%v", followUpNote, err)
	}
	var tasks, contacts int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM tasks WHERE episode_id=$1`, id).Scan(&tasks)
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM contact_attempts WHERE episode_id=$1`, id).Scan(&contacts)
	if tasks != 1 || contacts != 1 {
		t.Fatal("duplicate callback produced duplicate work")
	}
	request("POST", path+"/calls", coordinator, nil, 409)
	detail := request("GET", path, doctor, nil, 200)
	if detail["consent"] != "withdrawn" {
		t.Fatal("opt-out did not withdraw consent")
	}
	for _, item := range detail["tests"].([]any) {
		test := item.(map[string]any)
		action(coordinator, map[string]string{"action": "verify_test", "testId": test["id"].(string), "source": "staff_confirmed"}, 204)
	}
	action(coordinator, map[string]string{"action": "report", "reference": "DEMO-REPORT"}, 204)
	action(doctor, map[string]string{"action": "review", "nextAction": "Fictional clinician follow-up recorded."}, 204)
	action(doctor, map[string]string{"action": "close"}, 409)
	task := detail["tasks"].([]any)[0].(map[string]any)
	action(coordinator, map[string]string{"action": "resolve_task", "taskId": task["id"].(string), "note": "Opt-out acknowledged."}, 204)
	action(doctor, map[string]string{"action": "close"}, 204)
	request("GET", path, doctor, nil, 200)
	action(coordinator, map[string]string{"action": "consent", "consent": "granted", "source": "test"}, 409)
	var auditCount int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM audit_events WHERE episode_id=$1`, id).Scan(&auditCount)
	if auditCount < 9 {
		t.Fatal("workflow missing audit history")
	}
	if _, err := pool.Exec(ctx, `UPDATE audit_events SET event_type='tampered' WHERE episode_id=$1`, id); err == nil {
		t.Fatal("audit history was mutable")
	}
	// A manually reconciled call remains terminal even if its callback is delayed.
	create["reference"] = "DEMO-RECONCILE"
	id2 := request("POST", "/api/v1/episodes", doctor, create, 201)["id"].(string)
	secret := token()
	var reconciled string
	_ = pool.QueryRow(ctx, `INSERT INTO call_attempts(episode_id,requested_by,callback_hash,agent_phone_number,status,created_at) VALUES($1,$2,$3,'+918000000000','uncertain',$4) RETURNING id::text`, id2, owner, digest(secret), time.Now().Add(-3*time.Minute)).Scan(&reconciled)
	request("POST", "/api/v1/episodes/"+id2+"/calls/"+reconciled+"/reconcile", coordinator, map[string]string{"note": "Confirmed in provider logs: no active call."}, 204)
	event["attempt_id"] = "delayed"
	request("POST", "/api/v1/sarvam/webhooks/"+secret, "", event, 204)
	var status string
	_ = pool.QueryRow(ctx, `SELECT status FROM call_attempts WHERE id=$1`, reconciled).Scan(&status)
	if status != "reconciled" {
		t.Fatal("late callback overwrote reconciliation")
	}
	request("DELETE", "/api/v1/session", doctor, nil, 204)
	request("GET", path, doctor, nil, 401)
}
