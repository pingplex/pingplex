package start

import (
	"fmt"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/bot/middlewares/userauth"
	"go.uber.org/zap"
)

const (
	welcomeText = "Welcome! Your account is ready."
	failureText = "Something went wrong. Please try again later."
)

type Handler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) handler.Handler {
	return &Handler{
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

	_, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	return h.reply(ctx, update.Message.Chat.ID, welcomeText)
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text))
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}
