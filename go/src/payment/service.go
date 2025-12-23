package payment

import (
	"errors"
	"fmt"
	"payment-microservice/src/database"
	"payment-microservice/src/models"
)

// PaymentService handles payment processing
type PaymentService struct{}

// ProcessPayment processes a payment
// Issues:
// - No idempotency check
// - No transaction handling
// - Race conditions possible
// - Tight coupling to concrete classes
// - No error recovery
func (s *PaymentService) ProcessPayment(userID string, amount float64, paymentType string, paymentData map[string]string) (string, error) {
	// Create payment object based on type - violates Open/Closed Principle
	var payment interface{}
	
	if paymentType == "credit_card" {
		payment = models.NewCreditCardPayment(
			"", userID, amount,
			paymentData["card_number"],
			paymentData["cvv"],
			paymentData["expiry"],
		)
	} else if paymentType == "paypal" {
		payment = models.NewPayPalPayment(
			"", userID, amount,
			paymentData["email"],
		)
	} else {
		return "", errors.New(fmt.Sprintf("unknown payment type: %s", paymentType))
	}
	
	// Save payment BEFORE processing - if processing fails, we have orphaned record
	paymentID, err := database.SavePayment(payment)
	if err != nil {
		return "", err
	}
	
	// Process payment - if this fails, payment is saved but not processed
	var result bool
	var processErr error
	
	switch p := payment.(type) {
	case *models.CreditCardPayment:
		result, processErr = p.Process()
	case *models.PayPalPayment:
		result, processErr = p.Process()
	}
	
	if processErr != nil {
		// Error handling but payment already saved
		database.UpdatePaymentStatus(paymentID, "failed")
		return "", processErr
	}
	
	if result {
		database.UpdatePaymentStatus(paymentID, "processed")
		// Update user balance - NO TRANSACTION with payment save
		database.UpdateUserBalance(userID, -amount)
	} else {
		database.UpdatePaymentStatus(paymentID, "failed")
	}
	
	return paymentID, nil
}

// ProcessRefund processes a refund
// Issues:
// - No idempotency check (can refund multiple times)
// - No validation that payment exists or was processed
// - No transaction handling
// - Race conditions possible
func (s *PaymentService) ProcessRefund(paymentID string, amount float64) (string, error) {
	// Get payment - but don't validate status
	paymentData, err := database.GetPayment(paymentID)
	if err != nil {
		return "", errors.New("payment not found")
	}
	
	// Create refund - no duplicate check
	refund := models.NewRefund(fmt.Sprintf("refund_%s", paymentID), paymentID, amount)
	
	// Process refund - no validation
	refund.Process()
	
	// Update payment status - NO TRANSACTION
	database.UpdatePaymentStatus(paymentID, "refunded")
	
	// Update user balance - NO TRANSACTION with payment update
	database.UpdateUserBalance(paymentData.UserID, amount)
	
	return refund.RefundID, nil
}

// GetUserTransactions gets all transactions for a user
func (s *PaymentService) GetUserTransactions(userID string) []*database.Transaction {
	return database.GetUserTransactions(userID)
}

// GetPaymentStatus gets payment status
func (s *PaymentService) GetPaymentStatus(paymentID string) (string, error) {
	payment, err := database.GetPayment(paymentID)
	if err != nil {
		return "", err
	}
	return payment.Status, nil
}

