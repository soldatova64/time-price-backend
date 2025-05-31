package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (ns NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Time)
}

func (ns *NullTime) UnmarshalJSON(data []byte) error {
	// Если пришел null
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}

	// Пробуем распарсить время
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	ns.Time = t
	ns.Valid = true
	return nil
}
