import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { Journal, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, FileText, Eye } from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function JournalsPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('all');

  const { data, isLoading } = useQuery({
    queryKey: ['journals', page, search, status],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Journal>>>(
        '/journals',
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
      case 'approved': return 'success';
      case 'pending': return 'warning';
      case 'rejected': return 'danger';
      case 'draft': return 'default';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      draft: 'Draft',
      pending: 'Pending',
      approved: 'Disetujui',
      rejected: 'Ditolak',
    };
    return labels[status] || status;
  };

  const columns = [
    {
      key: 'journal_number',
      label: 'No. Jurnal',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'journal_date',
      label: 'Tanggal',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'description',
      label: 'Deskripsi',
      render: (value: string) => (
        <div className="max-w-md truncate">{value}</div>
      ),
    },
    {
      key: 'total_debit',
      label: 'Total Debit',
      render: (value: number) => (
        <span className="font-medium">{formatCurrency(value)}</span>
      ),
    },
    {
      key: 'total_credit',
      label: 'Total Kredit',
      render: (value: number) => (
        <span className="font-medium">{formatCurrency(value)}</span>
      ),
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
    {
      key: 'id',
      label: 'Aksi',
      render: (value: string) => (
        <Button
          size="sm"
          variant="ghost"
          onClick={() => navigate(`/finance/journals/${value}`)}
        >
          <Eye className="w-4 h-4" />
        </Button>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Jurnal Umum</h1>
          <p className="text-gray-600 mt-1">Kelola pencatatan jurnal akuntansi</p>
        </div>
        <Button variant="primary" onClick={() => navigate('/finance/journals/new')}>
          <Plus className="w-4 h-4 mr-2" />
          Buat Jurnal
        </Button>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <div className="text-sm text-gray-600">Total Jurnal</div>
          <div className="text-2xl font-bold text-gray-900 mt-1">
            {data?.total || 0}
          </div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Draft</div>
          <div className="text-2xl font-bold text-gray-600 mt-1">0</div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Pending</div>
          <div className="text-2xl font-bold text-yellow-600 mt-1">0</div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Disetujui</div>
          <div className="text-2xl font-bold text-green-600 mt-1">0</div>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <div className="flex gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Cari no. jurnal atau deskripsi..."
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
            <option value="draft">Draft</option>
            <option value="pending">Pending</option>
            <option value="approved">Disetujui</option>
            <option value="rejected">Ditolak</option>
          </select>
        </div>
      </Card>

      {/* Table */}
      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada data jurnal"
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
