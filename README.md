# Yayasan As-Salam Joglo - ERP System

Enterprise Resource Planning system untuk Yayasan Masjid dan Perguruan As-Salam Joglo.

## 📊 Project Status

**Current Phase:** Phase 1 - Foundation ✅ **COMPLETE!**  
**Progress:** 100% of Authentication Module  
**Next:** Phase 2 - Finance & Accounting Module

## 🎯 Features

### ✅ Implemented (Phase 1)
- ✅ Multi-branch/multi-unit support
- ✅ Role-based access control (RBAC)
- ✅ JWT Authentication
- ✅ User management
- ✅ Branch management
- ✅ Role & Permission management
- ✅ Audit logging
- ✅ API documentation ready

### 🔄 In Development (Phase 2)
- Chart of Accounts (COA) management
- Journal Entry system (Maker-Checker-Approver)
- Fund accounting (restricted/unrestricted)
- Donor tracking
- Budget management
- Financial reports

### ⏳ Planned (Phase 3+)
- Asset management & depreciation
- Inventory management
- Purchase management
- CRM/Donor management
- Advanced reporting & analytics

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL 15+
- **Authentication:** JWT (golang-jwt)
- **Validation:** validator/v10
- **Password Hashing:** bcrypt

## 🚀 Quick Start

### 1. Setup Database

```bash
# Create database
createdb yayasan_erp

# Create schema
psql yayasan_erp < database/schema/complete_schema.sql
```

### 2. Generate Password Hash

```bash
cd backend
go run cmd/hashpass/main.go "Admin123!"
# Copy the generated hash
```

### 3. Update Seed Data

Edit `database/seeds/initial_data.sql`:
- Replace `$2a$10$YourBcryptHashHere` with actual hash

### 4. Seed Initial Data

```bash
psql yayasan_erp < database/seeds/initial_data.sql
```

### 5. Configure Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your settings
```

### 6. Run Server

```bash
go mod download
go run cmd/api/main.go
```

Server starts at http://localhost:8080

## 🧪 Testing

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"Admin123!"}'
```

## 🔐 Default Users

| Username | Password | Role |
|----------|----------|------|
| admin | Admin123! | Super Administrator |
| bendahara.umum | Admin123! | Bendahara Umum |

**⚠️ Change these in production!**

## 📚 API Documentation

See full API docs at: [API Documentation](docs/api/)

### Key Endpoints

- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user
- `GET /api/v1/users` - List users
- `GET /api/v1/branches` - List branches
- `GET /api/v1/roles` - List roles

## 🏢 Default Branches

- YAY - Sekretariat
- SDI - Sekolah Dasar Islam
- TKI - Taman Kanak-kanak
- TPA - Taman Pendidikan Al-Quran
- MAD - Madrasah
- BKM - Kemakmuran Masjid
- AYD - Anak Yatim & Dhuafa

## 🔧 Troubleshooting

**Port in use:**
```bash
lsof -i :8080
kill -9 <PID>
```

**DB connection error:**
```bash
pg_isready
psql -l | grep yayasan_erp
```

## 📝 License

Proprietary - Yayasan As-Salam Joglo

---

**Built with ❤️ for Yayasan As-Salam Joglo**
