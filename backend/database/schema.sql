-- ============================================
-- YAYASAN ERP - COMPLETE DATABASE SCHEMA
-- PostgreSQL 14+
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 1. BRANCHES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'school', 'office', 'warehouse', etc
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    manager_id UUID,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_branches_code ON branches(code);
CREATE INDEX idx_branches_type ON branches(type);
CREATE INDEX idx_branches_is_active ON branches(is_active);

-- ============================================
-- 2. USERS & AUTHENTICATION
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL, -- 'admin', 'manager', 'staff', 'viewer'
    branch_id UUID REFERENCES branches(id),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- ============================================
-- 3. STUDENTS
-- ============================================
CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_number VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    nik VARCHAR(16),
    nisn VARCHAR(20),
    gender VARCHAR(10) NOT NULL, -- 'male', 'female'
    birth_place VARCHAR(100),
    birth_date DATE,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    branch_id UUID REFERENCES branches(id) NOT NULL,
    grade VARCHAR(20),
    class VARCHAR(20),
    academic_year VARCHAR(20),
    enrollment_date DATE,
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'graduated', 'dropped'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_students_student_number ON students(student_number);
CREATE INDEX idx_students_branch_id ON students(branch_id);
CREATE INDEX idx_students_status ON students(status);
CREATE INDEX idx_students_grade ON students(grade);

-- ============================================
-- 4. STUDENT PARENTS/GUARDIANS
-- ============================================
CREATE TABLE IF NOT EXISTS student_parents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    relationship VARCHAR(20) NOT NULL, -- 'father', 'mother', 'guardian'
    full_name VARCHAR(255) NOT NULL,
    nik VARCHAR(16),
    phone VARCHAR(50),
    email VARCHAR(255),
    occupation VARCHAR(100),
    address TEXT,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_student_parents_student_id ON student_parents(student_id);

-- ============================================
-- 5. INVOICES
-- ============================================
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    student_id UUID REFERENCES students(id) NOT NULL,
    branch_id UUID REFERENCES branches(id) NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount DECIMAL(15,2) DEFAULT 0,
    tax DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'partial', 'paid', 'overdue', 'cancelled'
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_student_id ON invoices(student_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);

-- ============================================
-- 6. INVOICE ITEMS
-- ============================================
CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);

-- ============================================
-- 7. PAYMENTS
-- ============================================
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_number VARCHAR(50) UNIQUE NOT NULL,
    invoice_id UUID REFERENCES invoices(id) NOT NULL,
    student_id UUID REFERENCES students(id) NOT NULL,
    payment_date DATE NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- 'cash', 'transfer', 'card', 'ewallet'
    reference_no VARCHAR(100),
    notes TEXT,
    attachment_url TEXT,
    status VARCHAR(20) DEFAULT 'completed', -- 'completed', 'pending', 'cancelled'
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_payments_payment_number ON payments(payment_number);
CREATE INDEX idx_payments_invoice_id ON payments(invoice_id);
CREATE INDEX idx_payments_student_id ON payments(student_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);

-- ============================================
-- 8. CHART OF ACCOUNTS (COA)
-- ============================================
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- 'asset', 'liability', 'equity', 'revenue', 'expense'
    parent_id UUID REFERENCES accounts(id),
    is_header BOOLEAN DEFAULT false,
    normal_balance VARCHAR(10) NOT NULL, -- 'debit', 'credit'
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_accounts_code ON accounts(code);
CREATE INDEX idx_accounts_type ON accounts(account_type);
CREATE INDEX idx_accounts_parent_id ON accounts(parent_id);

-- ============================================
-- 9. JOURNAL ENTRIES
-- ============================================
CREATE TABLE IF NOT EXISTS journals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_number VARCHAR(50) UNIQUE NOT NULL,
    journal_date DATE NOT NULL,
    description TEXT NOT NULL,
    reference_no VARCHAR(100),
    status VARCHAR(20) DEFAULT 'draft', -- 'draft', 'posted', 'approved', 'cancelled'
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_journals_journal_number ON journals(journal_number);
CREATE INDEX idx_journals_journal_date ON journals(journal_date);
CREATE INDEX idx_journals_status ON journals(status);

-- ============================================
-- 10. JOURNAL ITEMS (Lines)
-- ============================================
CREATE TABLE IF NOT EXISTS journal_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_id UUID REFERENCES journals(id) ON DELETE CASCADE,
    account_id UUID REFERENCES accounts(id) NOT NULL,
    description VARCHAR(255),
    debit DECIMAL(15,2) DEFAULT 0,
    credit DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_journal_items_journal_id ON journal_items(journal_id);
