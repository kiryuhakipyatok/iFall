package scheduler

import (
	"fmt"
	"iFall/internal/config"
	"iFall/internal/domain/services"
	"iFall/pkg/logger"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Cron                *cron.Cron
	Logger              *logger.Logger
	IPhoneService       services.IPhoneService
	IPhoneReportService services.IphoneReportService
	SchedulerConfig     config.SchedulerConfig
}

func NewScheduler(is services.IPhoneService, irs services.IphoneReportService, l *logger.Logger, scfg config.SchedulerConfig) *Scheduler {
	cr := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	return &Scheduler{
		Cron:                cr,
		IPhoneService:       is,
		IPhoneReportService: irs,
		Logger:              l,
		SchedulerConfig:     scfg,
	}
}

func (s *Scheduler) Start() {
	if _, err := s.Cron.AddFunc(fmt.Sprintf("%d %d * * *", s.SchedulerConfig.Minute, s.SchedulerConfig.Hour), func() {
		op := "scheduler.GetIPhoneData"
		log := s.Logger.AddOp(op)
		log.Info("updating iphone data")
		iphones, err := s.IPhoneService.UpdateAll()
		if err != nil {
			log.Error("failed to updated iphones info", logger.Err(err))
		} else if len(iphones) > 0 {
			log.Info("sending iphones info")
			if err := s.IPhoneReportService.SendIPhonesInfo(iphones); err != nil {
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

func (s *Scheduler) Stop() {
	ctx := s.Cron.Stop()
	<-ctx.Done()
}
