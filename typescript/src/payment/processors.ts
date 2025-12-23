/**
 * Payment processors - External service integrations
 */

export function creditCardProcessor(cardNumber: string, cvv: string, expiry: string, amount: number): boolean {
  // Simulate network delay
  // Note: In real code this would be async
  
  // Simulate random failures
  if (Math.random() < 0.1) { // 10% failure rate
    throw new Error("Credit card processing failed");
  }
  
  return true;
}

export function paypalProcessor(email: string, amount: number): boolean {
  // Simulate network delay
  // Note: In real code this would be async
  
  // Simulate random failures
  if (Math.random() < 0.1) { // 10% failure rate
    throw new Error("PayPal processing failed");
  }
  
  return true;
}

