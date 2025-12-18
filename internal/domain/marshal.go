package domain

import (
	"database/sql"
	"encoding/json"
)

// NullTime wraps sql.NullTime with custom JSON marshaling
type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

// NullBool wraps sql.NullBool with custom JSON marshaling
type NullBool struct {
	sql.NullBool
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

func (nb *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nb.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nb.Bool)
	if err != nil {
		return err
	}
	nb.Valid = true
	return nil
}
