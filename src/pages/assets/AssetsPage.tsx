import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { Asset, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, Filter, Box } from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function AssetsPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('all');

  const { data, isLoading } = useQuery({
    queryKey: ['assets', page, search, status],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Asset>>>(
        '/assets',
        {
          params: {
            page,
            page_size: 20,
            search,
            status: status !== 'all' ? status : undefined,
          },
        }
      );
      return response.data.data;
    },
  });

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'maintenance': return 'warning';
      case 'disposed': return 'danger';
      case 'inactive': return 'default';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      active: 'Aktif',
      maintenance: 'Maintenance',
      disposed: 'Dihapus',
      inactive: 'Nonaktif',
    };
    return labels[status] || status;
  };

  const columns = [
    {
      key: 'asset_code',
      label: 'Kode Aset',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'name',
      label: 'Nama Aset',
      render: (value: string, row: Asset) => (
        <div>
          <div className="font-medium text-gray-900">{value}</div>
          <div className="text-xs text-gray-500">{row.category_name}</div>
        </div>
      ),
    },
    {
      key: 'acquisition_date',
      label: 'Tanggal Perolehan',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'acquisition_cost',
      label: 'Nilai Perolehan',
      render: (value: number) => (
        <span className="font-medium">{formatCurrency(value)}</span>
      ),
    },
    {
      key: 'book_value',
      label: 'Nilai Buku',
      render: (value: number) => formatCurrency(value),
    },
    {
      key: 'location',
      label: 'Lokasi',
      render: (value: string) => value || '-',
    },
    {
      key: 'status',
      label: 'Status',
      render: (value: string) => (
        <Badge variant={getStatusVariant(value)}>
          {getStatusLabel(value)}
        </Badge>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Manajemen Aset</h1>
          <p className="text-gray-600 mt-1">Kelola aset tetap yayasan</p>
        </div>
        <Button variant="primary" onClick={() => navigate('/assets/new')}>
          <Plus className="w-4 h-4 mr-2" />
          Tambah Aset
        </Button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-600">Total Aset</div>
              <div className="text-2xl font-bold text-gray-900 mt-1">
                {data?.total || 0}
              </div>
            </div>
            <div className="p-3 bg-blue-100 rounded-lg">
              <Box className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </Card>

        <Card>
          <div className="text-sm text-gray-600">Nilai Perolehan</div>
          <div className="text-xl font-bold text-gray-900 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>

        <Card>
          <div className="text-sm text-gray-600">Nilai Buku</div>
          <div className="text-xl font-bold text-gray-900 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>

        <Card>
          <div className="text-sm text-gray-600">Akumulasi Penyusutan</div>
          <div className="text-xl font-bold text-red-600 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Cari kode atau nama aset..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
          
          <select
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
          >
            <option value="all">Semua Status</option>
            <option value="active">Aktif</option>
            <option value="maintenance">Maintenance</option>
            <option value="inactive">Nonaktif</option>
            <option value="disposed">Dihapus</option>
          </select>

          <Button variant="secondary">
            <Filter className="w-4 h-4 mr-2" />
            Filter Lanjut
          </Button>
        </div>
      </Card>

      {/* Table */}
      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada data aset"
          onRowClick={(row) => navigate(`/assets/${row.id}`)}
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
