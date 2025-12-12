package sqlite

import (
	"database/sql"

	"github.com/kosumoff/secret-santa-bot/internal/domain"
	"github.com/kosumoff/secret-santa-bot/internal/usecase"
)

type AssignmentRepo struct {
	db *sql.DB
}

var _ usecase.AssignmentRepository = (*AssignmentRepo)(nil)

func NewAssignmentRepo(db *sql.DB) *AssignmentRepo {
	return &AssignmentRepo{db: db}
}

func (r *AssignmentRepo) Save(assignments []domain.Assignment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, _ := tx.Prepare("INSERT INTO assignments(chat_id, giver_id, receiver_id) VALUES(?,?,?)")
	defer stmt.Close()

	for _, a := range assignments {
		if _, err := stmt.Exec(a.Giver.ChatID, a.Giver.UserID, a.Receiver.UserID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
