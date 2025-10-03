package scheduler

import (
	"context"
	"fmt"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/internal/domain/services"
	"iFall/internal/email"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Cron            *cron.Cron
	Logger          *logger.Logger
	IPhoneService   services.IPhoneService
	UserRepository  repositories.UserRepository
	EmailSender     *email.EmailSender
	SchedulerConfig config.SchedulerConfig
	IPhonesConfig   config.IPhonesConfig
}

func NewScheduler(is services.IPhoneService, ur repositories.UserRepository, es *email.EmailSender, l *logger.Logger, scfg config.SchedulerConfig, icfg config.IPhonesConfig) *Scheduler {
	cr := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	return &Scheduler{
		Cron:            cr,
		IPhoneService:   is,
		UserRepository:  ur,
		Logger:          l,
		EmailSender:     es,
		SchedulerConfig: scfg,
		IPhonesConfig:   icfg,
	}
}

func (s *Scheduler) Start() {
	if _, err := s.Cron.AddFunc(fmt.Sprintf("%d %d * * *", s.SchedulerConfig.Minute, s.SchedulerConfig.Hour), func() {
		op := "scheduler.GetIPhoneData"
		log := s.Logger.AddOp(op)
		log.Info("updating iphone data")
		iphones, err := s.updateAll()
		if err != nil {
			log.Error("failed to updated iphones info", logger.Err(err))
		} else if len(iphones) > 0 {
			log.Info("sending iphones info")
			if err := s.sendIPhonesInfo(iphones); err != nil {
				log.Error("failed to send iphones info", logger.Err(err))
			} else {
				log.Info("iphones info sended")
			}
		}
	}); err != nil {
		panic(fmt.Errorf("failed to start scheduler: %w", err))
	}
	s.Cron.Start()
}

func (s *Scheduler) updateAll() ([]models.IPhone, error) {
	op := "scheduler.updateAll"
	iphones := []string{
		s.IPhonesConfig.Black,
		s.IPhonesConfig.Green,
		s.IPhonesConfig.White,
		s.IPhonesConfig.Blue,
		s.IPhonesConfig.Pink,
	}

	fmt.Println(iphones)

	type errStruct struct {
		err error
		id  string
	}
	errChan := make(chan errStruct, len(iphones))
	iphoneChan := make(chan models.IPhone, len(iphones))

	var wg sync.WaitGroup

	for _, v := range iphones {
		if len(v) == 0 {
			continue
		}
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), s.SchedulerConfig.Timeout)
			defer cancel()
			iphone, err := s.IPhoneService.Update(ctx, id)
			if err != nil {
				errChan <- errStruct{err: err, id: id}
			} else {
				iphoneChan <- *iphone
			}
		}(v)
	}
	go func() {
		wg.Wait()
		close(errChan)
		close(iphoneChan)
	}()
	if len(errChan) > 0 {
		for err := range errChan {
			return nil, errs.NewAppError(op, err.err)
		}
	}

	iphonesData := []models.IPhone{}
	for i := range iphoneChan {
		iphonesData = append(iphonesData, i)
	}
	return iphonesData, nil
}

func (s *Scheduler) sendIPhonesInfo(iphones []models.IPhone) error {
	op := "scheduler.sendIPhonesInfo"
	ctx, cancel := context.WithTimeout(context.Background(), s.SchedulerConfig.Timeout)
	defer cancel()
	emails, err := s.UserRepository.FetchEmails(ctx)
	if err != nil {
		return errs.NewAppError(op, err)
	}
	if len(emails) > 0 {
		emailCtx, cancel := context.WithTimeout(context.Background(), s.SchedulerConfig.Timeout)
		defer cancel()
		subject := "цена говнофона семнадцатого 17"
		content, err := email.BuildEmailLetter(iphones)
		if err != nil {
			return errs.NewAppError(op, err)
		}
		if err := s.EmailSender.SendMessage(emailCtx, subject, []byte(content), emails, nil); err != nil {
			return errs.NewAppError(op, err)
		}
	}
	return nil
}

func (s *Scheduler) Stop() {
	ctx := s.Cron.Stop()
	<-ctx.Done()
}
