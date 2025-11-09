# Yayasan ERP - Backend API

Complete ERP system backend for Indonesian educational foundations (Yayasan).

## 🚀 Features

- **Multi-Branch Management** - Support multiple school branches
- **Finance & Accounting** - Complete accounting with 363 COA (Indonesia standard)
- **Student Management** - Registration, billing, payments
- **HR & Payroll** - Employee management, payroll with PPh 21 calculation
- **Asset Management** - Fixed assets with depreciation
- **Inventory Management** - Stock tracking, opname

## 📋 Requirements

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git (optional)

## 🛠️ Installation

### 1. Install Go

**Linux (Ubuntu/Debian):**
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

**Windows:**
- Download from: https://go.dev/dl/
- Run installer
- Verify: `go version`

**MacOS:**
```bash
brew install go
go version
```

### 2. Install PostgreSQL

See [SETUP-AND-TEST-GUIDE.md](../SETUP-AND-TEST-GUIDE.md) for detailed PostgreSQL installation.

**Quick Setup:**
```bash
# Create database
sudo -u postgres psql
CREATE DATABASE yayasan_erp;
CREATE USER yayasan_user WITH PASSWORD 'StrongPassword123!';
GRANT ALL PRIVILEGES ON DATABASE yayasan_erp TO yayasan_user;
\q
```

### 3. Clone & Setup Project

```bash
cd yayasan-erp/backend

# Install dependencies
go mod download

# Copy config
cp config/config.yaml.example config/config.yaml

# Edit config (set your database password)
nano config/config.yaml
```

### 4. Run Application

```bash
# Development mode
go run cmd/api/main.go

# Or build and run
go build -o bin/server cmd/api/main.go
./bin/server
```

Server will start on **http://localhost:8080**

### 5. Seed Initial Data

```bash
# In another terminal
go run cmd/seed/main.go
```

This creates:
- 3 branches
- 4 roles
- 2 users (admin/finance)
- 19 basic accounts

**Login:** admin / admin123

## 📁 Project Structure

```
backend/
├── cmd/
│   ├── api/           # Main application
│   └── seed/          # Database seeder
├── config/            # Configuration files
├── internal/
│   ├── config/        # Config loader
│   ├── database/      # Database connection
│   ├── handler/       # HTTP handlers (13 files)
│   ├── middleware/    # Middlewares (4 files)
│   ├── models/        # Data models (15 files)
│   ├── repository/    # Data access layer (15 files)
│   ├── routes/        # Route definitions
│   ├── service/       # Business logic (15 files)
│   └── utils/         # Utility functions
├── go.mod             # Go dependencies
└── README.md
```

## 🔌 API Endpoints

Total: **116 endpoints**

### Authentication
```
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
```

### Students
```
GET    /api/v1/students
GET    /api/v1/students/:id
POST   /api/v1/students
PUT    /api/v1/students/:id
DELETE /api/v1/students/:id
GET    /api/v1/students/search
```

### Finance
```
GET    /api/v1/accounts
GET    /api/v1/journals
POST   /api/v1/journals
POST   /api/v1/journals/:id/submit
POST   /api/v1/journals/:id/approve
GET    /api/v1/reports/trial-balance
GET    /api/v1/reports/balance-sheet
GET    /api/v1/reports/income-statement
```

### HR & Payroll
```
GET    /api/v1/employees
POST   /api/v1/employees
GET    /api/v1/payrolls
POST   /api/v1/payrolls
POST   /api/v1/payrolls/bulk
```

### Assets
```
GET    /api/v1/assets
POST   /api/v1/assets
POST   /api/v1/assets/maintenance
POST   /api/v1/assets/transfer
```

### Inventory
```
GET    /api/v1/inventory/items
POST   /api/v1/inventory/stock-in
POST   /api/v1/inventory/stock-out
POST   /api/v1/inventory/opname
```

[Full API documentation available in Swagger/OpenAPI format]

## ⚙️ Configuration

Edit `config/config.yaml`:

```yaml
database:
  host: "localhost"
  port: 5432
  name: "yayasan_erp"
  user: "yayasan_user"
  password: "YOUR_PASSWORD"

jwt:
  secret: "YOUR_SECRET_KEY"
  expiration: 24

cors:
  allowed_origins:
    - "http://localhost:3000"
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/service/...
```

## 🔐 Security

- JWT authentication
- Role-based access control (RBAC)
- Permission-based authorization
- Password hashing (bcrypt)
- SQL injection prevention (GORM)
- CORS protection
- Rate limiting (optional)

## 📊 Database Schema

42 tables including:
- users, roles, permissions
- branches
- accounts (363 COA)
- journals, journal_items
- students, parents, student_parents
- invoices, invoice_items, payments
- employees, employment_contracts
- payrolls, payroll_items, attendances
- assets, asset_categories, asset_transfers
- inventory_items, stock_transactions
- And more...

## 🚀 Production Deployment

### Build for Production

```bash
# Build
go build -o bin/server cmd/api/main.go

# Run
./bin/server
```

### Using Docker (Optional)

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o server cmd/api/main.go
CMD ["./server"]
```

### Systemd Service (Linux)

```bash
sudo nano /etc/systemd/system/yayasan-erp.service
```

```ini
[Unit]
Description=Yayasan ERP Backend
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/yayasan-erp/backend
ExecStart=/opt/yayasan-erp/backend/bin/server
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable yayasan-erp
sudo systemctl start yayasan-erp
```

## 📝 Environment Variables

Alternative to config.yaml:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=yayasan_erp
export DB_USER=yayasan_user
export DB_PASSWORD=password
export JWT_SECRET=your-secret-key
export APP_PORT=8080
```

## 🐛 Troubleshooting

### Database connection failed
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U yayasan_user -d yayasan_erp -h localhost
```

### Port already in use
```bash
# Find process on port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Module not found
```bash
go mod tidy
go mod download
```

## 📚 Documentation

- [Setup Guide](../SETUP-AND-TEST-GUIDE.md)
- [API Documentation](./API.md) (coming soon)
- [Database Schema](./SCHEMA.md) (coming soon)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

Private - Yayasan ERP System

## 📞 Support

For issues and questions, please open an issue on GitHub.

---

**Built with ❤️ for Indonesian educational foundations**
