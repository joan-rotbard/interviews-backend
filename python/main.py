"""
Main entry point - Simple CLI to test the payment service
"""

from src.payment.service import PaymentService
from src.database import db

def main():
    print("Payment Microservice - Test Interface")
    print("=" * 50)
    
    service = PaymentService()
    
    # Initialize some test data
    db.get_connection()["users"]["user_123"] = {"balance": 1000.0}
    
    while True:
        print("\nOptions:")
        print("1. Process Credit Card Payment")
        print("2. Process PayPal Payment")
        print("3. Process Refund")
        print("4. Get User Transactions")
        print("5. Get Payment Status")
        print("6. Exit")
        
        choice = input("\nSelect option: ").strip()
        
        if choice == "1":
            try:
                user_id = input("User ID: ").strip()
                amount = float(input("Amount: ").strip())
                card_number = input("Card Number: ").strip()
                cvv = input("CVV: ").strip()
                expiry = input("Expiry (MM/YY): ").strip()
                
                payment_id = service.process_payment(
                    user_id, amount, "credit_card",
                    {"card_number": card_number, "cvv": cvv, "expiry": expiry}
                )
                print(f"Payment processed: {payment_id}")
            except Exception as e:
                print(f"Error: {e}")
        
        elif choice == "2":
            try:
                user_id = input("User ID: ").strip()
                amount = float(input("Amount: ").strip())
                email = input("PayPal Email: ").strip()
                
                payment_id = service.process_payment(
                    user_id, amount, "paypal",
                    {"email": email}
                )
                print(f"Payment processed: {payment_id}")
            except Exception as e:
                print(f"Error: {e}")
        
        elif choice == "3":
            try:
                payment_id = input("Payment ID: ").strip()
                amount = float(input("Refund Amount: ").strip())
                
                refund_id = service.process_refund(payment_id, amount)
                print(f"Refund processed: {refund_id}")
            except Exception as e:
                print(f"Error: {e}")
        
        elif choice == "4":
            try:
                user_id = input("User ID: ").strip()
                transactions = service.get_user_transactions(user_id)
                print(f"\nTransactions for {user_id}:")
                for txn in transactions:
                    print(f"  {txn['payment_id']}: {txn['amount']} {txn['currency']} - {txn['status']}")
            except Exception as e:
                print(f"Error: {e}")
        
        elif choice == "5":
            try:
                payment_id = input("Payment ID: ").strip()
                status = service.get_payment_status(payment_id)
                if status:
                    print(f"Payment status: {status}")
                else:
                    print("Payment not found")
            except Exception as e:
                print(f"Error: {e}")
        
        elif choice == "6":
            break
        
        else:
            print("Invalid option")

if __name__ == "__main__":
    main()

