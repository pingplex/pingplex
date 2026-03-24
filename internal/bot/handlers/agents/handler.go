package agents

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/pingplex/pingplex/internal/agents"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/users"
	"go.uber.org/zap"
)

const (
	listEmptyText = "У вас пока нет зарегистрированных агентов."
	failureText   = "Не удалось получить список агентов. Попробуйте позже."
)

var (
	errUserIsNil = errors.New("user is nil")
)

type Handler struct {
	users  *users.Service
	agents *agents.Service
	logger *zap.Logger
}

func New(users *users.Service, agents *agents.Service, logger *zap.Logger) handler.Handler {
	return &Handler{
		users:  users,
		agents: agents,
		logger: logger,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleAgents,
		th.CommandEqual("agents"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleAgents(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		h.logger.Warn("received /agents update without message sender")
		return nil
	}

	identity, err := identityFromTelegramUser(update.Message.From)
	if err != nil {
		h.logger.Error("failed to map telegram user to identity", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	user, err := h.users.RegisterOrLogin(ctx, *identity)
	if err != nil {
		h.logger.Error("failed to register or login telegram user", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	ownedAgents, err := h.agents.ListOwned(ctx, user.ID)
	if err != nil {
		h.logger.Error("failed to list user agents", zap.String("user_id", user.ID), zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	if len(ownedAgents) == 0 {
		return h.reply(ctx, update.Message.Chat.ID, listEmptyText)
	}

	lines := make([]string, 0, len(ownedAgents)+1)
	lines = append(lines, "Ваши агенты:")

	for _, agent := range ownedAgents {
		lines = append(lines, fmt.Sprintf("• %s (%s)", agent.Name, agent.Visibility))
	}

	return h.reply(ctx, update.Message.Chat.ID, strings.Join(lines, "\n"))
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

func identityFromTelegramUser(user *telego.User) (*users.Identity, error) {
	if user == nil {
		return nil, errUserIsNil
	}

	return &users.Identity{
		Provider:     users.ProviderTelegram,
		ProviderID:   strconv.FormatInt(user.ID, 10),
		ProviderData: "",
	}, nil
}
