import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api } from '@/services/api';
import { InventoryItem, PaginationResponse, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Table } from '@/components/common/Table';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Plus, Search, Package, TrendingUp, TrendingDown, AlertTriangle } from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function InventoryPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['inventory', page, search],
    queryFn: async () => {
      const response = await api.get<ApiResponse<PaginationResponse<InventoryItem>>>(
        '/inventory/items',
        {
          params: { page, page_size: 20, search },
        }
      );
      return response.data.data;
    },
  });

  const getStockStatus = (current: number, min: number) => {
    if (current === 0) return { variant: 'danger' as const, label: 'Habis', icon: AlertTriangle };
    if (current <= min) return { variant: 'warning' as const, label: 'Rendah', icon: TrendingDown };
    return { variant: 'success' as const, label: 'Normal', icon: TrendingUp };
  };

  const columns = [
    {
      key: 'item_code',
      label: 'Kode Item',
      render: (value: string) => (
        <span className="font-mono font-medium">{value}</span>
      ),
    },
    {
      key: 'name',
      label: 'Nama Item',
      render: (value: string, row: InventoryItem) => (
        <div>
          <div className="font-medium text-gray-900">{value}</div>
          <div className="text-xs text-gray-500">{row.category}</div>
        </div>
      ),
    },
    {
      key: 'unit',
      label: 'Satuan',
    },
    {
      key: 'current_stock',
      label: 'Stok',
      render: (value: number, row: InventoryItem) => {
        const status = getStockStatus(value, row.minimum_stock);
        const Icon = status.icon;
        return (
          <div className="flex items-center gap-2">
            <span className="font-medium">{value}</span>
            <Icon className={`w-4 h-4 ${
              status.variant === 'danger' ? 'text-red-600' :
              status.variant === 'warning' ? 'text-yellow-600' :
              'text-green-600'
            }`} />
          </div>
        );
      },
    },
    {
      key: 'minimum_stock',
      label: 'Min. Stok',
      render: (value: number) => (
        <span className="text-gray-600">{value}</span>
      ),
    },
    {
      key: 'unit_price',
      label: 'Harga Satuan',
      render: (value: number) => formatCurrency(value),
    },
    {
      key: 'total_value',
      label: 'Nilai Total',
      render: (value: number, row: InventoryItem) => (
        <span className="font-medium">{formatCurrency(row.current_stock * row.unit_price)}</span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (value: number, row: InventoryItem) => {
        const status = getStockStatus(row.current_stock, row.minimum_stock);
        return <Badge variant={status.variant}>{status.label}</Badge>;
      },
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Manajemen Inventory</h1>
          <p className="text-gray-600 mt-1">Kelola stok barang dan persediaan</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => navigate('/inventory/transactions')}>
            Transaksi
          </Button>
          <Button variant="primary" onClick={() => navigate('/inventory/new')}>
            <Plus className="w-4 h-4 mr-2" />
            Tambah Item
          </Button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-600">Total Item</div>
              <div className="text-2xl font-bold text-gray-900 mt-1">
                {data?.total || 0}
              </div>
            </div>
            <div className="p-3 bg-blue-100 rounded-lg">
              <Package className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-600">Stok Normal</div>
              <div className="text-2xl font-bold text-green-600 mt-1">0</div>
            </div>
            <TrendingUp className="w-8 h-8 text-green-600" />
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-600">Stok Rendah</div>
              <div className="text-2xl font-bold text-yellow-600 mt-1">0</div>
            </div>
            <TrendingDown className="w-8 h-8 text-yellow-600" />
          </div>
        </Card>

        <Card>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-600">Stok Habis</div>
              <div className="text-2xl font-bold text-red-600 mt-1">0</div>
            </div>
            <AlertTriangle className="w-8 h-8 text-red-600" />
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
                placeholder="Cari kode atau nama item..."
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
          
          <Button variant="secondary">
            Stok Rendah
          </Button>
          <Button variant="secondary">
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
          emptyMessage="Tidak ada data inventory"
          onRowClick={(row) => navigate(`/inventory/${row.id}`)}
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
