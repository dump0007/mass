package controllers

import (
	// "context"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mass/database"
	"mass/models"
	// "math"
	"math/big"
	"mime/multipart"
	"net/http"
	"path/filepath"
	// "sync"
	"time"
	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const maxAddresses = 1000000
const chunkSize = 100000

// UploadExcel handles the uploading and processing of an Excel file with wallet addresses
// func UploadExcel() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("1111111111111")
// 		// Create a context with timeout for the database operation
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		fmt.Println("2222222222222")

// 		// Parse the request body for additional fields
// 		var requestBody struct {
// 			WalletAddress string `json:"wallet_address" binding:"required"`
// 		}
// 		fmt.Println("333333333333333")

// 		// Bind the form data to requestBody struct
// 		if err := c.ShouldBind(&requestBody); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 			return
// 		}
// 		fmt.Println("44444444444444444")

// 		fmt.Println(requestBody.WalletAddress)



// func EstimateGas() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
// 		defer cancel()

// 		fmt.Println("11111111111111111")
// 		client, err := ethclient.Dial("https://ethereum-holesky-rpc.publicnode.com")
// 		if err != nil {
// 			log.Printf("Failed to connect to the Ethereum client: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the Ethereum client"})
// 			return
// 		}
// 		fmt.Println("2222222222222222222")

// 		// Parse request body
// 		var requestBody struct {
// 			Message       string `json:"Message" binding:"required"`
// 			WalletAddress string `json:"WalletAddress" binding:"required"`
// 		}
// 		if err := c.ShouldBindJSON(&requestBody); err != nil {
// 			log.Printf("Invalid request payload: %v\n", err)
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 			return
// 		}
// 		fmt.Println("333333333333333333", requestBody)



// 		// toAddress := common.HexToAddress(addresses[0])
// 		toAddress := common.HexToAddress(requestBody.WalletAddress)

// 		// fmt.Println("88888888888888888888888888888")

// 		// Convert message to hex
// 		data := []byte(requestBody.Message)
// 		fmt.Println("9999999999999999", data)

// 		// Prepare the transaction
// 		fromAddress := common.HexToAddress(requestBody.WalletAddress)
// 		msg := ethereum.CallMsg{
// 			From: fromAddress,
// 			To:   &toAddress,
// 			Data: data,
// 		}
// 		fmt.Println("10101010101010101010101", msg)

// 		// Estimate the gas required for the transaction
// 		gasLimit, err := client.EstimateGas(ctx, msg)
// 		if err != nil {
// 			log.Printf("Failed to estimate gas: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to estimate gas", "details": err.Error()})
// 			return
// 		}
// 		fmt.Println("1011101110110101010110101011111110000011111111")

// 		// Get the current gas price
// 		gasPrice, err := client.SuggestGasPrice(ctx)
// 		if err != nil {
// 			log.Printf("Failed to get gas price: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price", "details": err.Error()})
// 			return
// 		}
// 		fmt.Println("12121212121212121212121212")

// 		// Calculate the total gas cost
// 		totalGasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))

// 		// Calculate the total cost for the number of addresses
// 		numAddresses := 1000000
// 		totalCost := new(big.Int).Mul(totalGasCost, big.NewInt(int64(numAddresses)))
// 		// Step 2: Create a big integer representing 10^18
// 		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// 		// Convert to big.Float for precise division
// 	totalCostFloat := new(big.Float).SetInt(totalCost)
// 	divisorFloat := new(big.Float).SetInt(divisor)

// 	// Step 3: Perform the division with big.Float
// 	finalCostFloat := new(big.Float).Quo(totalCostFloat, divisorFloat)

// 	// Set precision to 4 decimal places
// 	finalCostFloat = finalCostFloat.SetPrec(4 * 8) // 4 decimal places

// 	// Convert the result to a string with 4 decimal places
// 	finalCostString := fmt.Sprintf("%.4f", finalCostFloat)
// 		// finalCost := new(big.Int).Mul(totalGasCost, big.NewInt(int64(numAddresses)))
// 		fmt.Println("totalCost",totalCost,"divisor",divisor,"finalCost",finalCostString)

// 		c.JSON(http.StatusOK, gin.H{
// 			"gasLimit":     gasLimit,
// 			"gasPrice":     gasPrice.String(),
// 			"totalGasCost": totalGasCost.String(),
// 			"totalCost":    totalCost.String(),
// 			"numAddresses": numAddresses,
// 			"finalCost":	finalCostString,
// 		})
// 	}
// }


func EstimateGas() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		client, err := ethclient.Dial("https://ethereum-holesky-rpc.publicnode.com")
		if err != nil {
			log.Printf("Failed to connect to the Ethereum client: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the Ethereum client"})
			return
		}
		// Parse request body
		var requestBody struct {
			OrderID       float64 `json:"order_id" binding:"required"`
			Message       string  `json:"Message" binding:"required"`
			WalletAddress string  `json:"WalletAddress" binding:"required"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			log.Printf("Invalid request payload: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}


	

		orderCollection := database.OpenCollection(database.Client, "orders")

		var order models.Order
		filter := bson.M{"order_id": requestBody.OrderID}
		update := bson.M{"$set": bson.M{"message": requestBody.Message}}

		err = orderCollection.FindOneAndUpdate(ctx, filter, update).Decode(&order)
		if err != nil {
			log.Printf("Failed to find and update order: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find and update order"})
			return
		}

		// Prepare the transaction
		toAddress := common.HexToAddress(requestBody.WalletAddress)
		data := []byte(requestBody.Message)
		fromAddress := common.HexToAddress(requestBody.WalletAddress)
		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: data,
		}

		// Estimate the gas required for the transaction
		gasLimit, err := client.EstimateGas(ctx, msg)
		if err != nil {
			log.Printf("Failed to estimate gas: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to estimate gas", "details": err.Error()})
			return
		}
		// Get the current gas price
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			log.Printf("Failed to get gas price: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price", "details": err.Error()})
			return
		}
		// Calculate the total gas cost
		totalGasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
		// Multiply the total gas cost by the count value from the order
		totalCost := new(big.Int).Mul(totalGasCost, big.NewInt(int64(order.Count)))
		// Convert the total cost to Ether
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
		totalCostFloat := new(big.Float).SetInt(totalCost)
		divisorFloat := new(big.Float).SetInt(divisor)
		finalCostFloat := new(big.Float).Quo(totalCostFloat, divisorFloat)
		finalCostString := fmt.Sprintf("%.4f", finalCostFloat)

			// finalCost := new(big.Int).Mul(totalGasCost, big.NewInt(int64(numAddresses)))
			fmt.Println("totalCost",totalCost,"divisor",divisor,"finalCost",finalCostString)

			c.JSON(http.StatusOK, gin.H{
				"gasLimit":     gasLimit,
				"gasPrice":     gasPrice.String(),
				"totalGasCost": totalGasCost.String(),
				"totalCost":    totalCost.String(),
				"finalCost":	finalCostString,
			})
		}
	}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout for the database operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Define a struct to bind the request body
		var requestBody struct {
			WalletAddress string `json:"walletAddress" binding:"required"`
		}

		// Bind the JSON body to requestBody struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		fmt.Println("wallet ->", requestBody.WalletAddress)

		// Access the users collection
		usersCollection := database.OpenCollection(database.Client, "users")

		// Check if the wallet address exists in the database
		var existingUser models.User
		err := usersCollection.FindOne(ctx, bson.M{"walletaddress": requestBody.WalletAddress}).Decode(&existingUser)
		fmt.Println("err -> ", requestBody.WalletAddress, err)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Wallet address not found, register the new user
				newUser := models.User{
					WalletAddress: requestBody.WalletAddress,
					Created_at:    time.Now(),
				}

				// Insert the new user document
				_, err := usersCollection.InsertOne(ctx, newUser)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register new user"})
					return
				}

				// Respond with registration success
				c.JSON(http.StatusOK, gin.H{"message": "Wallet registered successfully"})
				return
			}

			// Handle other potential database errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user database"})
			return
		}

		// Respond with success if the wallet address is already connected
		c.JSON(http.StatusOK, gin.H{"message": "Wallet connected successfully"})
	}
}



func generateOrderID() (float64, error) {
	// Find the highest existing OrderID in the collection
	opts := options.FindOne().SetSort(primitive.D{{Key: "order_id", Value: -1}})
	orderCollection := database.OpenCollection(database.Client, "orders")
	var lastOrder models.Order
	err := orderCollection.FindOne(context.Background(), primitive.M{}, opts).Decode(&lastOrder)
	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	// Increment the OrderID by 1
	return lastOrder.OrderID + 1, nil
}

func saveFile(c *gin.Context, file *multipart.FileHeader, filename string) error {
	// Create the file path
	filePath := filepath.Join(".", filename)

	// Save the file to the specified path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return err
	}

	return nil
}

// generateNewEVMAccount generates a new EVM private key and corresponding public address
func generateNewEVMAccount() (string, string, error) {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// Convert the private key to a hexadecimal string
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)

	// Derive the public address from the private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return privateKeyHex, publicAddress, nil
}

// ExecuteTxns checks the order and user wallet address and updates the order status if applicable.
func ExecuteTxns() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout for the database operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var requestBody struct {
			WalletAddress string  `json:"WalletAddress" binding:"required"`
			OrderID       float64 `json:"OrderID" binding:"required"`
		}

		// Bind the JSON request body to the Order struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		fmt.Println("requestOrder ->", requestBody.WalletAddress)
		userCollection := database.OpenCollection(database.Client, "users")
		orderCollection := database.OpenCollection(database.Client, "orders")

		// Check if the wallet address exists in the User collection
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"walletaddress": requestBody.WalletAddress}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with the given wallet address not found"})
			return
		}

		// Check if the order exists in the Order collection and matches the wallet address and order ID
		var order models.Order
		filter := bson.M{"wallet_address": requestBody.WalletAddress, "order_id": requestBody.OrderID}
		err = orderCollection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found with the given wallet address and order ID"})
			return
		}

		// Check if the order status is inactive
		if order.OrderStatus == "inactive" {
			// Update the order status to active
			update := bson.M{
				"$set": bson.M{
					"order_status": "active",
					"updated_at":   time.Now(),
				},
			}
			_, err = orderCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Order status updated to active"})
		} else {
			// Return an error if the order status is not inactive
			c.JSON(http.StatusConflict, gin.H{"error": "Order has already been activated or is in another stage of being complete"})
		}
	}
}

func GetOrders() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Create a context with timeout for the database operation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var requestBody struct {
			WalletAddress string `json:"WalletAddress" binding:"required"`
		}

		// Bind the JSON request body to the Order struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		fmt.Println("requestOrder ->", requestBody.WalletAddress)
		userCollection := database.OpenCollection(database.Client, "users")
		orderCollection := database.OpenCollection(database.Client, "orders")

		// Check if the wallet address exists in the User collection
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"walletaddress": requestBody.WalletAddress}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User with the given wallet address not found"})
			return
		}

		// Check if the order exists in the Order collection and matches the wallet address and order ID
		// var order models.Order
		// var order struct {
		// 	OrderID       float64 `json:"OrderID" binding:"required"`
		// 	OrderStatus   string  `json:"OrderStatus" binding:"required"`
		// 	WalletAddress string  `json:"WalletAddress" binding:"required"`
		// }
		filter := bson.M{"wallet_address": requestBody.WalletAddress}
		projection := bson.M{
			// "OrderID": 1,
			// "OrderStatus": 1,
			// "WalletAddress": 1,
			"id":            0,
			"gas_collected": 0,
			"gas_used":      0,
			"nonce":         0,
			"file_name":     0,
			"pvt_key":       0,
			"addresses":     0,
			"created_at":    0,
			"updated_at":    0,
		}
		// err = orderCollection.FindOne(ctx, filter,options.FindOne().SetProjection(projection)).Decode(&order)
		cursor, err := orderCollection.Find(ctx, filter, options.Find().SetProjection(projection))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Orders not found with the given wallet address"})
			return
		}
		// defer cursor.Close(ctx)
		var orders []models.Order

		if err = cursor.All(context.TODO(), &orders); err != nil {
			panic(err)
		}
		fmt.Println("err", orders)
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"error": false,
			"data":  orders,
			// "walletADdress":order.WalletAddress,
			// "orderId":order.OrderID,
			// "orderStatus":order.OrderStatus,
			"message": "Order fetched successfully",
		})
	}
}


func isValidChecksumAddress(address string) bool {
	return common.IsHexAddress(address) && address == common.HexToAddress(address).Hex()
}



func SendMsg() gin.HandlerFunc {
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
			OrderType: 		"msg",
			Count:			len(addresses),
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

// parseAndCleanWalletAddresses reads wallet addresses from an Excel file, removes duplicates and invalid addresses.
func parseAndCleanWalletAddresses(file multipart.File) ([]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read the excel file: %w", err)
	}

	var addresses []string
	addressMap := make(map[string]bool)

	// Get the name of the first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("unable to get rows from the sheet: %w", err)
	}

	for _, row := range rows {
		if len(row) > 0 {
			address := row[0]
			if !isValidChecksumAddress(address) {
				return nil, fmt.Errorf("invalid wallet address: %s", address)
			}
			if !addressMap[address] {
				addressMap[address] = true
				addresses = append(addresses, address)
			}
		}
	}

	return addresses, nil
}

// saveAddressesToFile saves the cleaned wallet addresses to a new Excel file.
func saveAddressesToFile(addresses []string, filename string) error {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Write the cleaned addresses to the new Excel file
	for i, address := range addresses {
		cell := fmt.Sprintf("A%d", i+2) // Starting from the second row to leave space for headers
		if err := f.SetCellValue(sheetName, cell, address); err != nil {
			return fmt.Errorf("failed to write address to Excel file: %w", err)
		}
	}

	// Insert headers in the first row
	headers := []string{"wallet_address", "txn_hash", "gas_consumed", "txn_status", "order_id"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i) // A1, B1, C1, etc.
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set header: %w", err)
		}
	}

	// Save the modified file to the root directory
	savePath := filepath.Join(".", filename)
	if err := f.SaveAs(savePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}



// CheckAddresses handles the API request to validate and clean wallet addresses from an Excel file.
func CheckAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the wallet address from the request body (though not used here, this is an example of handling it)
		walletAddress := c.PostForm("WalletAddress")
		fmt.Println("Received Wallet Address:", walletAddress)

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

		// Check and clean the wallet addresses
		validAddresses,err := parseAndCleanWalletAddresses(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Return the results
		c.JSON(http.StatusOK, gin.H{
			"valid_address_count": len(validAddresses),
		})
	}
}
