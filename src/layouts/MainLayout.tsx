import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth.store';
import { authService } from '@/services/auth.service';
import {
  LayoutDashboard,
  Users,
  Wallet,
  UserCog,
  Package,
  Box,
  LogOut,
  Menu,
  X,
  ChevronDown,
  Building2,
  FileText,
  Receipt,
  DollarSign,
  TrendingUp,
  Calendar,
  LucideIcon,
} from 'lucide-react';
import { useState } from 'react';

interface MenuItem {
  name: string;
  icon: LucideIcon;
  path: string;
  children?: { name: string; path: string }[];
}

const menuItems: MenuItem[] = [
  {
    name: 'Dashboard',
    icon: LayoutDashboard,
    path: '/dashboard',
  },
  {
    name: 'Siswa',
    icon: Users,
    path: '/students',
    children: [
      { name: 'Daftar Siswa', path: '/students' },
      { name: 'Pendaftaran Baru', path: '/students/register' },
      { name: 'Tagihan & Invoice', path: '/students/invoices' },
      { name: 'Pembayaran', path: '/students/payments' },
    ],
  },
  {
    name: 'Keuangan',
    icon: Wallet,
    path: '/finance',
    children: [
      { name: 'Chart of Accounts', path: '/finance/accounts' },
      { name: 'Jurnal Umum', path: '/finance/journals' },
      { name: 'Laporan Keuangan', path: '/finance/reports' },
      { name: 'Anggaran', path: '/finance/budgets' },
    ],
  },
  {
    name: 'SDM & Payroll',
    icon: UserCog,
    path: '/employees',
    children: [
      { name: 'Data Karyawan', path: '/employees' },
      { name: 'Absensi', path: '/employees/attendance' },
      { name: 'Gaji/Payroll', path: '/employees/payroll' },
      { name: 'Cuti', path: '/employees/leave' },
    ],
  },
  {
    name: 'Aset',
    icon: Package,
    path: '/assets',
    children: [
      { name: 'Daftar Aset', path: '/assets' },
      { name: 'Depresiasi', path: '/assets/depreciation' },
      { name: 'Maintenance', path: '/assets/maintenance' },
      { name: 'Transfer Aset', path: '/assets/transfer' },
    ],
  },
  {
    name: 'Inventori',
    icon: Box,
    path: '/inventory',
    children: [
      { name: 'Daftar Barang', path: '/inventory/items' },
      { name: 'Stock In/Out', path: '/inventory/transactions' },
      { name: 'Stock Opname', path: '/inventory/opname' },
      { name: 'Laporan Stock', path: '/inventory/reports' },
    ],
  },
];

export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuthStore();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [expandedMenu, setExpandedMenu] = useState<string | null>(null);

  const handleLogout = async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      logout();
      navigate('/login');
    }
  };

  const isActive = (path: string) => {
    return location.pathname === path || location.pathname.startsWith(path + '/');
  };

  const toggleMenu = (name: string) => {
    setExpandedMenu(expandedMenu === name ? null : name);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside
        className={`fixed top-0 left-0 z-40 h-screen transition-transform ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } bg-white border-r border-gray-200 w-64`}
      >
        {/* Logo */}
        <div className="h-16 flex items-center justify-between px-4 border-b border-gray-200">
          <div className="flex items-center gap-2">
            <Building2 className="w-8 h-8 text-blue-600" />
            <span className="font-bold text-lg text-gray-900">Yayasan ERP</span>
          </div>
        </div>

        {/* User Info */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center text-white font-semibold">
              {user?.full_name?.charAt(0) || 'U'}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-semibold text-gray-900 truncate">
                {user?.full_name}
              </p>
              <p className="text-xs text-gray-500 truncate">{user?.role_name}</p>
            </div>
          </div>
        </div>

        {/* Navigation */}
        <nav className="p-4 space-y-1 overflow-y-auto h-[calc(100vh-180px)]">
          {menuItems.map((item) => (
            <div key={item.name}>
              {item.children ? (
                <div>
                  <button
                    onClick={() => toggleMenu(item.name)}
                    className={`w-full flex items-center justify-between px-3 py-2 text-sm font-medium rounded-lg transition-colors ${
                      isActive(item.path)
                        ? 'bg-blue-50 text-blue-700'
                        : 'text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <item.icon className="w-5 h-5" />
                      <span>{item.name}</span>
                    </div>
                    <ChevronDown
                      className={`w-4 h-4 transition-transform ${
                        expandedMenu === item.name ? 'transform rotate-180' : ''
                      }`}
                    />
                  </button>
                  {expandedMenu === item.name && (
                    <div className="ml-4 mt-1 space-y-1">
                      {item.children.map((child) => (
                        <button
                          key={child.path}
                          onClick={() => navigate(child.path)}
                          className={`w-full flex items-center gap-3 px-3 py-2 text-sm rounded-lg transition-colors ${
                            location.pathname === child.path
                              ? 'bg-blue-50 text-blue-700'
                              : 'text-gray-600 hover:bg-gray-100'
                          }`}
                        >
                          <div className="w-1.5 h-1.5 bg-gray-400 rounded-full" />
                          <span>{child.name}</span>
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              ) : (
                <button
                  onClick={() => navigate(item.path)}
                  className={`w-full flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-lg transition-colors ${
                    isActive(item.path)
                      ? 'bg-blue-50 text-blue-700'
                      : 'text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  <span>{item.name}</span>
                </button>
              )}
            </div>
          ))}
        </nav>

        {/* Logout Button */}
        <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 bg-white">
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors"
          >
            <LogOut className="w-5 h-5" />
            <span>Keluar</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <div className={`${sidebarOpen ? 'ml-64' : 'ml-0'} transition-all`}>
        {/* Header */}
        <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
          >
            {sidebarOpen ? (
              <X className="w-6 h-6 text-gray-600" />
            ) : (
              <Menu className="w-6 h-6 text-gray-600" />
            )}
          </button>

          <div className="flex items-center gap-4">
            <div className="text-sm text-gray-600">
              <span className="font-medium">{user?.branch_name || 'Semua Cabang'}</span>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="p-6">
          <Outlet />
        </main>
      </div>

      {/* Mobile Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
    </div>
  );
}
