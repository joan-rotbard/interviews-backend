"""
Payment processors - External service integrations
"""

import random
import time

def credit_card_processor(card_number, cvv, expiry, amount):
    """Simulate credit card processing"""
    # Simulate network delay
    time.sleep(0.1)
    
    # Simulate random failures
    if random.random() < 0.1:  # 10% failure rate
        raise Exception("Credit card processing failed")
    
    return True


def paypal_processor(email, amount):
    """Simulate PayPal processing"""
    # Simulate network delay
    time.sleep(0.1)
    
    # Simulate random failures
    if random.random() < 0.1:  # 10% failure rate
        raise Exception("PayPal processing failed")
    
    return True

