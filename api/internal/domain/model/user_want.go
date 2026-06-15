package model

type UserWant struct {
	UserID  uint `gorm:"primaryKey"`
	ThingID uint `gorm:"primaryKey"`
}
