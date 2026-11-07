package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func token() string          { return rand.Text() }
func digest(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func fail(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"message": message}})
}

// shortcut: process-wide demo login limit; use shared per-account throttling before scaling beyond one API instance.
var loginLimit struct {
	sync.Mutex
	start time.Time
	count int
}

func (r *Router) login(c *gin.Context) {
	loginLimit.Lock()
	if time.Since(loginLimit.start) > time.Minute {
		loginLimit.start = time.Now()
		loginLimit.count = 0
	}
	loginLimit.count++
	allowed := loginLimit.count <= 20
	loginLimit.Unlock()
	if !allowed {
		fail(c, 429, "Too many sign-in attempts. Try again in a minute.")
		return
	}
	var in struct {
		Login    string `json:"login" binding:"required,max=80"`
		Password string `json:"password" binding:"required,min=12,max=72"`
	}
	if c.ShouldBindJSON(&in) != nil {
		fail(c, 400, "Enter a login and a password of 12–72 characters.")
		return
	}
	var u User
	var hash string
	err := r.db.QueryRow(c.Request.Context(), `SELECT id::text,display_name,role,password_hash FROM users WHERE login=$1 AND active`, in.Login).Scan(&u.ID, &u.Name, &u.Role, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password)) != nil {
		fail(c, 401, "Invalid credentials.")
		return
	}
	t := token()
	if _, err = r.db.Exec(c.Request.Context(), `INSERT INTO sessions(token_hash,user_id,expires_at) VALUES($1,$2,now()+interval '8 hours')`, digest(t), u.ID); err != nil {
		writeUnavailable(c, "session_unavailable", err)
		return
	}
	c.JSON(200, gin.H{"token": t, "user": u})
}
func (r *Router) authenticate(c *gin.Context) {
	raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if raw == "" || raw == c.GetHeader("Authorization") {
		fail(c, 401, "Sign in to continue.")
		return
	}
	var u User
	err := r.db.QueryRow(c.Request.Context(), `SELECT u.id::text,u.display_name,u.role FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.expires_at>now() AND u.active AND u.role IN ('clinician','coordinator')`, digest(raw)).Scan(&u.ID, &u.Name, &u.Role)
	if err != nil {
		fail(c, 401, "Your session expired. Sign in again.")
		return
	}
	c.Set("user", u)
	c.Next()
}
func actor(c *gin.Context) User { return c.MustGet("user").(User) }
func clinician(c *gin.Context) bool {
	if actor(c).Role != "clinician" {
		fail(c, http.StatusForbidden, "A clinician must perform this action.")
		return false
	}
	return true
}
func (r *Router) logout(c *gin.Context) {
	_, err := r.db.Exec(c.Request.Context(), `DELETE FROM sessions WHERE token_hash=$1`, digest(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")))
	if err != nil {
		writeUnavailable(c, "logout_unavailable", err)
		return
	}
	c.Status(204)
}
func (r *Router) users(c *gin.Context) {
	rows, err := r.db.Query(c.Request.Context(), `SELECT id::text,display_name,role FROM users WHERE active AND role IN ('clinician','coordinator') ORDER BY display_name`)
	if err != nil {
		writeUnavailable(c, "users_unavailable", err)
		return
	}
	defer rows.Close()
	out := []User{}
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name, &u.Role); err != nil {
			writeUnavailable(c, "users_unavailable", err)
			return
		}
		out = append(out, u)
	}
	if rows.Err() != nil {
		writeUnavailable(c, "users_unavailable", rows.Err())
		return
	}
	c.JSON(200, out)
}
