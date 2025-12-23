package models

import (
	"payment-microservice/src/processors"
)

// Payment represents a payment transaction
type Payment struct {
	PaymentID  string
	UserID     string
	Amount     float64
	Currency   string
	Status     string
	CreatedAt  *string
	ProcessedAt *string
}

// NewPayment creates a new payment
func NewPayment(paymentID string, userID string, amount float64, currency string) *Payment {
	if currency == "" {
		currency = "USD"
	}
	return &Payment{
		PaymentID:  paymentID,
		UserID:     userID,
		Amount:     amount,
		Currency:   currency,
		Status:     "pending",
		CreatedAt:  nil,
		ProcessedAt: nil,
	}
}

// Process processes the payment
func (p *Payment) Process() (bool, error) {
	p.Status = "processed"
	now := "now"
	p.ProcessedAt = &now
	return true, nil
}

// Refund refunds the payment
func (p *Payment) Refund() (bool, error) {
	p.Status = "refunded"
	return true, nil
}

// CreditCardPayment represents a credit card payment
type CreditCardPayment struct {
	*Payment
	CardNumber string
	CVV        string
	Expiry     string
}

// NewCreditCardPayment creates a new credit card payment
func NewCreditCardPayment(paymentID string, userID string, amount float64, cardNumber string, cvv string, expiry string) *CreditCardPayment {
	return &CreditCardPayment{
		Payment:    NewPayment(paymentID, userID, amount, "USD"),
		CardNumber: cardNumber,
		CVV:        cvv,
		Expiry:     expiry,
	}
}

// Process processes credit card payment
func (p *CreditCardPayment) Process() (bool, error) {
	// Directly accessing external service - tight coupling
	result, err := processors.CreditCardProcessor(p.CardNumber, p.CVV, p.Expiry, p.Amount)
	if err != nil {
		return false, err
	}
	if result {
		p.Status = "processed"
	}
	return result, nil
}

// PayPalPayment represents a PayPal payment
type PayPalPayment struct {
	*Payment
	PayPalEmail string
}

// NewPayPalPayment creates a new PayPal payment
func NewPayPalPayment(paymentID string, userID string, amount float64, paypalEmail string) *PayPalPayment {
	return &PayPalPayment{
		Payment:    NewPayment(paymentID, userID, amount, "USD"),
		PayPalEmail: paypalEmail,
	}
}

// Process processes PayPal payment
func (p *PayPalPayment) Process() (bool, error) {
	// Directly accessing external service - tight coupling
	result, err := processors.PayPalProcessor(p.PayPalEmail, p.Amount)
	if err != nil {
		return false, err
	}
	if result {
		p.Status = "processed"
	}
	return result, nil
}

// Refund represents a refund transaction
type Refund struct {
	RefundID  string
	PaymentID string
	Amount    float64
	Status    string
}

// NewRefund creates a new refund
func NewRefund(refundID string, paymentID string, amount float64) *Refund {
	return &Refund{
		RefundID:  refundID,
		PaymentID: paymentID,
		Amount:    amount,
		Status:    "pending",
	}
}

// Process processes the refund
func (r *Refund) Process() (bool, error) {
	// No validation, no idempotency check
	r.Status = "processed"
	return true, nil
}

// PaymentProcessor interface - NOT USED (violates Open/Closed Principle)
type PaymentProcessor interface {
	Process() (bool, error)
}

