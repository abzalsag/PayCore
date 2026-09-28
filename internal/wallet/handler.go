package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type TransferRequest struct {
	FromWalletID uint  `json:"from_wallet_id" binding:"required"`
	ToWalletID   uint  `json:"to_wallet_id" binding:"required"`
	Amount       int64 `json:"amount" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateWalletRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required,len=3"`
}

type DepositRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

type WalletResponse struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateWalletRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	wallet, err := h.service.CreateWallet(
		req.UserID,
		req.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Currency: wallet.Currency,
		Balance:  wallet.Balance,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid wallet id",
		})
		return
	}

	wallet, err := h.service.GetWalletByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "wallet not found",
		})
		return
	}

	c.JSON(http.StatusOK, WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Currency: wallet.Currency,
		Balance:  wallet.Balance,
	})
}

func (h *Handler) Deposit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid wallet id",
		})
		return
	}

	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	wallet, err := h.service.Deposit(uint(id), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Currency: wallet.Currency,
		Balance:  wallet.Balance,
	})
}

func (h *Handler) Withdraw(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid wallet id",
		})
		return
	}

	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	wallet, err := h.service.Withdraw(uint(id), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Currency: wallet.Currency,
		Balance:  wallet.Balance,
	})
}

func (h *Handler) Transfer(c *gin.Context) {
	var req TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Transfer(
		req.FromWalletID,
		req.ToWalletID,
		req.Amount,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transfer successful",
	})
}
