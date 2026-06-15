package model

type UserDoing struct {
	UserID  uint `gorm:"primaryKey"`
	ThingID uint `gorm:"primaryKey"`
}
