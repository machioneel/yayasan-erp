-- ============================================
-- YAYASAN ERP DATABASE SCHEMA
-- PostgreSQL 15+
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- CORE TABLES
-- ============================================

-- Organizations/Branches
CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'headquarters', 'school', 'clinic', 'mosque', 'office'
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    avatar_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    is_super_admin BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP,
    password_changed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    module VARCHAR(50) NOT NULL, -- 'finance', 'inventory', 'sales', etc.
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Role Permissions (Many-to-Many)
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- User Roles (Many-to-Many)
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(id),
    PRIMARY KEY (user_id, role_id, branch_id)
);

-- Audit Log
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    branch_id UUID REFERENCES branches(id),
    action VARCHAR(50) NOT NULL, -- 'create', 'update', 'delete', 'view'
    entity_type VARCHAR(100) NOT NULL, -- 'journal', 'invoice', 'product', etc.
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- FINANCE & ACCOUNTING MODULE
-- ============================================

-- Chart of Accounts
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- 'asset', 'liability', 'equity', 'revenue', 'expense'
    category VARCHAR(100), -- 'current_asset', 'fixed_asset', 'operating_expense', etc.
    parent_id UUID REFERENCES accounts(id),
    level INTEGER DEFAULT 0,
    is_header BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    balance_type VARCHAR(20) NOT NULL, -- 'debit', 'credit'
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Fiscal Years
CREATE TABLE fiscal_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, start_date)
);

