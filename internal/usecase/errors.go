package usecase

import "errors"

var (
	ErrNotEnoughParticipants = errors.New("at least 3 participants are required")
)
