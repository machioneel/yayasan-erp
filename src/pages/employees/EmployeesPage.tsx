import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { Employee, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, Filter } from 'lucide-react';
import { formatDate } from '@/utils/format';

export default function EmployeesPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['employees', page, search],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Employee>>>(
        '/employees',
        {
          params: { page, page_size: 20, search },
        }
      );
      return response.data.data;
    },
  });

  const columns = [
    {
      key: 'employee_number',
      label: 'NIP',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'full_name',
      label: 'Nama Lengkap',
      render: (value: string, row: Employee) => (
        <div>
          <div className="font-medium">{value}</div>
          {row.is_teacher && (
            <Badge variant="info" className="text-xs mt-1">Guru</Badge>
          )}
        </div>
      ),
    },
    {
      key: 'position',
      label: 'Jabatan',
    },
    {
      key: 'department',
      label: 'Departemen',
      render: (value: string) => value || '-',
    },
    {
      key: 'branch_name',
      label: 'Cabang',
    },
    {
      key: 'employment_type',
      label: 'Tipe',
      render: (value: string) => {
        const labels: Record<string, string> = {
          permanent: 'Tetap',
          contract: 'Kontrak',
          intern: 'Magang',
          freelance: 'Freelance',
        };
        return labels[value] || value;
      },
    },
    {
      key: 'status',
      label: 'Status',
      render: (value: string) => (
        <Badge variant={value === 'active' ? 'success' : 'default'}>
          {value === 'active' ? 'Aktif' : 'Tidak Aktif'}
        </Badge>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Data Karyawan</h1>
          <p className="text-gray-600 mt-1">Kelola data karyawan dan guru</p>
        </div>
        <Button variant="primary" onClick={() => navigate('/employees/new')}>
          <Plus className="w-4 h-4 mr-2" />
          Tambah Karyawan
        </Button>
      </div>

      <Card>
        <div className="flex gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Cari nama, NIP, atau email..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
          <Button variant="secondary">
            <Filter className="w-4 h-4 mr-2" />
            Filter
          </Button>
        </div>
      </Card>

      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada data karyawan"
          onRowClick={(row) => navigate(`/employees/${row.id}`)}
        />

        {data && data.total > 0 && (
          <div className="mt-4 flex items-center justify-between border-t pt-4">
            <div className="text-sm text-gray-700">
              Menampilkan {(page - 1) * data.page_size + 1} - {Math.min(page * data.page_size, data.total)} dari {data.total}
            </div>
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" disabled={page === 1} onClick={() => setPage(page - 1)}>
                Previous
              </Button>
              <Button variant="secondary" size="sm" disabled={page >= data.total_pages} onClick={() => setPage(page + 1)}>
                Next
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}
