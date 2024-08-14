package helpers

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"mass/database"
	"mass/models"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/xuri/excelize/v2"
)

func ProcessExcelAndExecuteTransactions(orderID float64) error {
	orderCollection := database.OpenCollection(database.Client, "orders")
	var lastOrder models.Order
	err := orderCollection.FindOne(context.Background(), bson.M{"order_id": orderID}).Decode(&lastOrder)

	if err != nil {
		return err
	}
	

	// fileName := "modified_Aaddresses.xlsx"
	fileName := lastOrder.FileName
	fmt.Println("start")
	// Open the Excel file
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	fmt.Println("here1")

	// Connect to the Ethereum client
	client, err := ethclient.Dial("https://ethereum-holesky-rpc.publicnode.com")
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()
	fmt.Println("here2")

	// Iterate over the rows to find active transactions
	rows, err := f.GetRows("Sheet1")
	fmt.Println("here3")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	fmt.Println("1111111111111111111111111111111111111111111")

	// Load your private key (replace with your actual private key)
	privateKey, err := crypto.HexToECDSA(lastOrder.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}
	fmt.Println("222222222222222222222222222222222222222222222")

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	fmt.Println("333333333333333333333333333333333333333333333")
	fmt.Println("----------Entering The Matrix----------------")
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	fmt.Println("nonce", nonce)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Iterate over the rows starting from row 2 (skip header)
	for i := int(nonce+1); i < len(rows); i++ {
		fmt.Println("i->", i, "rows[i][0]", rows[i][0])
		walletAddress := rows[i][0]
		fmt.Println("First Record", rows[0][0])
		// Create the transaction
		fmt.Println("Current nonce", nonce)

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get gas price: %w", err)
		}
		fmt.Println("gasPrice", gasPrice)

		toAddress := common.HexToAddress(walletAddress)
		fmt.Println("toAddress", toAddress)

		message := "hello world"
		data := []byte(message)
		fmt.Println("data", data)
		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: data,
		}

		gasLimit, err := client.EstimateGas(context.Background(), msg)
		if err != nil {
			log.Fatalf("Failed to estimate gas: %v", err)
		}

		fmt.Println("Gas Estimate:", gasLimit)
		fmt.Println("Gas Price:", gasPrice)

		// tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), uint64(21000), gasPrice, []byte(data))
		tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gasLimit, gasPrice, data)

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get network ID: %w", err)
		}
		fmt.Println("chainId", chainID)

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			return fmt.Errorf("failed to sign transaction: %w", err)
		}
		fmt.Println("111111111111111111111")

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return fmt.Errorf("failed to send transaction: %w", err)
		}

		fmt.Println("2222222222222222222222222222")
		// Wait for transaction confirmation
		// receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		// if err != nil {
		// 	return fmt.Errorf("failed to wait for transaction confirmation: %w", err)
		// }

		// // Ensure the transaction was successful
		// if receipt.Status != 1 {
		// 	return fmt.Errorf("transaction failed: %s", signedTx.Hash().Hex())
		// }

		// Update txn_hash and txn_status in the Excel file
		txnHash := signedTx.Hash().Hex()
		// f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+1), txnHash)    // Update txn_hash
		// f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+1), "executed") // Update txn_status to "executed"

		// // Save the updated Excel file after each transaction
		// if err := f.SaveAs(fileName); err != nil {
		// 	return fmt.Errorf("failed to save Excel file: %w", err)
		// }

		fmt.Printf("Transaction sent and confirmed: %s\n", txnHash)
		nonce = nonce + 1
		fmt.Println("New nonce", nonce)

	}

	return nil
}
