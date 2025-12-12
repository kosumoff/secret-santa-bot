package usecase

import "github.com/kosumoff/secret-santa-bot/internal/domain"

type ParticipantRepository interface {
	Add(p domain.Participant) error
	GetByChat(chatID int64) ([]domain.Participant, error)
}

type AssignmentRepository interface {
	Save(assignments []domain.Assignment) error
}
