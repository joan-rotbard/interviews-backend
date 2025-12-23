"""
Demo script to demonstrate the issues in the payment service
Run this to see the problems in action
"""

import threading
import time
from src.payment.service import PaymentService
from src.database import db

def demo_race_condition():
    """Demonstrates race condition with concurrent payments"""
    print("\n" + "="*60)
    print("DEMO 1: Race Condition - Concurrent Payments")
    print("="*60)
    
    # Reset
    db.get_connection()["users"]["demo_user"] = {"balance": 1000.0}
    db.get_connection()["transactions"].clear()
    
    service = PaymentService()
    payment_ids = []
    errors = []
    
    def make_payment():
        try:
            pid = service.process_payment(
                "demo_user", 100.0, "paypal",
                {"email": "demo@example.com"}
            )
            payment_ids.append(pid)
        except Exception as e:
            errors.append(str(e))
    
    # Simulate 3 concurrent payments
    threads = []
    for i in range(3):
        t = threading.Thread(target=make_payment)
        threads.append(t)
        t.start()
    
    for t in threads:
        t.join()
    
    print(f"\nPayment IDs created: {payment_ids}")
    print(f"Unique payment IDs: {len(set(payment_ids))}")
    print(f"Expected: 3 unique IDs")
    
    balance = db.get_user_balance("demo_user")
    print(f"\nUser balance: {balance}")
    print(f"Expected: 700.0 (1000 - 3*100)")
    print(f"Actual: {balance}")
    
    if len(set(payment_ids)) < 3:
        print("\n⚠️  ISSUE DETECTED: Duplicate payment IDs (race condition)")
    if balance != 700.0:
        print("⚠️  ISSUE DETECTED: Incorrect balance (race condition)")


def demo_idempotency():
    """Demonstrates lack of idempotency in refunds"""
    print("\n" + "="*60)
    print("DEMO 2: Lack of Idempotency - Multiple Refunds")
    print("="*60)
    
    # Reset
    db.get_connection()["users"]["demo_user2"] = {"balance": 1000.0}
    db.get_connection()["transactions"].clear()
    
    service = PaymentService()
    
    # Process a payment
    payment_id = service.process_payment(
        "demo_user2", 100.0, "paypal",
        {"email": "demo@example.com"}
    )
    print(f"\nPayment created: {payment_id}")
    
    balance_after_payment = db.get_user_balance("demo_user2")
    print(f"Balance after payment: {balance_after_payment}")
    
    # Process refund 3 times (should be idempotent)
    print("\nProcessing refund 3 times...")
    for i in range(3):
        try:
            refund_id = service.process_refund(payment_id, 100.0)
            print(f"  Refund {i+1}: {refund_id}")
        except Exception as e:
            print(f"  Refund {i+1} failed: {e}")
    
    balance_after_refunds = db.get_user_balance("demo_user2")
    print(f"\nBalance after 3 refunds: {balance_after_refunds}")
    print(f"Expected: 1000.0 (should only refund once)")
    
    if balance_after_refunds != 1000.0:
        print("⚠️  ISSUE DETECTED: Multiple refunds processed (lack of idempotency)")


def demo_transaction_issues():
    """Demonstrates transaction/consistency issues"""
    print("\n" + "="*60)
    print("DEMO 3: Transaction Issues - Partial Failures")
    print("="*60)
    
    print("\nNote: This demo shows that if payment processing fails,")
    print("the payment record is already saved in the database.")
    print("This creates orphaned records and inconsistent state.")
    
    # The issue is visible in the code:
    # 1. Payment is saved FIRST
    # 2. Then processing happens
    # 3. If processing fails, payment record exists but wasn't processed
    
    print("\nCheck the code in src/payment/service.py:")
    print("  - process_payment() saves payment BEFORE processing")
    print("  - If processing fails, payment record remains")
    print("  - No transaction rollback mechanism")


if __name__ == "__main__":
    print("\n" + "="*60)
    print("PAYMENT SERVICE - ISSUE DEMONSTRATION")
    print("="*60)
    print("\nThis script demonstrates the issues in the payment service.")
    print("Run the tests to see them in action.\n")
    
    try:
        demo_race_condition()
        demo_idempotency()
        demo_transaction_issues()
        
        print("\n" + "="*60)
        print("DEMO COMPLETE")
        print("="*60)
        print("\nThese are some of the issues you should identify and fix.")
        print("There are more issues in the codebase - explore and find them!")
        
    except Exception as e:
        print(f"\nError during demo: {e}")
        import traceback
        traceback.print_exc()

