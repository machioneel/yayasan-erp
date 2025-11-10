import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { Invoice, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, FileText, Eye, DollarSign } from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function InvoicesPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('all');

  const { data, isLoading } = useQuery({
    queryKey: ['invoices', page, search, status],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Invoice>>>(
        '/invoices',
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
      case 'paid': return 'success';
      case 'partial': return 'warning';
      case 'unpaid': return 'danger';
      case 'cancelled': return 'default';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'paid': return 'Lunas';
      case 'partial': return 'Sebagian';
      case 'unpaid': return 'Belum Bayar';
      case 'cancelled': return 'Dibatalkan';
      default: return status;
    }
  };

  const columns = [
    {
      key: 'invoice_number',
      label: 'No. Invoice',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'student_name',
      label: 'Nama Siswa',
    },
    {
      key: 'invoice_date',
      label: 'Tanggal',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'due_date',
      label: 'Jatuh Tempo',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'total_amount',
      label: 'Total',
      render: (value: number) => (
        <span className="font-medium">{formatCurrency(value)}</span>
      ),
    },
    {
      key: 'paid_amount',
      label: 'Terbayar',
      render: (value: number) => formatCurrency(value),
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
      render: (value: string, row: Invoice) => (
        <div className="flex gap-2">
          <Button
            size="sm"
            variant="ghost"
            onClick={() => navigate(`/invoices/${value}`)}
          >
            <Eye className="w-4 h-4" />
          </Button>
          {row.status !== 'paid' && (
            <Button
              size="sm"
              variant="success"
              onClick={() => navigate(`/payments/new?invoice_id=${value}`)}
            >
              <DollarSign className="w-4 h-4" />
            </Button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Invoice & Tagihan</h1>
          <p className="text-gray-600 mt-1">Kelola invoice dan tagihan siswa</p>
        </div>
        <Button
          variant="primary"
          onClick={() => navigate('/invoices/new')}
        >
          <Plus className="w-4 h-4 mr-2" />
          Buat Invoice
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
                placeholder="Cari no. invoice atau nama siswa..."
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
            <option value="unpaid">Belum Bayar</option>
            <option value="partial">Sebagian</option>
            <option value="paid">Lunas</option>
            <option value="cancelled">Dibatalkan</option>
          </select>
        </div>
      </Card>

      {/* Table */}
      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada invoice"
        />

        {/* Pagination */}
        {data && data.total > 0 && (
          <div className="mt-4 flex items-center justify-between border-t border-gray-200 pt-4">
            <div className="text-sm text-gray-700">
              Menampilkan {(page - 1) * data.page_size + 1} -{' '}
              {Math.min(page * data.page_size, data.total)} dari {data.total} invoice
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
