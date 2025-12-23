/**
 * Payment Service - Main business logic with multiple issues
 */

import { Payment, CreditCardPayment, PayPalPayment, Refund } from '../models/payment';
import * as db from '../database/db';

export class PaymentService {
  /**
   * Process a payment
   * 
   * Issues:
   * - No idempotency check
   * - No transaction handling
   * - Race conditions possible
   * - Tight coupling to concrete classes
   * - No error recovery
   */
  processPayment(userId: string, amount: number, paymentType: string, paymentData: any): string {
    // Create payment object based on type - violates Open/Closed Principle
    let payment: Payment;
    
    if (paymentType === "credit_card") {
      payment = new CreditCardPayment(
        null, userId, amount,
        paymentData.card_number,
        paymentData.cvv,
        paymentData.expiry
      );
    } else if (paymentType === "paypal") {
      payment = new PayPalPayment(
        null, userId, amount,
        paymentData.email
      );
    } else {
      throw new Error(`Unknown payment type: ${paymentType}`);
    }
    
    // Save payment BEFORE processing - if processing fails, we have orphaned record
    const paymentId = db.savePayment(payment);
    
    // Process payment - if this fails, payment is saved but not processed
    try {
      const result = payment.process();
      if (result) {
        db.updatePaymentStatus(paymentId, "processed");
        // Update user balance - NO TRANSACTION with payment save
        db.updateUserBalance(userId, -amount);
      } else {
        db.updatePaymentStatus(paymentId, "failed");
      }
    } catch (e) {
      // Error handling but payment already saved
      db.updatePaymentStatus(paymentId, "failed");
      throw e;
    }
    
    return paymentId;
  }

  /**
   * Process a refund
   * 
   * Issues:
   * - No idempotency check (can refund multiple times)
   * - No validation that payment exists or was processed
   * - No transaction handling
   * - Race conditions possible
   */
  processRefund(paymentId: string, amount: number): string {
    // Get payment - but don't validate status
    const paymentData = db.getPayment(paymentId);
    if (!paymentData) {
      throw new Error("Payment not found");
    }
    
    // Create refund - no duplicate check
    const refund = new Refund(`refund_${paymentId}`, paymentId, amount);
    
    // Process refund - no validation
    refund.process();
    
    // Update payment status - NO TRANSACTION
    db.updatePaymentStatus(paymentId, "refunded");
    
    // Update user balance - NO TRANSACTION with payment update
    db.updateUserBalance(paymentData.userId, amount);
    
    return refund.refundId;
  }

  getUserTransactions(userId: string): any[] {
    return db.getUserTransactions(userId);
  }

  getPaymentStatus(paymentId: string): string | null {
    const payment = db.getPayment(paymentId);
    if (!payment) {
      return null;
    }
    return payment.status;
  }
}

