# 🔧 FIX DATABASE SEED ERROR

## ❌ Error Yang Terjadi

```
ERROR: column "updated_at" of relation "permissions" does not exist
```

## ✅ Solusi

Table `permissions`, `roles`, `role_permissions`, dan `user_roles` tidak ada di schema asli.

### **OPTION 1: Gunakan Seed Data yang Sudah Diperbaiki (RECOMMENDED)**

```powershell
# Drop dan recreate database
dropdb -U postgres yayasan_erp
createdb -U postgres yayasan_erp

# Run schema
Get-Content database/schema.sql | psql -U postgres yayasan_erp

# Run FIXED seed data (simplified)
Get-Content database/seeds/fixed_initial_data.sql | psql -U postgres yayasan_erp
```

### **OPTION 2: Tambahkan RBAC Tables (Jika butuh permissions system)**

```powershell
# 1. Drop dan recreate database
dropdb -U postgres yayasan_erp
createdb -U postgres yayasan_erp

# 2. Run schema
Get-Content database/schema.sql | psql -U postgres yayasan_erp

# 3. Add RBAC tables
Get-Content database/fix_schema_add_rbac.sql | psql -U postgres yayasan_erp

# 4. Run original seed (with RBAC)
Get-Content database/seed_data.sql | psql -U postgres yayasan_erp
```

---

## 📊 Verifikasi

Setelah run seed data, harusnya muncul:

```
✅ SEED DATA COMPLETED!

 Branches: | 4
 Users:    | 5
 Accounts: | 18
 Settings: | 11

Default Users (password: admin123):
  superadmin | superadmin@yayasan.org | Super Administrator | admin
  admin      | admin@yayasan.org      | System Administrator | admin
  manager    | manager@yayasan.org    | Finance Manager      | manager
  staff      | staff@yayasan.org      | Admin Staff          | staff
  viewer     | viewer@yayasan.org     | Data Viewer          | viewer
```

---

## 🎯 Perbedaan 2 Versi

### **fixed_initial_data.sql** (Simple - RECOMMENDED):
- ✅ Branches (4)
- ✅ Users (5)
- ✅ Chart of Accounts (18)
- ✅ Settings (11)
- ✅ NO RBAC complexity
- ✅ Langsung bisa dipakai

### **Original seed_data.sql** (With RBAC):
- ✅ Semua yang di atas
- ✅ + Roles (8)
- ✅ + Permissions (46)
- ✅ + Role Permissions (mappings)
- ✅ + User Roles (mappings)
- ❌ Butuh RBAC tables dulu

---

## 🚀 Quick Fix (1 Command)

```powershell
# RECOMMENDED: Simple seed
dropdb -U postgres yayasan_erp; createdb -U postgres yayasan_erp; Get-Content database/schema.sql | psql -U postgres yayasan_erp; Get-Content database/seeds/fixed_initial_data.sql | psql -U postgres yayasan_erp
```

---

## ✅ Default Login

Setelah seed data berhasil:

```
Username: admin
Password: admin123

atau

Username: superadmin  
Password: admin123
```

---

**Status:** ✅ Fixed!  
**File:** `database/seeds/fixed_initial_data.sql`  
**Ready:** Use the command above!
