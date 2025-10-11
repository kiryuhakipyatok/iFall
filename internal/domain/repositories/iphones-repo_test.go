package repositories

import (
	"context"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"iFall/pkg/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIPhoneRepository_Get(t *testing.T) {
	type testIphoneData struct {
		id     string
		name   string
		price  float64
		color  string
		change float64
	}
	tests := []struct {
		id             string
		testName       string
		expectedError  error
		testIphoneData testIphoneData
		expectedResult *models.IPhone
	}{
		{
			testName:      "success receiving",
			expectedError: nil,
			id:            "iphone-1-id",
			testIphoneData: testIphoneData{
				id:     "iphone-1-id",
				name:   "iphone1",
				price:  1000.0,
				change: 0.0,
				color:  "ffffff",
			},
			expectedResult: &models.IPhone{
				Id:     "iphone-1-id",
				Name:   "iphone1",
				Price:  1000.0,
				Change: 0.0,
				Color:  "ffffff",
			},
		},
		{
			testName:       "not found",
			expectedError:  errs.ErrNotFoundBase,
			testIphoneData: testIphoneData{},
			id:             "iphone-1-id",
			expectedResult: nil,
		},
	}
	storage := storage.MustConnect(config.StorageConfig{Path: ":memory:", PingTimeout: time.Second})

	schema := `
   				CREATE TABLE IF NOT EXISTS iphones (
    				id TEXT PRIMARY KEY,
    				name TEXT NOT NULL UNIQUE,
    				price NUMERIC NOT NULL,
    				change NUMERIC NOT NULL DEFAULT 0,
    				color TEXT NOT NULL DEFAULT 'ffffff'
				);				
    		`
	if _, err := storage.DB.Exec(schema); err != nil {
		t.Fatalf("failed to create test iphones table: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			if _, err := storage.DB.Exec("DELETE FROM iphones"); err != nil {
				t.Fatalf("failed to delete iphones data from table: %v", err)
			}

			if _, err := storage.DB.Exec("INSERT INTO iphones (id, name, price, color, change) VALUES ($1, $2, $3, $4, $5)", tt.testIphoneData.id, tt.testIphoneData.name, tt.testIphoneData.price, tt.testIphoneData.color, tt.testIphoneData.change); err != nil {
				t.Fatalf("failed to insert test iphone data in the table: %v", err)
			}
			repo := NewIPhoneRepository(storage)
			iphone, err := repo.Get(context.Background(), tt.id)
			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, iphone)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
