package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type UintSlice []uint

func (u UintSlice) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

func (u *UintSlice) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return json.Unmarshal(b, u)
}

type UserPickupCache struct {
	UserID        uint      `gorm:"primaryKey"`
	CacheDate     time.Time `gorm:"primaryKey;type:date"`
	PickedUserIDs UintSlice `gorm:"type:json;not null"`
	CreatedAt     time.Time
}
