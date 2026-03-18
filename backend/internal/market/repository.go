package market

import (
	"database/sql"
	"time"

	"trading-dashboard/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// history ko upsert karna hai, kyunki agar same date ke liye data aa raha hai toh update karna hai, nahi toh insert karna hai. Isse hum ensure karte hain ki humare paas hamesha latest data ho aur duplicate entries na ho.
func (r *Repository) UpsertPriceHistory(symbol string, date time.Time, open, high, low, close float64, volume int64) error {
	_, err := r.db.Exec(`
		INSERT INTO price_history (symbol, trade_date, open_price, high_price, low_price, close_price, volume, fetched_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (symbol, trade_date) DO UPDATE SET
			open_price  = EXCLUDED.open_price,
			high_price  = EXCLUDED.high_price,
			low_price   = EXCLUDED.low_price,
			close_price = EXCLUDED.close_price,
			volume      = EXCLUDED.volume,
			fetched_at  = NOW()
	`, symbol, date, open, high, low, close, volume)
	return err
}

func (r *Repository) GetHistory(symbol string, days int) ([]models.PriceHistory, error) {
	rows, err := r.db.Query(`
		SELECT id, symbol, trade_date, open_price, high_price, low_price, close_price, volume
		FROM price_history
		WHERE symbol = $1
		  AND trade_date >= NOW() - ($2 || ' days')::INTERVAL
		ORDER BY trade_date ASC
	`, symbol, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var p models.PriceHistory
		if err := rows.Scan(&p.ID, &p.Symbol, &p.TradeDate, &p.Open, &p.High, &p.Low, &p.Close, &p.Volume); err != nil {
			return nil, err
		}
		history = append(history, p)
	}
	return history, rows.Err()
}

// update the symbol metadata (company name, exchange, currency) - this is called when we fetch new data from Yahoo
func (r *Repository) UpsertSymbolMetadata(meta models.SymbolMetadata) error {
	_, err := r.db.Exec(`
		INSERT INTO symbol_metadata (symbol, company_name, currency, exchange, last_fetched)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (symbol) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			currency     = EXCLUDED.currency,
			exchange     = EXCLUDED.exchange,
			last_fetched = NOW()
	`, meta.Symbol, meta.CompanyName, meta.Currency, meta.Exchange)
	return err
}

// get the symbol metadata - this is used to display company name, exchange, etc. in the UI
func (r *Repository) GetSymbolMetadata(symbol string) (*models.SymbolMetadata, error) {
	var meta models.SymbolMetadata
	err := r.db.QueryRow(`
		SELECT symbol, company_name, currency, exchange, last_fetched
		FROM symbol_metadata WHERE symbol = $1
	`, symbol).Scan(&meta.Symbol, &meta.CompanyName, &meta.Currency, &meta.Exchange, &meta.LastFetched)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

// check if we have recent data for this symbol - if not, we should fetch from Yahoo
func (r *Repository) IsDataStale(symbol string) bool {
	var lastFetched time.Time
	err := r.db.QueryRow(
		`SELECT MAX(fetched_at) FROM price_history WHERE symbol = $1`, symbol,
	).Scan(&lastFetched)
	if err != nil || lastFetched.IsZero() {
		return true
	}
	return time.Since(lastFetched) > time.Hour
}
