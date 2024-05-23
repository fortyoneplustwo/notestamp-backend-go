package auth

import (
	"database/sql"
	"log"
	"time"
)

type RevokedDB struct {
	db *sql.DB
}

// Constructor
func NewRevokedDB(db *sql.DB) *RevokedDB {
	return &RevokedDB{db: db}
}

// Implements revokedStore interface
func (r *RevokedDB) Add(token string, id int, exp time.Time) error {
	stmt, err := r.db.Prepare("INSERT INTO revoked (token, id, exp) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(token, id, exp); err != nil {
		return err
	}

	return nil
}

func (r *RevokedDB) IsRevoked(token string, uid int) (bool, error) {
	return false, nil
}

func (r *RevokedDB) Cleanup(id int) {
	stmt, err := r.db.Prepare("DELETE FROM revoked WHERE id = $1 AND exp < NOW()")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id); err != nil {
		log.Println(err)
	}
}
