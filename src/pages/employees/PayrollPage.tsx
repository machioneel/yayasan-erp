import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Input } from '@/components/common/Input';
import { Employee, ApiResponse } from '@/types';
import { 
  DollarSign, 
  Download, 
  Calendar,
  CheckCircle,
  Clock,
  Filter,
  FileText
} from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

interface PayrollItem {
  employee_id: string;
  employee_name: string;
  position: string;
  base_salary: number;
  allowances: number;
  deductions: number;
  net_salary: number;
  status: 'pending' | 'approved' | 'paid';
}

export default function PayrollPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [selectedMonth, setSelectedMonth] = useState(new Date().toISOString().slice(0, 7));
  const [selectedEmployees, setSelectedEmployees] = useState<Set<string>>(new Set());

  const { data: employees } = useQuery({
    queryKey: ['employees'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Employee[]>>('/employees');
      return response.data.data;
    },
  });

  // Mock payroll data - in real app, fetch from API
  const payrollItems: PayrollItem[] = employees?.map(emp => ({
    employee_id: emp.id,
    employee_name: emp.full_name,
    position: emp.position,
    base_salary: emp.salary || 0,
    allowances: Math.floor((emp.salary || 0) * 0.15), // 15% transport + meals
    deductions: Math.floor((emp.salary || 0) * 0.08), // 8% tax + BPJS
    net_salary: (emp.salary || 0) + Math.floor((emp.salary || 0) * 0.15) - Math.floor((emp.salary || 0) * 0.08),
    status: 'pending',
  })) || [];

  const processMutation = useMutation({
    mutationFn: async (data: { month: string; employee_ids: string[] }) => {
      return api.post('/payroll/process', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payroll'] });
      alert('Payroll berhasil diproses!');
      setSelectedEmployees(new Set());
    },
  });

  const handleSelectAll = () => {
    if (selectedEmployees.size === payrollItems.length) {
      setSelectedEmployees(new Set());
    } else {
      setSelectedEmployees(new Set(payrollItems.map(item => item.employee_id)));
    }
  };

  const handleSelectEmployee = (employeeId: string) => {
    const newSelected = new Set(selectedEmployees);
    if (newSelected.has(employeeId)) {
      newSelected.delete(employeeId);
    } else {
      newSelected.add(employeeId);
    }
    setSelectedEmployees(newSelected);
  };

  const handleProcessPayroll = () => {
    if (selectedEmployees.size === 0) {
      alert('Pilih minimal 1 karyawan');
      return;
    }

    if (confirm(`Proses payroll untuk ${selectedEmployees.size} karyawan?`)) {
      processMutation.mutate({
        month: selectedMonth,
        employee_ids: Array.from(selectedEmployees),
      });
    }
  };

  const selectedItems = payrollItems.filter(item => selectedEmployees.has(item.employee_id));
  const totalBaseSalary = selectedItems.reduce((sum, item) => sum + item.base_salary, 0);
  const totalAllowances = selectedItems.reduce((sum, item) => sum + item.allowances, 0);
  const totalDeductions = selectedItems.reduce((sum, item) => sum + item.deductions, 0);
  const totalNetSalary = selectedItems.reduce((sum, item) => sum + item.net_salary, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Payroll Processing</h1>
          <p className="text-gray-600">Generate dan proses gaji karyawan</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            Export Excel
          </Button>
          <Button 
            variant="primary"
            onClick={handleProcessPayroll}
            disabled={selectedEmployees.size === 0 || processMutation.isPending}
            loading={processMutation.isPending}
          >
            <CheckCircle className="w-4 h-4 mr-2" />
            Process Payroll ({selectedEmployees.size})
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              <Calendar className="w-4 h-4 inline mr-2" />
              Periode
            </label>
            <input
              type="month"
              value={selectedMonth}
              onChange={(e) => setSelectedMonth(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              <Filter className="w-4 h-4 inline mr-2" />
              Status
            </label>
            <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
              <option value="all">Semua Status</option>
              <option value="pending">Pending</option>
              <option value="approved">Approved</option>
              <option value="paid">Paid</option>
            </select>
          </div>

          <div className="flex items-end">
            <Button variant="secondary">
              <Filter className="w-4 h-4 mr-2" />
              Filter
            </Button>
          </div>
        </div>
      </Card>

      {/* Summary Cards */}
      {selectedEmployees.size > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <div className="text-sm text-gray-600 mb-1">Karyawan Dipilih</div>
            <div className="text-2xl font-bold text-blue-600">{selectedEmployees.size}</div>
          </Card>
          <Card>
            <div className="text-sm text-gray-600 mb-1">Total Gaji Pokok</div>
            <div className="text-2xl font-bold text-gray-900">{formatCurrency(totalBaseSalary)}</div>
          </Card>
          <Card>
            <div className="text-sm text-gray-600 mb-1">Total Potongan</div>
            <div className="text-2xl font-bold text-red-600">{formatCurrency(totalDeductions)}</div>
          </Card>
          <Card>
            <div className="text-sm text-gray-600 mb-1">Total Nett</div>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(totalNetSalary)}</div>
          </Card>
        </div>
      )}

      {/* Payroll Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-3 text-left">
                  <input
                    type="checkbox"
                    checked={selectedEmployees.size === payrollItems.length}
                    onChange={handleSelectAll}
                    className="rounded border-gray-300"
                  />
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Karyawan
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Jabatan
                </th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase">
                  Gaji Pokok
                </th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase">
                  Tunjangan
                </th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase">
                  Potongan
                </th>
                <th className="px-4 py-3 text-right text-xs font-medium text-gray-700 uppercase">
                  Gaji Nett
                </th>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase">
                  Status
                </th>
                <th className="px-4 py-3 text-center text-xs font-medium text-gray-700 uppercase">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {payrollItems.map((item) => (
                <tr 
                  key={item.employee_id}
                  className={`hover:bg-gray-50 ${
                    selectedEmployees.has(item.employee_id) ? 'bg-blue-50' : ''
                  }`}
                >
                  <td className="px-4 py-3">
                    <input
                      type="checkbox"
                      checked={selectedEmployees.has(item.employee_id)}
                      onChange={() => handleSelectEmployee(item.employee_id)}
                      className="rounded border-gray-300"
                    />
                  </td>
                  <td className="px-4 py-3">
                    <div className="font-medium text-gray-900">{item.employee_name}</div>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-600">
                    {item.position}
                  </td>
                  <td className="px-4 py-3 text-right font-medium">
                    {formatCurrency(item.base_salary)}
                  </td>
                  <td className="px-4 py-3 text-right text-green-600">
                    +{formatCurrency(item.allowances)}
                  </td>
                  <td className="px-4 py-3 text-right text-red-600">
                    -{formatCurrency(item.deductions)}
                  </td>
                  <td className="px-4 py-3 text-right font-bold text-gray-900">
                    {formatCurrency(item.net_salary)}
                  </td>
                  <td className="px-4 py-3 text-center">
                    <Badge variant={
                      item.status === 'paid' ? 'success' :
                      item.status === 'approved' ? 'info' :
                      'warning'
                    }>
                      {item.status === 'paid' ? 'Paid' :
                       item.status === 'approved' ? 'Approved' :
                       'Pending'}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <Button variant="ghost" size="sm">
                      <FileText className="w-4 h-4" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
            <tfoot className="bg-gray-100 font-bold border-t-2">
              <tr>
                <td colSpan={3} className="px-4 py-3 text-right">
                  Total:
                </td>
                <td className="px-4 py-3 text-right">
                  {formatCurrency(payrollItems.reduce((sum, item) => sum + item.base_salary, 0))}
                </td>
                <td className="px-4 py-3 text-right text-green-600">
                  +{formatCurrency(payrollItems.reduce((sum, item) => sum + item.allowances, 0))}
                </td>
                <td className="px-4 py-3 text-right text-red-600">
                  -{formatCurrency(payrollItems.reduce((sum, item) => sum + item.deductions, 0))}
                </td>
                <td className="px-4 py-3 text-right text-gray-900">
                  {formatCurrency(payrollItems.reduce((sum, item) => sum + item.net_salary, 0))}
                </td>
                <td colSpan={2}></td>
              </tr>
            </tfoot>
          </table>
        </div>

        {payrollItems.length === 0 && (
          <div className="text-center py-12">
            <DollarSign className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500">Tidak ada data karyawan</p>
          </div>
        )}
      </Card>

      {/* Instructions */}
      <Card title="Panduan Payroll">
        <div className="text-sm text-gray-600 space-y-2">
          <p><strong>Langkah-langkah:</strong></p>
          <ol className="list-decimal list-inside space-y-1 ml-2">
            <li>Pilih periode bulan yang akan diproses</li>
            <li>Centang karyawan yang akan di-generate payroll-nya</li>
            <li>Review total gaji pokok, tunjangan, dan potongan</li>
            <li>Klik "Process Payroll" untuk generate slip gaji</li>
            <li>Slip gaji dapat di-download dan dikirim ke karyawan</li>
          </ol>
          <p className="mt-3"><strong>Komponen Gaji:</strong></p>
          <ul className="list-disc list-inside space-y-1 ml-2">
            <li>Gaji Pokok: Sesuai data karyawan</li>
            <li>Tunjangan: Transport (10%) + Makan (5%)</li>
            <li>Potongan: PPh 21 (5%) + BPJS (3%)</li>
          </ul>
        </div>
      </Card>
    </div>
  );
}
