package usecase

import (
	"math/rand"

	"github.com/kosumoff/secret-santa-bot/internal/domain"
)

// Draw assigns a recipient to each participant for the gift exchange.
func Draw(participants []domain.Participant) ([]domain.Assignment, error) {
	if len(participants) < 3 {
		return nil, ErrNotEnoughParticipants
	}

	shuffled := make([]domain.Participant, len(participants))
	copy(shuffled, participants)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	assignments := make([]domain.Assignment, len(shuffled))
	for i := range shuffled {
		assignments[i] = domain.Assignment{
			Giver:    shuffled[i],
			Receiver: shuffled[(i+1)%len(shuffled)],
		}
	}
	return assignments, nil
}
