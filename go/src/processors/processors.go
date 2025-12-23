package processors

import (
	"errors"
	"math/rand"
	"time"
)

// CreditCardProcessor simulates credit card processing
func CreditCardProcessor(cardNumber string, cvv string, expiry string, amount float64) (bool, error) {
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)
	
	// Simulate random failures
	if rand.Float64() < 0.1 { // 10% failure rate
		return false, errors.New("credit card processing failed")
	}
	
	return true, nil
}

// PayPalProcessor simulates PayPal processing
func PayPalProcessor(email string, amount float64) (bool, error) {
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)
	
	// Simulate random failures
	if rand.Float64() < 0.1 { // 10% failure rate
		return false, errors.New("PayPal processing failed")
	}
	
	return true, nil
}