-- Journals
CREATE TABLE journals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    fiscal_year_id UUID REFERENCES fiscal_years(id),
    journal_number VARCHAR(50) UNIQUE NOT NULL,
    journal_date DATE NOT NULL,
    reference_type VARCHAR(50), -- 'invoice', 'payment', 'adjustment', etc.
    reference_id UUID,
    reference_number VARCHAR(100),
    description TEXT,
    total_debit DECIMAL(20,2) DEFAULT 0,
    total_credit DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(30) DEFAULT 'draft', -- 'draft', 'posted', 'cancelled'
    posted_at TIMESTAMP,
    posted_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- Journal Entries
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journal_id UUID REFERENCES journals(id) ON DELETE CASCADE,
    account_id UUID REFERENCES accounts(id),
    description TEXT,
    debit DECIMAL(20,2) DEFAULT 0,
    credit DECIMAL(20,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Funds (for Non-Profit Fund Accounting)
CREATE TABLE funds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    fund_type VARCHAR(50) NOT NULL, -- 'unrestricted', 'temporarily_restricted', 'permanently_restricted'
    purpose TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Programs (for program-based accounting)
CREATE TABLE programs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    program_type VARCHAR(50), -- 'education', 'social', 'religious'
    start_date DATE,
    end_date DATE,
    budget DECIMAL(20,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Budgets
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    fiscal_year_id UUID REFERENCES fiscal_years(id),
    program_id UUID REFERENCES programs(id),
    fund_id UUID REFERENCES funds(id),
    account_id UUID REFERENCES accounts(id),
    period_type VARCHAR(20) NOT NULL, -- 'yearly', 'quarterly', 'monthly'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    budget_amount DECIMAL(20,2) NOT NULL,
    actual_amount DECIMAL(20,2) DEFAULT 0,
    variance DECIMAL(20,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Donations (Zakat, Infaq, Sedekah, etc.)
CREATE TABLE donations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    donor_id UUID, -- references donors in CRM
    donation_number VARCHAR(50) UNIQUE NOT NULL,
    donation_date DATE NOT NULL,
    donation_type VARCHAR(50) NOT NULL, -- 'zakat', 'infaq', 'sedekah', 'wakaf', 'general'
    fund_id UUID REFERENCES funds(id),
    program_id UUID REFERENCES programs(id),
    amount DECIMAL(20,2) NOT NULL,
    payment_method VARCHAR(50), -- 'cash', 'transfer', 'check', 'online'
    receipt_number VARCHAR(100),
    receipt_issued BOOLEAN DEFAULT false,
    is_tax_deductible BOOLEAN DEFAULT true,
    notes TEXT,
    status VARCHAR(30) DEFAULT 'received', -- 'received', 'acknowledged', 'receipted'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- ============================================
-- INVENTORY MANAGEMENT MODULE
-- ============================================

-- Product Categories
CREATE TABLE product_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES product_categories(id),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Units of Measure
CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20),
    unit_type VARCHAR(50), -- 'length', 'weight', 'volume', 'quantity', etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    category_id UUID REFERENCES product_categories(id),
    sku VARCHAR(100) UNIQUE NOT NULL,
    barcode VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    product_type VARCHAR(50) DEFAULT 'goods', -- 'goods', 'service', 'donation_item'
    unit_id UUID REFERENCES units(id),
    cost_price DECIMAL(20,2),
    selling_price DECIMAL(20,2),
    reorder_level INTEGER DEFAULT 0,
    reorder_quantity INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Warehouses
CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    location TEXT,
    warehouse_type VARCHAR(50), -- 'main', 'branch', 'temporary'
    manager_id UUID REFERENCES users(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Stock
CREATE TABLE stock (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id),
    warehouse_id UUID REFERENCES warehouses(id),
    quantity DECIMAL(20,4) DEFAULT 0,
    reserved_quantity DECIMAL(20,4) DEFAULT 0,
    available_quantity DECIMAL(20,4) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, warehouse_id)
);

-- Stock Movements
CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    movement_number VARCHAR(50) UNIQUE NOT NULL,
    movement_date DATE NOT NULL,
    movement_type VARCHAR(50) NOT NULL, -- 'in', 'out', 'transfer', 'adjustment'
    reference_type VARCHAR(50), -- 'purchase', 'sale', 'donation', 'adjustment'
    reference_id UUID,
    reference_number VARCHAR(100),
    product_id UUID REFERENCES products(id),
    from_warehouse_id UUID REFERENCES warehouses(id),
    to_warehouse_id UUID REFERENCES warehouses(id),
    quantity DECIMAL(20,4) NOT NULL,
    unit_cost DECIMAL(20,2),
    total_cost DECIMAL(20,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- ============================================
-- SALES & CRM MODULE (Donor Management)
-- ============================================

-- Donors (Customers in typical CRM)
CREATE TABLE donors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    donor_number VARCHAR(50) UNIQUE NOT NULL,
    donor_type VARCHAR(50) NOT NULL, -- 'individual', 'corporate', 'foundation', 'government'
    salutation VARCHAR(20),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    organization_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Indonesia',
    tax_id VARCHAR(100), -- NPWP
    donor_category VARCHAR(50), -- 'major', 'regular', 'occasional', 'lapsed'
    total_donations DECIMAL(20,2) DEFAULT 0,
    last_donation_date DATE,
    communication_preference VARCHAR(50), -- 'email', 'phone', 'whatsapp', 'mail'
    tags TEXT[], -- Array of tags
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Grants (for grant management)
CREATE TABLE grants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    grant_number VARCHAR(50) UNIQUE NOT NULL,
    donor_id UUID REFERENCES donors(id),
    program_id UUID REFERENCES programs(id),
    fund_id UUID REFERENCES funds(id),
    grant_name VARCHAR(255) NOT NULL,
    grant_type VARCHAR(50), -- 'restricted', 'unrestricted', 'challenge'
    amount DECIMAL(20,2) NOT NULL,
    start_date DATE,
    end_date DATE,
    payment_schedule VARCHAR(50), -- 'lump_sum', 'quarterly', 'milestone'
    amount_received DECIMAL(20,2) DEFAULT 0,
    amount_spent DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'active', 'completed', 'cancelled'
    terms_conditions TEXT,
    reporting_requirements TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Donor Interactions
CREATE TABLE donor_interactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    donor_id UUID REFERENCES donors(id),
    interaction_date TIMESTAMP NOT NULL,
    interaction_type VARCHAR(50) NOT NULL, -- 'call', 'email', 'meeting', 'event', 'thank_you'
    subject VARCHAR(255),
    description TEXT,
    outcome VARCHAR(255),
    follow_up_date DATE,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- PURCHASE MANAGEMENT MODULE
-- ============================================

-- Vendors/Suppliers
CREATE TABLE vendors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    vendor_number VARCHAR(50) UNIQUE NOT NULL,
    vendor_name VARCHAR(255) NOT NULL,
    vendor_type VARCHAR(50), -- 'supplier', 'contractor', 'service_provider'
    contact_person VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    tax_id VARCHAR(100), -- NPWP
    payment_terms VARCHAR(100), -- 'net_30', 'net_60', 'cod', etc.
    bank_name VARCHAR(255),
    bank_account_number VARCHAR(100),
    bank_account_name VARCHAR(255),
    rating INTEGER, -- 1-5 stars
    is_active BOOLEAN DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Purchase Requisitions
CREATE TABLE purchase_requisitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    pr_number VARCHAR(50) UNIQUE NOT NULL,
    pr_date DATE NOT NULL,
    requested_by UUID REFERENCES users(id),
    department VARCHAR(100),
    program_id UUID REFERENCES programs(id),
    required_date DATE,
    priority VARCHAR(20) DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'
    justification TEXT,
    total_amount DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'pending_approval', 'approved', 'rejected', 'po_created'
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase Requisition Items
CREATE TABLE purchase_requisition_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pr_id UUID REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    description TEXT,
    quantity DECIMAL(20,4) NOT NULL,
    unit_id UUID REFERENCES units(id),
    estimated_price DECIMAL(20,2),
    total_estimated DECIMAL(20,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase Orders
CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    pr_id UUID REFERENCES purchase_requisitions(id),
    vendor_id UUID REFERENCES vendors(id),
    po_number VARCHAR(50) UNIQUE NOT NULL,
    po_date DATE NOT NULL,
    expected_delivery_date DATE,
    payment_terms VARCHAR(100),
    subtotal DECIMAL(20,2) DEFAULT 0,
    tax_percentage DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(20,2) DEFAULT 0,
    shipping_cost DECIMAL(20,2) DEFAULT 0,
    other_costs DECIMAL(20,2) DEFAULT 0,
    discount_amount DECIMAL(20,2) DEFAULT 0,
    total_amount DECIMAL(20,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'sent', 'confirmed', 'partially_received', 'received', 'cancelled'
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Purchase Order Items
CREATE TABLE purchase_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    po_id UUID REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    description TEXT,
    quantity DECIMAL(20,4) NOT NULL,
    unit_id UUID REFERENCES units(id),
    unit_price DECIMAL(20,2) NOT NULL,
    discount_percentage DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(20,2) DEFAULT 0,
    tax_percentage DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(20,2) DEFAULT 0,
    total_amount DECIMAL(20,2) NOT NULL,
    received_quantity DECIMAL(20,4) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Goods Receipts
CREATE TABLE goods_receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    po_id UUID REFERENCES purchase_orders(id),
    gr_number VARCHAR(50) UNIQUE NOT NULL,
    gr_date DATE NOT NULL,
    warehouse_id UUID REFERENCES warehouses(id),
    received_by UUID REFERENCES users(id),
    delivery_note_number VARCHAR(100),
    notes TEXT,
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'received', 'inspected', 'accepted', 'rejected'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Goods Receipt Items
CREATE TABLE goods_receipt_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    gr_id UUID REFERENCES goods_receipts(id) ON DELETE CASCADE,
    po_item_id UUID REFERENCES purchase_order_items(id),
    product_id UUID REFERENCES products(id),
    ordered_quantity DECIMAL(20,4),
    received_quantity DECIMAL(20,4) NOT NULL,
    rejected_quantity DECIMAL(20,4) DEFAULT 0,
    accepted_quantity DECIMAL(20,4) DEFAULT 0,
    quality_status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'passed', 'failed'
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- ASSET MANAGEMENT MODULE
-- ============================================

-- Asset Categories
CREATE TABLE asset_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES asset_categories(id),
    depreciation_method VARCHAR(50), -- 'straight_line', 'declining_balance', 'none'
    depreciation_rate DECIMAL(5,2), -- percentage
    useful_life_years INTEGER,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Assets
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    category_id UUID REFERENCES asset_categories(id),
    asset_number VARCHAR(50) UNIQUE NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    description TEXT,
    acquisition_date DATE NOT NULL,
    acquisition_cost DECIMAL(20,2) NOT NULL,
    residual_value DECIMAL(20,2) DEFAULT 0,
    useful_life_years INTEGER,
    depreciation_method VARCHAR(50),
    current_value DECIMAL(20,2),
    accumulated_depreciation DECIMAL(20,2) DEFAULT 0,
    location VARCHAR(255),
    custodian_id UUID REFERENCES users(id),
    serial_number VARCHAR(100),
    manufacturer VARCHAR(255),
    model VARCHAR(255),
    warranty_expiry_date DATE,
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'in_maintenance', 'disposed', 'retired'
    condition VARCHAR(50) DEFAULT 'good', -- 'excellent', 'good', 'fair', 'poor'
    po_id UUID REFERENCES purchase_orders(id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Asset Transfers
CREATE TABLE asset_transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID REFERENCES assets(id),
    from_branch_id UUID REFERENCES branches(id),
    to_branch_id UUID REFERENCES branches(id),
    from_custodian_id UUID REFERENCES users(id),
    to_custodian_id UUID REFERENCES users(id),
    from_location VARCHAR(255),
    to_location VARCHAR(255),
    transfer_date DATE NOT NULL,
    reason TEXT,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Asset Maintenance
CREATE TABLE asset_maintenance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID REFERENCES assets(id),
    maintenance_number VARCHAR(50) UNIQUE NOT NULL,
    maintenance_type VARCHAR(50) NOT NULL, -- 'preventive', 'corrective', 'inspection'
    scheduled_date DATE,
    completed_date DATE,
    vendor_id UUID REFERENCES vendors(id),
    cost DECIMAL(20,2),
    description TEXT,
    work_performed TEXT,
    next_maintenance_date DATE,
    status VARCHAR(50) DEFAULT 'scheduled', -- 'scheduled', 'in_progress', 'completed', 'cancelled'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Asset Depreciation
CREATE TABLE asset_depreciation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID REFERENCES assets(id),
    fiscal_year_id UUID REFERENCES fiscal_years(id),
    depreciation_date DATE NOT NULL,
    depreciation_amount DECIMAL(20,2) NOT NULL,
    accumulated_depreciation DECIMAL(20,2) NOT NULL,
    book_value DECIMAL(20,2) NOT NULL,
    journal_id UUID REFERENCES journals(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- ============================================
-- REPORTING & SETTINGS
-- ============================================

-- Report Templates
CREATE TABLE report_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    module VARCHAR(50) NOT NULL,
    report_type VARCHAR(50), -- 'financial', 'operational', 'analytical'
    description TEXT,
    query_template TEXT,
    parameters JSONB,
    is_system_report BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- System Settings
CREATE TABLE system_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id),
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT,
    setting_type VARCHAR(50), -- 'string', 'number', 'boolean', 'json'
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES users(id),
    UNIQUE(branch_id, setting_key)
);

-- File Attachments (Generic)
CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    uploaded_by UUID REFERENCES users(id),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Users
CREATE INDEX idx_users_branch_id ON users(branch_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Audit Logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Accounts
CREATE INDEX idx_accounts_branch_id ON accounts(branch_id);
CREATE INDEX idx_accounts_code ON accounts(code);
CREATE INDEX idx_accounts_parent_id ON accounts(parent_id);

-- Journals
CREATE INDEX idx_journals_branch_id ON journals(branch_id);
CREATE INDEX idx_journals_journal_date ON journals(journal_date);
CREATE INDEX idx_journals_status ON journals(status);
CREATE INDEX idx_journals_fiscal_year_id ON journals(fiscal_year_id);

-- Journal Entries
CREATE INDEX idx_journal_entries_journal_id ON journal_entries(journal_id);
CREATE INDEX idx_journal_entries_account_id ON journal_entries(account_id);

-- Donations
CREATE INDEX idx_donations_branch_id ON donations(branch_id);
CREATE INDEX idx_donations_donor_id ON donations(donor_id);
CREATE INDEX idx_donations_donation_date ON donations(donation_date);
CREATE INDEX idx_donations_donation_type ON donations(donation_type);

-- Products
CREATE INDEX idx_products_branch_id ON products(branch_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_barcode ON products(barcode);

-- Stock
CREATE INDEX idx_stock_product_id ON stock(product_id);
CREATE INDEX idx_stock_warehouse_id ON stock(warehouse_id);

-- Stock Movements
CREATE INDEX idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX idx_stock_movements_movement_date ON stock_movements(movement_date);
CREATE INDEX idx_stock_movements_branch_id ON stock_movements(branch_id);

-- Donors
CREATE INDEX idx_donors_branch_id ON donors(branch_id);
CREATE INDEX idx_donors_email ON donors(email);
CREATE INDEX idx_donors_donor_type ON donors(donor_type);

-- Purchase Orders
CREATE INDEX idx_po_branch_id ON purchase_orders(branch_id);
CREATE INDEX idx_po_vendor_id ON purchase_orders(vendor_id);
CREATE INDEX idx_po_status ON purchase_orders(status);
CREATE INDEX idx_po_date ON purchase_orders(po_date);

-- Assets
CREATE INDEX idx_assets_branch_id ON assets(branch_id);
CREATE INDEX idx_assets_category_id ON assets(category_id);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_custodian_id ON assets(custodian_id);

-- ============================================
-- FUNCTIONS AND TRIGGERS
-- ============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply update_updated_at trigger to relevant tables
CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_journals_updated_at BEFORE UPDATE ON journals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_donors_updated_at BEFORE UPDATE ON donors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vendors_updated_at BEFORE UPDATE ON vendors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON assets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update stock available quantity
CREATE OR REPLACE FUNCTION update_stock_available_quantity()
RETURNS TRIGGER AS $$
BEGIN
    NEW.available_quantity = NEW.quantity - NEW.reserved_quantity;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stock_available BEFORE INSERT OR UPDATE ON stock
    FOR EACH ROW EXECUTE FUNCTION update_stock_available_quantity();

-- ============================================
-- INITIAL DATA / SEED DATA
-- ============================================

-- Insert default roles
INSERT INTO roles (code, name, description, is_system_role) VALUES
('super_admin', 'Super Administrator', 'Full system access', true),
('branch_admin', 'Branch Administrator', 'Full access to branch', true),
('finance_manager', 'Finance Manager', 'Finance & accounting access', true),
('inventory_manager', 'Inventory Manager', 'Inventory & purchase access', true),
('donor_manager', 'Donor/Sales Manager', 'CRM & donor management access', true),
('asset_manager', 'Asset Manager', 'Asset management access', true),
('staff', 'Staff', 'Limited access', true),
('viewer', 'Viewer', 'View-only access', true);

-- Insert default permissions (examples)
INSERT INTO permissions (code, name, module) VALUES
-- Finance
('finance.view', 'View Finance Data', 'finance'),
('finance.create', 'Create Finance Transactions', 'finance'),
('finance.update', 'Update Finance Transactions', 'finance'),
('finance.delete', 'Delete Finance Transactions', 'finance'),
('finance.approve', 'Approve Finance Transactions', 'finance'),
('finance.reports', 'View Finance Reports', 'finance'),

-- Inventory
('inventory.view', 'View Inventory', 'inventory'),
('inventory.create', 'Create Inventory Items', 'inventory'),
('inventory.update', 'Update Inventory Items', 'inventory'),
('inventory.delete', 'Delete Inventory Items', 'inventory'),
('inventory.adjust', 'Adjust Stock', 'inventory'),
('inventory.reports', 'View Inventory Reports', 'inventory'),

-- Sales/CRM
('crm.view', 'View Donors/CRM', 'crm'),
('crm.create', 'Create Donors/Contacts', 'crm'),
('crm.update', 'Update Donors/Contacts', 'crm'),
('crm.delete', 'Delete Donors/Contacts', 'crm'),
('crm.reports', 'View CRM Reports', 'crm'),

-- Purchase
('purchase.view', 'View Purchase Orders', 'purchase'),
('purchase.create', 'Create Purchase Orders', 'purchase'),
('purchase.update', 'Update Purchase Orders', 'purchase'),
('purchase.delete', 'Delete Purchase Orders', 'purchase'),
('purchase.approve', 'Approve Purchase Orders', 'purchase'),

-- Assets
('assets.view', 'View Assets', 'assets'),
('assets.create', 'Create Assets', 'assets'),
('assets.update', 'Update Assets', 'assets'),
('assets.delete', 'Delete Assets', 'assets'),
('assets.transfer', 'Transfer Assets', 'assets'),

-- System
('system.settings', 'Manage System Settings', 'system'),
('system.users', 'Manage Users', 'system'),
('system.branches', 'Manage Branches', 'system');

-- Insert default units
INSERT INTO units (code, name, symbol, unit_type) VALUES
('pcs', 'Pieces', 'pcs', 'quantity'),
('unit', 'Unit', 'unit', 'quantity'),
('box', 'Box', 'box', 'quantity'),
('kg', 'Kilogram', 'kg', 'weight'),
('g', 'Gram', 'g', 'weight'),
('l', 'Liter', 'l', 'volume'),
('ml', 'Milliliter', 'ml', 'volume'),
('m', 'Meter', 'm', 'length'),
('cm', 'Centimeter', 'cm', 'length'),
('set', 'Set', 'set', 'quantity'),
('dozen', 'Dozen', 'dz', 'quantity');

-- ============================================
-- VIEWS FOR REPORTING
-- ============================================

-- View: Account Balances
CREATE OR REPLACE VIEW v_account_balances AS
SELECT 
    a.id as account_id,
    a.branch_id,
    a.code as account_code,
    a.name as account_name,
    a.account_type,
    COALESCE(SUM(je.debit), 0) as total_debit,
    COALESCE(SUM(je.credit), 0) as total_credit,
    CASE 
        WHEN a.balance_type = 'debit' THEN COALESCE(SUM(je.debit - je.credit), 0)
        ELSE COALESCE(SUM(je.credit - je.debit), 0)
    END as balance
FROM accounts a
LEFT JOIN journal_entries je ON a.id = je.account_id
LEFT JOIN journals j ON je.journal_id = j.id AND j.status = 'posted'
WHERE a.is_active = true AND a.is_header = false
GROUP BY a.id, a.branch_id, a.code, a.name, a.account_type, a.balance_type;

-- View: Stock Levels
CREATE OR REPLACE VIEW v_stock_levels AS
SELECT 
    p.id as product_id,
    p.branch_id,
    p.sku,
    p.name as product_name,
    w.id as warehouse_id,
    w.name as warehouse_name,
    s.quantity,
    s.reserved_quantity,
    s.available_quantity,
    p.reorder_level,
    CASE 
        WHEN s.available_quantity <= p.reorder_level THEN 'low_stock'
        WHEN s.available_quantity = 0 THEN 'out_of_stock'
        ELSE 'in_stock'
    END as stock_status
FROM products p
JOIN stock s ON p.id = s.product_id
JOIN warehouses w ON s.warehouse_id = w.id
WHERE p.is_active = true;

-- ============================================
-- COMMENTS ON TABLES
-- ============================================

COMMENT ON TABLE branches IS 'Multi-branch organization structure';
COMMENT ON TABLE users IS 'System users with authentication';
COMMENT ON TABLE accounts IS 'Chart of accounts for financial tracking';
COMMENT ON TABLE journals IS 'Financial journal entries';
COMMENT ON TABLE funds IS 'Non-profit fund accounting (restricted/unrestricted)';
COMMENT ON TABLE programs IS 'Programs for program-based accounting';
COMMENT ON TABLE donations IS 'Donations, Zakat, Infaq, Sedekah tracking';
COMMENT ON TABLE products IS 'Inventory products/items';
COMMENT ON TABLE donors IS 'Donor/customer management (CRM)';
COMMENT ON TABLE grants IS 'Grant management for restricted funds';
COMMENT ON TABLE vendors IS 'Suppliers and vendors';
COMMENT ON TABLE assets IS 'Fixed asset registry';

-- ============================================
-- END OF SCHEMA
-- ============================================
