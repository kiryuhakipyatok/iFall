package services

import (
	"context"
	"iFall/internal/bot"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/internal/email"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"sync"
)

type IphoneReportService interface {
	SendIPhonesInfo(emailSupp bool, iphones []models.IPhone) error
}

type iPhoneReportService struct {
	IPhoneService  IPhoneService
	UserRepository repositories.UserRepository
	IPonesConfig   config.IPhonesConfig
	EmailSender    *email.EmailSender
	Bot            *bot.TelegramBot
	Logger         *logger.Logger
}

func NewIPhoneReportService(is IPhoneService, ur repositories.UserRepository, l *logger.Logger, b *bot.TelegramBot, es *email.EmailSender, cfg config.IPhonesConfig) IphoneReportService {
	return &iPhoneReportService{
		IPhoneService:  is,
		UserRepository: ur,
		EmailSender:    es,
		Bot:            b,
		IPonesConfig:   cfg,
		Logger:         l,
	}
}

func (irs *iPhoneReportService) SendIPhonesInfo(emailSupp bool, iphones []models.IPhone) error {
	op := "iPhoneReportService.sendIPhonesInfo"
	log := irs.Logger.AddOp(op)
	log.Info("sending iphones info")
	ctx, cancel := context.WithTimeout(context.Background(), irs.IPonesConfig.Timeout)
	defer cancel()
	contacts, err := irs.UserRepository.FetchContacts(ctx)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	emails := []string{}
	chatIds := []int64{}
	if len(contacts) > 0 {
		for _, c := range contacts {
			if emailSupp {
				emails = append(emails, c.Email)
			}
			if c.ChatId != nil {
				chatIds = append(chatIds, *c.ChatId)
			}
		}
	} else {
		return nil
	}

	errlen := len(chatIds) + 1
	if emailSupp {
		errlen = len(contacts) + 1
	}

	errChan := make(chan error, errlen)
	var wg sync.WaitGroup
	if emailSupp {
		if len(emails) > 0 {
			wg.Add(1)
			log.Info("sending on emails")
			emailCtx, cancel := context.WithTimeout(context.Background(), irs.IPonesConfig.Timeout)
			defer cancel()
			subject := "цена говнофона семнадцатого 17"
			content, err := email.BuildEmailLetter(iphones)
			if err != nil {
				errChan <- err
			}
			go func() {
				defer wg.Done()
				if err := irs.EmailSender.SendMessage(emailCtx, subject, []byte(content), emails, nil); err != nil {
					errChan <- err
				}
			}()
		}
	}
	if len(chatIds) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("sending on telegrams")
			if err := irs.Bot.SendIPhonesInfo(chatIds, iphones); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Error("failed to send iphones info", logger.Err(err))
		return errs.NewAppError(op, err)
	}

	log.Info("iphones info sended")
	return nil
}
