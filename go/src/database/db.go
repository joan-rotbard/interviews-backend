package database

import (
	"errors"
	"fmt"
	"payment-microservice/src/models"
)

// Transaction represents a payment transaction in the database
type Transaction struct {
	PaymentID  string
	UserID     string
	Amount     float64
	Currency   string
	Status     string
	CreatedAt  *string
	ProcessedAt *string
}

// User represents a user in the database
type User struct {
	Balance float64
}

// Simulated in-memory database
var (
	transactions = make(map[string]*Transaction)
	users        = make(map[string]*User)
	paymentCounter int
	// NO MUTEX - Race conditions possible
)

// GetConnection returns database connection (simulated)
func GetConnection() (map[string]*Transaction, map[string]*User) {
	return transactions, users
}

// SavePayment saves payment to database - NO TRANSACTION, NO LOCKING
func SavePayment(payment interface{}) (string, error) {
	// Race condition: multiple goroutines can read same counter value
	paymentCounter++
	paymentID := ""
	
	// Type assertion based on payment type - violates Open/Closed Principle
	switch p := payment.(type) {
	case *models.CreditCardPayment:
		paymentID = p.PaymentID
		if paymentID == "" {
			paymentID = fmt.Sprintf("pay_%d", paymentCounter)
			p.PaymentID = paymentID
		}
		transactions[paymentID] = &Transaction{
			PaymentID:  paymentID,
			UserID:     p.UserID,
			Amount:     p.Amount,
			Currency:   p.Currency,
			Status:     p.Status,
			CreatedAt:  p.CreatedAt,
			ProcessedAt: p.ProcessedAt,
		}
	case *models.PayPalPayment:
		paymentID = p.PaymentID
		if paymentID == "" {
			paymentID = fmt.Sprintf("pay_%d", paymentCounter)
			p.PaymentID = paymentID
		}
		transactions[paymentID] = &Transaction{
			PaymentID:  paymentID,
			UserID:     p.UserID,
			Amount:     p.Amount,
			Currency:   p.Currency,
			Status:     p.Status,
			CreatedAt:  p.CreatedAt,
			ProcessedAt: p.ProcessedAt,
		}
	default:
		return "", errors.New("unknown payment type")
	}
	
	// No duplicate check, no transaction isolation
	return paymentID, nil
}

// GetPayment gets payment from database
func GetPayment(paymentID string) (*Transaction, error) {
	txn, exists := transactions[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}
	return txn, nil
}

// UpdatePaymentStatus updates payment status - NO TRANSACTION
func UpdatePaymentStatus(paymentID string, status string) error {
	if txn, exists := transactions[paymentID]; exists {
		txn.Status = status
		return nil
	}
	return errors.New("payment not found")
}

// GetUserBalance gets user balance - NO LOCKING, RACE CONDITION
func GetUserBalance(userID string) float64 {
	user, exists := users[userID]
	if !exists {
		return 0.0
	}
	return user.Balance
}

// UpdateUserBalance updates user balance - NO TRANSACTION, RACE CONDITION
func UpdateUserBalance(userID string, amount float64) (float64, error) {
	if _, exists := users[userID]; !exists {
		users[userID] = &User{Balance: 0.0}
	}
	
	// Race condition: read-modify-write without locking
	currentBalance := users[userID].Balance
	users[userID].Balance = currentBalance + amount
	
	return users[userID].Balance, nil
}

// GetUserTransactions gets all transactions for a user
func GetUserTransactions(userID string) []*Transaction {
	var userTransactions []*Transaction
	for _, txn := range transactions {
		if txn.UserID == userID {
			userTransactions = append(userTransactions, txn)
		}
	}
	return userTransactions
}

