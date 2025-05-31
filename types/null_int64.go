package types

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	ni.Int64 = i
	ni.Valid = true
	return nil
}
