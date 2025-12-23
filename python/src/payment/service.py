"""
Payment Service - Main business logic with multiple issues
"""

from src.models.payment import Payment, CreditCardPayment, PayPalPayment, Refund
from src.database import db


class PaymentService:
    """Main payment service - Contains multiple design and concurrency issues"""
    
    def process_payment(self, user_id, amount, payment_type, payment_data):
        """
        Process a payment
        
        Issues:
        - No idempotency check
        - No transaction handling
        - Race conditions possible
        - Tight coupling to concrete classes
        - No error recovery
        """
        
        # Create payment object based on type - violates Open/Closed Principle
        if payment_type == "credit_card":
            payment = CreditCardPayment(
                None, user_id, amount,
                payment_data["card_number"],
                payment_data["cvv"],
                payment_data["expiry"]
            )
        elif payment_type == "paypal":
            payment = PayPalPayment(
                None, user_id, amount,
                payment_data["email"]
            )
        else:
            raise ValueError(f"Unknown payment type: {payment_type}")
        
        # Save payment BEFORE processing - if processing fails, we have orphaned record
        payment_id = db.save_payment(payment)
        
        # Process payment - if this fails, payment is saved but not processed
        try:
            result = payment.process()
            if result:
                db.update_payment_status(payment_id, "processed")
                # Update user balance - NO TRANSACTION with payment save
                db.update_user_balance(user_id, -amount)
            else:
                db.update_payment_status(payment_id, "failed")
        except Exception as e:
            # Error handling but payment already saved
            db.update_payment_status(payment_id, "failed")
            raise e
        
        return payment_id
    
    def process_refund(self, payment_id, amount):
        """
        Process a refund
        
        Issues:
        - No idempotency check (can refund multiple times)
        - No validation that payment exists or was processed
        - No transaction handling
        - Race conditions possible
        """
        
        # Get payment - but don't validate status
        payment_data = db.get_payment(payment_id)
        if not payment_data:
            raise ValueError("Payment not found")
        
        # Create refund - no duplicate check
        refund = Refund(f"refund_{payment_id}", payment_id, amount)
        
        # Process refund - no validation
        refund.process()
        
        # Update payment status - NO TRANSACTION
        db.update_payment_status(payment_id, "refunded")
        
        # Update user balance - NO TRANSACTION with payment update
        db.update_user_balance(payment_data["user_id"], amount)
        
        return refund.refund_id
    
    def get_user_transactions(self, user_id):
        """Get all transactions for a user"""
        return db.get_user_transactions(user_id)
    
    def get_payment_status(self, payment_id):
        """Get payment status"""
        payment = db.get_payment(payment_id)
        if not payment:
            return None
        return payment["status"]

