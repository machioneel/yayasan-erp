-- Database Seeder for Yayasan As-Salam Joglo ERP
-- Run this after creating database schema
-- Usage: psql yayasan_erp < database/seeds/initial_data.sql

BEGIN;

-- ============================================================================
-- 1. INSERT BRANCHES
-- ============================================================================

INSERT INTO branches (id, code, name, type, address, city, province, postal_code, phone, email, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'YAY', 'Sekretariat Yayasan', 'sekretariat', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345678', 'sekretariat@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'SDI', 'Sekolah Dasar Islam As-Salam', 'pendidikan', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345679', 'sdi@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'TKI', 'Taman Kanak-kanak Islam As-Salam', 'pendidikan', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345680', 'tki@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'TPA', 'Taman Pendidikan Al-Quran As-Salam', 'pendidikan', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345681', 'tpa@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'MAD', 'Madrasah As-Salam', 'pendidikan', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345682', 'mad@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'BKM', 'Bidang Kemakmuran Masjid', 'masjid', 'Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345683', 'bkm@assalamjoglo.or.id', true, NOW(), NOW()),
    (gen_random_uuid(), 'AYD', 'Anak Yatim & Dhuafa', 'sosial', 'Jl. Masjid As-Salam Joglo', 'Jakarta Barat', 'DKI Jakarta', '11640', '021-12345684', 'ayd@assalamjoglo.or.id', true, NOW(), NOW());

-- ============================================================================
-- 2. INSERT PERMISSIONS
-- ============================================================================

INSERT INTO permissions (id, code, name, description, module, created_at, updated_at)
VALUES
    -- User Management
    (gen_random_uuid(), 'users.view', 'View Users', 'Can view user list', 'users', NOW(), NOW()),
    (gen_random_uuid(), 'users.create', 'Create User', 'Can create new users', 'users', NOW(), NOW()),
    (gen_random_uuid(), 'users.update', 'Update User', 'Can update user information', 'users', NOW(), NOW()),
    (gen_random_uuid(), 'users.delete', 'Delete User', 'Can delete users', 'users', NOW(), NOW()),
    
    -- Branch Management
    (gen_random_uuid(), 'branches.view', 'View Branches', 'Can view branch list', 'branches', NOW(), NOW()),
    (gen_random_uuid(), 'branches.create', 'Create Branch', 'Can create new branches', 'branches', NOW(), NOW()),
    (gen_random_uuid(), 'branches.update', 'Update Branch', 'Can update branch information', 'branches', NOW(), NOW()),
    (gen_random_uuid(), 'branches.delete', 'Delete Branch', 'Can delete branches', 'branches', NOW(), NOW()),
    
    -- Role Management
    (gen_random_uuid(), 'roles.view', 'View Roles', 'Can view role list', 'roles', NOW(), NOW()),
    (gen_random_uuid(), 'roles.create', 'Create Role', 'Can create new roles', 'roles', NOW(), NOW()),
    (gen_random_uuid(), 'roles.update', 'Update Role', 'Can update role information', 'roles', NOW(), NOW()),
    (gen_random_uuid(), 'roles.delete', 'Delete Role', 'Can delete roles', 'roles', NOW(), NOW()),
    
    -- Finance - Accounts
    (gen_random_uuid(), 'accounts.view', 'View Accounts', 'Can view chart of accounts', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'accounts.create', 'Create Account', 'Can create new accounts', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'accounts.update', 'Update Account', 'Can update account information', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'accounts.delete', 'Delete Account', 'Can delete accounts', 'finance', NOW(), NOW()),
    
    -- Finance - Journal Entries
    (gen_random_uuid(), 'journals.view', 'View Journals', 'Can view journal entries', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'journals.create', 'Create Journal', 'Can create journal entries', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'journals.update', 'Update Journal', 'Can update journal entries', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'journals.delete', 'Delete Journal', 'Can delete journal entries', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'journals.approve', 'Approve Journal', 'Can approve journal entries', 'finance', NOW(), NOW()),
    (gen_random_uuid(), 'journals.post', 'Post Journal', 'Can post journal entries', 'finance', NOW(), NOW()),
    
    -- Assets
    (gen_random_uuid(), 'assets.view', 'View Assets', 'Can view asset list', 'assets', NOW(), NOW()),
    (gen_random_uuid(), 'assets.create', 'Create Asset', 'Can create new assets', 'assets', NOW(), NOW()),
    (gen_random_uuid(), 'assets.update', 'Update Asset', 'Can update asset information', 'assets', NOW(), NOW()),
    (gen_random_uuid(), 'assets.delete', 'Delete Asset', 'Can delete assets', 'assets', NOW(), NOW()),
    (gen_random_uuid(), 'assets.approve', 'Approve Asset', 'Can approve asset transactions', 'assets', NOW(), NOW()),
    
    -- Reports
    (gen_random_uuid(), 'reports.finance', 'Financial Reports', 'Can view financial reports', 'reports', NOW(), NOW()),
    (gen_random_uuid(), 'reports.assets', 'Asset Reports', 'Can view asset reports', 'reports', NOW(), NOW()),
    (gen_random_uuid(), 'reports.donor', 'Donor Reports', 'Can view donor reports', 'reports', NOW(), NOW());

