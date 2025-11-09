-- ============================================
-- FIXED SEED DATA FOR YAYASAN ERP
-- No RBAC tables, simple and working!
-- ============================================

BEGIN;

-- ============================================
-- 1. BRANCHES (4 branches)
-- ============================================
INSERT INTO branches (id, code, name, type, address, phone, email, is_active) VALUES
(uuid_generate_v4(), 'HQ', 'Kantor Pusat', 'office', 'Jl. Pendidikan No. 123, Jakarta', '+62 21 1234567', 'hq@yayasan.org', true),
(uuid_generate_v4(), 'SD01', 'SD Yayasan ABC 1', 'school', 'Jl. Sekolah No. 45, Jakarta Selatan', '+62 21 2345678', 'sd1@yayasan.org', true),
(uuid_generate_v4(), 'SMP01', 'SMP Yayasan ABC 1', 'school', 'Jl. Pendidikan No. 67, Jakarta Timur', '+62 21 3456789', 'smp1@yayasan.org', true),
(uuid_generate_v4(), 'SMA01', 'SMA Yayasan ABC 1', 'school', 'Jl. Pendidikan No. 89, Jakarta Barat', '+62 21 4567890', 'sma1@yayasan.org', true)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 2. USERS (5 users - password: admin123)
-- ============================================
INSERT INTO users (id, username, email, password_hash, full_name, role, is_active) VALUES
(uuid_generate_v4(), 'superadmin', 'superadmin@yayasan.org', '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', 'Super Administrator', 'admin', true),
(uuid_generate_v4(), 'admin', 'admin@yayasan.org', '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', 'System Administrator', 'admin', true),
(uuid_generate_v4(), 'manager', 'manager@yayasan.org', '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', 'Finance Manager', 'manager', true),
(uuid_generate_v4(), 'staff', 'staff@yayasan.org', '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', 'Admin Staff', 'staff', true),
(uuid_generate_v4(), 'viewer', 'viewer@yayasan.org', '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', 'Data Viewer', 'viewer', true)
ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 3. CHART OF ACCOUNTS (18 accounts)
-- ============================================
INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
-- Assets
('1-0000', 'ASET', 'asset', NULL, true, 'debit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('1-1000', 'Aset Lancar', 'asset', (SELECT id FROM accounts WHERE code = '1-0000' LIMIT 1), true, 'debit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('1-1100', 'Kas', 'asset', (SELECT id FROM accounts WHERE code = '1-1000' LIMIT 1), false, 'debit', true),
('1-1200', 'Bank', 'asset', (SELECT id FROM accounts WHERE code = '1-1000' LIMIT 1), false, 'debit', true),
('1-1300', 'Piutang', 'asset', (SELECT id FROM accounts WHERE code = '1-1000' LIMIT 1), false, 'debit', true)
ON CONFLICT (code) DO NOTHING;

-- Liabilities
INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('2-0000', 'KEWAJIBAN', 'liability', NULL, true, 'credit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('2-1000', 'Kewajiban Lancar', 'liability', (SELECT id FROM accounts WHERE code = '2-0000' LIMIT 1), true, 'credit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('2-1100', 'Utang Usaha', 'liability', (SELECT id FROM accounts WHERE code = '2-1000' LIMIT 1), false, 'credit', true)
ON CONFLICT (code) DO NOTHING;

-- Equity
INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('3-0000', 'EKUITAS', 'equity', NULL, true, 'credit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('3-1000', 'Modal', 'equity', (SELECT id FROM accounts WHERE code = '3-0000' LIMIT 1), false, 'credit', true)
ON CONFLICT (code) DO NOTHING;

-- Revenue
INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('4-0000', 'PENDAPATAN', 'revenue', NULL, true, 'credit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('4-1000', 'Pendapatan Operasional', 'revenue', (SELECT id FROM accounts WHERE code = '4-0000' LIMIT 1), true, 'credit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('4-1100', 'Pendapatan SPP', 'revenue', (SELECT id FROM accounts WHERE code = '4-1000' LIMIT 1), false, 'credit', true),
('4-1200', 'Pendapatan Donasi', 'revenue', (SELECT id FROM accounts WHERE code = '4-1000' LIMIT 1), false, 'credit', true)
ON CONFLICT (code) DO NOTHING;

-- Expenses
INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('5-0000', 'BEBAN', 'expense', NULL, true, 'debit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('5-1000', 'Beban Operasional', 'expense', (SELECT id FROM accounts WHERE code = '5-0000' LIMIT 1), true, 'debit', true)
ON CONFLICT (code) DO NOTHING;

INSERT INTO accounts (code, name, account_type, parent_id, is_header, normal_balance, is_active) VALUES
('5-1100', 'Beban Gaji', 'expense', (SELECT id FROM accounts WHERE code = '5-1000' LIMIT 1), false, 'debit', true),
('5-1200', 'Beban Utilitas', 'expense', (SELECT id FROM accounts WHERE code = '5-1000' LIMIT 1), false, 'debit', true)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 4. SETTINGS (11 settings)
-- ============================================
INSERT INTO settings (setting_key, setting_value, setting_type, category, description) VALUES
('company_name', 'Yayasan Pendidikan ABC', 'string', 'company', 'Official company name'),
('company_email', 'info@yayasan.org', 'string', 'company', 'Company email address'),
('company_phone', '+62 21 1234567', 'string', 'company', 'Company phone number'),
('invoice_prefix', 'INV', 'string', 'system', 'Invoice number prefix'),
('payment_prefix', 'PAY', 'string', 'system', 'Payment number prefix'),
('employee_prefix', 'EMP', 'string', 'system', 'Employee number prefix'),
('student_prefix', 'STD', 'string', 'system', 'Student number prefix'),
('asset_prefix', 'AST', 'string', 'system', 'Asset code prefix'),
('inventory_prefix', 'ITM', 'string', 'system', 'Inventory item prefix'),
('low_stock_alert', 'true', 'boolean', 'notification', 'Enable low stock alerts'),
('email_notifications', 'true', 'boolean', 'notification', 'Enable email notifications')
ON CONFLICT (setting_key) DO NOTHING;

COMMIT;

-- ============================================
-- VERIFICATION
-- ============================================
\echo ''
\echo '✅ SEED DATA COMPLETED!'
\echo ''
\echo 'Data Summary:'
SELECT 'Branches:' as info, COUNT(*)::text as total FROM branches
UNION ALL
SELECT 'Users:', COUNT(*)::text FROM users
UNION ALL
SELECT 'Accounts:', COUNT(*)::text FROM accounts
UNION ALL
SELECT 'Settings:', COUNT(*)::text FROM settings;

\echo ''
\echo 'Default Users (password: admin123):'
SELECT username, email, full_name, role FROM users ORDER BY username;
\echo ''

-- ============================================
-- END OF SEED DATA
-- ============================================
