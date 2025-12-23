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
func (s *PaymentService) ProcessPayment(userID string, amount float64, paymentType string, paymentData map[string]string) (string, error) {
	// Create payment object based on type
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
		return "", fmt.Errorf("unknown payment type: %s", paymentType)
	}
	
	// Save payment
	paymentID, err := database.SavePayment(payment)
	if err != nil {
		return "", err
	}
	
	// Process payment
	var result bool
	var processErr error
	
	switch p := payment.(type) {
	case *models.CreditCardPayment:
		result, processErr = p.Process()
	case *models.PayPalPayment:
		result, processErr = p.Process()
	}
	
	if processErr != nil {
		database.UpdatePaymentStatus(paymentID, "failed")
		return "", processErr
	}
	
	if result {
		database.UpdatePaymentStatus(paymentID, "processed")
		database.UpdateUserBalance(userID, -amount)
	} else {
		database.UpdatePaymentStatus(paymentID, "failed")
	}
	
	return paymentID, nil
}

// ProcessRefund processes a refund
func (s *PaymentService) ProcessRefund(paymentID string, amount float64) (string, error) {
	// Get payment
	paymentData, err := database.GetPayment(paymentID)
	if err != nil {
		return "", errors.New("payment not found")
	}
	
	// Create refund
	refund := models.NewRefund(fmt.Sprintf("refund_%s", paymentID), paymentID, amount)
	
	// Process refund
	refund.Process()
	
	// Update payment status
	database.UpdatePaymentStatus(paymentID, "refunded")
	
	// Update user balance
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

