package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"PayCore/config"
	"PayCore/internal/transaction"
	"PayCore/internal/user"
	"PayCore/internal/wallet"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	db := config.ConnectDB()
	redisClient := config.ConnectRedis()
	defer redisClient.Close()

	err := db.AutoMigrate(
		&user.User{},
		&wallet.Wallet{},
		&transaction.Transaction{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("database migration completed successfully")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	walletRepository := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepository, db, redisClient)
	walletHandler := wallet.NewHandler(walletService)

	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)
	transactionHandler := transaction.NewHandler(transactionService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    "paycore",
		})
	})

	router.POST("/api/v1/users", userHandler.Create)

	router.POST("/api/v1/wallets", walletHandler.Create)
	router.GET("/api/v1/wallets/:id", walletHandler.GetByID)
	router.POST("/api/v1/wallets/:id/deposit", walletHandler.Deposit)
	router.POST("/api/v1/wallets/:id/withdraw", walletHandler.Withdraw)
	router.POST("/api/v1/transfers", walletHandler.Transfer)

	router.GET("/api/v1/transactions/:id", transactionHandler.GetByID)
	router.GET("/api/v1/wallets/:id/transactions", transactionHandler.GetByWalletID)

	log.Println("PayCore API running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
