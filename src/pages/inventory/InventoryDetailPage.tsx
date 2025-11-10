import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { InventoryItem, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  ArrowLeft, 
  Edit,
  Package,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
  Plus,
  Minus,
  RefreshCw
} from 'lucide-react';
import { formatCurrency } from '@/utils/format';

export default function InventoryDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: item, isLoading } = useQuery({
    queryKey: ['inventory-item', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<InventoryItem>>(`/inventory/items/${id}`);
      return response.data.data;
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!item) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Item tidak ditemukan</p>
      </div>
    );
  }

  const getStockStatus = () => {
    if (item.current_stock === 0) {
      return { variant: 'danger' as const, label: 'Stok Habis', icon: AlertTriangle };
    }
    if (item.current_stock <= item.minimum_stock) {
      return { variant: 'warning' as const, label: 'Stok Rendah', icon: TrendingDown };
    }
    return { variant: 'success' as const, label: 'Stok Normal', icon: TrendingUp };
  };

  const status = getStockStatus();
  const StatusIcon = status.icon;
  const totalValue = item.current_stock * item.unit_price;
  const stockPercentage = Math.min((item.current_stock / (item.minimum_stock * 3)) * 100, 100);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/inventory')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{item.name}</h1>
            <p className="text-gray-600 font-mono">{item.item_code}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="success" onClick={() => navigate(`/inventory/stock-in?item_id=${id}`)}>
            <Plus className="w-4 h-4 mr-2" />
            Stok Masuk
          </Button>
          <Button variant="danger" onClick={() => navigate(`/inventory/stock-out?item_id=${id}`)}>
            <Minus className="w-4 h-4 mr-2" />
            Stok Keluar
          </Button>
          <Button variant="primary" onClick={() => navigate(`/inventory/${id}/edit`)}>
            <Edit className="w-4 h-4 mr-2" />
            Edit
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Item Information */}
          <Card title="Informasi Item">
            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Kode Item</label>
                <p className="text-base font-mono font-medium text-gray-900 mt-1">
                  {item.item_code}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Kategori</label>
                <p className="text-base text-gray-900 mt-1">{item.category}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Nama Item</label>
                <p className="text-base font-medium text-gray-900 mt-1">{item.name}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Satuan</label>
                <p className="text-base text-gray-900 mt-1">{item.unit}</p>
              </div>

              {item.brand && (
                <div>
                  <label className="text-sm font-medium text-gray-500">Brand</label>
                  <p className="text-base text-gray-900 mt-1">{item.brand}</p>
                </div>
              )}

              {item.description && (
                <div className="col-span-2">
                  <label className="text-sm font-medium text-gray-500">Deskripsi</label>
                  <p className="text-base text-gray-900 mt-1">{item.description}</p>
                </div>
              )}
            </div>
          </Card>

          {/* Stock Information */}
          <Card title="Informasi Stok">
            <div className="space-y-6">
              {/* Current Stock */}
              <div className="p-4 bg-blue-50 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">Stok Saat Ini</span>
                  <Badge variant={status.variant}>
                    <StatusIcon className="w-3 h-3 mr-1" />
                    {status.label}
                  </Badge>
                </div>
                <div className="flex items-end gap-2">
                  <span className="text-4xl font-bold text-blue-900">
                    {item.current_stock}
                  </span>
                  <span className="text-lg text-blue-700 mb-1">{item.unit}</span>
                </div>
                <div className="mt-3">
                  <div className="w-full bg-blue-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${
                        item.current_stock === 0 ? 'bg-red-600' :
                        item.current_stock <= item.minimum_stock ? 'bg-yellow-600' :
                        'bg-green-600'
                      }`}
                      style={{ width: `${stockPercentage}%` }}
                    ></div>
                  </div>
                </div>
              </div>

              {/* Stock Limits */}
              <div className="grid grid-cols-2 gap-4">
                <div className="p-3 border border-gray-200 rounded-lg">
                  <div className="text-sm text-gray-500 mb-1">Minimum Stok</div>
                  <div className="text-2xl font-bold text-gray-900">
                    {item.minimum_stock} {item.unit}
                  </div>
                </div>

                {item.maximum_stock && (
                  <div className="p-3 border border-gray-200 rounded-lg">
                    <div className="text-sm text-gray-500 mb-1">Maximum Stok</div>
                    <div className="text-2xl font-bold text-gray-900">
                      {item.maximum_stock} {item.unit}
                    </div>
                  </div>
                )}
              </div>

              {/* Stock Alert */}
              {item.current_stock <= item.minimum_stock && (
                <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <div className="flex items-start gap-3">
                    <AlertTriangle className="w-5 h-5 text-yellow-600 mt-0.5" />
                    <div>
                      <div className="font-medium text-yellow-900">
                        {item.current_stock === 0 ? 'Stok Habis!' : 'Peringatan Stok Rendah'}
                      </div>
                      <div className="text-sm text-yellow-700 mt-1">
                        {item.current_stock === 0 
                          ? 'Segera lakukan pembelian atau stok masuk.'
                          : 'Stok sudah mencapai batas minimum. Pertimbangkan untuk melakukan restock.'}
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </Card>

          {/* Pricing Information */}
          <Card title="Informasi Harga & Nilai">
            <div className="space-y-4">
              <div className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                <span className="text-sm font-medium text-gray-700">Harga Satuan</span>
                <span className="text-lg font-bold text-gray-900">
                  {formatCurrency(item.unit_price)}
                </span>
              </div>

              <div className="flex justify-between items-center p-3 bg-green-50 rounded-lg">
                <span className="text-sm font-medium text-gray-700">Total Nilai Stok</span>
                <span className="text-xl font-bold text-green-900">
                  {formatCurrency(totalValue)}
                </span>
              </div>
            </div>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Status Card */}
          <Card>
            <div className="flex items-center justify-between mb-4">
              <span className="text-sm font-medium text-gray-500">Status Stok</span>
              <Badge variant={status.variant}>
                <StatusIcon className="w-3 h-3 mr-1" />
                {status.label}
              </Badge>
            </div>

            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Tersedia</span>
                <span className="font-medium">{item.current_stock} {item.unit}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Min. Stok</span>
                <span className="font-medium">{item.minimum_stock} {item.unit}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Selisih</span>
                <span className={`font-medium ${
                  item.current_stock > item.minimum_stock ? 'text-green-600' : 'text-red-600'
                }`}>
                  {item.current_stock > item.minimum_stock ? '+' : ''}
                  {item.current_stock - item.minimum_stock} {item.unit}
                </span>
              </div>
            </div>
          </Card>

          {/* Quick Actions */}
          <Card title="Aksi Cepat">
            <div className="space-y-2">
              <Button variant="success" className="w-full justify-start">
                <Plus className="w-4 h-4 mr-2" />
                Stok Masuk
              </Button>
              <Button variant="danger" className="w-full justify-start">
                <Minus className="w-4 h-4 mr-2" />
                Stok Keluar
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <RefreshCw className="w-4 h-4 mr-2" />
                Stock Opname
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <Package className="w-4 h-4 mr-2" />
                Riwayat Transaksi
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
