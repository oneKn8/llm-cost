package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type UsageEntry struct {
	ID           int64     `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	CachedTokens int       `json:"cached_tokens"`
	Cost         float64   `json:"cost"`
}

type QueryFilters struct {
	Since    time.Time
	Provider string
	Model    string
}

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := migrate(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func migrate(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS usage (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			provider TEXT NOT NULL,
			model TEXT NOT NULL,
			input_tokens INTEGER NOT NULL,
			output_tokens INTEGER NOT NULL,
			cached_tokens INTEGER NOT NULL DEFAULT 0,
			cost REAL NOT NULL
		);
		CREATE TABLE IF NOT EXISTS budget (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			monthly_limit REAL NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_usage_timestamp ON usage(timestamp);
		CREATE INDEX IF NOT EXISTS idx_usage_provider ON usage(provider);
	`)
	return err
}

func (d *DB) RecordUsage(e UsageEntry) error {
	_, err := d.conn.Exec(
		`INSERT INTO usage (timestamp, provider, model, input_tokens, output_tokens, cached_tokens, cost) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.Timestamp.Format(time.RFC3339), e.Provider, e.Model, e.InputTokens, e.OutputTokens, e.CachedTokens, e.Cost,
	)
	return err
}

func (d *DB) QueryUsage(f QueryFilters) ([]UsageEntry, error) {
	query := `SELECT id, timestamp, provider, model, input_tokens, output_tokens, cached_tokens, cost FROM usage WHERE 1=1`
	args := []interface{}{}

	if !f.Since.IsZero() {
		query += ` AND timestamp >= ?`
		args = append(args, f.Since.Format(time.RFC3339))
	}
	if f.Provider != "" {
		query += ` AND provider = ?`
		args = append(args, f.Provider)
	}
	if f.Model != "" {
		query += ` AND model = ?`
		args = append(args, f.Model)
	}

	query += ` ORDER BY timestamp DESC`

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []UsageEntry
	for rows.Next() {
		var e UsageEntry
		var ts string
		if err := rows.Scan(&e.ID, &ts, &e.Provider, &e.Model, &e.InputTokens, &e.OutputTokens, &e.CachedTokens, &e.Cost); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse(time.RFC3339, ts)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (d *DB) SetBudget(limit float64) error {
	_, err := d.conn.Exec(
		`INSERT INTO budget (id, monthly_limit) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET monthly_limit = excluded.monthly_limit`,
		limit,
	)
	return err
}

func (d *DB) GetBudget() (float64, error) {
	var limit float64
	err := d.conn.QueryRow(`SELECT monthly_limit FROM budget WHERE id = 1`).Scan(&limit)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return limit, err
}
