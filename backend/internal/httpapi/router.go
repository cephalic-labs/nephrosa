package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) http.Handler {
	router := &Router{db: pool}
	engine := gin.New()
	engine.Use(gin.Recovery(), responseHeaders())
	_ = engine.SetTrustedProxies(nil)

	engine.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/health/ready", router.ready)
	engine.GET("/api/v1/dashboard/overview", router.dashboardOverview)
	return engine
}

func responseHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
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
	WHERE e.status <> 'closed'
	GROUP BY e.id, p.external_reference, p.display_name, u.display_name
	ORDER BY e.due_at, p.display_name
	LIMIT 10
	`)
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
