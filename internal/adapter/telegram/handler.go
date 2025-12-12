// Package telegram implements Telegram bot handlers.
package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kosumoff/secret-santa-bot/internal/domain"
	"github.com/kosumoff/secret-santa-bot/internal/usecase"
)

// Handler handles Telegram bot updates.
type Handler struct {
	participantRepo usecase.ParticipantRepository
	assignmentRepo  usecase.AssignmentRepository
}

func NewHandler(pr usecase.ParticipantRepository, ar usecase.AssignmentRepository) *Handler {
	return &Handler{
		participantRepo: pr,
		assignmentRepo:  ar,
	}
}

func (h *Handler) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/santa", bot.MatchTypeExact, h.santa)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/draw", bot.MatchTypeExact, h.draw)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, h.start)
}

func (h *Handler) santa(ctx context.Context, b *bot.Bot, upd *models.Update) {
	chatID := upd.Message.Chat.ID

	payload := fmt.Sprintf("santa_%d", chatID)

	botUser, err := b.GetMe(ctx)
	if err != nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   "Tap the button below to join the exchange!",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Join Secret Santa",
						URL:  "https://t.me/" + botUser.Username + "?start=" + payload},
				},
			},
		},
	})
}

func (h *Handler) start(ctx context.Context, b *bot.Bot, upd *models.Update) {
	if upd.Message.Chat.Type != "private" {
		return
	}

	parts := strings.Split(upd.Message.Text, " ")
	if len(parts) < 2 {
		return
	}

	payload := parts[1]
	if !strings.HasPrefix(payload, "santa_") {
		return
	}

	chatID, _ := strconv.ParseInt(strings.TrimPrefix(payload, "santa_"), 10, 64)
	user := domain.Participant{
		ChatID:   chatID,
		UserID:   upd.Message.From.ID,
		Username: displayName(upd.Message.From),
	}

	_ = h.participantRepo.Add(user)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: upd.Message.Chat.ID,
		Text:   "You're in! I'll message you once the draw is complete.",
	})
}

func (h *Handler) draw(ctx context.Context, b *bot.Bot, upd *models.Update) {
	chatID := upd.Message.Chat.ID
	userID := upd.Message.From.ID

	isAdmin, err := IsAdmin(ctx, b, chatID, userID)
	if err != nil {
		return
	}

	if !isAdmin {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Only group administrators can start the draw.",
		})
		return
	}

	users, _ := h.participantRepo.GetByChat(chatID)

	assignments, err := usecase.Draw(users)
	if err != nil {
		if errors.Is(err, usecase.ErrNotEnoughParticipants) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "We need at least 3 people to start the Secret Santa.",
			})
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Something went wrong. Please try again later.",
			})
		}
		return
	}

	_ = h.assignmentRepo.Save(assignments)

	for _, a := range assignments {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: a.Giver.UserID,
			Text:   "Your Secret Santa recipient is: " + a.Receiver.Username,
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "The draw is complete! Everyone has received their assignments via private message.",
	})
}

func displayName(u *models.User) string {
	if u.Username != "" {
		return "@" + u.Username
	}
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// IsAdmin checks if a user is an admin in a chat.
func IsAdmin(ctx context.Context, b *bot.Bot, chatID int64, userID int64) (bool, error) {

	member, err := b.GetChatMember(ctx, &bot.GetChatMemberParams{
		ChatID: chatID,
		UserID: userID,
	})
	if err != nil {
		return false, err
	}

	switch member.Type {
	case models.ChatMemberTypeOwner,
		models.ChatMemberTypeAdministrator:
		return true, nil
	default:
		return false, nil
	}
}
