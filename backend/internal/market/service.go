package market

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"trading-dashboard/internal/models"
)

type Service struct {
	repo   *Repository
	client *http.Client
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) Search(query string) ([]models.SearchResult, error) {
	endpoint := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v1/finance/search?q=%s&quotesCount=10&lang=en-US",
		url.QueryEscape(query),
	)

	// my backend is calling external API, toh user-agent set karna zaruri hai taaki Yahoo Finance request ko block na kare. Isse hum apne app ko identify kar sakte hain.
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TradingDashboard/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch search: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// yahoo finance response
	/*
				{
		            "symbol": "AAPL",
		            "name": "Apple Inc.",
		            "exchange": "NMS",
		            "quote_type": "EQUITY"
		        }
	*/
	var yfResp models.YFSearchResponse
	if err := json.Unmarshal(body, &yfResp); err != nil {
		return nil, fmt.Errorf("parse search: %w", err)
	}

	var results []models.SearchResult
	for _, q := range yfResp.Quotes {
		if q.QuoteType == "" {
			continue
		}
		name := q.Longname
		if name == "" {
			name = q.Shortname
		}
		results = append(results, models.SearchResult{
			Symbol:    q.Symbol,
			Name:      name,
			Exchange:  q.Exchange,
			QuoteType: q.QuoteType,
		})
	}
	return results, nil
}

func (s *Service) FetchAndStoreHistory(symbol, period string) (*models.HistoryResponse, error) {
	symbol = strings.ToUpper(symbol)
	days := 30
	yahooRange := "1mo"
	if period == "1d" {
		days = 7 // 7 days fetch karo taaki weekend/holiday miss na ho
		yahooRange = "5d"
	}

	if s.repo.IsDataStale(symbol) {
		log.Printf("📡 Fetching %s from Yahoo Finance...", symbol)
		chart, err := s.fetchYahooChart(symbol, yahooRange)
		if err != nil {
			log.Printf("⚠️  Yahoo fetch failed: %v", err)
		} else {
			if err := s.storeChartData(symbol, chart); err != nil {
				log.Printf("⚠️  Store failed: %v", err)
			}
		}
	}

	history, err := s.repo.GetHistory(symbol, days)
	if err != nil {
		return nil, err
	}

	meta, _ := s.repo.GetSymbolMetadata(symbol)
	company := symbol
	if meta != nil {
		company = meta.CompanyName
	}

	// 1d ke liye sirf last trading day ka data return karo
	if period == "1d" && len(history) > 0 {
		history = history[len(history)-1:]
	}

	return &models.HistoryResponse{
		Symbol:  symbol,
		Company: company,
		Period:  period,
		Data:    history,
	}, nil
}

func (s *Service) GetQuote(symbol string) (*models.QuoteResponse, error) {
	symbol = strings.ToUpper(symbol)
	chart, err := s.fetchYahooChart(symbol, "5d")
	if err != nil {
		return nil, err
	}
	if len(chart.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data for %s", symbol)
	}

	meta := chart.Chart.Result[0].Meta
	currentPrice := meta.RegularMarketPrice
	prevClose := meta.PreviousClose
	change := currentPrice - prevClose
	changePct := 0.0
	if prevClose != 0 {
		changePct = (change / prevClose) * 100
	}

	_ = s.repo.UpsertSymbolMetadata(models.SymbolMetadata{
		Symbol:      symbol,
		CompanyName: meta.LongName,
		Currency:    meta.Currency,
		Exchange:    meta.ExchangeName,
	})

	return &models.QuoteResponse{
		Symbol:        symbol,
		CompanyName:   meta.LongName,
		CurrentPrice:  currentPrice,
		Change:        change,
		ChangePercent: changePct,
		Currency:      meta.Currency,
		Exchange:      meta.ExchangeName,
	}, nil
}

