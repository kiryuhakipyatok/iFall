package dto

type IPhoneResponse struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Prices struct {
		PriceMin struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"price_min"`
	} `json:"prices"`
	Color string `json:"color_code"`
}
