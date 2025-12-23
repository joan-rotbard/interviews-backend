# Python Implementation

## Setup

1. **Install Dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

2. **Run the Service:**
   ```bash
   python main.py
   ```

3. **Test the Service:**
   ```bash
   pytest tests/
   ```

4. **See Issues in Action (Optional):**
   ```bash
   python demo_issues.py
   ```
   This script demonstrates some of the problems you should identify and fix.

## Structure

- `src/`: Main source code
  - `payment/`: Payment processing logic
  - `models/`: Data models and entities
  - `database/`: Database access layer
- `tests/`: Test files
- `main.py`: Entry point
- `demo_issues.py`: Script to demonstrate issues

