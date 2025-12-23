/**
 * Test file for payment service
 * Run with: npm test
 */

import * as db from '../src/database/db';
import { PaymentService } from '../src/payment/service';

describe('Payment Service Tests', () => {
  beforeEach(() => {
    // Reset database before each test
    const connection = db.getConnection();
    Object.keys(connection.transactions).forEach(key => delete connection.transactions[key]);
    Object.keys(connection.users).forEach(key => delete connection.users[key]);
  });

  test('Concurrent Payments', () => {
    // Reset database
    db.getConnection().users['user_test'] = { balance: 1000.0 };
    
    const service = new PaymentService();
    
    // Simulate concurrent payments
    const results: string[] = [];
    const promises: Promise<void>[] = [];
    
    const processPayment = async () => {
      try {
        const paymentId = service.processPayment(
          'user_test', 100.0, 'paypal',
          { email: 'test@example.com' }
        );
        results.push(paymentId);
      } catch (e: any) {
        results.push(`Error: ${e.message}`);
      }
    };
    
    // Create multiple promises to process payments simultaneously
    for (let i = 0; i < 5; i++) {
      promises.push(processPayment());
    }
    
    return Promise.all(promises).then(() => {
      // Check for duplicate payment IDs
      console.log(`Results: ${results}`);
      const uniqueResults = new Set(results);
      console.log(`Unique payments: ${uniqueResults.size}`);
      
      // Check user balance
      const balance = db.getUserBalance('user_test');
      console.log(`Final balance: ${balance}`);
      console.log(`Expected balance: 500.0 (1000 - 5*100)`);
      console.log(`Actual balance: ${balance}`);
    });
  });

  test('Multiple Refunds', () => {
    // Reset database
    db.getConnection().users['user_test2'] = { balance: 1000.0 };
    
    const service = new PaymentService();
    
    // Process a payment
    const paymentId = service.processPayment(
      'user_test2', 100.0, 'paypal',
      { email: 'test@example.com' }
    );
    
    const initialBalance = db.getUserBalance('user_test2');
    console.log(`Balance after payment: ${initialBalance}`);
    
    // Process refund multiple times
    service.processRefund(paymentId, 100.0);
    service.processRefund(paymentId, 100.0);
    service.processRefund(paymentId, 100.0);
    
    const finalBalance = db.getUserBalance('user_test2');
    console.log(`Balance after 3 refunds: ${finalBalance}`);
    console.log(`Expected balance: 1000.0`);
    console.log(`Actual balance: ${finalBalance}`);
  });
});

