package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"PayCore/internal/transaction"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repository    *Repository
	transactionDB *gorm.DB
	redis         *redis.Client
}

func NewService(repository *Repository, db *gorm.DB, redisClient *redis.Client) *Service {
	return &Service{
		repository:    repository,
		transactionDB: db,
		redis:         redisClient,
	}
}

func (s *Service) CreateWallet(userID uint, currency string) (*Wallet, error) {
	wallet := &Wallet{
		UserID:   userID,
		Currency: currency,
		Balance:  0,
	}

	err := s.repository.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Service) GetWalletByID(id uint) (*Wallet, error) {
	ctx := context.Background()
	key := "wallet:" + strconv.Itoa(int(id))

	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var wallet Wallet

		if err := json.Unmarshal([]byte(cached), &wallet); err == nil {
			return &wallet, nil
		}
	}

	wallet, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(wallet)
	if err == nil {
		s.redis.Set(ctx, key, data, 5*time.Minute)
	}

	return wallet, nil
}

func (s *Service) Deposit(walletID uint, amount int64) (*Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	var wallet *Wallet

	err := s.transactionDB.Transaction(func(tx *gorm.DB) error {
		var err error

		wallet, err = s.repository.GetByID(walletID)
		if err != nil {
			return err
		}

		wallet.Balance += amount

		if err := s.repository.UpdateBalance(tx, wallet); err != nil {
			return err
		}

		transactionRepository := transaction.NewRepository(tx)

		newTransaction := &transaction.Transaction{
			ToWalletID: &wallet.ID,
			Amount:     amount,
			Type:       "DEPOSIT",
			Status:     "SUCCESS",
		}

		if err := transactionRepository.Create(tx, newTransaction); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.deleteWalletCache(walletID)

	return wallet, nil
}

func (s *Service) Withdraw(walletID uint, amount int64) (*Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	var wallet *Wallet

	err := s.transactionDB.Transaction(func(tx *gorm.DB) error {
		var err error

		wallet, err = s.repository.GetByID(walletID)
		if err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		wallet.Balance -= amount

		if err := s.repository.UpdateBalance(tx, wallet); err != nil {
			return err
		}

		transactionRepository := transaction.NewRepository(tx)

		newTransaction := &transaction.Transaction{
			FromWalletID: &wallet.ID,
			Amount:       amount,
			Type:         "WITHDRAW",
			Status:       "SUCCESS",
		}

		if err := transactionRepository.Create(tx, newTransaction); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.deleteWalletCache(walletID)

	return wallet, nil
}

func (s *Service) Transfer(fromWalletID, toWalletID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if fromWalletID == toWalletID {
		return errors.New("cannot transfer to the same wallet")
	}

	err := s.transactionDB.Transaction(func(tx *gorm.DB) error {
		var fromWallet Wallet
		var toWallet Wallet

		if err := tx.First(&fromWallet, fromWalletID).Error; err != nil {
			return errors.New("sender wallet not found")
		}

		if err := tx.First(&toWallet, toWalletID).Error; err != nil {
			return errors.New("receiver wallet not found")
		}

		if fromWallet.Currency != toWallet.Currency {
			return errors.New("wallet currencies must be the same")
		}

		if fromWallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		fromWallet.Balance -= amount
		toWallet.Balance += amount

		if err := tx.Model(&fromWallet).Update("balance", fromWallet.Balance).Error; err != nil {
			return err
		}

		if err := tx.Model(&toWallet).Update("balance", toWallet.Balance).Error; err != nil {
			return err
		}

		newTransaction := &transaction.Transaction{
			FromWalletID: &fromWallet.ID,
			ToWalletID:   &toWallet.ID,
			Amount:       amount,
			Type:         "TRANSFER",
			Status:       "SUCCESS",
		}

		return tx.Create(newTransaction).Error
	})

	if err != nil {
		return err
	}

	s.deleteWalletCache(fromWalletID)
	s.deleteWalletCache(toWalletID)

	return nil
}

func (s *Service) deleteWalletCache(walletID uint) {
	ctx := context.Background()
	key := "wallet:" + strconv.Itoa(int(walletID))

	s.redis.Del(ctx, key)
}