-- ============================================================================
-- 3. INSERT ROLES
-- ============================================================================

-- Get branch IDs
DO $$ 
DECLARE 
    v_yay_branch_id UUID;
BEGIN
    SELECT id INTO v_yay_branch_id FROM branches WHERE code = 'YAY';

    -- Super Admin Role
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'SUPER_ADMIN', 'Super Administrator', 'Full system access', true, NOW(), NOW());

    -- Bendahara Umum
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'BENDAHARA_UMUM', 'Bendahara Umum', 'Chief Financial Officer', true, NOW(), NOW());

    -- Bendahara Unit
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'BENDAHARA_UNIT', 'Bendahara Unit', 'Unit Financial Officer', true, NOW(), NOW());

    -- Staff Keuangan
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'STAFF_KEUANGAN', 'Staff Keuangan', 'Finance Staff (Maker)', true, NOW(), NOW());

    -- Sekretaris Umum
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'SEKRETARIS_UMUM', 'Sekretaris Umum', 'General Secretary - Asset Management', true, NOW(), NOW());

    -- Staff Inventaris
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'STAFF_INVENTARIS', 'Staff Inventaris', 'Asset Management Staff', true, NOW(), NOW());

    -- Ketua Bidang
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'KETUA_BIDANG', 'Ketua Bidang', 'Department Head', true, NOW(), NOW());

    -- Kepala Sekolah
    INSERT INTO roles (id, code, name, description, is_system_role, created_at, updated_at)
    VALUES (gen_random_uuid(), 'KEPALA_SEKOLAH', 'Kepala Sekolah', 'School Principal', true, NOW(), NOW());

END $$;

-- ============================================================================
-- 4. ASSIGN PERMISSIONS TO ROLES
-- ============================================================================

DO $$
DECLARE
    v_super_admin_id UUID;
    v_bendahara_umum_id UUID;
    v_bendahara_unit_id UUID;
    v_staff_keuangan_id UUID;
    v_perm_id UUID;
