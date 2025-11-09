// API Response Types
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginationResponse<T = any> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface PaginationParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

// Auth Types
export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  branch_id: string;
  branch_name?: string;
  role_id: string;
  role_name?: string;
  is_active: boolean;
  last_login?: string;
  created_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// Branch Types
export interface Branch {
  id: string;
  code: string;
  name: string;
  address?: string;
  phone?: string;
  email?: string;
  is_active: boolean;
  created_at: string;
}

// Student Types
export interface Student {
  id: string;
  registration_number: string;
  nisn?: string;
  full_name: string;
  nick_name?: string;
  gender: 'male' | 'female';
  birth_place?: string;
  birth_date?: string;
  religion?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  branch_id: string;
  branch_name?: string;
  current_class_id?: string;
  class_name?: string;
  class_level?: string;
  status: string;
  status_name?: string;
  registration_date: string;
  photo_url?: string;
  parents?: ParentSummary[];
  created_at: string;
}

export interface ParentSummary {
  id: string;
  full_name: string;
  relationship: string;
  phone: string;
  is_primary: boolean;
  is_financial: boolean;
}

export interface CreateStudentRequest {
  full_name: string;
  gender: 'male' | 'female';
  birth_date?: string;
  branch_id: string;
  registration_date: string;
  email?: string;
  phone?: string;
  address?: string;
  parents?: ParentRequest[];
}

export interface ParentRequest {
  parent_id?: string;
  full_name?: string;
  nik?: string;
  gender?: 'male' | 'female';
  phone?: string;
  email?: string;
  address?: string;
  occupation?: string;
  relationship: string;
  is_primary: boolean;
  is_financial: boolean;
}

// Invoice & Payment Types
export interface Invoice {
  id: string;
  student_id: string;
  student_name?: string;
  branch_id: string;
  invoice_number: string;
  invoice_date: string;
  due_date: string;
  description: string;
  total_amount: number;
  paid_amount: number;
  status: string;
  items?: InvoiceItem[];
  created_at: string;
}

export interface InvoiceItem {
  id: string;
  description: string;
  quantity: number;
  unit_price: number;
  amount: number;
}

export interface Payment {
  id: string;
  invoice_id: string;
  student_id: string;
  payment_number: string;
  payment_date: string;
  amount: number;
  payment_method: string;
  reference_no?: string;
  notes?: string;
  created_at: string;
}

// Employee Types
export interface Employee {
  id: string;
  employee_number: string;
  nik?: string;
  npwp?: string;
  full_name: string;
  gender: 'male' | 'female';
  birth_place?: string;
  birth_date?: string;
  email: string;
  phone: string;
  address?: string;
  branch_id: string;
  branch_name?: string;
  department?: string;
  position: string;
  employment_type: string;
  join_date: string;
  status: string;
  is_teacher: boolean;
  created_at: string;
}

// Payroll Types
export interface Payroll {
  id: string;
  employee_id: string;
  employee_name?: string;
  payroll_number: string;
  period: string;
  payment_date: string;
  base_salary: number;
  allowances: number;
  overtime: number;
  bonus: number;
  tax_pph21: number;
  bpjs: number;
  loan: number;
  other_deductions: number;
  gross_salary: number;
  total_deductions: number;
  net_salary: number;
  payment_status: string;
  created_at: string;
}

// Account (COA) Types
export interface Account {
  id: string;
  code: string;
  name: string;
  account_type: string;
  normal_balance: string;
  parent_id?: string;
  level: number;
  is_header: boolean;
  is_active: boolean;
  balance?: number;
}

// Journal Types
export interface Journal {
  id: string;
  journal_number: string;
  journal_date: string;
  description: string;
  branch_id: string;
  branch_name?: string;
  total_debit: number;
  total_credit: number;
  status: string;
  created_by_name?: string;
  items?: JournalItem[];
  created_at: string;
}

export interface JournalItem {
  id: string;
  account_id: string;
  account_code?: string;
  account_name?: string;
  description: string;
  debit: number;
  credit: number;
}

// Asset Types
export interface Asset {
  id: string;
  asset_number: string;
  name: string;
  category_id: string;
  category_name?: string;
  branch_id: string;
  branch_name?: string;
  purchase_date: string;
  purchase_price: number;
  depreciation_method?: string;
  useful_life?: number;
  book_value?: number;
  condition: string;
  status: string;
  location?: string;
  created_at: string;
}

// Inventory Types
export interface InventoryItem {
  id: string;
  item_code: string;
  name: string;
  category_id: string;
  category_name?: string;
  branch_id: string;
  unit: string;
  current_stock: number;
  minimum_stock: number;
  purchase_price: number;
  selling_price: number;
  is_active: boolean;
  is_saleable: boolean;
  created_at: string;
}

export interface StockTransaction {
  id: string;
  transaction_number: string;
  item_id: string;
  item_name?: string;
  transaction_type: string;
  transaction_date: string;
  quantity: number;
  unit_price: number;
  total_value: number;
  stock_before: number;
  stock_after: number;
  supplier?: string;
  customer?: string;
  reason?: string;
  created_at: string;
}

// Dashboard Stats
export interface DashboardStats {
  total_students: number;
  total_employees: number;
  total_assets: number;
  monthly_revenue: number;
  monthly_expenses: number;
  pending_invoices: number;
  overdue_invoices: number;
  low_stock_items: number;
}

// Common Types
export interface SelectOption {
  value: string;
  label: string;
}

export interface FilterOptions {
  branches?: SelectOption[];
  statuses?: SelectOption[];
  categories?: SelectOption[];
}
