/**
 * Database layer
 */

interface Transaction {
  paymentId: string;
  userId: string;
  amount: number;
  currency: string;
  status: string;
  createdAt: string | null;
  processedAt: string | null;
}

interface User {
  balance: number;
}

// Simulated in-memory database
const transactions: Record<string, Transaction> = {};
const users: Record<string, User> = {};
let paymentCounter = 0;

export function getConnection(): { transactions: Record<string, Transaction>; users: Record<string, User> } {
  return { transactions, users };
}

export function savePayment(payment: any): string {
  paymentCounter++;
  
  const paymentId = `pay_${paymentCounter}`;
  payment.paymentId = paymentId;
  transactions[paymentId] = {
    paymentId: paymentId,
    userId: payment.userId,
    amount: payment.amount,
    currency: payment.currency,
    status: payment.status,
    createdAt: payment.createdAt,
    processedAt: payment.processedAt
  };
  
  return paymentId;
}

export function getPayment(paymentId: string): Transaction | undefined {
  return transactions[paymentId];
}

export function updatePaymentStatus(paymentId: string, status: string): boolean {
  if (paymentId in transactions) {
    transactions[paymentId].status = status;
    return true;
  }
  return false;
}

export function getUserBalance(userId: string): number {
  const user = users[userId] || { balance: 0.0 };
  return user.balance;
}

export function updateUserBalance(userId: string, amount: number): number {
  if (!users[userId]) {
    users[userId] = { balance: 0.0 };
  }
  
  const currentBalance = users[userId].balance;
  users[userId].balance = currentBalance + amount;
  
  return users[userId].balance;
}

export function getUserTransactions(userId: string): Transaction[] {
  const userTransactions: Transaction[] = [];
  for (const paymentId in transactions) {
    if (transactions[paymentId].userId === userId) {
      userTransactions.push(transactions[paymentId]);
    }
  }
  return userTransactions;
}

