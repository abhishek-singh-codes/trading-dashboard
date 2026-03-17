package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

type PriceHistory struct {
	ID        int       `json:"id"`
	Symbol    string    `json:"symbol"`
	TradeDate time.Time `json:"trade_date"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

type SymbolMetadata struct {
	Symbol      string    `json:"symbol"`
	CompanyName string    `json:"company_name"`
	Currency    string    `json:"currency"`
	Exchange    string    `json:"exchange"`
	LastFetched time.Time `json:"last_fetched"`
}

type RegisterRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"     binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SearchResult struct {
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Exchange  string `json:"exchange"`
	QuoteType string `json:"quote_type"`
}

type QuoteResponse struct {
	Symbol        string  `json:"symbol"`
	CompanyName   string  `json:"company_name"`
	CurrentPrice  float64 `json:"current_price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Currency      string  `json:"currency"`
	Exchange      string  `json:"exchange"`
}

type HistoryResponse struct {
	Symbol  string         `json:"symbol"`
	Company string         `json:"company"`
	Period  string         `json:"period"`
	Data    []PriceHistory `json:"data"`
}

type YFSearchResponse struct {
	Quotes []struct {
		Symbol    string `json:"symbol"`
		Shortname string `json:"shortname"`
		Longname  string `json:"longname"`
		Exchange  string `json:"exchange"`
		QuoteType string `json:"quoteType"`
	} `json:"quotes"`
}

type YFChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol             string  `json:"symbol"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				PreviousClose      float64 `json:"chartPreviousClose"`
				Currency           string  `json:"currency"`
				ExchangeName       string  `json:"exchangeName"`
				LongName           string  `json:"longName"`
				ShortName          string  `json:"shortName"`
			} `json:"meta"`
			Timestamps []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}