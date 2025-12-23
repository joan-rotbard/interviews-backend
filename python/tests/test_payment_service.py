"""
Test file for payment service
Run with: pytest tests/test_payment_service.py
"""

import threading
from src.payment.service import PaymentService
from src.database import db

def test_concurrent_payments():
    """Test concurrent payment processing"""
    # Reset database
    db.get_connection()["users"]["user_test"] = {"balance": 1000.0}
    
    service = PaymentService()
    
    # Simulate concurrent payments
    results = []
    
    def process_payment():
        try:
            payment_id = service.process_payment(
                "user_test", 100.0, "paypal",
                {"email": "test@example.com"}
            )
            results.append(payment_id)
        except Exception as e:
            results.append(f"Error: {e}")
    
    # Create multiple threads to process payments simultaneously
    threads = []
    for i in range(5):
        thread = threading.Thread(target=process_payment)
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()
    
    # Check for duplicate payment IDs
    print(f"Results: {results}")
    print(f"Unique payments: {len(set(results))}")
    
    # Check user balance
    balance = db.get_user_balance("user_test")
    print(f"Final balance: {balance}")
    print(f"Expected balance: 500.0 (1000 - 5*100)")
    print(f"Actual balance: {balance}")


def test_duplicate_refund():
    """Test refund processing"""
    # Reset database
    db.get_connection()["users"]["user_test2"] = {"balance": 1000.0}
    
    service = PaymentService()
    
    # Process a payment
    payment_id = service.process_payment(
        "user_test2", 100.0, "paypal",
        {"email": "test@example.com"}
    )
    
    initial_balance = db.get_user_balance("user_test2")
    print(f"Balance after payment: {initial_balance}")
    
    # Process refund multiple times
    service.process_refund(payment_id, 100.0)
    service.process_refund(payment_id, 100.0)
    service.process_refund(payment_id, 100.0)
    
    final_balance = db.get_user_balance("user_test2")
    print(f"Balance after 3 refunds: {final_balance}")
    print(f"Expected balance: 1000.0")
    print(f"Actual balance: {final_balance}")


if __name__ == "__main__":
    print("=" * 50)
    print("Test 1: Concurrent Payments")
    print("=" * 50)
    test_concurrent_payments()
    
    print("\n" + "=" * 50)
    print("Test 2: Multiple Refunds")
    print("=" * 50)
    test_duplicate_refund()

