package services

import (
	"context"
	"iFall/internal/bot"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/internal/email"
	"iFall/pkg/errs"
	"sync"
)

type IphoneReportService interface {
	SendIPhonesInfo(iphones []models.IPhone) error
}

type iPhoneReportService struct {
	IPhoneService  IPhoneService
	UserRepository repositories.UserRepository
	IPonesConfig   config.IPhonesConfig
	EmailSender    *email.EmailSender
	Bot            *bot.TelegramBot
}

func NewIPhoneReportService(is IPhoneService, ur repositories.UserRepository, b *bot.TelegramBot, es *email.EmailSender, cfg config.IPhonesConfig) IphoneReportService {
	return &iPhoneReportService{
		IPhoneService:  is,
		UserRepository: ur,
		EmailSender:    es,
		Bot:            b,
		IPonesConfig:   cfg,
	}
}

func (irs *iPhoneReportService) SendIPhonesInfo(iphones []models.IPhone) error {
	op := "scheduler.sendIPhonesInfo"
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
			emails = append(emails, c.Email)
			if c.ChatId != nil {
				chatIds = append(chatIds, *c.ChatId)
			}
		}
	} else {
		return nil
	}

	emailCtx, cancel := context.WithTimeout(context.Background(), irs.IPonesConfig.Timeout)
	defer cancel()
	subject := "цена говнофона семнадцатого 17"
	content, err := email.BuildEmailLetter(iphones)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	errChan := make(chan error, len(contacts))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := irs.EmailSender.SendMessage(emailCtx, subject, []byte(content), emails, nil); err != nil {
			errChan <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := irs.Bot.SendIPhonesInfo(chatIds, iphones); err != nil {
			errChan <- err
		}
	}()
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
