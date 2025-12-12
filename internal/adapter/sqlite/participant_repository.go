package sqlite

import (
	"database/sql"

	"github.com/kosumoff/secret-santa-bot/internal/domain"
	"github.com/kosumoff/secret-santa-bot/internal/usecase"
)

type ParticipantRepo struct {
	db *sql.DB
}

var _ usecase.ParticipantRepository = (*ParticipantRepo)(nil)

func NewParticipantRepo(db *sql.DB) *ParticipantRepo {
	return &ParticipantRepo{db: db}
}

func (r *ParticipantRepo) Add(p domain.Participant) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO participants(chat_id, user_id, username) VALUES(?,?,?)",
		p.ChatID, p.UserID, p.Username)
	return err
}

func (r *ParticipantRepo) GetByChat(chatID int64) ([]domain.Participant, error) {
	rows, err := r.db.Query("SELECT chat_id, user_id, username FROM participants WHERE chat_id=?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		if err := rows.Scan(&p.ChatID, &p.UserID, &p.Username); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}
