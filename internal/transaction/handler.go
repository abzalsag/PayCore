package transaction

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type TransactionResponse struct {
	ID           uint   `json:"id"`
	FromWalletID *uint  `json:"from_wallet_id"`
	ToWalletID   *uint  `json:"to_wallet_id"`
	Amount       int64  `json:"amount"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction id",
		})
		return
	}

	transaction, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		ID:           transaction.ID,
		FromWalletID: transaction.FromWalletID,
		ToWalletID:   transaction.ToWalletID,
		Amount:       transaction.Amount,
		Type:         transaction.Type,
		Status:       transaction.Status,
	})
}

func (h *Handler) GetByWalletID(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid wallet id",
		})
		return
	}

	transactions, err := h.service.GetByWalletID(uint(walletID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		response = append(response, TransactionResponse{
			ID:           transaction.ID,
			FromWalletID: transaction.FromWalletID,
			ToWalletID:   transaction.ToWalletID,
			Amount:       transaction.Amount,
			Type:         transaction.Type,
			Status:       transaction.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}
