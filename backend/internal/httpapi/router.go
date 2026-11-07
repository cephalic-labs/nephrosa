package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cephalic-labs/nephrosa/backend/internal/sarvam"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	db    *pgxpool.Pool
	voice *sarvam.Client
}

func New(pool *pgxpool.Pool, voice *sarvam.Client) http.Handler {
	router := &Router{db: pool, voice: voice}
	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, _ any) { fail(c, 500, "The request could not be completed.") }), responseHeaders(), func(c *gin.Context) {
		limit := int64(16384)
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/sarvam/webhooks/") {
			limit = 1 << 20
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	})
	_ = engine.SetTrustedProxies(nil)

	engine.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/health/ready", router.ready)
	engine.POST("/api/v1/sessions", router.login)
	engine.POST("/api/v1/sarvam/webhooks/:secret", router.webhook)
	api := engine.Group("/api/v1", router.authenticate)
	api.GET("/session", func(c *gin.Context) { c.JSON(200, actor(c)) })
	api.DELETE("/session", router.logout)
	api.GET("/users", router.users)
	api.GET("/dashboard/overview", router.dashboardOverview)
	api.GET("/episodes/:id", router.episodeDetail)
	api.POST("/episodes", router.createEpisode)
	api.POST("/episodes/:id/actions", router.workflow)
	api.POST("/episodes/:id/calls", router.startCall)
	api.GET("/episodes/:id/calls", router.callList)
	api.POST("/episodes/:id/calls/:call/reconcile", router.reconcile)
	return engine
}

func responseHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		c.Next()
	}
}

