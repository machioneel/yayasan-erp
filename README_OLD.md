# Yayasan ERP System

Enterprise Resource Planning system untuk Yayasan (Pendidikan, Sosial, dan Keagamaan)

## 🎯 Project Overview

Sistem ERP terintegrasi untuk mengelola operasional yayasan meliputi:
- Finance & Accounting (Fund Accounting)
- Inventory Management
- Sales & CRM (Donor Management)
- Purchase Management
- Asset Management
- Reporting & Analytics

## 🏗️ Architecture

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP framework)
- **Database**: PostgreSQL 15+
- **ORM**: GORM
- **Authentication**: JWT
- **API Style**: RESTful

### Frontend
- **Framework**: React 18+
- **State Management**: Redux Toolkit / Zustand
- **UI Library**: Ant Design / Material-UI
- **HTTP Client**: Axios
- **Routing**: React Router v6

### Database
- **Primary DB**: PostgreSQL
- **Features**: Multi-branch, Multi-tenant ready
- **Backup**: Automated daily backups

## 📁 Project Structure

```
yayasan-erp/
├── backend/                    # Go backend
│   ├── cmd/
│   │   └── api/               # Main application
│   ├── internal/
│   │   ├── config/            # Configuration
│   │   ├── database/          # Database connection
│   │   ├── middleware/        # HTTP middlewares
│   │   ├── models/            # Data models
│   │   ├── repository/        # Data access layer
│   │   ├── service/           # Business logic
│   │   ├── handler/           # HTTP handlers
│   │   └── utils/             # Utilities
│   ├── migrations/            # Database migrations
│   ├── scripts/               # Utility scripts
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # React frontend
│   ├── public/
│   ├── src/
│   │   ├── api/               # API clients
│   │   ├── assets/            # Static assets
│   │   ├── components/        # Reusable components
│   │   ├── features/          # Feature modules
│   │   │   ├── auth/
│   │   │   ├── finance/
│   │   │   ├── inventory/
│   │   │   ├── sales/
│   │   │   ├── purchase/
│   │   │   ├── assets/
│   │   │   └── reports/
│   │   ├── hooks/             # Custom hooks
│   │   ├── layouts/           # Layout components
│   │   ├── routes/            # Route definitions
│   │   ├── store/             # State management
│   │   ├── utils/             # Utilities
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
│
├── database/
│   ├── schema/                # Database schema files
│   ├── migrations/            # Migration files
│   └── seeds/                 # Seed data
│
├── docs/                      # Documentation
│   ├── api/                   # API documentation
│   ├── user-guide/            # User manual
│   └── technical/             # Technical docs
│
├── docker/                    # Docker configs
│   ├── backend.Dockerfile
│   ├── frontend.Dockerfile
│   └── postgres.Dockerfile
│
├── scripts/                   # Deployment scripts
│   ├── deploy.sh
│   ├── backup.sh
│   └── restore.sh
│
└── docker-compose.yml
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Git

### Backend Setup

```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your configurations
nano .env

# Run migrations
go run cmd/api/main.go migrate

# Start server
go run cmd/api/main.go
```

### Frontend Setup

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Edit .env with your API URL
nano .env

# Start development server
npm run dev
```

### Database Setup

```bash
# Create database
createdb yayasan_erp

# Run migrations (from backend directory)
go run cmd/api/main.go migrate

# Seed initial data (optional)
go run cmd/api/main.go seed
```

## 🔐 Default Credentials

**Super Admin:**
- Email: admin@yayasan.org
- Password: Admin123!

**Branch Admin:**
- Email: branch1@yayasan.org
- Password: Branch123!

⚠️ **IMPORTANT**: Change these passwords immediately after first login!

## 📚 Documentation

- [API Documentation](docs/api/README.md)
- [User Guide](docs/user-guide/README.md)
- [Technical Documentation](docs/technical/README.md)
- [Deployment Guide](docs/deployment/README.md)

## 🎯 Module Status

| Module | Status | Version |
|--------|--------|---------|
| Authentication | ✅ Complete | 1.0 |
| Multi-Branch | ✅ Complete | 1.0 |
| Finance & Accounting | 🚧 In Progress | 0.8 |
| Inventory Management | 📋 Planned | - |
| Sales & CRM | 📋 Planned | - |
| Purchase Management | 📋 Planned | - |
| Asset Management | 📋 Planned | - |
| Reporting & Analytics | 📋 Planned | - |

