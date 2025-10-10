package services

import (
	"context"
	"iFall/internal/client"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/internal/email"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
	"sync"
)

type IPhoneService interface {
	Get(ctx context.Context, id string) (*models.IPhone, error)
	UpdateAll() ([]models.IPhone, error)
}

type iPhoneService struct {
	IPhoneRepository repositories.IPhoneRepository
	ApiClient        client.ApiClient
	IPhonesConfig    config.IPhonesConfig
	EmailSendler     email.EmailSender
	Logger           *logger.Logger
	Mutex            sync.Mutex
}

func NewIPhoneService(ir repositories.IPhoneRepository, ac client.ApiClient, l *logger.Logger, es email.EmailSender, cfg config.IPhonesConfig) IPhoneService {
	return &iPhoneService{
		IPhoneRepository: ir,
		ApiClient:        ac,
		Logger:           l,
		IPhonesConfig:    cfg,
		EmailSendler:     es,
	}
}

const place = "iPhoneService."

func (is *iPhoneService) Get(ctx context.Context, id string) (*models.IPhone, error) {
	op := place + "Get"
	log := is.Logger.AddOp(op)
	log.Info("iphone receiving")
	iphone, err := is.IPhoneRepository.Get(ctx, id)
	if err != nil {
		log.Error("failed to receive iphone", logger.Err(err))
		return nil, errs.NewAppError(op, err)
	}
	log.Info("iphone received", "id", iphone.Id)
	return iphone, nil
}

func (is *iPhoneService) update(ctx context.Context, id string) (*models.IPhone, error) {
	op := place + "update"
	log := is.Logger.AddOp(op)
	log.Info("iphone updating", "id", id)
	iphoneData, err := is.ApiClient.GetIPhoneData(id)
	if err != nil {
		log.Error("failed to receive iphone data", logger.Err(err))
		return nil, errs.NewAppError(op, err)
	}
	is.Mutex.Lock()
	iphone, err := is.IPhoneRepository.Update(ctx, id, iphoneData.Price)
	if err != nil {
		log.Error("failed to update iphone", logger.Err(err))
		return nil, errs.NewAppError(op, err)
	}
	is.Mutex.Unlock()
	log.Info("iphone updated", "id", id)
	return iphone, nil
}

func (is *iPhoneService) UpdateAll() ([]models.IPhone, error) {
	op := place + "updateAll"
	log := is.Logger.AddOp(op)
	log.Info("updating all iphones")
	iphones := []string{
		is.IPhonesConfig.Black,
		is.IPhonesConfig.Green,
		is.IPhonesConfig.White,
		is.IPhonesConfig.Blue,
		is.IPhonesConfig.Pink,
	}

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
			ctx, cancel := context.WithTimeout(context.Background(), is.IPhonesConfig.Timeout)
			defer cancel()
			iphone, err := is.update(ctx, id)
			if err != nil {
				errChan <- errStruct{err: err, id: id}
			} else {
				iphoneChan <- *iphone
			}
		}(v)
	}

	wg.Wait()
	close(errChan)
	close(iphoneChan)

	for err := range errChan {
		return nil, errs.NewAppError(op, err.err)
	}

	iphonesData := []models.IPhone{}
	for i := range iphoneChan {
		iphonesData = append(iphonesData, i)
	}
	log.Info("all iphones updated")
	return iphonesData, nil
}
