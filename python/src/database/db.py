"""
Database layer
"""

# Simulated in-memory database
_transactions = {}
_users = {}
_payment_counter = 0

def get_connection():
    """Get database connection - simulated"""
    return {"transactions": _transactions, "users": _users}


def save_payment(payment):
    """Save payment to database"""
    global _payment_counter
    _payment_counter += 1
    
    payment_id = f"pay_{_payment_counter}"
    payment.payment_id = payment_id
    _transactions[payment_id] = {
        "payment_id": payment_id,
        "user_id": payment.user_id,
        "amount": payment.amount,
        "currency": payment.currency,
        "status": payment.status,
        "created_at": payment.created_at,
        "processed_at": payment.processed_at
    }
    
    return payment_id


def get_payment(payment_id):
    """Get payment from database"""
    return _transactions.get(payment_id)


def update_payment_status(payment_id, status):
    """Update payment status"""
    if payment_id in _transactions:
        _transactions[payment_id]["status"] = status
        return True
    return False


def get_user_balance(user_id):
    """Get user balance"""
    user = _users.get(user_id, {"balance": 0.0})
    return user["balance"]


def update_user_balance(user_id, amount):
    """Update user balance"""
    if user_id not in _users:
        _users[user_id] = {"balance": 0.0}
    
    current_balance = _users[user_id]["balance"]
    _users[user_id]["balance"] = current_balance + amount
    
    return _users[user_id]["balance"]


def get_user_transactions(user_id):
    """Get all transactions for a user"""
    user_transactions = []
    for payment_id, transaction in _transactions.items():
        if transaction["user_id"] == user_id:
            user_transactions.append(transaction)
    return user_transactions

