package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"trading-dashboard/internal/config"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var database *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		database, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Attempt %d: failed to open DB: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if err = database.Ping(); err != nil {
			log.Printf("Attempt %d: DB not ready: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("✅ Database connected successfully")
		return database
	}

	log.Fatalf("❌ Could not connect to database: %v", err)
	return nil
}

func Migrate(database *sql.DB) {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id            SERIAL PRIMARY KEY,
			email         VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name          VARCHAR(255) NOT NULL,
			created_at    TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS price_history (
			id          SERIAL PRIMARY KEY,
			symbol      VARCHAR(20)   NOT NULL,
			trade_date  DATE          NOT NULL,
			open_price  DECIMAL(14,4),
			high_price  DECIMAL(14,4),
			low_price   DECIMAL(14,4),
			close_price DECIMAL(14,4),
			volume      BIGINT,
			fetched_at  TIMESTAMP DEFAULT NOW(),
			UNIQUE(symbol, trade_date)
		)`,
		`CREATE TABLE IF NOT EXISTS symbol_metadata (
			symbol       VARCHAR(20) PRIMARY KEY,
			company_name VARCHAR(255),
			currency     VARCHAR(10),
			exchange     VARCHAR(50),
			last_fetched TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_price_symbol     ON price_history(symbol)`,
		`CREATE INDEX IF NOT EXISTS idx_price_trade_date ON price_history(trade_date)`,
	}

	for _, stmt := range statements {
		if _, err := database.Exec(stmt); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	}
	log.Println("✅ Database migrations applied")
}