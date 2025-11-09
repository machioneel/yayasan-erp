# 📊 DATABASE SCHEMA - Yayasan ERP

## 📋 Overview

Complete PostgreSQL database schema for Yayasan ERP system.

## 📁 Files

- `schema.sql` - Complete database schema (all tables, indexes, triggers)
- `seed_data.sql` - Initial/sample data
- `test_schema.sql` - Comprehensive test suite
- `migrations/` - Incremental migration files

## 🗃️ Database Tables (18 total)

### Core Tables:
1. **branches** - Organization branches
2. **users** - System users & authentication
3. **students** - Student master data
4. **student_parents** - Parent/guardian info
5. **invoices** - Student invoices
6. **invoice_items** - Invoice line items
7. **payments** - Payment transactions
8. **accounts** - Chart of accounts
9. **journals** - Journal entries
10. **journal_items** - Journal line items
11. **employees** - Employee master data
12. **attendance** - Employee attendance
13. **payroll** - Payroll processing
14. **assets** - Fixed assets
15. **inventory_items** - Inventory master
16. **inventory_transactions** - Stock movements
17. **settings** - System settings
18. **audit_logs** - Audit trail

## 🚀 Quick Start

### 1. Create Database

```bash
createdb yayasan_erp
```

### 2. Run Schema

```bash
psql -d yayasan_erp -f schema.sql
```

### 3. Load Seed Data

```bash
psql -d yayasan_erp -f seed_data.sql
```

### 4. Test Schema

```bash
psql -d yayasan_erp -f test_schema.sql
```

## 🔧 Features

### ✅ Implemented:

- **UUID Primary Keys** - All tables use UUID
- **Soft Deletes** - deleted_at column where applicable
- **Timestamps** - created_at, updated_at on all tables
- **Foreign Keys** - Proper relationships with cascades
- **Indexes** - Optimized queries on common fields
- **Triggers** - Auto-update updated_at timestamps
- **Views** - Convenient data access
- **Constraints** - Data integrity enforced
- **Comments** - Table documentation

### 🔒 Security:

- Password hashing (bcrypt)
- Role-based access control
- Audit logging
- Referential integrity

## 📊 Default Users

| Username | Password | Role | Email |
|----------|----------|------|-------|
| admin | admin123 | admin | admin@yayasan.org |
| manager | admin123 | manager | manager@yayasan.org |
| staff | admin123 | staff | staff@yayasan.org |
| viewer | admin123 | viewer | viewer@yayasan.org |

⚠️ **IMPORTANT:** Change passwords in production!

## 🔄 Migrations

Migrations are in `migrations/` directory:
- `001_initial_schema.sql` - Initial database setup

To run migrations:

```bash
psql -d yayasan_erp -f migrations/001_initial_schema.sql
```

## 📈 Views

### Available Views:

1. **v_students** - Students with branch info
2. **v_invoice_summary** - Invoices with student & outstanding
3. **v_employees** - Employees with branch info
4. **v_assets** - Assets with branch info

## 🧪 Testing

Run comprehensive tests:

```bash
psql -d yayasan_erp -f test_schema.sql
```

Tests include:
- ✅ Table creation
- ✅ Primary keys
- ✅ Foreign keys
- ✅ Indexes
- ✅ INSERT/UPDATE/DELETE operations
- ✅ Triggers
- ✅ Views
- ✅ Referential integrity

## 📝 Schema Validation

### Check if schema is valid:

```sql
-- Count tables
SELECT COUNT(*) FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
-- Should return: 18

-- Count foreign keys
SELECT COUNT(*) FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';
-- Should return: 20+

-- Count indexes
SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public';
-- Should return: 30+
```

## 🔍 Common Queries

### Check table sizes:

```sql
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Check all indexes:

```sql
SELECT tablename, indexname 
FROM pg_indexes 
WHERE schemaname = 'public'
ORDER BY tablename, indexname;
```

## ⚙️ Configuration

### PostgreSQL Requirements:

- PostgreSQL 14+
- Extensions: uuid-ossp
- Encoding: UTF8

### Recommended Settings:

```sql
-- In postgresql.conf
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
```

## 🐛 Troubleshooting

### Error: Extension uuid-ossp does not exist

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Error: Permission denied

```bash
# Grant permissions
GRANT ALL PRIVILEGES ON DATABASE yayasan_erp TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
```

## 📚 Documentation

- PostgreSQL Docs: https://www.postgresql.org/docs/
- UUID Functions: https://www.postgresql.org/docs/current/uuid-ossp.html

## ✅ Validation Checklist

- [x] All tables created
- [x] Primary keys defined
- [x] Foreign keys with proper cascades
- [x] Indexes on common query fields
- [x] Triggers for updated_at
- [x] Views for reporting
- [x] Seed data loaded
- [x] Default admin user
- [x] Chart of accounts structure
- [x] Audit logging ready
- [x] Soft delete implemented
- [x] UUID primary keys
- [x] Timestamps on all tables
- [x] Referential integrity enforced

## 🎯 Next Steps

1. ✅ Run schema.sql
2. ✅ Load seed_data.sql
3. ✅ Test with test_schema.sql
4. ✅ Verify all tests pass
5. ✅ Change default passwords
6. ✅ Configure backups
7. ✅ Set up monitoring

## 📞 Support

For issues or questions:
- Check test_schema.sql results
- Review PostgreSQL logs
- Verify all prerequisites met

---

**Status:** ✅ Production Ready  
**Version:** 1.0  
**Last Updated:** November 8, 2025
