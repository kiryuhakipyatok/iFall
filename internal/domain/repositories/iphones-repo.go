package repositories

import (
	"context"
	"errors"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"iFall/pkg/storage"
)

//go:generate mockgen -source=iphones-repo.go -destination=mocks/iphones-repo-mock.go
type IPhoneRepository interface {
	Get(ctx context.Context, id string) (*models.IPhone, error)
	Update(ctx context.Context, id string, price float64) (*models.IPhone, error)
}

type iPhoneRepository struct {
	Storage *storage.Storage
}

func NewIPhoneRepository(s *storage.Storage) IPhoneRepository {
	return &iPhoneRepository{
		Storage: s,
	}
}

const iphonesRepo = "iPhoneRepository."

func (ir *iPhoneRepository) Get(ctx context.Context, id string) (*models.IPhone, error) {
	op := iphonesRepo + "Get"
	query := "SELECT * FROM iphones WHERE id = $1"
	iphone := &models.IPhone{}
	if err := ir.Storage.DB.QueryRowContext(ctx, query, id).Scan(
		&iphone.Id,
		&iphone.Name,
		&iphone.Price,
		&iphone.Change,
		&iphone.Color,
	); err != nil {
		if errors.Is(err, storage.ErrNotFound()) {
			return nil, errs.ErrNotFound(op)
		}
	}
	return iphone, nil
}

func (ir *iPhoneRepository) Update(ctx context.Context, id string, price float64) (*models.IPhone, error) {
	op := iphonesRepo + "Update"
	query := "UPDATE iphones SET price=$1, change=$1-iphones.price WHERE id=$2 RETURNING name, price, color, change"
	iphone := &models.IPhone{}
	if err := ir.Storage.DB.QueryRowContext(ctx, query, price, id).Scan(
		&iphone.Name,
		&iphone.Price,
		&iphone.Color,
		&iphone.Change,
	); err != nil {
		if errors.Is(err, storage.ErrNotFound()) {
			return nil, errs.ErrNotFound(op)
		}
		return nil, errs.NewAppError(op, err)
	}

	return iphone, nil
}
