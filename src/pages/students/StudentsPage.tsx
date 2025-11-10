import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '@/services/api';
import { Student, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, Download, Filter } from 'lucide-react';
import { formatDate } from '@/utils/format';

export default function StudentsPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['students', page, search],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Student>>>(
        '/students',
        {
          params: {
            page,
            page_size: 20,
            search,
          },
        }
      );
      return response.data.data;
    },
  });

  const columns = [
    {
      key: 'registration_number',
      label: 'No. Registrasi',
    },
    {
      key: 'full_name',
      label: 'Nama Lengkap',
      render: (value: string, row: Student) => (
        <div>
          <div className="font-medium text-gray-900">{value}</div>
          {row.nick_name && (
            <div className="text-sm text-gray-500">({row.nick_name})</div>
          )}
        </div>
      ),
    },
    {
      key: 'gender',
      label: 'Jenis Kelamin',
      render: (value: string) => (
        <span className="capitalize">{value === 'male' ? 'L' : 'P'}</span>
      ),
    },
    {
      key: 'class_name',
      label: 'Kelas',
      render: (value: string) => value || '-',
    },
    {
      key: 'branch_name',
      label: 'Cabang',
    },
    {
      key: 'status',
      label: 'Status',
      render: (value: string) => {
        const variants: any = {
          active: 'success',
          inactive: 'default',
          graduated: 'info',
          dropped: 'danger',
        };
        return (
          <Badge variant={variants[value] || 'default'}>
            {value === 'active' ? 'Aktif' : value === 'inactive' ? 'Tidak Aktif' : value}
          </Badge>
        );
      },
    },
    {
      key: 'registration_date',
      label: 'Tgl. Daftar',
      render: (value: string) => formatDate(value),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Data Siswa</h1>
          <p className="text-gray-600 mt-1">Kelola data siswa dan registrasi</p>
        </div>
        <Button variant="primary">
          <Plus className="w-4 h-4 mr-2" />
          Tambah Siswa
        </Button>
      </div>

      {/* Filters */}
      <Card>
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Cari nama, no. registrasi, NISN..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
          <Button variant="secondary">
            <Filter className="w-4 h-4 mr-2" />
            Filter
          </Button>
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
        </div>
      </Card>

      {/* Table */}
      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada data siswa"
        />

        {/* Pagination */}
        {data && data.total > 0 && (
          <div className="mt-4 flex items-center justify-between border-t border-gray-200 pt-4">
            <div className="text-sm text-gray-700">
              Menampilkan{' '}
              <span className="font-medium">
                {(page - 1) * data.page_size + 1}
              </span>{' '}
              -{' '}
              <span className="font-medium">
                {Math.min(page * data.page_size, data.total)}
              </span>{' '}
              dari <span className="font-medium">{data.total}</span> siswa
            </div>
            <div className="flex gap-2">
              <Button
                variant="secondary"
                size="sm"
                disabled={page === 1}
                onClick={() => setPage(page - 1)}
              >
                Previous
              </Button>
              <Button
                variant="secondary"
                size="sm"
                disabled={page >= data.total_pages}
                onClick={() => setPage(page + 1)}
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}
