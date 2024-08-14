package controllers

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"bytes"
	"github.com/xuri/excelize/v2"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mass/database"
	"mass/models"
	"path/filepath"
	"strconv"
	"time"
)

const (
	contractAddress = "0x7CFDFf5b6b8D450682b1FC2182B81231E3A2468c"
	CONTRACT_ABI    = `[
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"internalType": "address[]",
					"name": "_to",
					"type": "address[]"
				},
				{
					"internalType": "uint256",
					"name": "id",
					"type": "uint256"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "BatchsafeTransferFrom",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "initialOwner",
					"type": "address"
				}
			],
			"stateMutability": "nonpayable",
			"type": "constructor"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "sender",
					"type": "address"
				},
				{
					"internalType": "uint256",
					"name": "balance",
					"type": "uint256"
				},
				{
					"internalType": "uint256",
					"name": "needed",
					"type": "uint256"
				},
				{
					"internalType": "uint256",
					"name": "tokenId",
					"type": "uint256"
				}
			],
			"name": "ERC1155InsufficientBalance",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "approver",
					"type": "address"
				}
			],
			"name": "ERC1155InvalidApprover",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "idsLength",
					"type": "uint256"
				},
				{
					"internalType": "uint256",
					"name": "valuesLength",
					"type": "uint256"
				}
			],
			"name": "ERC1155InvalidArrayLength",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "operator",
					"type": "address"
				}
			],
			"name": "ERC1155InvalidOperator",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "receiver",
					"type": "address"
				}
			],
			"name": "ERC1155InvalidReceiver",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "sender",
					"type": "address"
				}
			],
			"name": "ERC1155InvalidSender",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "operator",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "owner",
					"type": "address"
				}
			],
			"name": "ERC1155MissingApprovalForAll",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "account",
					"type": "address"
				},
				{
					"internalType": "uint256",
					"name": "_tokenId",
					"type": "uint256"
				},
				{
					"internalType": "uint256",
					"name": "amount",
					"type": "uint256"
				},
				{
					"internalType": "string",
					"name": "_tokenURI",
					"type": "string"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "mint",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"internalType": "uint256[]",
					"name": "ids",
					"type": "uint256[]"
				},
				{
					"internalType": "uint256[]",
					"name": "amounts",
					"type": "uint256[]"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "mintBatch",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "owner",
					"type": "address"
				}
			],
			"name": "OwnableInvalidOwner",
			"type": "error"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "account",
					"type": "address"
				}
			],
			"name": "OwnableUnauthorizedAccount",
			"type": "error"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "account",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "operator",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "bool",
					"name": "approved",
					"type": "bool"
				}
			],
			"name": "ApprovalForAll",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "previousOwner",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "newOwner",
					"type": "address"
				}
			],
			"name": "OwnershipTransferred",
			"type": "event"
		},
		{
			"inputs": [],
			"name": "renounceOwnership",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"internalType": "uint256[]",
					"name": "ids",
					"type": "uint256[]"
				},
				{
					"internalType": "uint256[]",
					"name": "values",
					"type": "uint256[]"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "safeBatchTransferFrom",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"internalType": "uint256",
					"name": "id",
					"type": "uint256"
				},
				{
					"internalType": "uint256",
					"name": "value",
					"type": "uint256"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "safeTransferFrom",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "operator",
					"type": "address"
				},
				{
					"internalType": "bool",
					"name": "approved",
					"type": "bool"
				}
			],
			"name": "setApprovalForAll",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "operator",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256[]",
					"name": "ids",
					"type": "uint256[]"
				},
				{
					"indexed": false,
					"internalType": "uint256[]",
					"name": "values",
					"type": "uint256[]"
				}
			],
			"name": "TransferBatch",
			"type": "event"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "newOwner",
					"type": "address"
				}
			],
			"name": "transferOwnership",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"internalType": "address",
					"name": "operator",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "from",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "to",
					"type": "address"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "id",
					"type": "uint256"
				},
				{
					"indexed": false,
					"internalType": "uint256",
					"name": "value",
					"type": "uint256"
				}
			],
			"name": "TransferSingle",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": false,
					"internalType": "string",
					"name": "value",
					"type": "string"
				},
				{
					"indexed": true,
					"internalType": "uint256",
					"name": "id",
					"type": "uint256"
				}
			],
			"name": "URI",
			"type": "event"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "account",
					"type": "address"
				},
				{
					"internalType": "uint256",
					"name": "id",
					"type": "uint256"
				}
			],
			"name": "balanceOf",
			"outputs": [
				{
					"internalType": "uint256",
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address[]",
					"name": "accounts",
					"type": "address[]"
				},
				{
					"internalType": "uint256[]",
					"name": "ids",
					"type": "uint256[]"
				}
			],
			"name": "balanceOfBatch",
			"outputs": [
				{
					"internalType": "uint256[]",
					"name": "",
					"type": "uint256[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "address",
					"name": "account",
					"type": "address"
				},
				{
					"internalType": "address",
					"name": "operator",
					"type": "address"
				}
			],
			"name": "isApprovedForAll",
			"outputs": [
				{
					"internalType": "bool",
					"name": "",
					"type": "bool"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "owner",
			"outputs": [
				{
					"internalType": "address",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes4",
					"name": "interfaceId",
					"type": "bytes4"
				}
			],
			"name": "supportsInterface",
			"outputs": [
				{
					"internalType": "bool",
					"name": "",
					"type": "bool"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint256",
					"name": "tokenId",
					"type": "uint256"
				}
			],
			"name": "uri",
			"outputs": [
				{
					"internalType": "string",
					"name": "",
					"type": "string"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`
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
			OrderType:     "nft",
			Count:         len(addresses),
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

// UploadMetadata handles the image upload and updates the order in MongoDB
func UploadMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract order ID from the request
		orderID := c.PostForm("order_id")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		// Extract image file from the request
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			return
		}

		// Validate image format
		ext := filepath.Ext(file.Filename)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format. Only PNG, JPG, and JPEG are allowed."})
			return
		}

		// Save the file locally (temporarily)
		filePath := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
			return
		}
		defer os.Remove(filePath) // Remove the file after the function completes

		// Upload image to Pinata
		openedFile, err := readFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer openedFile.Close()

		resp, err := uploadToPinata(openedFile, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Pinata"})
			return
		}

		// Extract IPFS hash from Pinata response
		var pinataResp struct {
			IpfsHash string `json:"IpfsHash"`
		}
		if err := json.NewDecoder(bytes.NewReader(resp.Body())).Decode(&pinataResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Pinata response"})
			return
		}

		// Create JSON metadata
		metadata := map[string]string{
			"name":        "My First NFT",
			"description": "This is my first NFT on the ERC1155 standard.",
			"image":       fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", pinataResp.IpfsHash),
		}

		// Save JSON metadata to file
		jsonFilePath := "metadata.json"
		metadataFile, err := os.Create(jsonFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JSON file"})
			return
		}
		defer metadataFile.Close()

		if err := json.NewEncoder(metadataFile).Encode(metadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON to file"})
			return
		}

		// Upload JSON file to Pinata
		jsonFile, err := readFile(jsonFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open JSON file"})
			return
		}
		defer jsonFile.Close()

		jsonResp, err := uploadToPinata(jsonFile, jsonFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload JSON to Pinata"})
			return
		}

		// Extract IPFS hash of the JSON file
		var jsonPinataResp struct {
			IpfsHash string `json:"IpfsHash"`
		}
		if err := json.NewDecoder(bytes.NewReader(jsonResp.Body())).Decode(&jsonPinataResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Pinata JSON response"})
			return
		}

		// Update MongoDB order entry with the IPFS hash
		oid, err := strconv.ParseFloat(orderID, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
			return
		}

		filter := bson.M{"order_id": oid}
		update := bson.M{"$set": bson.M{"ipfs": fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", jsonPinataResp.IpfsHash)}}
		orderCollection := database.OpenCollection(database.Client, "orders")

		if _, err := orderCollection.UpdateOne(c, filter, update); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order with IPFS hash"})
			return
		}

		// Delete the local JSON file
		if err := os.Remove(jsonFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete local JSON file"})
			return
		}

		// Respond with the updated IPFS hash
		c.JSON(http.StatusOK, gin.H{
			"message":                 "Metadata uploaded and order updated successfully",
			"json_metadata_ipfs_hash": jsonPinataResp.IpfsHash,
		})
	}
}

