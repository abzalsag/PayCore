package transaction

import "time"

type Transaction struct {
	ID           uint `gorm:"primaryKey"`
	FromWalletID *uint
	ToWalletID   *uint
	Amount       int64  `gorm:"not null"`
	Type         string `gorm:"not null"`
	Status       string `gorm:"not null"`
	CreatedAt    time.Time
}