CREATE INDEX idx_journal_items_account_id ON journal_items(account_id);

-- ============================================
-- 11. EMPLOYEES
-- ============================================
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_number VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    nik VARCHAR(16),
    npwp VARCHAR(20),
    gender VARCHAR(10) NOT NULL,
    birth_place VARCHAR(100),
    birth_date DATE,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    branch_id UUID REFERENCES branches(id) NOT NULL,
    position VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    employment_type VARCHAR(20) NOT NULL, -- 'permanent', 'contract', 'intern', 'freelance'
    join_date DATE NOT NULL,
    resign_date DATE,
    salary DECIMAL(15,2),
    is_teacher BOOLEAN DEFAULT false,
    education_level VARCHAR(50),
    major VARCHAR(100),
    certification VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive', 'resigned'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_employees_employee_number ON employees(employee_number);
CREATE INDEX idx_employees_branch_id ON employees(branch_id);
CREATE INDEX idx_employees_status ON employees(status);
CREATE INDEX idx_employees_is_teacher ON employees(is_teacher);

-- ============================================
-- 12. ATTENDANCE
-- ============================================
CREATE TABLE IF NOT EXISTS attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID REFERENCES employees(id) NOT NULL,
    attendance_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'present', 'absent', 'late', 'sick', 'permission', 'leave'
    check_in TIME,
    check_out TIME,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, attendance_date)
);

CREATE INDEX idx_attendance_employee_id ON attendance(employee_id);
CREATE INDEX idx_attendance_date ON attendance(attendance_date);
CREATE INDEX idx_attendance_status ON attendance(status);

-- ============================================
-- 13. PAYROLL
-- ============================================
CREATE TABLE IF NOT EXISTS payroll (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payroll_number VARCHAR(50) UNIQUE NOT NULL,
    employee_id UUID REFERENCES employees(id) NOT NULL,
    period_month INTEGER NOT NULL, -- 1-12
    period_year INTEGER NOT NULL,
    base_salary DECIMAL(15,2) NOT NULL,
    allowances DECIMAL(15,2) DEFAULT 0,
    deductions DECIMAL(15,2) DEFAULT 0,
    net_salary DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'approved', 'paid'
    payment_date DATE,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, period_month, period_year)
);

CREATE INDEX idx_payroll_employee_id ON payroll(employee_id);
CREATE INDEX idx_payroll_period ON payroll(period_year, period_month);
CREATE INDEX idx_payroll_status ON payroll(status);

-- ============================================
-- 14. ASSETS
-- ============================================
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id UUID, -- Can reference asset_categories if needed
    category_name VARCHAR(100),
    branch_id UUID REFERENCES branches(id) NOT NULL,
    acquisition_date DATE NOT NULL,
    acquisition_cost DECIMAL(15,2) NOT NULL,
    useful_life INTEGER, -- in years
    salvage_value DECIMAL(15,2) DEFAULT 0,
    book_value DECIMAL(15,2) NOT NULL,
    accumulated_depreciation DECIMAL(15,2) DEFAULT 0,
    brand VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(100),
    location VARCHAR(255),
    condition VARCHAR(50), -- 'excellent', 'good', 'fair', 'poor'
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'maintenance', 'disposed', 'inactive'
    pic_name VARCHAR(255), -- Person in charge
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_assets_asset_code ON assets(asset_code);
CREATE INDEX idx_assets_branch_id ON assets(branch_id);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_category_name ON assets(category_name);

-- ============================================
-- 15. INVENTORY ITEMS
-- ============================================
CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    unit VARCHAR(20) NOT NULL, -- 'pcs', 'box', 'pack', etc
    current_stock INTEGER DEFAULT 0,
    minimum_stock INTEGER DEFAULT 0,
    maximum_stock INTEGER,
    unit_price DECIMAL(15,2) NOT NULL DEFAULT 0,
    brand VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_inventory_items_item_code ON inventory_items(item_code);
CREATE INDEX idx_inventory_items_category ON inventory_items(category);
CREATE INDEX idx_inventory_items_current_stock ON inventory_items(current_stock);

-- ============================================
-- 16. INVENTORY TRANSACTIONS
-- ============================================
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_number VARCHAR(50) UNIQUE NOT NULL,
    transaction_date DATE NOT NULL,
    transaction_type VARCHAR(20) NOT NULL, -- 'in', 'out', 'adjustment', 'transfer'
    item_id UUID REFERENCES inventory_items(id) NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(15,2),
    total_value DECIMAL(15,2),
    source VARCHAR(255), -- For 'in' transactions
    destination VARCHAR(255), -- For 'out' transactions
    reference_no VARCHAR(100),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inventory_transactions_item_id ON inventory_transactions(item_id);
