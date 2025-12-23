package tests

import (
	"fmt"
	"payment-microservice/src/database"
	"payment-microservice/src/payment"
	"sync"
	"testing"
)

// TestConcurrentPayments demonstrates race conditions
func TestConcurrentPayments(t *testing.T) {
	// Reset database
	transactions, users := database.GetConnection()
	users["user_test"] = &database.User{Balance: 1000.0}
	_ = transactions // suppress unused variable
	
	service := &payment.PaymentService()
	
	// Simulate concurrent payments
	var wg sync.WaitGroup
	results := make([]string, 0)
	var mu sync.Mutex
	
	processPayment := func() {
		defer wg.Done()
		paymentID, err := service.ProcessPayment(
			"user_test", 100.0, "paypal",
			map[string]string{"email": "test@example.com"},
		)
		mu.Lock()
		if err != nil {
			results = append(results, fmt.Sprintf("Error: %v", err))
		} else {
			results = append(results, paymentID)
		}
		mu.Unlock()
	}
	
	// Create multiple goroutines to process payments simultaneously
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go processPayment()
	}
	
	wg.Wait()
	
	// Check for duplicate payment IDs (race condition)
	fmt.Printf("Results: %v\n", results)
	uniqueResults := make(map[string]bool)
	for _, r := range results {
		uniqueResults[r] = true
	}
	fmt.Printf("Unique payments: %d\n", len(uniqueResults))
	
	// Check user balance
	balance := database.GetUserBalance("user_test")
	fmt.Printf("Final balance: %.2f\n", balance)
	fmt.Printf("Expected balance: 500.0 (1000 - 5*100)\n")
	fmt.Printf("Actual balance: %.2f\n", balance)
}

// TestDuplicateRefund demonstrates lack of idempotency
func TestDuplicateRefund(t *testing.T) {
	// Reset database
	transactions, users := database.GetConnection()
	users["user_test2"] = &database.User{Balance: 1000.0}
	_ = transactions // suppress unused variable
	
	service := &payment.PaymentService()
	
	// Process a payment
	paymentID, err := service.ProcessPayment(
		"user_test2", 100.0, "paypal",
		map[string]string{"email": "test@example.com"},
	)
	if err != nil {
		t.Fatalf("Failed to process payment: %v", err)
	}
	
	initialBalance := database.GetUserBalance("user_test2")
	fmt.Printf("Balance after payment: %.2f\n", initialBalance)
	
	// Process refund multiple times (should be idempotent but isn't)
	service.ProcessRefund(paymentID, 100.0)
	service.ProcessRefund(paymentID, 100.0)
	service.ProcessRefund(paymentID, 100.0)
	
	finalBalance := database.GetUserBalance("user_test2")
	fmt.Printf("Balance after 3 refunds: %.2f\n", finalBalance)
	fmt.Printf("Expected balance: 1000.0 (should only refund once)\n")
	fmt.Printf("Actual balance: %.2f\n", finalBalance)
}

