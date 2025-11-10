import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAuthStore } from './store/auth.store';
import { useEffect } from 'react';

// Pages
import LoginPage from './pages/auth/LoginPage';
import DashboardPage from './pages/dashboard/DashboardPage';
import StudentsPage from './pages/students/StudentsPage';
import StudentDetailPage from './pages/students/StudentDetailPage';
import StudentFormPage from './pages/students/StudentFormPage';
import InvoicesPage from './pages/invoices/InvoicesPage';
import InvoiceDetailPage from './pages/invoices/InvoiceDetailPage';
import InvoiceFormPage from './pages/invoices/InvoiceFormPage';
import PaymentsPage from './pages/payments/PaymentsPage';
import PaymentFormPage from './pages/payments/PaymentFormPage';
import PaymentDetailPage from './pages/payments/PaymentDetailPage';
import EmployeesPage from './pages/employees/EmployeesPage';
import EmployeeDetailPage from './pages/employees/EmployeeDetailPage';
import EmployeeFormPage from './pages/employees/EmployeeFormPage';
import PayrollPage from './pages/employees/PayrollPage';
import AttendancePage from './pages/employees/AttendancePage';
import AccountsPage from './pages/finance/AccountsPage';
import ReportsPage from './pages/finance/ReportsPage';
import JournalsPage from './pages/finance/JournalsPage';
import JournalFormPage from './pages/finance/JournalFormPage';
import AssetsPage from './pages/assets/AssetsPage';
import AssetDetailPage from './pages/assets/AssetDetailPage';
import AssetFormPage from './pages/assets/AssetFormPage';
import InventoryPage from './pages/inventory/InventoryPage';
import InventoryDetailPage from './pages/inventory/InventoryDetailPage';
import InventoryFormPage from './pages/inventory/InventoryFormPage';
import StockInPage from './pages/inventory/StockInPage';
import StockOutPage from './pages/inventory/StockOutPage';
import SettingsPage from './pages/settings/SettingsPage';
import UsersPage from './pages/users/UsersPage';
import UserFormPage from './pages/users/UserFormPage';
import RolesPage from './pages/users/RolesPage';
import MainLayout from './layouts/MainLayout';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
}

function App() {
  const { initAuth } = useAuthStore();

  useEffect(() => {
    initAuth();
  }, [initAuth]);

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          
          <Route
            path="/"
            element={
              <PrivateRoute>
                <MainLayout />
              </PrivateRoute>
            }
          >
            <Route index element={<DashboardPage />} />
            <Route path="dashboard" element={<DashboardPage />} />
            
            {/* Student Routes */}
            <Route path="students" element={<StudentsPage />} />
            <Route path="students/new" element={<StudentFormPage />} />
            <Route path="students/:id" element={<StudentDetailPage />} />
            <Route path="students/:id/edit" element={<StudentFormPage />} />
            
            {/* Invoice Routes */}
            <Route path="invoices" element={<InvoicesPage />} />
            <Route path="invoices/new" element={<InvoiceFormPage />} />
            <Route path="invoices/:id" element={<InvoiceDetailPage />} />
            <Route path="invoices/:id/edit" element={<InvoiceFormPage />} />
            
            {/* Payment Routes */}
            <Route path="payments" element={<PaymentsPage />} />
            <Route path="payments/new" element={<PaymentFormPage />} />
            <Route path="payments/:id" element={<PaymentDetailPage />} />
            <Route path="payments/:id" element={<PaymentDetailPage />} />
            
            {/* Finance Routes */}
            <Route path="finance/accounts" element={<AccountsPage />} />
            <Route path="finance/journals" element={<JournalsPage />} />
            <Route path="finance/journals/new" element={<JournalFormPage />} />
            <Route path="finance/reports" element={<ReportsPage />} />
            <Route path="finance/budgets" element={<div className="p-6"><h1 className="text-2xl font-bold">Budgets - Coming Soon</h1></div>} />
            
            {/* Assets Routes */}
            <Route path="assets" element={<AssetsPage />} />
            <Route path="assets/new" element={<div className="p-6"><h1 className="text-2xl font-bold">Add Asset - Coming Soon</h1></div>} />
            <Route path="assets/:id" element={<div className="p-6"><h1 className="text-2xl font-bold">Asset Detail - Coming Soon</h1></div>} />
            
            {/* Inventory Routes */}
            <Route path="inventory" element={<InventoryPage />} />
            <Route path="inventory/new" element={<div className="p-6"><h1 className="text-2xl font-bold">Add Item - Coming Soon</h1></div>} />
            <Route path="inventory/:id" element={<div className="p-6"><h1 className="text-2xl font-bold">Item Detail - Coming Soon</h1></div>} />
            <Route path="inventory/transactions" element={<div className="p-6"><h1 className="text-2xl font-bold">Stock Transactions - Coming Soon</h1></div>} />
            
            {/* HR Routes */}
            <Route path="employees" element={<EmployeesPage />} />
            <Route path="employees/new" element={<EmployeeFormPage />} />
            <Route path="employees/:id" element={<EmployeeDetailPage />} />
            <Route path="employees/:id/edit" element={<EmployeeFormPage />} />
            <Route path="employees/payroll" element={<PayrollPage />} />
            <Route path="employees/attendance" element={<AttendancePage />} />
            <Route path="employees/leave" element={<div className="p-6"><h1 className="text-2xl font-bold">Leave Management - Coming Soon</h1></div>} />
            
            {/* Assets Routes */}
            <Route path="assets" element={<AssetsPage />} />
            <Route path="assets/new" element={<AssetFormPage />} />
            <Route path="assets/:id" element={<AssetDetailPage />} />
            <Route path="assets/:id/edit" element={<AssetFormPage />} />
            
            {/* Inventory Routes */}
            <Route path="inventory" element={<InventoryPage />} />
            <Route path="inventory/new" element={<InventoryFormPage />} />
            <Route path="inventory/:id" element={<InventoryDetailPage />} />
            <Route path="inventory/:id/edit" element={<InventoryFormPage />} />
            <Route path="inventory/stock-in" element={<StockInPage />} />
            <Route path="inventory/stock-out" element={<StockOutPage />} />
            
            {/* Settings Route */}
            <Route path="settings" element={<SettingsPage />} />
            
            {/* User Management Routes */}
            <Route path="users" element={<UsersPage />} />
            <Route path="users/new" element={<UserFormPage />} />
            <Route path="users/:id/edit" element={<UserFormPage />} />
            <Route path="users/roles" element={<RolesPage />} />
          </Route>
          
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