CREATE INDEX idx_inventory_transactions_date ON inventory_transactions(transaction_date);
CREATE INDEX idx_inventory_transactions_type ON inventory_transactions(transaction_type);

-- ============================================
-- 17. SYSTEM SETTINGS
-- ============================================
CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    setting_type VARCHAR(50) NOT NULL, -- 'string', 'number', 'boolean', 'json'
    category VARCHAR(50), -- 'company', 'system', 'notification', etc
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_settings_key ON settings(setting_key);
CREATE INDEX idx_settings_category ON settings(category);

-- ============================================
-- 18. AUDIT LOGS
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'create', 'update', 'delete', 'login', etc
    entity_type VARCHAR(50) NOT NULL, -- 'student', 'invoice', 'payment', etc
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- ============================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables with updated_at
CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_students_updated_at BEFORE UPDATE ON students
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_journals_updated_at BEFORE UPDATE ON journals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_attendance_updated_at BEFORE UPDATE ON attendance
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payroll_updated_at BEFORE UPDATE ON payroll
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON assets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inventory_items_updated_at BEFORE UPDATE ON inventory_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_settings_updated_at BEFORE UPDATE ON settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- INITIAL DATA (SEED)
-- ============================================

-- Insert default admin user (password: admin123)
INSERT INTO users (id, username, email, password_hash, full_name, role, is_active)
VALUES (
    uuid_generate_v4(),
    'admin',
    'admin@yayasan.org',
    '$2a$10$rZ7YfJKqPr9r8vL4Wy5Lg.6xQ2PzB5zJYX4VqZKL9AxFc8W7vCj2G', -- admin123
    'System Administrator',
    'admin',
    true
) ON CONFLICT (username) DO NOTHING;

-- Insert default branch
INSERT INTO branches (id, code, name, type, is_active)
VALUES (
    uuid_generate_v4(),
    'HQ',
    'Head Office',
    'office',
    true
) ON CONFLICT (code) DO NOTHING;

-- ============================================
-- VIEWS FOR REPORTING
-- ============================================

-- View: Student with Branch
CREATE OR REPLACE VIEW v_students AS
SELECT 
    s.*,
    b.name as branch_name,
    b.code as branch_code
FROM students s
LEFT JOIN branches b ON s.branch_id = b.id
WHERE s.deleted_at IS NULL;

-- View: Invoice Summary
CREATE OR REPLACE VIEW v_invoice_summary AS
SELECT 
    i.*,
    s.full_name as student_name,
    s.student_number,
    b.name as branch_name,
    (i.total_amount - i.paid_amount) as outstanding_amount
FROM invoices i
JOIN students s ON i.student_id = s.id
JOIN branches b ON i.branch_id = b.id
WHERE i.deleted_at IS NULL;

-- View: Employee with Branch
CREATE OR REPLACE VIEW v_employees AS
SELECT 
    e.*,
    b.name as branch_name,
    b.code as branch_code
FROM employees e
LEFT JOIN branches b ON e.branch_id = b.id
WHERE e.deleted_at IS NULL;

-- View: Asset with Branch
CREATE OR REPLACE VIEW v_assets AS
SELECT 
    a.*,
    b.name as branch_name,
    b.code as branch_code
FROM assets a
LEFT JOIN branches b ON a.branch_id = b.id
WHERE a.deleted_at IS NULL;

-- ============================================
-- COMMENTS
-- ============================================
COMMENT ON TABLE branches IS 'Organization branches (schools, offices, warehouses)';
COMMENT ON TABLE users IS 'System users with role-based access control';
COMMENT ON TABLE students IS 'Student master data';
COMMENT ON TABLE student_parents IS 'Student parent/guardian information';
COMMENT ON TABLE invoices IS 'Student invoices (tuition, fees, etc)';
COMMENT ON TABLE payments IS 'Payment transactions for invoices';
COMMENT ON TABLE accounts IS 'Chart of accounts for accounting';
COMMENT ON TABLE journals IS 'General journal entries';
COMMENT ON TABLE employees IS 'Employee master data';
COMMENT ON TABLE attendance IS 'Daily employee attendance records';
COMMENT ON TABLE payroll IS 'Monthly employee payroll';
COMMENT ON TABLE assets IS 'Fixed assets with depreciation tracking';
COMMENT ON TABLE inventory_items IS 'Inventory item master data';
COMMENT ON TABLE inventory_transactions IS 'Stock in/out transactions';
COMMENT ON TABLE audit_logs IS 'System audit trail';

-- ============================================
-- END OF SCHEMA
-- ============================================
