package client

import (
	"fmt"
	"iFall/internal/config"
	"iFall/internal/domain/models"
	"iFall/pkg/errs"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ApiClient struct {
	Client  *http.Client
	BaseURL string
}

func NewClient(cfg config.ApiClientConfig) *ApiClient {
	client := http.Client{
		Timeout: cfg.Timeout,
	}
	return &ApiClient{
		Client:  &client,
		BaseURL: cfg.BaseURL,
	}
}

var nonDigit = regexp.MustCompile(`\D`)

func parseMajorPrice(s string) (float64, error) {
	s = strings.ReplaceAll(s, "\u00A0", " ")
	digits := nonDigit.ReplaceAllString(s, "")
	if digits == "" {
		return 0, fmt.Errorf("no digits in price %q", s)
	}
	return strconv.ParseFloat(digits, 32)
}

func (ac *ApiClient) GetIPhoneData(id string) (*models.IPhone, error) {
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
	priceSel.Find("i").Remove()
	rawPrice := strings.TrimSpace(priceSel.Text())

	price, err := parseMajorPrice(rawPrice)
	if err != nil {
		return nil, errs.NewAppError(op, err)
	}
	iphone := &models.IPhone{
		Name:  name,
		Price: price,
	}
	return iphone, nil
}
