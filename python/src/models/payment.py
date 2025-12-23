"""
Payment models
"""

class Payment:
    """Represents a payment transaction"""
    
    def __init__(self, payment_id, user_id, amount, currency="USD"):
        self.payment_id = payment_id
        self.user_id = user_id
        self.amount = amount
        self.currency = currency
        self.status = "pending"
        self.created_at = None
        self.processed_at = None
    
    def process(self):
        """Process the payment"""
        self.status = "processed"
        self.processed_at = "now"
        return True
    
    def refund(self):
        """Refund the payment"""
        self.status = "refunded"
        return True


class CreditCardPayment(Payment):
    """Credit card payment"""
    
    def __init__(self, payment_id, user_id, amount, card_number, cvv, expiry):
        super().__init__(payment_id, user_id, amount)
        self.card_number = card_number
        self.cvv = cvv
        self.expiry = expiry
    
    def process(self):
        """Process credit card payment"""
        from src.payment.processors import credit_card_processor
        result = credit_card_processor(self.card_number, self.cvv, self.expiry, self.amount)
        if result:
            self.status = "processed"
        return result


class PayPalPayment(Payment):
    """PayPal payment"""
    
    def __init__(self, payment_id, user_id, amount, paypal_email):
        super().__init__(payment_id, user_id, amount)
        self.paypal_email = paypal_email
    
    def process(self):
        """Process PayPal payment"""
        from src.payment.processors import paypal_processor
        result = paypal_processor(self.paypal_email, self.amount)
        if result:
            self.status = "processed"
        return result


class Refund:
    """Represents a refund transaction"""
    
    def __init__(self, refund_id, payment_id, amount):
        self.refund_id = refund_id
        self.payment_id = payment_id
        self.amount = amount
        self.status = "pending"
    
    def process(self):
        """Process the refund"""
        self.status = "processed"
        return True

