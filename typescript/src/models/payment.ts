/**
 * Payment models
 */

import { creditCardProcessor, paypalProcessor } from '../payment/processors';

export class Payment {
  paymentId: string | null;
  userId: string;
  amount: number;
  currency: string;
  status: string;
  createdAt: string | null;
  processedAt: string | null;

  constructor(paymentId: string | null, userId: string, amount: number, currency: string = "USD") {
    this.paymentId = paymentId;
    this.userId = userId;
    this.amount = amount;
    this.currency = currency;
    this.status = "pending";
    this.createdAt = null;
    this.processedAt = null;
  }

  process(): boolean {
    this.status = "processed";
    this.processedAt = "now";
    return true;
  }

  refund(): boolean {
    this.status = "refunded";
    return true;
  }
}

export class CreditCardPayment extends Payment {
  cardNumber: string;
  cvv: string;
  expiry: string;

  constructor(paymentId: string | null, userId: string, amount: number, cardNumber: string, cvv: string, expiry: string) {
    super(paymentId, userId, amount);
    this.cardNumber = cardNumber;
    this.cvv = cvv;
    this.expiry = expiry;
  }

  process(): boolean {
    const result = creditCardProcessor(this.cardNumber, this.cvv, this.expiry, this.amount);
    if (result) {
      this.status = "processed";
    }
    return result;
  }
}

export class PayPalPayment extends Payment {
  paypalEmail: string;

  constructor(paymentId: string | null, userId: string, amount: number, paypalEmail: string) {
    super(paymentId, userId, amount);
    this.paypalEmail = paypalEmail;
  }

  process(): boolean {
    const result = paypalProcessor(this.paypalEmail, this.amount);
    if (result) {
      this.status = "processed";
    }
    return result;
  }
}

export class Refund {
  refundId: string;
  paymentId: string;
  amount: number;
  status: string;

  constructor(refundId: string, paymentId: string, amount: number) {
    this.refundId = refundId;
    this.paymentId = paymentId;
    this.amount = amount;
    this.status = "pending";
  }

  process(): boolean {
    this.status = "processed";
    return true;
  }
}

