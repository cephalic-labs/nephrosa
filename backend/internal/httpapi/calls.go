package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cephalic-labs/nephrosa/backend/internal/sarvam"
	"github.com/gin-gonic/gin"
)

func (r *Router) startCall(c *gin.Context) {
	if r.voice == nil || !r.voice.Config.Enabled {
		fail(c, 503, "Samvaad calls are disabled until configuration and callback security are confirmed.")
		return
	}
	if !uuidPattern.MatchString(c.Param("id")) {
		fail(c, 400, "Invalid episode.")
		return
	}
	ctx := c.Request.Context()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		writeUnavailable(c, "call_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())
	var patient, status, phone, name, due, consent string
	err = tx.QueryRow(ctx, `SELECT e.patient_id::text,e.status,COALESCE(p.phone_number,''),p.display_name,e.due_at::text FROM episodes e JOIN patients p ON p.id=e.patient_id WHERE e.id=$1 FOR UPDATE OF e,p`, c.Param("id")).Scan(&patient, &status, &phone, &name, &due)
	if err != nil {
		fail(c, 404, "Episode not found.")
		return
	}
	if status == "closed" || !phonePattern.MatchString(phone) {
		fail(c, 409, "An open episode and valid phone number are required.")
		return
	}
	err = tx.QueryRow(ctx, `SELECT status FROM consents WHERE patient_id=$1 AND channel='voice' ORDER BY created_at DESC,id DESC LIMIT 1`, patient).Scan(&consent)
	if err != nil || consent != "granted" {
		fail(c, 409, "Record voice consent before calling.")
		return
	}
	secret := token()
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO call_attempts(episode_id,requested_by,callback_hash,agent_phone_number,status) VALUES($1,$2,$3,$4,'requesting') RETURNING id::text`, c.Param("id"), actor(c).ID, digest(secret), r.voice.Config.Phone).Scan(&id)
	if err != nil {
		fail(c, 409, "A call is already active or uncertain. Reconcile it before redialing.")
		return
	}
	if err = audit(ctx, tx, actor(c).ID, c.Param("id"), "call_requested"); err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeUnavailable(c, "call_unavailable", err)
		return
	}
	// Do not bind provider dispatch to browser disconnects after the durable reservation.
	dispatch, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	attempt, callErr := r.voice.Start(dispatch, phone, strings.TrimRight(r.voice.Config.Webhook, "/")+"/"+secret, map[string]string{"patient_name": name, "due_date": due, "episode_reference": c.Param("id")})
	state := "pending"
	if callErr != nil {
		state = "uncertain"
	}
	_, err = r.db.Exec(dispatch, `UPDATE call_attempts SET attempt_id=COALESCE(attempt_id,NULLIF($1,'')),status=CASE WHEN finished_at IS NULL THEN $2 ELSE status END WHERE id=$3`, attempt, state, id)
	if err != nil {
		writeUnavailable(c, "call_record_unavailable", err)
		return
	}
	if callErr != nil {
		slog.Warn("Samvaad call acceptance failed", "error", callErr)
		fail(c, 502, "Call acceptance is uncertain. Check Samvaad call logs and reconcile before redialing.")
		return
	}
	c.JSON(202, gin.H{"id": id, "attemptId": attempt})
}
func (r *Router) callList(c *gin.Context) {
	rows, err := r.db.Query(c.Request.Context(), `SELECT id::text,COALESCE(attempt_id,''),status,outcome,created_at,reconciliation_note,follow_up_note FROM call_attempts WHERE episode_id::text=$1 ORDER BY created_at DESC`, c.Param("id"))
	if err != nil {
		writeUnavailable(c, "calls_unavailable", err)
		return
	}
	defer rows.Close()
	out := []gin.H{}
	for rows.Next() {
		var id, attempt, status, outcome, note, followUp string
		var at time.Time
		if err = rows.Scan(&id, &attempt, &status, &outcome, &at, &note, &followUp); err != nil {
			writeUnavailable(c, "calls_unavailable", err)
			return
		}
		out = append(out, gin.H{"id": id, "attemptId": attempt, "status": status, "outcome": outcome, "createdAt": at, "note": note, "followUpNote": followUp})
	}
	if rows.Err() != nil {
		writeUnavailable(c, "calls_unavailable", rows.Err())
		return
	}
	c.JSON(200, out)
}
func (r *Router) webhook(c *gin.Context) {
	if r.voice == nil || !r.voice.Config.Enabled {
		c.Status(404)
		return
	}
	var result sarvam.Result
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	if json.NewDecoder(c.Request.Body).Decode(&result) != nil || !result.Valid() {
		fail(c, 400, "Invalid instant-call event.")
		return
	}
	ctx := c.Request.Context()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		writeUnavailable(c, "callback_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())
	var id, episode, patient, attempt, phone string
	var finished *time.Time
	// Use the same episode→patient→call lock order as workflow mutations.
	err = tx.QueryRow(ctx, `SELECT e.id::text,e.patient_id::text FROM episodes e JOIN call_attempts a ON a.episode_id=e.id WHERE a.callback_hash=$1 FOR UPDATE OF e`, digest(c.Param("secret"))).Scan(&episode, &patient)
	if err != nil {
		c.Status(404)
		return
	}
	_, err = tx.Exec(ctx, `SELECT id FROM patients WHERE id=$1 FOR UPDATE`, patient)
	if err != nil {
		writeUnavailable(c, "callback_unavailable", err)
		return
	}
	err = tx.QueryRow(ctx, `SELECT id::text,COALESCE(attempt_id,''),agent_phone_number,finished_at FROM call_attempts WHERE callback_hash=$1 FOR UPDATE`, digest(c.Param("secret"))).Scan(&id, &attempt, &phone, &finished)
	if err != nil {
		c.Status(404)
		return
	}
	if (attempt != "" && attempt != result.AttemptID) || phone != result.Channel.Phone {
		fail(c, 409, "Call identifiers do not match.")
		return
	}
	if finished != nil {
		c.Status(204)
		return
	}
	outcome := result.Outcome()
	if result.Status != "connected" {
		outcome = result.Status
	}
	_, err = tx.Exec(ctx, `UPDATE call_attempts SET attempt_id=$1,status=$2,outcome=$3,interaction_id=$4,follow_up_note=$6,finished_at=now() WHERE id=$5`, result.AttemptID, result.Status, outcome, result.InteractionID, id, result.Note())
	contact := "other"
	switch outcome {
	case "opted_out", "wrong_contact", "callback_requested", "no_answer":
		contact = outcome
	default:
		if result.Status == "connected" {
			contact = "reached"
		}
	}
	if err == nil {
		_, err = tx.Exec(ctx, `INSERT INTO contact_attempts(episode_id,channel,outcome,language,attempted_at,recorded_by) SELECT episode_id,'voice',$1,'en',created_at,requested_by FROM call_attempts WHERE id=$2`, contact, id)
	}
	if err == nil && outcome == "opted_out" {
		_, err = tx.Exec(ctx, `INSERT INTO consents(patient_id,channel,status,recorded_at,source) VALUES($1,'voice','withdrawn',now(),'Samvaad patient-reported opt-out')`, patient)
	}
	if err == nil {
		_, err = tx.Exec(ctx, `INSERT INTO tasks(episode_id,assigned_to,task_type,summary) SELECT id,owner_id,'contact_patient',$1 FROM episodes WHERE id=$2`, "Review patient-reported call outcome: "+outcome, episode)
	}
	if err == nil && outcome == "barrier" {
		_, err = tx.Exec(ctx, `INSERT INTO barriers(episode_id,category,description,recorded_by) SELECT episode_id,'other','Patient reported a practical barrier during a voice call; coordinator review required.',requested_by FROM call_attempts WHERE id=$1`, id)
	}
	if err == nil && outcome == "booked" {
		_, err = tx.Exec(ctx, `UPDATE episodes SET status='booked',updated_at=now() WHERE id=$1 AND status IN ('ordered','outreach_pending','awaiting_booking')`, episode)
	}
	if err == nil {
		err = audit(ctx, tx, "", episode, "samvaad_call_"+outcome)
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeUnavailable(c, "callback_unavailable", err)
		return
	}
	c.Status(204)
}
func (r *Router) reconcile(c *gin.Context) {
	var in struct {
		Note    string `json:"note" binding:"required,max=240"`
		Attempt string `json:"attemptId" binding:"max=200"`
	}
	if c.ShouldBindJSON(&in) != nil || strings.TrimSpace(in.Note) == "" || !uuidPattern.MatchString(c.Param("id")) || !uuidPattern.MatchString(c.Param("call")) {
		fail(c, 400, "Record what you verified in Samvaad before resolving the call.")
		return
	}
	ctx := c.Request.Context()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		writeUnavailable(c, "reconcile_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(ctx, `SELECT id FROM episodes WHERE id=$1 FOR UPDATE`, c.Param("id"))
	if err == nil {
		var active int
		err = tx.QueryRow(ctx, `SELECT count(*) FROM call_attempts WHERE id=$1 AND episode_id=$2 AND status IN ('requesting','pending','uncertain') AND created_at<now()-interval '2 minutes'`, c.Param("call"), c.Param("id")).Scan(&active)
		if err == nil && active == 0 {
			fail(c, 409, "Wait at least two minutes and verify that the provider call is terminal before reconciling.")
			return
		}
	}
	if err == nil {
		_, err = tx.Exec(ctx, `UPDATE call_attempts SET attempt_id=COALESCE(attempt_id,NULLIF($1,'')),status='reconciled',finished_at=now(),reconciliation_note=$2 WHERE id=$3 AND episode_id=$4`, in.Attempt, in.Note, c.Param("call"), c.Param("id"))
	}
	if err == nil {
		err = audit(ctx, tx, actor(c).ID, c.Param("id"), "call_manually_reconciled")
	}
	if err == nil {
		err = tx.Commit(ctx)
	}
	if err != nil {
		writeUnavailable(c, "reconcile_unavailable", err)
		return
	}
	c.Status(204)
}
