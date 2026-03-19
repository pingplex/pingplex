package start

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/users"
	"go.uber.org/zap"
)

const (
	welcomeText = "Welcome! Your account is ready."
	failureText = "Something went wrong. Please try again in a moment."
)

var (
	errUserIsNil = errors.New("user is nil")
)

type Handler struct {
	users  *users.Service
	logger *zap.Logger
}

func New(users *users.Service, logger *zap.Logger) handler.Handler {
	return &Handler{
		users:  users,
		logger: logger,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleStart,
		th.CommandEqual("start"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleStart(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		h.logger.Warn("received /start update without message sender")
		return nil
	}

	identity, err := identityFromTelegramUser(update.Message.From)
	if err != nil {
		h.logger.Error("failed to map telegram user to identity", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	user, err := h.users.RegisterOrLogin(ctx, *identity)
	if err != nil {
		h.logger.Error(
			"failed to register or login telegram user",
			zap.String("provider", string(users.ProviderTelegram)),
			zap.Error(err),
		)
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	h.logger.Info(
		"telegram user registered or logged in",
		zap.String("user_id", user.ID),
		zap.String("provider", string(users.ProviderTelegram)),
	)

	return h.reply(ctx, update.Message.Chat.ID, welcomeText)
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: chatID, Username: ""},
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}

// identityFromTelegramUser maps a Telegram user into the users identity model.
func identityFromTelegramUser(user *telego.User) (*users.Identity, error) {
	if user == nil {
		return nil, errUserIsNil
	}

	providerData, err := toProviderDataJSON(user)
	if err != nil {
		return nil, err
	}

	return &users.Identity{
		Provider:     users.ProviderTelegram,
		ProviderID:   strconv.FormatInt(user.ID, 10),
		ProviderData: providerData,
	}, nil
}

func toProviderDataJSON(user *telego.User) (string, error) {
	payload := struct {
		Username     string `json:"username,omitempty"`
		FirstName    string `json:"first_name,omitempty"`
		LastName     string `json:"last_name,omitempty"`
		LanguageCode string `json:"language_code,omitempty"`
		IsBot        bool   `json:"is_bot"`
	}{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,
		IsBot:        user.IsBot,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal provider data: %w", err)
	}

	return string(encoded), nil
}
