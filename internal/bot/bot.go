package bot

import (
	"context"
	"errors"
	"fmt"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"strings"
	"sync"

	telebot "gopkg.in/telebot.v4"
)

type TelegramBot struct {
	Bot            *telebot.Bot
	Config         config.TelegramBotConfig
	Logger         *logger.Logger
	UserRepository repositories.UserRepository
}

func NewTelegramBot(cfg config.TelegramBotConfig, l *logger.Logger, ur repositories.UserRepository) *TelegramBot {
	pref := telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: cfg.Timeout},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		panic(fmt.Errorf("failed to create new telegram bot: %w", err))
	}
	return &TelegramBot{
		Bot:            bot,
		Config:         cfg,
		UserRepository: ur,
		Logger:         l,
	}
}

const place = "telegramBot"

func (tb *TelegramBot) StoreChatId() {
	op := place + "StoreChatId"
	log := tb.Logger.AddOp(op)
	log.Info("storing chat id")
	yes := telebot.Btn{Unique: "yes", Text: "✅ да"}
	no := telebot.Btn{Unique: "no", Text: "❌ Нет"}

	tb.Bot.Handle("/start", func(c telebot.Context) error {
		markup := &telebot.ReplyMarkup{}
		markup.Inline(markup.Row(yes, no))
		return c.Send("хотите получать обновления цены айфончика 17??", markup)
	})
	tb.Bot.Handle(&yes, func(c telebot.Context) error {
		chatId := c.Chat().ID
		username := c.Sender().Username
		ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
		defer cancel()
		if err := tb.UserRepository.SetChatId(ctx, username, chatId); err != nil {
			if errors.Is(err, errs.ErrAlreadyExistsBase) {
				return c.Edit("вы уже получаете обновления")
			}
			log.Error("failed to store chat id", logger.Err(err))
			return c.Edit("произошла ошибка, возможно вы не зарегестрированы((")
		}
		return c.Edit("ждите обновления))")
	})
	tb.Bot.Handle(&no, func(c telebot.Context) error {
		chatId := c.Chat().ID
		username := c.Sender().Username
		ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
		defer cancel()
		if err := tb.UserRepository.DropChatId(ctx, username, chatId); err != nil {
			if errors.Is(err, errs.ErrNotFoundBase) {
				return c.Edit("вы не подписаны на обновления")
			}
			log.Error("failed to drop chat id", logger.Err(err))
			return c.Edit("произошла ошибка((")
		}
		return c.Edit("обновлений не ждите((")
	})
}

func (tb *TelegramBot) SendIPhonesInfo(chatIds []int64, iphones []models.IPhone) error {
	op := place + "SendIphoneInfo"
	log := tb.Logger.AddOp(op)
	log.Info("sending iphone info")
	msgArr := []string{}
	for _, iphone := range iphones {
		msgArr = append(msgArr, fmt.Sprintf("%s:\n 💰 цена: %.2f | 📉 разница: %.2f\n", iphone.Name, iphone.Price, iphone.Change))
	}
	msg := strings.Join(msgArr, "\n")
	errChan := make(chan error, len(chatIds))
	var wg sync.WaitGroup
	for _, cid := range chatIds {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			if _, err := tb.Bot.Send(&telebot.Chat{ID: id}, msg, telebot.ModeMarkdown); err != nil {
				errChan <- err
			}
		}(cid)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	if len(errChan) > 0 {
		for err := range errChan {
			log.Error("failed to send iphones info", logger.Err(err))
			return errs.NewAppError(op, err)
		}
	}
	return nil
}

func (tb *TelegramBot) Start() {
	tb.Bot.Start()
}

func (tb *TelegramBot) Stop() {
	tb.Bot.Stop()
}
