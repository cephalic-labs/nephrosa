package httpapi

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var phonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func audit(ctx context.Context, tx pgx.Tx, user, episode, event string) error {
	_, err := tx.Exec(ctx, `INSERT INTO audit_events(actor_id,episode_id,event_type) VALUES(NULLIF($1,'')::uuid,$2,$3)`, user, episode, event)
	return err
}
func (r *Router) createEpisode(c *gin.Context) {
	if !clinician(c) {
		return
	}
	var in struct {
		Name          string   `json:"name" binding:"required,max=120"`
		Reference     string   `json:"reference" binding:"required,max=80"`
		Phone         string   `json:"phone" binding:"required"`
		Language      string   `json:"language"`
		Consent       bool     `json:"consent"`
		ConsentSource string   `json:"consentSource" binding:"required,max=120"`
		Owner         string   `json:"owner"`
		Due           string   `json:"due"`
		Tests         []string `json:"tests"`
	}
	if c.ShouldBindJSON(&in) != nil || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Reference) == "" || strings.TrimSpace(in.ConsentSource) == "" || !phonePattern.MatchString(in.Phone) || !uuidPattern.MatchString(in.Owner) || in.Language != "en" || len(in.Tests) == 0 || len(in.Tests) > 2 {
		fail(c, 400, "Provide patient details, an E.164 phone number, English, an owner, consent source, and tests.")
		return
	}
	if _, err := time.Parse("2006-01-02", in.Due); err != nil {
		fail(c, 400, "Provide a valid due date.")
		return
	}
	for i, t := range in.Tests {
		if !slices.Contains([]string{"uacr", "egfr"}, t) || slices.Contains(in.Tests[:i], t) {
			fail(c, 400, "Select uACR and/or eGFR once.")
			return
		}
	}
	ctx := c.Request.Context()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		writeUnavailable(c, "create_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())
	var valid bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 AND active AND role IN ('clinician','coordinator'))`, in.Owner).Scan(&valid)
	if err != nil || !valid {
		fail(c, 400, "Choose an active clinician or coordinator.")
		return
	}
	var patient, id string
	err = tx.QueryRow(ctx, `INSERT INTO patients(external_reference,display_name,preferred_language,phone_number) VALUES($1,$2,$3,$4) RETURNING id::text`, in.Reference, in.Name, in.Language, in.Phone).Scan(&patient)
	if err != nil {
		fail(c, 409, "Patient reference already exists or patient could not be created.")
		return
	}
	err = tx.QueryRow(ctx, `INSERT INTO episodes(patient_id,owner_id,status,due_at,created_by) VALUES($1,$2,'ordered',$3,$4) RETURNING id::text`, patient, in.Owner, in.Due, actor(c).ID).Scan(&id)
	if err == nil {
		for _, test := range in.Tests {
			_, err = tx.Exec(ctx, `INSERT INTO test_orders(episode_id,test_type,due_at) VALUES($1,$2,$3)`, id, test, in.Due)
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		status := "pending"
		var recorded *time.Time
		if in.Consent {
			status = "granted"
			now := time.Now()
			recorded = &now
		}
		_, err = tx.Exec(ctx, `INSERT INTO consents(patient_id,channel,status,recorded_at,source) VALUES($1,'voice',$2,$3,$4)`, patient, status, recorded, in.ConsentSource)
	}
	if err == nil {
		err = audit(ctx, tx, actor(c).ID, id, "episode_created")
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeUnavailable(c, "create_unavailable", err)
		return
	}
	c.JSON(201, gin.H{"id": id})
}

type workflowInput struct {
	Action     string `json:"action"`
	TestID     string `json:"testId"`
	TaskID     string `json:"taskId"`
	Source     string `json:"source"`
	Reference  string `json:"reference"`
	NextAction string `json:"nextAction"`
	Note       string `json:"note"`
	Consent    string `json:"consent"`
	Owner      string `json:"owner"`
}

func closureReady(tests, reports, reviews, work int) bool {
	return tests == 0 && reports > 0 && reviews > 0 && work == 0
}
func (r *Router) workflow(c *gin.Context) {
	var in workflowInput
	if c.ShouldBindJSON(&in) != nil || !uuidPattern.MatchString(c.Param("id")) {
		fail(c, 400, "Invalid workflow request.")
		return
	}
	if slices.Contains([]string{"review", "close", "assign"}, in.Action) && !clinician(c) {
		return
	}
	ctx := c.Request.Context()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		writeUnavailable(c, "workflow_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())
	var status, patient string
	err = tx.QueryRow(ctx, `SELECT status,patient_id::text FROM episodes WHERE id=$1 FOR UPDATE`, c.Param("id")).Scan(&status, &patient)
	if errors.Is(err, pgx.ErrNoRows) {
		fail(c, 404, "Episode not found.")
		return
	}
	if err != nil {
		writeUnavailable(c, "workflow_unavailable", err)
		return
	}
	if status == "closed" {
		fail(c, 409, "This episode is closed.")
		return
	}
	invalid := func(message string) { fail(c, 409, message) }
	switch in.Action {
	case "consent":
		if !slices.Contains([]string{"granted", "withdrawn"}, in.Consent) || len(strings.TrimSpace(in.Source)) == 0 || len(in.Source) > 120 {
			fail(c, 400, "Provide consent status and its source.")
			return
		}
		// Lock the patient as well: consent applies to calls across all their episodes.
		_, err = tx.Exec(ctx, `SELECT id FROM patients WHERE id=$1 FOR UPDATE`, patient)
		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO consents(patient_id,channel,status,recorded_at,source) VALUES($1,'voice',$2,now(),$3)`, patient, in.Consent, in.Source)
		}
	case "assign":
		if !uuidPattern.MatchString(in.Owner) {
			fail(c, 400, "Select a care owner.")
			return
		}
		tag, e := tx.Exec(ctx, `UPDATE episodes SET owner_id=$1 WHERE id=$2 AND EXISTS(SELECT 1 FROM users WHERE id=$1 AND active AND role IN ('clinician','coordinator'))`, in.Owner, c.Param("id"))
		err = e
		if err == nil && tag.RowsAffected() != 1 {
			invalid("Choose an active care owner.")
			return
		}
	case "verify_test":
		if !uuidPattern.MatchString(in.TestID) || !slices.Contains([]string{"lab_document", "lab_portal", "staff_confirmed"}, in.Source) {
			fail(c, 400, "Select a test and verification source.")
			return
		}
		tag, e := tx.Exec(ctx, `INSERT INTO test_completions(test_order_id,completed_at,verification_source,recorded_by) SELECT id,now(),$1,$2 FROM test_orders WHERE id=$3 AND episode_id=$4 ON CONFLICT(test_order_id) DO NOTHING`, in.Source, actor(c).ID, in.TestID, c.Param("id"))
		err = e
		if err == nil && tag.RowsAffected() != 1 {
			invalid("Test is already verified or does not belong to this episode.")
			return
		}
	case "report":
		if len(strings.TrimSpace(in.Reference)) == 0 || len(in.Reference) > 120 {
			fail(c, 400, "Provide a report reference (up to 120 characters).")
			return
		}
		var missing int
		err = tx.QueryRow(ctx, `SELECT count(*) FROM test_orders o LEFT JOIN test_completions t ON t.test_order_id=o.id WHERE o.episode_id=$1 AND t.id IS NULL`, c.Param("id")).Scan(&missing)
		if err == nil && missing > 0 {
			invalid("Verify all ordered tests before recording report receipt.")
			return
		}
		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO reports(episode_id,received_at,reference,recorded_by) VALUES($1,now(),$2,$3)`, c.Param("id"), in.Reference, actor(c).ID)
		}
	case "review":
		if len(strings.TrimSpace(in.NextAction)) == 0 || len(in.NextAction) > 1000 {
			fail(c, 400, "Record the next action (up to 1000 characters).")
			return
		}
		var reports int
		err = tx.QueryRow(ctx, `SELECT count(*) FROM reports WHERE episode_id=$1`, c.Param("id")).Scan(&reports)
		if err == nil && reports == 0 {
			invalid("Record report receipt before clinician acknowledgement.")
			return
		}
		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO clinical_reviews(episode_id,clinician_id,acknowledged_at,next_action) VALUES($1,$2,now(),$3)`, c.Param("id"), actor(c).ID, in.NextAction)
		}
	case "resolve_task":
		if !uuidPattern.MatchString(in.TaskID) || len(strings.TrimSpace(in.Note)) == 0 || len(in.Note) > 240 {
			fail(c, 400, "Choose a task and document the resolution (up to 240 characters).")
			return
		}
		tag, e := tx.Exec(ctx, `UPDATE tasks SET status='completed',completed_at=now() WHERE id=$1 AND episode_id=$2 AND status IN ('open','in_progress')`, in.TaskID, c.Param("id"))
		err = e
		if err == nil && tag.RowsAffected() != 1 {
			invalid("Task is already resolved or does not belong to this episode.")
			return
		}

		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO audit_events(actor_id,episode_id,event_type,metadata) VALUES($1,$2,'task_resolution',jsonb_build_object('task_id',$3::text,'note',$4::text))`, actor(c).ID, c.Param("id"), in.TaskID, in.Note)
		}
	case "resolve_barriers":
		if len(strings.TrimSpace(in.Note)) == 0 || len(in.Note) > 240 {
			fail(c, 400, "Document the barrier resolution.")
			return
		}
		_, err = tx.Exec(ctx, `UPDATE barriers SET status='resolved',resolved_at=now() WHERE episode_id=$1 AND status='open'`, c.Param("id"))
		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO audit_events(actor_id,episode_id,event_type,metadata) VALUES($1,$2,'barrier_resolution',jsonb_build_object('note',$3::text))`, actor(c).ID, c.Param("id"), in.Note)
		}
	case "close":
		var missing, reports, reviews, work int
		err = tx.QueryRow(ctx, `SELECT (SELECT count(*) FROM test_orders o LEFT JOIN test_completions t ON t.test_order_id=o.id WHERE o.episode_id=$1 AND t.id IS NULL),(SELECT count(*) FROM reports WHERE episode_id=$1),(SELECT count(*) FROM clinical_reviews WHERE episode_id=$1),(SELECT count(*) FROM tasks WHERE episode_id=$1 AND status IN ('open','in_progress'))+(SELECT count(*) FROM barriers WHERE episode_id=$1 AND status='open')+(SELECT count(*) FROM call_attempts WHERE episode_id=$1 AND status IN ('requesting','pending','uncertain'))`, c.Param("id")).Scan(&missing, &reports, &reviews, &work)
		if err == nil && !closureReady(missing, reports, reviews, work) {
			invalid("Verify every test, record report receipt and clinician next action, and resolve outstanding tasks, barriers, and calls before closure.")
			return
		}
		if err == nil {
			_, err = tx.Exec(ctx, `UPDATE episodes SET status='closed' WHERE id=$1`, c.Param("id"))
		}
	default:
		fail(c, 400, "Unknown workflow action.")
		return
	}
	if err == nil && in.Action != "close" {
		_, err = tx.Exec(ctx, `UPDATE episodes SET status=CASE WHEN EXISTS(SELECT 1 FROM clinical_reviews WHERE episode_id=$1) THEN 'awaiting_clinician_review' WHEN EXISTS(SELECT 1 FROM reports WHERE episode_id=$1) THEN 'awaiting_clinician_review' WHEN NOT EXISTS(SELECT 1 FROM test_orders o LEFT JOIN test_completions t ON t.test_order_id=o.id WHERE o.episode_id=$1 AND t.id IS NULL) THEN 'awaiting_report' ELSE status END,updated_at=now() WHERE id=$1`, c.Param("id"))
	}
	if err == nil {
		err = audit(ctx, tx, actor(c).ID, c.Param("id"), in.Action)
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeUnavailable(c, "workflow_unavailable", err)
		return
	}
	c.Status(http.StatusNoContent)
}
