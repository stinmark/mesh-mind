package store

import (
	"database/sql"
	"time"
)

type HistoryRecord struct {
	ID        int       `json:"id"`
	InputText string    `json:"input_text"`
	Intent    string    `json:"intent"`
	Sentiment string    `json:"sentiment"`
	ModelUsed string    `json:"model_used"`
	LatencyMs float64   `json:"latency_ms"`
	CreatedAt time.Time `json:"created_at"`
}

type HistoryStore struct {
	db *sql.DB
}

func NewHistoryStore(db *sql.DB) *HistoryStore {
	return &HistoryStore{db: db}
}

func (s *HistoryStore) AddRecord(input, intent, sentiment, model string, latency float64) error {
	query := `INSERT INTO history (input_text, intent, sentiment, model_used, latency_ms)
	          VALUES (?, ?, ?, ?, ?);`
	_, err := s.db.Exec(query, input, intent, sentiment, model, latency)
	return err
}

func (s *HistoryStore) GetRecent(limit int) ([]HistoryRecord, error) {
	query := `SELECT id, input_text, intent, sentiment, model_used, latency_ms, created_at 
	          FROM history ORDER BY id DESC LIMIT ?;`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []HistoryRecord
	for rows.Next() {
		var r HistoryRecord
		if err := rows.Scan(&r.ID, &r.InputText, &r.Intent, &r.Sentiment, &r.ModelUsed, &r.LatencyMs, &r.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}
