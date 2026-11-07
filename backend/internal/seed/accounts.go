package seed

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func accounts(ctx context.Context, tx pgx.Tx) error {
	for _, role := range []string{"clinician", "coordinator"} {
		path := os.Getenv("DEMO_" + strings.ToUpper(role) + "_PASSWORD_FILE")
		b, err := os.ReadFile(path)
		if err != nil {
			return errors.New("demo account password files are required for seeding")
		}
		password := strings.TrimRight(string(b), "\r\n")
		if len(password) < 12 || len(password) > 72 {
			return errors.New("demo passwords must contain 12–72 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE users SET login=$1,password_hash=$2 WHERE display_name=$3 AND role=$1`, role, string(hash), "Demo "+strings.ToUpper(role[:1])+role[1:])
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `DELETE FROM sessions WHERE user_id IN (SELECT id FROM users WHERE login=$1)`, role)
		if err != nil {
			return err
		}
	}
	return nil
}
