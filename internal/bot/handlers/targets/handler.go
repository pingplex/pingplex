package targets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/pingplex/pingplex/internal/bot/handler"
	"github.com/pingplex/pingplex/internal/bot/middlewares/userauth"
	"github.com/pingplex/pingplex/internal/targets"
	"go.uber.org/zap"
)

const (
	menuHeaderText        = "*Управление целями мониторинга*"
	listEmptyText         = "*У вас пока нет зарегистрированных целей мониторинга.*"
	failureText           = "*Не удалось обработать команду для целей. Попробуйте позже.*"
	targetCreateUsageText = "Использование: `/target_create <name> <endpoint> [visibility] [param=value ...]`"
	targetDeleteUsageText = "Использование: `/target_delete <target_id>`"

	targetCreateSuccessFmt = "Цель создана: *%s* (%s)\n**ID: `%s`**\nПараметры: `%s`"
	targetDeleteSuccessFmt = "Цель удалена: *%s*"

	callbackPrefix     = "targets:"
	callbackList       = callbackPrefix + "list"
	callbackCreateHelp = callbackPrefix + "create_help"
	callbackDeleteHelp = callbackPrefix + "delete_help"
	callbackBack       = callbackPrefix + "menu"

	createHelpText = "Создать цель:\n`/target_create <name> <endpoint> [visibility] [param=value ...]`\n\nПримеры:\n`/target_create api-main https://example.com/health public timeout_seconds=5 retries=2`\n`/target_create dns-api https://example.com private dns_expected_ip=1.1.1.1`"
	deleteHelpText = "Удалить цель:\n`/target_delete <target_id>`\n\nID можно взять из списка /targets или меню «📋 Список целей»."
)

type Handler struct {
	targetsSvc *targets.Service

	logger *zap.Logger
}

func New(targets *targets.Service, logger *zap.Logger) handler.Handler {
	return &Handler{targetsSvc: targets, logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(h.handleTargets, th.CommandEqual("targets"), th.AnyMessageWithFrom())
	router.Handle(h.handleTargetCreate, th.CommandEqual("target_create"), th.AnyMessageWithFrom())
	router.Handle(h.handleTargetDelete, th.CommandEqual("target_delete"), th.AnyMessageWithFrom())
	router.HandleCallbackQuery(
		h.handleTargetsCallback,
		th.AnyCallbackQueryWithMessage(),
		th.CallbackDataPrefix(callbackPrefix),
	)
}

func (h *Handler) handleTargets(ctx *th.Context, update telego.Update) error {
	_, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	return h.replyWithMenu(ctx, update.Message.Chat.ID, menuHeaderText)
}

func (h *Handler) handleTargetsCallback(ctx *th.Context, query telego.CallbackQuery) error {
	if query.Message == nil {
		return nil
	}

	chatID := query.Message.GetChat().ID
	messageID := query.Message.GetMessageID()

	if err := h.answerCallback(ctx, query.ID); err != nil {
		h.logger.Warn("failed to answer callback", zap.Error(err))
	}

	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return h.editMessage(ctx, chatID, messageID, failureText, backKeyboard())
	}

	switch query.Data {
	case callbackList:
		items, listErr := h.targetsSvc.ListOwned(ctx, user.ID)
		if listErr != nil {
			h.logger.Error("failed to list targets from callback", zap.String("user_id", user.ID), zap.Error(listErr))
			return h.editMessage(ctx, chatID, messageID, failureText, menuKeyboard())
		}

		if len(items) == 0 {
			return h.editMessage(ctx, chatID, messageID, listEmptyText, backKeyboard())
		}

		lines := make([]string, 0, len(items)+1)
		lines = append(lines, "**Ваши цели:**")
		for _, t := range items {
			lines = append(
				lines,
				fmt.Sprintf(
					"• *%s* \\[`%s`] (*%s*) params=`%s`",
					t.Name,
					t.ID,
					t.Visibility,
					formatParams(t.CheckParams),
				),
			)
		}

		return h.editMessage(ctx, chatID, messageID, strings.Join(lines, "\n"), backKeyboard())
	case callbackCreateHelp:
		return h.editMessage(ctx, chatID, messageID, createHelpText, backKeyboard())
	case callbackDeleteHelp:
		return h.editMessage(ctx, chatID, messageID, deleteHelpText, backKeyboard())
	case callbackBack:
		return h.editMessage(ctx, chatID, messageID, menuHeaderText, menuKeyboard())
	default:
		return h.editMessage(ctx, chatID, messageID, menuHeaderText, menuKeyboard())
	}
}

