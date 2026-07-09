# Office File Sharing Application (OfficeFlow)

A secure, reactive document sharing and approval workflow application built using a Go (Echo) backend, an Angular frontend, and a PostgreSQL database.

## System Architecture Overview

- **Frontend**: Angular 17+ with Vanilla CSS (Dark Mode & Micro-animations) running on port `4200`.
- **Backend**: Go 1.26+ (Echo Framework + GORM) running on port `8080`.
- **Database**: PostgreSQL (Auto-migrations & Seeding via GORM) running on port `5432`.

---

## 1. Database Setup (PostgreSQL)

The application requires a PostgreSQL instance running locally.

### Installation
Make sure PostgreSQL is installed and running on your machine:
- **Windows**: Verify that the PostgreSQL service (e.g., `postgresql-x64-18`) is running.
- **Mac/Linux**: Install via Homebrew/APT and start the service:
  ```bash
  brew services start postgresql
  # or
  sudo systemctl start postgresql
  ```

### Creating the Database
1. Open your terminal or Command Prompt.
2. Run `createdb` using the default `postgres` superuser (enter the database password when prompted):
   ```bash
   createdb -U postgres office_files
   ```
   *(Note: The default password configured in the app is `postgres`. If your credentials differ, set the `DATABASE_URL` environment variable.)*

---

## 2. Backend Setup (Go)

The backend exposes a REST API, connects to PostgreSQL, runs migrations, and automatically seeds test data on start.

### Prerequisites
- Go 1.26 or higher.

### Configuration
By default, the backend connects using:
`host=localhost user=postgres password=postgres dbname=office_files port=5432 sslmode=disable`

To override the connection string, set the `DATABASE_URL` environment variable:
```bash
# Windows (PowerShell)
$env:DATABASE_URL="host=localhost user=your_user password=your_password dbname=office_files port=5432 sslmode=disable"

# Mac/Linux
export DATABASE_URL="host=localhost user=your_user password=your_password dbname=office_files port=5432 sslmode=disable"
```

### Installation & Run
1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```
2. Run the application:
   ```bash
   go run cmd/api/main.go
   ```
   On successful run, GORM will migrate database tables and seed test users. The API will listen on `http://localhost:8080`.

### Running Tests
To run backend unit tests (including the resubmission handler tests):
```bash
go test -v ./internal/handlers/...
```

---

## 3. Frontend Setup (Angular)

The Angular frontend provides dashboard controls for document actions, tracking logs, and a PDF viewer.

### Prerequisites
- Node.js v20+ and npm.
- If Node.js is not on your system path, you can run commands pointing directly to your Node path.

### Installation
1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install the package dependencies:
   ```bash
   npm install
   ```

### Run
1. Start the Angular local development server:
   ```bash
   npm start
   ```
2. Open your browser and navigate to [http://localhost:4200](http://localhost:4200).

---

## 4. Test Accounts

The database is pre-seeded with three mock users. You can log in on the login screen by entering one of these email addresses (no password required):

- **Alice Smith** (Uploader / Owner): `alice@office.com`
- **Bob Jones** (Approver): `bob@office.com`
- **Charlie Brown** (Approver): `charlie@office.com`

---

## Key Features Implemented

1. **Document Actions**: Approve, Reject, and Send Back for Revision.
2. **Resubmit or Replace**:
   - Uploaders can replace a document that was sent back.
   - Alternatively, they can **resubmit with comments** without modifying the original file.
3. **Workflow Timeline**: Complete action tracking shown chronologically on a vertical history timeline.
4. **PDF Document Preview**: Embeds an inline browser-native PDF viewer dynamically using `DomSanitizer` inside the document details page.
