# Windows Setup Guide for Stocky

## Prerequisites Installation

### 1. Install Go (Golang)

1. Download Go from: https://go.dev/dl/
2. Download the Windows installer (e.g., `go1.21.x.windows-amd64.msi`)
3. Run the installer and follow the wizard
4. Verify installation by opening a new PowerShell terminal:
   ```powershell
   go version
   ```

### 2. Install PostgreSQL

1. Download PostgreSQL from: https://www.postgresql.org/download/windows/
2. Run the installer and remember the password you set for the `postgres` user
3. Add PostgreSQL bin to PATH if not done automatically:
   ```
   C:\Program Files\PostgreSQL\15\bin
   ```
4. Verify installation:
   ```powershell
   psql --version
   ```

### 3. Setup Database

Open PowerShell as Administrator and run:

```powershell
# Connect to PostgreSQL
psql -U postgres

# In the PostgreSQL prompt, run:
CREATE DATABASE assignment;
CREATE USER stocky_user WITH PASSWORD 'stocky_password';
GRANT ALL PRIVILEGES ON DATABASE assignment TO stocky_user;
\q
```

Alternatively, use pgAdmin (installed with PostgreSQL) to create the database and user.

### 4. Configure Environment

```powershell
# Navigate to project directory
cd "c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky"

# Create .env file from example
Copy-Item .env.example .env

# Edit .env with your database credentials using Notepad
notepad .env
```

Update the database credentials in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=stocky_user
DB_PASSWORD=stocky_password
DB_NAME=assignment
```

### 5. Install Dependencies

```powershell
cd "c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky"
go mod download
```

### 6. Run Database Migrations

```powershell
# Connect to the database
psql -U stocky_user -d assignment

# In PostgreSQL prompt, run the migration:
\i migrations/001_initial_schema.sql
\q
```

Or run directly from PowerShell:
```powershell
Get-Content migrations\001_initial_schema.sql | psql -U stocky_user -d assignment
```

### 7. Start the Server

```powershell
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 8. Test the API

In a new PowerShell terminal:

```powershell
# Test health check
Invoke-RestMethod -Uri "http://localhost:8080/health" | ConvertTo-Json

# Run the full test script
.\test-api.ps1
```

---

## Quick Docker Setup (Alternative)

If you prefer using Docker (no Go installation needed):

### Prerequisites
- Docker Desktop for Windows: https://www.docker.com/products/docker-desktop/

### Steps

```powershell
# Navigate to project
cd "c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky"

# Build and start containers
docker-compose up --build

# Server will be available at http://localhost:8080
```

To stop:
```powershell
docker-compose down
```

---

## Troubleshooting

### Issue: "go: command not found"
**Solution**: Close and reopen PowerShell after installing Go, or add Go to PATH manually:
```powershell
$env:Path += ";C:\Go\bin"
```

### Issue: "psql: command not found"
**Solution**: Add PostgreSQL bin to PATH:
```powershell
$env:Path += ";C:\Program Files\PostgreSQL\15\bin"
```

### Issue: "connection refused" when connecting to database
**Solution**: 
1. Ensure PostgreSQL service is running:
   - Open Services (Win + R, type `services.msc`)
   - Find "postgresql-x64-15" and ensure it's running
2. Check firewall settings
3. Verify credentials in `.env` file

### Issue: Port 8080 already in use
**Solution**: Change the port in `.env`:
```env
PORT=8081
```

---

## Building Executable

To create a standalone executable:

```powershell
# Build for Windows
go build -o stocky.exe cmd/server/main.go

# Run the executable
.\stocky.exe
```

---

## Next Steps

1. Import `Stocky-API.postman_collection.json` into Postman for easy testing
2. Review the `README.md` for full API documentation
3. Check the database to see created tables and data
4. Modify fees configuration in `.env` as needed

---

## Support

For issues, check:
- Go installation: https://go.dev/doc/install
- PostgreSQL Windows guide: https://www.postgresqltutorial.com/postgresql-getting-started/install-postgresql/
- Gin framework docs: https://gin-gonic.com/docs/
