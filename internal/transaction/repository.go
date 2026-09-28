package transaction

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(tx *gorm.DB, transaction *Transaction) error {
	return tx.Create(transaction).Error
}

func (r *Repository) GetByID(id uint) (*Transaction, error) {
	var transaction Transaction

	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *Repository) GetByWalletID(walletID uint) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.
		Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
