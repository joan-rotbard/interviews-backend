# Go Implementation

## Setup

1. **Install Dependencies:**
   ```bash
   go mod download
   ```

2. **Run the Service:**
   ```bash
   go run src/main.go
   ```

3. **Test the Service:**
   ```bash
   go test ./...
   ```

4. **Build the Service:**
   ```bash
   go build -o payment-service src/main.go
   ```

## Structure

- `src/`: Main source code
  - `payment/`: Payment processing logic
  - `models/`: Data models and entities
  - `database/`: Database access layer
- `tests/`: Test files
- `src/main.go`: Entry point

## Requirements

- Go version: 1.21+