/* API Response for 5 day
{
  "chart": {
    "result": [
      {
        "meta": {
          "currency": "USD",
          "symbol": "RS",
          "exchangeName": "NYQ",
          "fullExchangeName": "NYSE",
          "instrumentType": "EQUITY",
          "firstTradeDate": 779722200,
          "regularMarketTime": 1773777602,
          "hasPrePostMarketData": true,
          "gmtoffset": -14400,
          "timezone": "EDT",
          "exchangeTimezoneName": "America/New_York",
          "regularMarketPrice": 299.04,
          "fiftyTwoWeekHigh": 365.59,
          "fiftyTwoWeekLow": 250.07,
          "regularMarketDayHigh": 301.77,
          "regularMarketDayLow": 297.94,
          "regularMarketVolume": 221426,
          "longName": "Reliance, Inc.",
          "shortName": "Reliance, Inc.",
          "chartPreviousClose": 300.31,
          "priceHint": 2,
          "currentTradingPeriod": {
            "pre": {
              "timezone": "EDT",
              "end": 1773840600,
              "start": 1773820800,
              "gmtoffset": -14400
            },
            "regular": {
              "timezone": "EDT",
              "end": 1773864000,
              "start": 1773840600,
              "gmtoffset": -14400
            },
            "post": {
              "timezone": "EDT",
              "end": 1773878400,
              "start": 1773864000,
              "gmtoffset": -14400
            }
          },
          "dataGranularity": "1d",
          "range": "5d",
          "validRanges": [
            "1d",
            "5d",
            "1mo",
            "3mo",
            "6mo",
            "1y",
            "2y",
            "5y",
            "10y",
            "ytd",
            "max"
          ]
        },
        "timestamp": [
          1773235800,
          1773322200,
          1773408600,
          1773667800,
          1773754200,
          1773777602
        ],
        "indicators": {
          "quote": [
            {
              "low": [
                297.95001220703125,
                298.7699890136719,
                294.70001220703125,
                296.6000061035156,
                297.6700134277344,
                297.94000244140625
              ],
              "close": [
                308.0899963378906,
                299.2799987792969,
                297.44000244140625,
                297.79998779296875,
                299.0400085449219,
                299.0400085449219
              ],
              "open": [
                299.6199951171875,
                306.19000244140625,
                300.32000732421875,
                296.6000061035156,
                299.7300109863281,
                299.7300109863281
              ],
              "high": [
                308.489990234375,
                307,
                300.32000732421875,
                303.239990234375,
                302.5400085449219,
                301.7699890136719
              ],
              "volume": [
                477600,
                384800,
                391100,
                368400,
                281200,
                221426
              ]
            }
          ],
          "adjclose": [
            {
              "adjclose": [
                308.0899963378906,
                299.2799987792969,
                297.44000244140625,
                297.79998779296875,
                299.0400085449219,
                299.0400085449219
              ]
            }
          ]
        }
      }
    ],
    "error": null
  }
}
*/

func (s *Service) fetchYahooChart(symbol, yahooRange string) (*models.YFChartResponse, error) {

	// symbol: RS (Reliance etc), range: 1d, 30d
	endpoint := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=%s",
		url.PathEscape(symbol), yahooRange,
	)

	// external api call to yahoo finance
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TradingDashboard/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo returned %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var chart models.YFChartResponse
	if err := json.Unmarshal(body, &chart); err != nil {
		return nil, err
	}
	return &chart, nil
}

func (s *Service) storeChartData(symbol string, chart *models.YFChartResponse) error {
	if len(chart.Chart.Result) == 0 {
		return fmt.Errorf("empty result")
	}
	result := chart.Chart.Result[0]
	meta := result.Meta

	companyName := meta.LongName
	if companyName == "" {
		companyName = meta.ShortName
	}
	_ = s.repo.UpsertSymbolMetadata(models.SymbolMetadata{
		Symbol:      symbol,
		CompanyName: companyName,
		Currency:    meta.Currency,
		Exchange:    meta.ExchangeName,
	})

	if len(result.Timestamps) == 0 || len(result.Indicators.Quote) == 0 {
		return fmt.Errorf("no quote data")
	}

	quotes := result.Indicators.Quote[0]
	for i, ts := range result.Timestamps {
		if i >= len(quotes.Open) || i >= len(quotes.Close) {
			continue
		}
		if quotes.Open[i] == 0 && quotes.Close[i] == 0 {
			continue
		}
		var vol int64
		if i < len(quotes.Volume) {
			vol = quotes.Volume[i]
		}
		_ = s.repo.UpsertPriceHistory(
			symbol,
			time.Unix(ts, 0).UTC(),
			quotes.Open[i], quotes.High[i], quotes.Low[i], quotes.Close[i],
			vol,
		)
	}
	return nil
}
