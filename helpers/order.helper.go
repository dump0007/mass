package helpers

import (
	"context"
	"fmt"
	"log"
	"mass/database"
	"mass/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// executeOrderProcessing processes orders based on their status
func ExecuteOrderProcessing() {
	var order models.Order
	orderCollection := database.OpenCollection(database.Client, "orders")

	// Step 1: Check for any pending orders
	err := orderCollection.FindOne(context.Background(), bson.M{"order_status": "stopped"}).Decode(&order)
	if err == nil {
		// Pending order found, pass the OrderID to another function
		handleOrder(order.OrderID)
		return
	}

	// Step 2: Check if there are no pending orders and find the oldest active order
	pendingCount, err := orderCollection.CountDocuments(context.Background(), bson.M{"order_status": "pending"})
	if err != nil {
		log.Fatalf("Failed to count pending orders: %v", err)
	}

	if pendingCount == 0 {
		// No pending orders, find the oldest active order
		opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: 1}})
		err = orderCollection.FindOne(context.Background(), bson.M{"order_status": "active"}, opts).Decode(&order)
		if err == nil {
			// Oldest active order found, pass the OrderID to another function
			handleOrder(order.OrderID)
			return
		}
	}

	// Step 3: No pending or active orders, exit the function
	fmt.Println("No pending or active orders found. Exiting function.")
}

// handleOrder is the function that processes the order with the given OrderID
func handleOrder(orderID float64) {
	fmt.Printf("Processing order with OrderID: %.0f\n", orderID)
	orderCollection := database.OpenCollection(database.Client, "orders")
	filter := bson.D{{"order_id", orderID}}
	update := bson.D{{"$set", bson.D{{"order_status", "pending"}}}}
	result, err := orderCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("mongo query dropped", err)
	}
	fmt.Println(result)
	// Add your order handling logic here
	ProcessExcelAndExecuteTransactions(orderID)
}
