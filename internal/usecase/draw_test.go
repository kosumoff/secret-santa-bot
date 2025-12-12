package usecase

import (
	"testing"

	"github.com/kosumoff/secret-santa-bot/internal/domain"
)

func TestDraw(t *testing.T) {
	participants := []domain.Participant{
		{UserID: 1, Username: "Nurhat"},
		{UserID: 2, Username: "Azat"},
		{UserID: 3, Username: "Ainel"},
		{UserID: 4, Username: "Zhanel"},
	}

	assignments, err := Draw(participants)
	if err != nil {
		t.Fatalf("Draw returned error: %v", err)
	}

	// Check if the number of assignments matches the number of participants.
	if len(assignments) != len(participants) {
		t.Fatalf("expected %d assignments, got %d", len(participants), len(assignments))
	}

	// Ensure no one is assigned to themselves.
	for _, a := range assignments {
		if a.Giver.UserID == a.Receiver.UserID {
			t.Errorf("giver %v assigned to themselves", a.Giver.Username)
		}
	}

	// Check that all receivers are unique.
	seen := make(map[int64]bool)
	for _, a := range assignments {
		if seen[a.Receiver.UserID] {
			t.Errorf("received %v appears more than ones", a.Receiver.Username)
		}
		seen[a.Receiver.UserID] = true
	}
}
