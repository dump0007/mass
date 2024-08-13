package controllers

import (
	// "context"
	"context"
	"fmt"
	// "github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	// "github.com/xuri/excelize/v2"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "log"
	"mass/database"
	"mass/models"
	// "math"
	// "math/big"
	// "mime/multipart"
	"net/http"
	// "path/filepath"
	// "sync"
	"time"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/ethclient"
)



func CreateNftOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the wallet address from the request
		walletAddress := c.PostForm("WalletAddress")

		// Parse the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
			return
		}

		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// Parse and clean wallet addresses from the file
		addresses, err := parseAndCleanWalletAddresses(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate the next OrderID
		orderID, err := generateOrderID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OrderID"})
			return
		}

		// Set the default OrderStatus to "inactive"
		orderStatus := "inactive"

		// Set the filename as "OrderID.xlsx"
		filename := fmt.Sprintf("%d.xlsx", int(orderID))

		// Save the cleaned addresses to a new Excel file in the root directory
		if err := saveAddressesToFile(addresses, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Generate a new private key and derive the public address
		privateKeyHex, publicAddress, err := generateNewEVMAccount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate EVM account"})
			return
		}

		// Create a new Order
		order := models.Order{
			ID:            primitive.NewObjectID(),
			OrderID:       orderID,
			OrderStatus:   orderStatus,
			OrderType: 		"nft",
			WalletAddress: walletAddress,
			FileName:      filename,
			PrivateKey:    privateKeyHex,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Insert the order into the database
		orderCollection := database.OpenCollection(database.Client, "orders")
		_, err = orderCollection.InsertOne(context.Background(), order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order into database"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": orderID, "public_address": publicAddress})
	}
}
