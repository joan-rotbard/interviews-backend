/**
 * Main entry point - Simple CLI to test the payment service
 */

import * as readline from 'readline';
import { PaymentService } from './payment/service';
import * as db from './database/db';

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

function question(query: string): Promise<string> {
  return new Promise(resolve => rl.question(query, resolve));
}

async function main() {
  console.log("Payment Microservice - Test Interface");
  console.log("=".repeat(50));
  
  const service = new PaymentService();
  
  // Initialize some test data
  db.getConnection().users["user_123"] = { balance: 1000.0 };
  
  while (true) {
    console.log("\nOptions:");
    console.log("1. Process Credit Card Payment");
    console.log("2. Process PayPal Payment");
    console.log("3. Process Refund");
    console.log("4. Get User Transactions");
    console.log("5. Get Payment Status");
    console.log("6. Exit");
    
    const choice = (await question("\nSelect option: ")).trim();
    
    if (choice === "1") {
      try {
        const userId = (await question("User ID: ")).trim();
        const amount = parseFloat((await question("Amount: ")).trim());
        const cardNumber = (await question("Card Number: ")).trim();
        const cvv = (await question("CVV: ")).trim();
        const expiry = (await question("Expiry (MM/YY): ")).trim();
        
        const paymentId = service.processPayment(
          userId, amount, "credit_card",
          { card_number: cardNumber, cvv: cvv, expiry: expiry }
        );
        console.log(`Payment processed: ${paymentId}`);
      } catch (e: any) {
        console.log(`Error: ${e.message}`);
      }
    } else if (choice === "2") {
      try {
        const userId = (await question("User ID: ")).trim();
        const amount = parseFloat((await question("Amount: ")).trim());
        const email = (await question("PayPal Email: ")).trim();
        
        const paymentId = service.processPayment(
          userId, amount, "paypal",
          { email: email }
        );
        console.log(`Payment processed: ${paymentId}`);
      } catch (e: any) {
        console.log(`Error: ${e.message}`);
      }
    } else if (choice === "3") {
      try {
        const paymentId = (await question("Payment ID: ")).trim();
        const amount = parseFloat((await question("Refund Amount: ")).trim());
        
        const refundId = service.processRefund(paymentId, amount);
        console.log(`Refund processed: ${refundId}`);
      } catch (e: any) {
        console.log(`Error: ${e.message}`);
      }
    } else if (choice === "4") {
      try {
        const userId = (await question("User ID: ")).trim();
        const transactions = service.getUserTransactions(userId);
        console.log(`\nTransactions for ${userId}:`);
        transactions.forEach(txn => {
          console.log(`  ${txn.paymentId}: ${txn.amount} ${txn.currency} - ${txn.status}`);
        });
      } catch (e: any) {
        console.log(`Error: ${e.message}`);
      }
    } else if (choice === "5") {
      try {
        const paymentId = (await question("Payment ID: ")).trim();
        const status = service.getPaymentStatus(paymentId);
        if (status) {
          console.log(`Payment status: ${status}`);
        } else {
          console.log("Payment not found");
        }
      } catch (e: any) {
        console.log(`Error: ${e.message}`);
      }
    } else if (choice === "6") {
      break;
    } else {
      console.log("Invalid option");
    }
  }
  
  rl.close();
}

if (require.main === module) {
  main().catch(console.error);
}