func (r *Router) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := r.db.Ping(ctx); err != nil {
		slog.Error("database readiness check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

type DashboardOverview struct {
	Summary  DashboardSummary   `json:"summary"`
	Episodes []DashboardEpisode `json:"episodes"`
}

type DashboardSummary struct {
	OpenEpisodes            int `json:"openEpisodes"`
	OverdueEpisodes         int `json:"overdueEpisodes"`
	AwaitingReport          int `json:"awaitingReport"`
	AwaitingClinicianReview int `json:"awaitingClinicianReview"`
	DueWithinSevenDays      int `json:"dueWithinSevenDays"`
}

type DashboardEpisode struct {
	ID                  string    `json:"id"`
	PatientReference    string    `json:"patientReference"`
	PatientDisplayName  string    `json:"patientDisplayName"`
	Status              string    `json:"status"`
	DueAt               time.Time `json:"dueAt"`
	OwnerDisplayName    string    `json:"ownerDisplayName"`
	Tests               []string  `json:"tests"`
	OpenBarrierCategory string    `json:"openBarrierCategory,omitempty"`
	NextTask            string    `json:"nextTask,omitempty"`
}

type EpisodeDetail struct {
	ID                 string           `json:"id"`
	PatientReference   string           `json:"patientReference"`
	PatientDisplayName string           `json:"patientDisplayName"`
	Status             string           `json:"status"`
	DueAt              time.Time        `json:"dueAt"`
	OwnerDisplayName   string           `json:"ownerDisplayName"`
	Tests              []EpisodeTest    `json:"tests"`
	Contacts           []EpisodeContact `json:"contacts"`
	Barriers           []EpisodeBarrier `json:"barriers"`
	Tasks              []EpisodeTask    `json:"tasks"`
	Reports            []EpisodeReport  `json:"reports"`
	ClinicalReviews    []EpisodeReview  `json:"clinicalReviews"`
	Consent            string           `json:"consent"`
	Language           string           `json:"language"`
	Phone              string           `json:"phone"`
	AuditEvents        []AuditEvent     `json:"auditEvents"`
}

type EpisodeTest struct {
	ID         string          `json:"id"`
	TestType   string          `json:"testType"`
	OrderedAt  time.Time       `json:"orderedAt"`
	DueAt      time.Time       `json:"dueAt"`
	Completion *TestCompletion `json:"completion,omitempty"`
}

type TestCompletion struct {
	CompletedAt        time.Time `json:"completedAt"`
	VerificationSource string    `json:"verificationSource"`
	RecordedBy         string    `json:"recordedBy"`
}

type EpisodeContact struct {
	Channel     string    `json:"channel"`
	Outcome     string    `json:"outcome"`
	Language    string    `json:"language"`
	AttemptedAt time.Time `json:"attemptedAt"`
	RecordedBy  string    `json:"recordedBy"`
}

type EpisodeBarrier struct {
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	ResolvedAt  *time.Time `json:"resolvedAt,omitempty"`
}

type EpisodeTask struct {
	ID                    string     `json:"id"`
	Type                  string     `json:"type"`
	Summary               string     `json:"summary"`
	Priority              string     `json:"priority"`
	Status                string     `json:"status"`
	DueAt                 *time.Time `json:"dueAt,omitempty"`
	AssignedTo            string     `json:"assignedTo,omitempty"`
	AssignedToDisplayName string     `json:"assignedToDisplayName,omitempty"`
}

type EpisodeReport struct {
	ReceivedAt time.Time `json:"receivedAt"`
	Reference  string    `json:"reference"`
	RecordedBy string    `json:"recordedBy"`
}

type EpisodeReview struct {
	AcknowledgedAt time.Time `json:"acknowledgedAt"`
	Clinician      string    `json:"clinician"`
	NextAction     string    `json:"nextAction"`
}

type AuditEvent struct {
	EventType  string    `json:"eventType"`
	OccurredAt time.Time `json:"occurredAt"`
	Actor      string    `json:"actor"`
	Note       string    `json:"note,omitempty"`
}

func (r *Router) episodeDetail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	defer tx.Rollback(context.Background())

	var detail EpisodeDetail
	err = tx.QueryRow(ctx, `
		SELECT e.id::text, p.external_reference, p.display_name, e.status, e.due_at, u.display_name
		FROM episodes e
		JOIN patients p ON p.id = e.patient_id
		JOIN users u ON u.id = e.owner_id
		WHERE e.id::text = $1
	`, c.Param("id")).Scan(
		&detail.ID,
		&detail.PatientReference,
		&detail.PatientDisplayName,
		&detail.Status,
		&detail.DueAt,
		&detail.OwnerDisplayName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "episode_not_found", "message": "Episode not found."}})
		return
	}
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}

	err = tx.QueryRow(ctx, `SELECT p.preferred_language,COALESCE(p.phone_number,''),COALESCE((SELECT status FROM consents WHERE patient_id=p.id AND channel='voice' ORDER BY created_at DESC,id DESC LIMIT 1),'pending') FROM patients p JOIN episodes e ON e.patient_id=p.id WHERE e.id::text=$1`, detail.ID).Scan(&detail.Language, &detail.Phone, &detail.Consent)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	detail.Tests = make([]EpisodeTest, 0)
	rows, err := tx.Query(ctx, `
		SELECT o.id::text, o.test_type, o.ordered_at, o.due_at,
		       c.completed_at, COALESCE(c.verification_source, ''), COALESCE(cu.display_name, '')
		FROM test_orders o
		LEFT JOIN test_completions c ON c.test_order_id = o.id
		LEFT JOIN users cu ON cu.id = c.recorded_by
		WHERE o.episode_id::text = $1
		ORDER BY o.test_type
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var test EpisodeTest
		var completedAt *time.Time
		var verificationSource string
		var recordedBy string
		if err := rows.Scan(&test.ID, &test.TestType, &test.OrderedAt, &test.DueAt, &completedAt, &verificationSource, &recordedBy); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		if completedAt != nil {
			test.Completion = &TestCompletion{CompletedAt: *completedAt, VerificationSource: verificationSource, RecordedBy: recordedBy}
		}
		detail.Tests = append(detail.Tests, test)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.Contacts = make([]EpisodeContact, 0)
	rows, err = tx.Query(ctx, `
		SELECT c.channel, c.outcome, c.language, c.attempted_at, u.display_name
		FROM contact_attempts c JOIN users u ON u.id = c.recorded_by
		WHERE c.episode_id::text = $1
		ORDER BY c.attempted_at DESC
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var contact EpisodeContact
		if err := rows.Scan(&contact.Channel, &contact.Outcome, &contact.Language, &contact.AttemptedAt, &contact.RecordedBy); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.Contacts = append(detail.Contacts, contact)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.Barriers = make([]EpisodeBarrier, 0)
	rows, err = tx.Query(ctx, `
		SELECT category, description, status, created_at, resolved_at
		FROM barriers WHERE episode_id::text = $1 ORDER BY created_at DESC
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var barrier EpisodeBarrier
		if err := rows.Scan(&barrier.Category, &barrier.Description, &barrier.Status, &barrier.CreatedAt, &barrier.ResolvedAt); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.Barriers = append(detail.Barriers, barrier)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.Tasks = make([]EpisodeTask, 0)
	rows, err = tx.Query(ctx, `
		SELECT t.id::text, t.task_type, t.summary, t.priority, t.status, t.due_at,
		       COALESCE(t.assigned_to::text, ''), COALESCE(u.display_name, '')
		FROM tasks t LEFT JOIN users u ON u.id = t.assigned_to
		WHERE t.episode_id::text = $1 ORDER BY t.created_at
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var task EpisodeTask
		if err := rows.Scan(&task.ID, &task.Type, &task.Summary, &task.Priority, &task.Status, &task.DueAt, &task.AssignedTo, &task.AssignedToDisplayName); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.Tasks = append(detail.Tasks, task)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.Reports = make([]EpisodeReport, 0)
	rows, err = tx.Query(ctx, `
		SELECT r.received_at, r.reference, u.display_name
		FROM reports r JOIN users u ON u.id = r.recorded_by
		WHERE r.episode_id::text = $1 ORDER BY r.received_at DESC
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var report EpisodeReport
		if err := rows.Scan(&report.ReceivedAt, &report.Reference, &report.RecordedBy); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.Reports = append(detail.Reports, report)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.ClinicalReviews = make([]EpisodeReview, 0)
	rows, err = tx.Query(ctx, `
		SELECT cr.acknowledged_at, u.display_name, cr.next_action
		FROM clinical_reviews cr JOIN users u ON u.id = cr.clinician_id
		WHERE cr.episode_id::text = $1 ORDER BY cr.acknowledged_at DESC
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var review EpisodeReview
		if err := rows.Scan(&review.AcknowledgedAt, &review.Clinician, &review.NextAction); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.ClinicalReviews = append(detail.ClinicalReviews, review)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	detail.AuditEvents = make([]AuditEvent, 0)
	rows, err = tx.Query(ctx, `
		SELECT ae.event_type, ae.occurred_at, COALESCE(u.display_name, 'System'),COALESCE(ae.metadata->>'note','')
		FROM audit_events ae LEFT JOIN users u ON u.id = ae.actor_id
		WHERE ae.episode_id::text = $1 ORDER BY ae.occurred_at DESC
	`, detail.ID)
	if err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	for rows.Next() {
		var event AuditEvent
		if err := rows.Scan(&event.EventType, &event.OccurredAt, &event.Actor, &event.Note); err != nil {
			rows.Close()
			writeUnavailable(c, "episode_detail_unavailable", err)
			return
		}
		detail.AuditEvents = append(detail.AuditEvents, event)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	rows.Close()

	if err := tx.Commit(ctx); err != nil {
		writeUnavailable(c, "episode_detail_unavailable", err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (r *Router) dashboardOverview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var overview DashboardOverview
	err := r.db.QueryRow(ctx, `
		SELECT
			count(*) FILTER (WHERE status <> 'closed'),
			count(*) FILTER (WHERE status <> 'closed' AND due_at < CURRENT_DATE),
			count(*) FILTER (WHERE status = 'awaiting_report'),
			count(*) FILTER (WHERE status = 'awaiting_clinician_review'),
			count(*) FILTER (
				WHERE status <> 'closed'
				  AND due_at >= CURRENT_DATE
				  AND due_at <= CURRENT_DATE + 7
			)
		FROM episodes
	`).Scan(
		&overview.Summary.OpenEpisodes,
		&overview.Summary.OverdueEpisodes,
		&overview.Summary.AwaitingReport,
		&overview.Summary.AwaitingClinicianReview,
		&overview.Summary.DueWithinSevenDays,
	)
	if err != nil {
		writeUnavailable(c, "dashboard_summary_unavailable", err)
		return
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			e.id::text,
			p.external_reference,
			p.display_name,
			e.status,
			e.due_at,
			u.display_name,
			COALESCE(
				array_agg(DISTINCT o.test_type ORDER BY o.test_type)
					FILTER (WHERE o.id IS NOT NULL),
				ARRAY[]::text[]
			),
			COALESCE((
				SELECT b.category
				FROM barriers b
				WHERE b.episode_id = e.id AND b.status = 'open'
				ORDER BY b.created_at DESC
				LIMIT 1
			), ''),
			COALESCE((
				SELECT t.summary
				FROM tasks t
				WHERE t.episode_id = e.id AND t.status IN ('open', 'in_progress')
				ORDER BY t.due_at NULLS LAST, t.created_at
				LIMIT 1
			), '')
		FROM episodes e
		JOIN patients p ON p.id = e.patient_id
		JOIN users u ON u.id = e.owner_id
		LEFT JOIN test_orders o ON o.episode_id = e.id
	WHERE ($1::boolean OR e.status <> 'closed')
	GROUP BY e.id, p.external_reference, p.display_name, u.display_name
	ORDER BY e.due_at, p.display_name
	LIMIT 100
	`, c.Query("history") == "true")
	if err != nil {
		writeUnavailable(c, "dashboard_episodes_unavailable", err)
		return
	}
	defer rows.Close()

	overview.Episodes = make([]DashboardEpisode, 0, 10)
	for rows.Next() {
		var episode DashboardEpisode
		var barrierCategory string
		var nextTask string
		if err := rows.Scan(
			&episode.ID,
			&episode.PatientReference,
			&episode.PatientDisplayName,
			&episode.Status,
			&episode.DueAt,
			&episode.OwnerDisplayName,
			&episode.Tests,
			&barrierCategory,
			&nextTask,
		); err != nil {
			writeUnavailable(c, "dashboard_episodes_unavailable", err)
			return
		}
		episode.OpenBarrierCategory = barrierCategory
		episode.NextTask = nextTask
		overview.Episodes = append(overview.Episodes, episode)
	}
	if err := rows.Err(); err != nil {
		writeUnavailable(c, "dashboard_episodes_unavailable", err)
		return
	}

	c.JSON(http.StatusOK, overview)
}

func writeUnavailable(c *gin.Context, code string, err error) {
	slog.Error("dashboard query failed", "code", code, "error", err)
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"code":    code,
			"message": "Dashboard data is temporarily unavailable.",
		},
	})
}
