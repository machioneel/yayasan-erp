import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { Payment, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, Download, Eye, Printer } from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function PaymentsPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['payments', page, search],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<Payment>>>(
        '/payments',
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

  const getMethodBadge = (method: string) => {
    switch (method) {
      case 'cash': return { variant: 'success' as const, label: 'Tunai' };
      case 'transfer': return { variant: 'info' as const, label: 'Transfer' };
      case 'card': return { variant: 'warning' as const, label: 'Kartu' };
      case 'ewallet': return { variant: 'info' as const, label: 'E-Wallet' };
      default: return { variant: 'default' as const, label: method };
    }
  };

  const columns = [
    {
      key: 'payment_number',
      label: 'No. Pembayaran',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'payment_date',
      label: 'Tanggal',
      render: (value: string) => formatDate(value),
    },
    {
      key: 'student_name',
      label: 'Nama Siswa',
    },
    {
      key: 'amount',
      label: 'Jumlah',
      render: (value: number) => (
        <span className="font-medium text-green-600">{formatCurrency(value)}</span>
      ),
    },
    {
      key: 'payment_method',
      label: 'Metode',
      render: (value: string) => {
        const badge = getMethodBadge(value);
        return <Badge variant={badge.variant}>{badge.label}</Badge>;
      },
    },
    {
      key: 'reference_no',
      label: 'Ref. No.',
      render: (value: string) => value || '-',
    },
    {
      key: 'id',
      label: 'Aksi',
      render: (value: string) => (
        <div className="flex gap-2">
          <Button
            size="sm"
            variant="ghost"
            onClick={() => navigate(`/payments/${value}`)}
          >
            <Eye className="w-4 h-4" />
          </Button>
          <Button
            size="sm"
            variant="secondary"
            onClick={() => window.open(`/api/v1/payments/${value}/receipt`, '_blank')}
          >
            <Printer className="w-4 h-4" />
          </Button>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Pembayaran</h1>
          <p className="text-gray-600 mt-1">Kelola pembayaran dari siswa</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button
            variant="primary"
            onClick={() => navigate('/payments/new')}
          >
            <Plus className="w-4 h-4 mr-2" />
            Input Pembayaran
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Cari no. pembayaran, nama siswa, atau ref. no..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </div>
      </Card>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <div className="text-sm text-gray-600">Hari Ini</div>
          <div className="text-2xl font-bold text-gray-900 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Minggu Ini</div>
          <div className="text-2xl font-bold text-gray-900 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Bulan Ini</div>
          <div className="text-2xl font-bold text-gray-900 mt-1">
            {formatCurrency(0)}
          </div>
        </Card>
        <Card>
          <div className="text-sm text-gray-600">Total</div>
          <div className="text-2xl font-bold text-gray-900 mt-1">
            {formatCurrency(data?.total || 0)}
          </div>
        </Card>
      </div>

      {/* Table */}
      <Card>
        <Table
          columns={columns}
          data={data?.data || []}
          loading={isLoading}
          emptyMessage="Tidak ada data pembayaran"
        />

        {/* Pagination */}
        {data && data.total > 0 && (
          <div className="mt-4 flex items-center justify-between border-t border-gray-200 pt-4">
            <div className="text-sm text-gray-700">
              Menampilkan {(page - 1) * data.page_size + 1} -{' '}
              {Math.min(page * data.page_size, data.total)} dari {data.total} pembayaran
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
