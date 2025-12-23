/**
 * Payment Service - Main business logic
 */

import { Payment, CreditCardPayment, PayPalPayment, Refund } from '../models/payment';
import * as db from '../database/db';

export class PaymentService {
  /**
   * Process a payment
   */
  processPayment(userId: string, amount: number, paymentType: string, paymentData: any): string {
    // Create payment object based on type
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
    
    // Save payment
    const paymentId = db.savePayment(payment);
    
    // Process payment
    try {
      const result = payment.process();
      if (result) {
        db.updatePaymentStatus(paymentId, "processed");
        db.updateUserBalance(userId, -amount);
      } else {
        db.updatePaymentStatus(paymentId, "failed");
      }
    } catch (e) {
      db.updatePaymentStatus(paymentId, "failed");
      throw e;
    }
    
    return paymentId;
  }

  /**
   * Process a refund
   */
  processRefund(paymentId: string, amount: number): string {
    // Get payment
    const paymentData = db.getPayment(paymentId);
    if (!paymentData) {
      throw new Error("Payment not found");
    }
    
    // Create refund
    const refund = new Refund(`refund_${paymentId}`, paymentId, amount);
    
    // Process refund
    refund.process();
    
    // Update payment status
    db.updatePaymentStatus(paymentId, "refunded");
    
    // Update user balance
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

