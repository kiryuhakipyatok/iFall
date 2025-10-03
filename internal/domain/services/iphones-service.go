package services

import (
	"context"
	"iFall/internal/client"
	"iFall/internal/domain/models"
	"iFall/internal/domain/repositories"
	"iFall/pkg/errs"
	"iFall/pkg/logger"
)

type IPhoneService interface {
	Get(ctx context.Context, id string) (*models.IPhone, error)
	Update(ctx context.Context, id string) (*models.IPhone, error)
}

type iPhoneService struct {
	IPhoneRepository repositories.IPhoneRepository
	ApiClient        *client.ApiClient
	Logger           *logger.Logger
}

func NewIPhoneService(ir repositories.IPhoneRepository, ac *client.ApiClient, l *logger.Logger) IPhoneService {
	return &iPhoneService{
		IPhoneRepository: ir,
		ApiClient:        ac,
		Logger:           l,
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
	return iphone, nil
}

func (is *iPhoneService) Update(ctx context.Context, id string) (*models.IPhone, error) {
	op := place + "Update"
	log := is.Logger.AddOp(op)
	log.Info("iphone updating", "id", id)
	// iphone, err := is.IPhoneRepository.Get(ctx, id)
	// if err != nil {
	// 	log.Error("failed to receive iphone", logger.Err(err))
	// 	return errs.NewAppError(op, err)
	// }
	iphoneData, err := is.ApiClient.GetIPhoneData(id)
	if err != nil {
		log.Error("failed to receive iphone data", logger.Err(err))
		return nil, errs.NewAppError(op, err)
	}
	iphone, err := is.IPhoneRepository.Update(ctx, id, iphoneData.Price)
	if err != nil {
		log.Error("failed to update iphone", logger.Err(err))
		return nil, errs.NewAppError(op, err)
	}

	log.Info("iphone updated", "id", id)
	return iphone, nil
}
