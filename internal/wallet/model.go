package wallet

import "time"

type Wallet struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Currency  string `gorm:"type:varchar(3);not null"`
	Balance   int64  `gorm:"not null;default:0"`
	CreatedAt time.Time
}