func (h *Handler) handleTargetCreate(ctx *th.Context, update telego.Update) error {
	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	chatID := update.Message.Chat.ID

	const minArgs = 2
	_, _, args := tu.ParseCommand(update.Message.Text)
	if len(args) < minArgs {
		return h.reply(ctx, chatID, targetCreateUsageText)
	}

	visibility := targets.VisibilityPrivate
	if len(args) >= minArgs+1 {
		visibility = targets.Visibility(args[2])
	}

	const keyValueParts = 2
	params := map[string]string{}
	for _, arg := range args[3:] {
		parts := strings.SplitN(arg, "=", keyValueParts)
		if len(parts) != keyValueParts || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return h.reply(ctx, chatID, targetCreateUsageText)
		}

		params[parts[0]] = parts[1]
	}

	created, err := h.targetsSvc.Create(ctx, user.ID, targets.CreateTargetInput{
		Name:        args[0],
		Endpoint:    args[1],
		Description: "",
		Visibility:  visibility,
		CheckParams: params,
	})
	if err != nil {
		h.logger.Error("failed to create target", zap.String("user_id", user.ID), zap.Error(err))
		if errors.Is(err, targets.ErrValidation) {
			return h.reply(ctx, chatID, targetCreateUsageText)
		}
		return h.reply(ctx, chatID, failureText)
	}

	return h.reply(
		ctx,
		chatID,
		fmt.Sprintf(
			targetCreateSuccessFmt,
			created.Name,
			created.Visibility,
			created.ID,
			formatParams(created.CheckParams),
		),
	)
}

func (h *Handler) handleTargetDelete(ctx *th.Context, update telego.Update) error {
	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, failureText)
	}

	chatID := update.Message.Chat.ID

	_, _, args := tu.ParseCommand(update.Message.Text)
	if len(args) < 1 {
		return h.reply(ctx, chatID, targetDeleteUsageText)
	}

	if err = h.targetsSvc.Delete(ctx, user.ID, args[0]); err != nil {
		h.logger.Error(
			"failed to delete target",
			zap.String("user_id", user.ID),
			zap.String("target_id", args[0]),
			zap.Error(err),
		)
		if errors.Is(err, targets.ErrNotFound) {
			return h.reply(ctx, chatID, "*Цель не найдена.*")
		}
		return h.reply(ctx, chatID, failureText)
	}

	return h.reply(ctx, chatID, fmt.Sprintf(targetDeleteSuccessFmt, args[0]))
}

func (h *Handler) replyWithMenu(ctx *th.Context, chatID int64, text string) error {
	return h.replyWithMarkup(ctx, chatID, text, menuKeyboard())
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	return h.replyWithMarkup(ctx, chatID, text, nil)
}

func (h *Handler) replyWithMarkup(
	ctx *th.Context,
	chatID int64,
	text string,
	keyboard *telego.InlineKeyboardMarkup,
) error {
	params := &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID, Username: ""},
		Text:      text,
		ParseMode: telego.ModeMarkdown,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard
	}

	_, err := ctx.Bot().SendMessage(ctx, params)
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}

func (h *Handler) editMessage(
	ctx *th.Context,
	chatID int64,
	messageID int,
	text string,
	keyboard *telego.InlineKeyboardMarkup,
) error {
	params := &telego.EditMessageTextParams{
		ChatID:    telego.ChatID{ID: chatID, Username: ""},
		MessageID: messageID,
		Text:      text,
		ParseMode: telego.ModeMarkdown,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard
	}

	_, err := ctx.Bot().EditMessageText(ctx, params)
	if err != nil {
		return fmt.Errorf("edit telegram message: %w", err)
	}

	return nil
}

func (h *Handler) answerCallback(ctx *th.Context, callbackID string) error {
	err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{CallbackQueryID: callbackID})
	if err != nil {
		return fmt.Errorf("answer callback query: %w", err)
	}

	return nil
}

func menuKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📋 Список целей").WithCallbackData(callbackList),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Как создать").WithCallbackData(callbackCreateHelp),
			tu.InlineKeyboardButton("🗑 Как удалить").WithCallbackData(callbackDeleteHelp),
		),
	)
}

func backKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅️ Назад в меню").WithCallbackData(callbackBack),
		),
	)
}

func formatParams(params map[string]string) string {
	if len(params) == 0 {
		return "-"
	}

	parts := make([]string, 0, len(params))
	for key, value := range params {
		parts = append(parts, key+"="+value)
	}

	return strings.Join(parts, ", ")
}
