package client

import (
	"fmt"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//go:generate mockgen -source=client.go -destination=mocks/client-mock.go
type ApiClient interface {
	GetIPhoneData(id string) (*models.IPhone, error)
}

type apiClient struct {
	Client  *http.Client
	BaseURL string
}

func NewClient(cfg config.ApiClientConfig) ApiClient {
	client := http.Client{
		Timeout: cfg.Timeout,
	}
	return &apiClient{
		Client:  &client,
		BaseURL: cfg.BaseURL,
	}
}

func (ac *apiClient) GetIPhoneData(id string) (*models.IPhone, error) {
	op := "apiClient.GetIPhoneData"
	url := fmt.Sprintf("%s/%s", ac.BaseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-http-client)")
	resp, err := ac.Client.Do(req)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	name := strings.TrimSpace(doc.Find(`h1[itemprop="name"]`).First().Text())
	priceSel := doc.Find(".price-block .price:not(.old)").First()
	rawPrice := strings.TrimSpace(priceSel.Text())
	strPrice := strings.ReplaceAll(rawPrice, " ", "")
	floatPrice, err := strconv.ParseFloat(strPrice, 32)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	price := floatPrice / 100
	iphone := &models.IPhone{
		Name:  name,
		Price: price,
	}
	return iphone, nil
}
