package model

type ThingAlias struct {
	ThingID uint   `gorm:"primaryKey"`
	Alias   string `gorm:"primaryKey"`
}