BEGIN
    -- Get role IDs
    SELECT id INTO v_super_admin_id FROM roles WHERE code = 'SUPER_ADMIN';
    SELECT id INTO v_bendahara_umum_id FROM roles WHERE code = 'BENDAHARA_UMUM';
    SELECT id INTO v_bendahara_unit_id FROM roles WHERE code = 'BENDAHARA_UNIT';
    SELECT id INTO v_staff_keuangan_id FROM roles WHERE code = 'STAFF_KEUANGAN';

    -- Super Admin gets ALL permissions
    FOR v_perm_id IN SELECT id FROM permissions
    LOOP
        INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
        VALUES (gen_random_uuid(), v_super_admin_id, v_perm_id, NOW(), NOW());
    END LOOP;

    -- Bendahara Umum - Finance permissions
    FOR v_perm_id IN SELECT id FROM permissions WHERE module IN ('finance', 'reports')
    LOOP
        INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
        VALUES (gen_random_uuid(), v_bendahara_umum_id, v_perm_id, NOW(), NOW());
    END LOOP;

    -- Bendahara Unit - Finance view and create
    FOR v_perm_id IN SELECT id FROM permissions WHERE code IN ('journals.view', 'journals.create', 'accounts.view', 'reports.finance')
    LOOP
        INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
        VALUES (gen_random_uuid(), v_bendahara_unit_id, v_perm_id, NOW(), NOW());
    END LOOP;

    -- Staff Keuangan - Finance create only
    FOR v_perm_id IN SELECT id FROM permissions WHERE code IN ('journals.view', 'journals.create', 'accounts.view')
    LOOP
        INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
        VALUES (gen_random_uuid(), v_staff_keuangan_id, v_perm_id, NOW(), NOW());
    END LOOP;

END $$;

-- ============================================================================
-- 5. CREATE SUPER ADMIN USER
-- ============================================================================

DO $$
DECLARE
    v_user_id UUID := gen_random_uuid();
    v_branch_id UUID;
    v_role_id UUID;
BEGIN
    -- Get Sekretariat branch
    SELECT id INTO v_branch_id FROM branches WHERE code = 'YAY';
    
    -- Create super admin user
    -- Password: Admin123! (bcrypt hash)
    INSERT INTO users (id, branch_id, username, email, password_hash, full_name, is_active, is_super_admin, created_at, updated_at)
    VALUES (
        v_user_id,
        v_branch_id,
        'admin',
        'admin@assalamjoglo.or.id',
        '$2a$10$YourBcryptHashHere', -- This will be replaced with actual hash
        'System Administrator',
        true,
        true,
        NOW(),
        NOW()
    );

    -- Assign Super Admin role
    SELECT id INTO v_role_id FROM roles WHERE code = 'SUPER_ADMIN';
    INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
    VALUES (gen_random_uuid(), v_user_id, v_role_id, NOW(), NOW());

    RAISE NOTICE 'Super Admin user created: admin / Admin123!';
END $$;

-- ============================================================================
-- 6. CREATE SAMPLE BENDAHARA UMUM USER
-- ============================================================================

DO $$
DECLARE
    v_user_id UUID := gen_random_uuid();
    v_branch_id UUID;
    v_role_id UUID;
BEGIN
    SELECT id INTO v_branch_id FROM branches WHERE code = 'YAY';
    
    INSERT INTO users (id, branch_id, username, email, password_hash, full_name, phone, is_active, created_at, updated_at)
    VALUES (
        v_user_id,
        v_branch_id,
        'bendahara.umum',
        'bendahara@assalamjoglo.or.id',
        '$2a$10$YourBcryptHashHere',
        'Bendahara Umum',
        '081234567890',
        true,
        NOW(),
        NOW()
    );

    SELECT id INTO v_role_id FROM roles WHERE code = 'BENDAHARA_UMUM';
    INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
    VALUES (gen_random_uuid(), v_user_id, v_role_id, NOW(), NOW());

    RAISE NOTICE 'Bendahara Umum user created: bendahara.umum / Admin123!';
END $$;

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- Check inserted data
SELECT 'Branches:', COUNT(*) FROM branches;
SELECT 'Permissions:', COUNT(*) FROM permissions;
SELECT 'Roles:', COUNT(*) FROM roles;
SELECT 'Role Permissions:', COUNT(*) FROM role_permissions;
SELECT 'Users:', COUNT(*) FROM users;
SELECT 'User Roles:', COUNT(*) FROM user_roles;

-- Show admin user
SELECT u.username, u.email, u.full_name, r.name as role
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.username = 'admin';