// Helper function to read a file
func readFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

// Uploads a file to Pinata and returns the response (placeholder)
func uploadToPinata(file *os.File, filePath string) (*resty.Response, error) {
	apiKey := "1652e12bd8f5d89ef5cc"
	secretKey := "968e840d7a0526e19dc8d8b090fab28dde7c40c2c254fbe09fd9442561c5e65e"

	client := resty.New()

	resp, err := client.R().
		SetHeader("pinata_api_key", apiKey).
		SetHeader("pinata_secret_api_key", secretKey).
		SetFileReader("file", filepath.Base(filePath), file).
		SetFormData(map[string]string{
			"pinataMetadata": `{"name": "` + filepath.Base(filePath) + `", "keyvalues": {"uploadedBy": "backend"}}`,
			"pinataOptions":  `{"cidVersion": 0}`,
		}).
		Post("https://api.pinata.cloud/pinning/pinFileToIPFS")

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func EstimateNftGas() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 	}
// }

// EstimateNftGas handles the gas estimation process
func EstimateNftGas() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			WalletAddress string  `json:"wallet_address" binding:"required"`
			OrderID       float64 `json:"order_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		fmt.Println("req.WalletAddress ->",req.WalletAddress)
		fmt.Println("req.OrderID ->",req.OrderID)


		// Verify that the order exists
		var order models.Order
		orderCollection := database.OpenCollection(database.Client, "orders")
		err := orderCollection.FindOne(context.Background(), bson.M{
			"wallet_address": req.WalletAddress,
			"order_id":       req.OrderID,
		}).Decode(&order)

		fmt.Println("order.FileName ->",order.FileName)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Connect to Ethereum client
		client, err := ethclient.Dial("https://ethereum-holesky-rpc.publicnode.com")
		if err != nil {
			log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		}

		// Load and parse the ABI
		contractABI, err := abi.JSON(strings.NewReader(string(CONTRACT_ABI)))
		if err != nil {
			log.Fatalf("Failed to parse contract ABI: %v", err)
		}

		// Read the addresses from the file specified in the order
		addresses, err := readAddressesFromFile(order.FileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read addresses from file"})
			return
		}
		// Estimate gas for mint
		mintGas := estimateGasForMint(client, contractABI, order.PrivateKey)
		fmt.Println("mintGas ->",mintGas)
		// Estimate gas for setApprovalForAll
		approvalGas := estimateGasForSetApprovalForAll(client, contractABI, order.PrivateKey)
		fmt.Println("approvalGas ->",approvalGas)

		// Batch process addresses in groups of 500
		totalGasEstimation := new(big.Int)
		for i := 0; i < len(addresses); i += 500 {
			end := i + 500
			if end > len(addresses) {
				end = len(addresses)
			}
			batch := addresses[i:end]
			fmt.Println("batch ->",batch)
			

			// Estimate gas for BatchsafeTransferFrom
			batchGas := estimateGasForBatchsafeTransferFrom(client, contractABI, batch, order.PrivateKey,15454438)

			// Accumulate total gas cost
			totalGasEstimation.Add(totalGasEstimation,batchGas)
		}
		totalGasEstimation.Add(totalGasEstimation,mintGas)
		totalGasEstimation.Add(totalGasEstimation,approvalGas)


		// Return the total gas estimation
		c.JSON(http.StatusOK, gin.H{
			"total_gas_estimation_wei": totalGasEstimation.String(),
		})
	}
}

// readAddressesFromFile reads and returns addresses from the given file
func readAddressesFromFile(fileName string) ([]common.Address, error) {
    // Open the Excel file
    f, err := excelize.OpenFile(fileName)
    if err != nil {
        return nil, fmt.Errorf("failed to open Excel file: %w", err)
    }
	fmt.Println("1111111111111111111")
    // Get the name of the first sheet
    sheetName := f.GetSheetName(0)
    if sheetName == "" {
        return nil, fmt.Errorf("failed to get sheet name from Excel file")
    }
	fmt.Println("22222222222222222222")


    // Get all the rows in the first sheet
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return nil, fmt.Errorf("failed to get rows from sheet: %w", err)
    }
	fmt.Println("3333333333333333333333")

    var addresses []common.Address
     // Iterate over rows, starting from index 1 to skip the headers
	 for i, row := range rows {
        if i == 0 {
            continue // Skip the header row
        }
        if len(row) > 0 {
            // Assuming the address is in the first column
            address := common.HexToAddress(row[0])
            addresses = append(addresses, address)
        }
    }
	fmt.Println("length ->",len(addresses))

    return addresses, nil
}


func estimateGasForMint(client *ethclient.Client, parsedABI abi.ABI, privateKey string) *big.Int {
	// Prepare the transaction data
	tokenID := big.NewInt(1)                                                                       // Replace with your token ID
	amount := big.NewInt(1000000)                                                                  // Replace with the amount
	tokenURI := "https://gateway.pinata.cloud/ipfs/QmUDMBgg9mbQytooovDUtJ3uCdS2THXNSBLzqpzsuSan7s" // Replace with your token URI
	data := []byte{}                                                                               // Replace with your data if any

	// ownerAddress := deriveAddressFromPrivateKey(privateKey)
	operator := common.HexToAddress("0x6F4Cd3C8D1cF68242a66479f2F9A3F976A5cd4dD") // Replace with the operator address

	// for _, address := range addresses {
	input, err := parsedABI.Pack("mint", operator, tokenID, amount, tokenURI, data)
	if err != nil {
		log.Fatalf("Failed to pack the mint function input: %v", err)
	}

	toAddress := common.HexToAddress(contractAddress)
	msg := ethereum.CallMsg{
		From: operator,
		To:   &toAddress,
		Data: input,
	}

	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to estimate gas for mint: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	totalGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	fmt.Printf("Total Gas Cost for Mint (in wei): %s\n", totalGasCost.String())

	return totalGasCost
	// }

	// return big.NewInt(0)
}

func estimateGasForSetApprovalForAll(client *ethclient.Client, parsedABI abi.ABI, privateKey string) *big.Int {

	// operator := common.HexToAddress(deriveAddressFromPrivateKey(privateKey)) // Replace with the operator address
	operator := common.HexToAddress("0x6F4Cd3C8D1cF68242a66479f2F9A3F976A5cd4dD") // Replace with the operator address

	approved := true

	input, err := parsedABI.Pack("setApprovalForAll", operator, approved)
	if err != nil {
		log.Fatalf("Failed to pack the setApprovalForAll function input: %v", err)
	}

	toAddress := common.HexToAddress(contractAddress)
	msg := ethereum.CallMsg{
		From: common.HexToAddress(operator.Hex()), // Use the operator address as the 'From' address
		To:   &toAddress,
		Data: input,
	}

	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to estimate gas for setApprovalForAll: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	totalGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	fmt.Printf("Total Gas Cost for setApprovalForAll (in wei): %s\n", totalGasCost.String())

	return totalGasCost
}

func estimateGasForBatchsafeTransferFrom(client *ethclient.Client, parsedABI abi.ABI, addresses []common.Address, privateKey string,gas int) *big.Int {
	// from := deriveAddressFromPrivateKey(privateKey) // Owner address
	// operator := common.HexToAddress("0x6F4Cd3C8D1cF68242a66479f2F9A3F976A5cd4dD") // Replace with the operator address
	// id := big.NewInt(1)                                                           // Token ID
	// data := []byte{}                                                              // Data, if any

	// input, err := parsedABI.Pack("BatchsafeTransferFrom", operator, addresses, id, data)
	// if err != nil {
	// 	log.Fatalf("Failed to pack the BatchsafeTransferFrom function input: %v", err)
	// }

	// toAddress := common.HexToAddress(contractAddress)
	// msg := ethereum.CallMsg{
	// 	From: operator,
	// 	To:   &toAddress,
	// 	Data: input,
	// }

	// gasLimit, err := client.EstimateGas(context.Background(), msg)
	gasLimit := gas
	// if err != nil {
	// 	log.Fatalf("Failed to estimate gas for BatchsafeTransferFrom: %v", err)
	// }

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	totalGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	fmt.Printf("Total Gas Cost for BatchsafeTransferFrom (in wei): %s\n", totalGasCost.String())

	return totalGasCost
}

func deriveAddressFromPrivateKey(privateKey string) common.Address {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
