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
	"strconv"
	"strings"
	"sync"

	telebot "gopkg.in/telebot.v4"
)

//go:generate mockgen -source=bot.go -destination=mocks/bot-mock.go
type TelegramBot interface {
	SetupTelegramBot()
	SendIPhonesInfo(chatIds []int64, iphones []models.IPhone) error
	Start()
	Stop()
}

type telegramBot struct {
	Bot             *telebot.Bot
	Config          config.TelegramBotConfig
	UserRepository  repositories.UserRepository
	usersStatements sync.Map
	Logger          *logger.Logger
}

func NewTelegramBot(cfg config.TelegramBotConfig, l *logger.Logger, ur repositories.UserRepository) TelegramBot {
	pref := telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: cfg.Timeout},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		panic(fmt.Errorf("failed to create new telegram bot: %w", err))
	}
	return &telegramBot{
		Bot:            bot,
		Config:         cfg,
		UserRepository: ur,
		Logger:         l,
	}
}

const place = "telegramBot."

func (tb *telegramBot) SetupTelegramBot() {
	tb.choosePrice()
	tb.storeChatId()
}

const (
	storingChatId  = "s"
	choosingPrice  = "c"
	askingForPrice = "a"
)

func (tb *telegramBot) storeChatId() {
	op := place + "storeChatId"
	log := tb.Logger.AddOp(op)
	yes := telebot.Btn{Unique: "store_chatid_yes", Text: "✅ да"}
	no := telebot.Btn{Unique: "store_chatid_no", Text: "❌ нет"}

	tb.Bot.Handle("/start", func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if !ok || state == storingChatId {
			tb.usersStatements.Store(c.Chat().ID, storingChatId)
			markup := &telebot.ReplyMarkup{}
			markup.Inline(markup.Row(yes, no))
			return c.Send("хотите получать обновления цены айфончика 17??", markup)
		}
		return nil
	})
	tb.Bot.Handle(&yes, func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if ok && state.(string) == storingChatId {
			defer tb.usersStatements.Delete(chatId)
			username := c.Sender().Username
			ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
			defer cancel()
			if err := tb.UserRepository.SetChatId(ctx, username, chatId); err != nil {
				if errors.Is(err, errs.ErrAlreadyExistsBase) {
					return c.Edit("вы уже получаете обновления")
				}
				log.Error("failed to set chat id", logger.Err(err))
				return c.Edit("произошла ошибка, возможно вы не зарегестрированы((")
			}
			return c.Edit("ждите обновления))")
		}
		return nil
	})
	tb.Bot.Handle(&no, func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if ok && state.(string) == storingChatId {
			defer tb.usersStatements.Delete(chatId)
			username := c.Sender().Username
			ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
			defer cancel()
			if err := tb.UserRepository.DropChatId(ctx, username, chatId); err != nil {
				if errors.Is(err, errs.ErrNotFoundBase) {
					return c.Edit("вы не подписаны на обновления")
				}
				log.Error("failed to delete chat id", logger.Err(err))
				return c.Edit("произошла ошибка((")
			}
			return c.Edit("обновлений не ждите((")
		}
		return nil
	})
}

func (tb *telegramBot) choosePrice() {
	op := place + "choosePrice"
	log := tb.Logger.AddOp(op)
	yes := telebot.Btn{Unique: "choose_price_yes", Text: "✅ да2"}
	no := telebot.Btn{Unique: "choose_price_no", Text: "❌ нет"}
	tb.Bot.Handle("/setprice", func(c telebot.Context) error {
		chatId := c.Chat().ID
		ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
		defer cancel()
		exist, err := tb.UserRepository.CheckChatId(ctx, op, c.Sender().Username, chatId)
		if err != nil {
			log.Error("failed to check chat id", logger.Err(err))
			return c.Send("произошла ошибка((")
		}
		if !exist {
			_, ok := tb.usersStatements.Load(chatId)
			if !ok {
				tb.usersStatements.Store(c.Chat().ID, askingForPrice)
				markup := &telebot.ReplyMarkup{}
				markup.Inline(markup.Row(yes, no))
				return c.Send("хотите установить цену айфончика при достижении которой жоско заспамлю??", markup)
			}
		} else {
			return c.Send("сначала на обновления подпишитесь")
		}
		return nil
	})
	tb.Bot.Handle(&no, func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if ok && state.(string) == askingForPrice {
			tb.usersStatements.Delete(chatId)
			ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
			defer cancel()
			if err := tb.UserRepository.DropDesiredPrice(ctx, chatId); err != nil {
				log.Error("failed to drop desired price", logger.Err(err))
				return c.Edit("произошла ошибка((")
			}
			return c.Edit("нет так нет")
		}
		return nil
	})
	tb.Bot.Handle(&yes, func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if ok && state.(string) == askingForPrice {
			tb.usersStatements.Store(chatId, choosingPrice)
			return c.Edit("напиши цену, например 2800.52")
		}
		return nil
	})
	tb.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		chatId := c.Chat().ID
		state, ok := tb.usersStatements.Load(chatId)
		if ok && state.(string) == choosingPrice {
			defer tb.usersStatements.Delete(chatId)
			strPrice := strings.TrimSpace(strings.ReplaceAll(c.Text(), ",", "."))
			price, err := strconv.ParseFloat(strPrice, 32)
			if err != nil {
				return c.Send("❌ неправильный формат цены!!")
			}
			ctx, cancel := context.WithTimeout(context.Background(), tb.Config.Timeout)
			defer cancel()
			if err := tb.UserRepository.SetDesiredPrice(ctx, chatId, price); err != nil {
				log.Error("failed to set desired price", logger.Err(err))
				return c.Send("произошла ошибка((")
			}
			return c.Send(fmt.Sprintf("✅ цена установлена: %.2f", price))
		}
		return nil
	})
}

const (
	grafUp   = "📈"
	grafDown = "📉"
	grafDef  = "0️⃣"
	blue     = "🟦"
	white    = "⬜"
	black    = "⬛"
	green    = "🟩"
	pink     = "🟪"
	plus     = "+"
	zero     = ""
)

func (tb *telegramBot) SendIPhonesInfo(chatIds []int64, iphones []models.IPhone) error {
	op := place + "SendIphoneInfo"
	msgArr := []string{}
	for _, iphone := range iphones {
		graf := grafDef
		color := white
		sign := zero
		switch iphone.Color {
		case "F5F5F5":
			color = white
		case "353839":
			color = black
		case "96AED1":
			color = blue
		case "A9B689":
			color = green
		case "DFCEEA":
			color = pink
		}

		if iphone.Change > 0 {
			graf = grafUp
			sign = plus
		} else if iphone.Change < 0 {
			graf = grafDown
		}
		msgArr = append(msgArr, fmt.Sprintf("%s %s:\n 💰 цена: %.2f | %s разница: %s%.2f\n", iphone.Name, color, iphone.Price, graf, sign, iphone.Change))
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

			return errs.NewAppError(op, err)
		}
	}
	return nil
}

func (tb *telegramBot) Start() {
	tb.Bot.Start()
}

func (tb *telegramBot) Stop() {
	tb.Bot.Stop()
}
