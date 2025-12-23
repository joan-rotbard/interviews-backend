package main

import (
	"bufio"
	"fmt"
	"os"
	"payment-microservice/src/database"
	"payment-microservice/src/payment"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Payment Microservice - Test Interface")
	fmt.Println(strings.Repeat("=", 50))
	
	service := &payment.PaymentService{}
	
	// Initialize some test data
	transactions, users := database.GetConnection()
	users["user_123"] = &database.User{Balance: 1000.0}
	_ = transactions // suppress unused variable
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Process Credit Card Payment")
		fmt.Println("2. Process PayPal Payment")
		fmt.Println("3. Process Refund")
		fmt.Println("4. Get User Transactions")
		fmt.Println("5. Get Payment Status")
		fmt.Println("6. Exit")
		
		fmt.Print("\nSelect option: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())
		
		if choice == "1" {
			fmt.Print("User ID: ")
			scanner.Scan()
			userID := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Amount: ")
			scanner.Scan()
			amountStr := strings.TrimSpace(scanner.Text())
			amount, _ := strconv.ParseFloat(amountStr, 64)
			
			fmt.Print("Card Number: ")
			scanner.Scan()
			cardNumber := strings.TrimSpace(scanner.Text())
			
			fmt.Print("CVV: ")
			scanner.Scan()
			cvv := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Expiry (MM/YY): ")
			scanner.Scan()
			expiry := strings.TrimSpace(scanner.Text())
			
			paymentID, err := service.ProcessPayment(userID, amount, "credit_card", map[string]string{
				"card_number": cardNumber,
				"cvv":         cvv,
				"expiry":      expiry,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Payment processed: %s\n", paymentID)
			}
		} else if choice == "2" {
			fmt.Print("User ID: ")
			scanner.Scan()
			userID := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Amount: ")
			scanner.Scan()
			amountStr := strings.TrimSpace(scanner.Text())
			amount, _ := strconv.ParseFloat(amountStr, 64)
			
			fmt.Print("PayPal Email: ")
			scanner.Scan()
			email := strings.TrimSpace(scanner.Text())
			
			paymentID, err := service.ProcessPayment(userID, amount, "paypal", map[string]string{
				"email": email,
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Payment processed: %s\n", paymentID)
			}
		} else if choice == "3" {
			fmt.Print("Payment ID: ")
			scanner.Scan()
			paymentID := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Refund Amount: ")
			scanner.Scan()
			amountStr := strings.TrimSpace(scanner.Text())
			amount, _ := strconv.ParseFloat(amountStr, 64)
			
			refundID, err := service.ProcessRefund(paymentID, amount)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Refund processed: %s\n", refundID)
			}
		} else if choice == "4" {
			fmt.Print("User ID: ")
			scanner.Scan()
			userID := strings.TrimSpace(scanner.Text())
			
			transactions := service.GetUserTransactions(userID)
			fmt.Printf("\nTransactions for %s:\n", userID)
			for _, txn := range transactions {
				fmt.Printf("  %s: %.2f %s - %s\n", txn.PaymentID, txn.Amount, txn.Currency, txn.Status)
			}
		} else if choice == "5" {
			fmt.Print("Payment ID: ")
			scanner.Scan()
			paymentID := strings.TrimSpace(scanner.Text())
			
			status, err := service.GetPaymentStatus(paymentID)
			if err != nil {
				fmt.Println("Payment not found")
			} else {
				fmt.Printf("Payment status: %s\n", status)
			}
		} else if choice == "6" {
			break
		} else {
			fmt.Println("Invalid option")
		}
	}
}

