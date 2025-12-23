# Backend Engineer Interview Task: Payment Microservice Refactoring

## Context
We built a payment processing microservice for our e-commerce platform. The service handles payment processing, refunds, and transaction history.

However, after deploying to production, we've discovered several critical issues:
1. **Race Conditions:** Multiple payments are being processed simultaneously, causing duplicate charges.
2. **Data Integrity:** Transactions are sometimes lost or duplicated in the database.
3. **Error Handling:** The service crashes unexpectedly and doesn't handle edge cases properly.
4. **Code Quality:** The codebase is tightly coupled, making it hard to test and maintain.
5. **Scalability:** The current design doesn't scale well and has performance bottlenecks.

## The Mission
Your goal is to **identify and fix** the most critical issues in the payment service.
You have **45-60 minutes**.

Focus on:
- **Object-Oriented Design:** Improve class structure, use proper abstractions, and reduce coupling.
- **Microservice Architecture:** Ensure proper transaction handling, error recovery, and idempotency.
- **Code Quality:** Make the code more maintainable and testable.

Prioritize what you think will have the biggest impact. We value pragmatic engineering over perfection.

## Setup

Choose your preferred language: **Python**, **TypeScript**, or **Go**. All three implementations have the same issues and structure.

### Python

See [python/README.md](python/README.md) for detailed instructions.

Quick start:
```bash
cd python
pip install -r requirements.txt
python main.py
```

### TypeScript

See [typescript/README.md](typescript/README.md) for detailed instructions.

Quick start:
```bash
cd typescript
npm install
npm start
```

### Go

See [go/README.md](go/README.md) for detailed instructions.

Quick start:
```bash
cd go
go mod download
go run src/main.go
```

## Repository Structure

Each language implementation follows the same structure:

- `src/`: Main source code
  - `payment/`: Payment processing logic
  - `models/`: Data models and entities
  - `database/`: Database access layer
- `tests/`: Test files
- `main.*`: Entry point

## What We're Looking For
- **Problem Identification:** Can you spot the issues quickly?
- **Design Thinking:** How do you approach refactoring?
- **Trade-offs:** What would you prioritize and why?
- **Code Quality:** Can you write clean, maintainable code?

## Deliverables
1. **Documentation:** List the issues you found and your proposed solutions (5-10 min)
2. **Code Changes:** Implement the most critical fixes
3. **Brief Explanation:** Why you chose these fixes and what trade-offs you considered

Send compressed folder with results to jrotbard@volans.world

## Good Luck!