## 🔧 Configuration

### Environment Variables

**Backend (.env)**
```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=yayasan_erp
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# CORS
CORS_ORIGINS=http://localhost:5173

# File Upload
MAX_UPLOAD_SIZE=10485760  # 10MB
UPLOAD_PATH=./uploads
```

**Frontend (.env)**
```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Yayasan ERP
```

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test

# E2E tests
npm run test:e2e
```

## 📦 Deployment

### On-Premise Deployment

```bash
# Build backend
cd backend
go build -o yayasan-erp cmd/api/main.go

# Build frontend
cd frontend
npm run build

# Deploy using provided scripts
./scripts/deploy.sh
```

### Docker Deployment

```bash
# Build and run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🔄 Backup & Restore

### Automated Backup
```bash
# Setup cron job for daily backup
0 2 * * * /path/to/scripts/backup.sh
```

### Manual Backup
```bash
./scripts/backup.sh
```

### Restore
```bash
./scripts/restore.sh backup_file.sql
```

## 🛡️ Security Features

- ✅ JWT-based authentication
- ✅ Role-based access control (RBAC)
- ✅ Branch-level data isolation
- ✅ Password hashing (bcrypt)
- ✅ SQL injection protection (parameterized queries)
- ✅ XSS protection
- ✅ CSRF protection
- ✅ Rate limiting
- ✅ Audit logging

## 👥 User Roles

1. **Super Admin**
   - Full system access
   - Manage all branches
   - System configuration

2. **Branch Admin**
   - Full access to assigned branch
   - User management (branch level)
   - All module access

3. **Finance Manager**
   - Finance & accounting access
   - Financial reports
   - Budget management

4. **Inventory Manager**
   - Inventory management
   - Purchase management
   - Stock reports

5. **Sales/Donor Manager**
   - CRM access
   - Donation management
   - Donor reports

6. **Staff**
   - Limited access based on assignment
   - View-only for most modules

## 📊 Key Features

### Finance & Accounting
- ✅ Fund accounting (restricted/unrestricted)
- ✅ Multi-branch accounting
- ✅ General ledger
- ✅ Accounts payable/receivable
- ✅ Budget tracking
- ✅ Financial reporting
- ✅ Donation/Zakat tracking

### Inventory Management
- ✅ Multi-location inventory
- ✅ Stock tracking
- ✅ Warehouse management
- ✅ Stock movements
- ✅ Inventory valuation

### Sales & CRM
- ✅ Donor management
- ✅ Donation tracking
- ✅ Grant management
- ✅ Communication history
- ✅ Donor reports

### Purchase Management
- ✅ Vendor management
- ✅ Purchase orders
- ✅ Purchase requisitions
- ✅ Receiving goods
- ✅ Vendor evaluation

### Asset Management
- ✅ Asset registry
- ✅ Depreciation tracking
- ✅ Maintenance scheduling
- ✅ Asset transfer
- ✅ Asset reports

### Reporting & Analytics
- ✅ Financial reports
- ✅ Operational reports
- ✅ Custom report builder
- ✅ Dashboard analytics
- ✅ Export to PDF/Excel

## 🤝 Contributing

This is an internal project for the Yayasan. For contributions:
1. Create a feature branch
2. Make your changes
3. Submit for review
4. Wait for approval before merging

## 📝 License

Proprietary - Internal use only for Yayasan operations

## 📞 Support

For technical support:
- Email: it@yayasan.org
- Phone: +62-xxx-xxxx-xxxx
- Internal ticketing system

## 🗺️ Roadmap

### Q1 2024
- ✅ Authentication & Authorization
- ✅ Multi-branch setup
- 🚧 Finance & Accounting module

### Q2 2024
- 📋 Inventory Management
- 📋 Sales & CRM
- 📋 Purchase Management

### Q3 2024
- 📋 Asset Management
- 📋 Advanced Reporting
- 📋 Mobile app (optional)

### Q4 2024
- 📋 API integrations
- 📋 Advanced analytics
- 📋 Performance optimization

## 📈 Version History

### v0.1.0 (Current)
- Initial project setup
- Authentication system
- Multi-branch foundation
- Database schema design

---

**Built with ❤️ for Yayasan operations**
