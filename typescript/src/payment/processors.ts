/**
 * Payment processors - External service integrations
 */

export function creditCardProcessor(cardNumber: string, cvv: string, expiry: string, amount: number): boolean {
  // Simulate network delay
  // Note: In real code this would be async
  const start = Date.now();
  while (Date.now() - start < 100) {
    // Busy wait to simulate delay (100ms)
  }
  
  // Simulate random failures
  if (Math.random() < 0.1) { // 10% failure rate
    throw new Error("Credit card processing failed");
  }
  
  return true;
}

export function paypalProcessor(email: string, amount: number): boolean {
  // Simulate network delay
  // Note: In real code this would be async
  const start = Date.now();
  while (Date.now() - start < 100) {
    // Busy wait to simulate delay (100ms)
  }
  
  // Simulate random failures
  if (Math.random() < 0.1) { // 10% failure rate
    throw new Error("PayPal processing failed");
  }
  
  return true;
}

