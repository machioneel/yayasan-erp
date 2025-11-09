import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Asset, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  ArrowLeft, 
  Edit, 
  Box,
  Calendar,
  DollarSign,
  MapPin,
  TrendingDown,
  Wrench,
  ArrowRightLeft
} from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function AssetDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: asset, isLoading } = useQuery({
    queryKey: ['asset', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Asset>>(`/assets/${id}`);
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

  if (!asset) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Aset tidak ditemukan</p>
      </div>
    );
  }

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'maintenance': return 'warning';
      case 'disposed': return 'danger';
      case 'inactive': return 'default';
      default: return 'default';
    }
  };

  const depreciation = asset.acquisition_cost - asset.book_value;
  const depreciationRate = (depreciation / asset.acquisition_cost) * 100;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/assets')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{asset.name}</h1>
            <p className="text-gray-600 font-mono">{asset.asset_code}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="primary" onClick={() => navigate(`/assets/${id}/edit`)}>
            <Edit className="w-4 h-4 mr-2" />
            Edit
          </Button>
          <Button variant="secondary">
            <Wrench className="w-4 h-4 mr-2" />
            Maintenance
          </Button>
          <Button variant="secondary">
            <ArrowRightLeft className="w-4 h-4 mr-2" />
            Transfer
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Asset Information */}
          <Card title="Informasi Aset">
            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Kode Aset</label>
                <p className="text-base font-mono font-medium text-gray-900 mt-1">
                  {asset.asset_code}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Kategori</label>
                <p className="text-base text-gray-900 mt-1">{asset.category_name}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Nama Aset</label>
                <p className="text-base font-medium text-gray-900 mt-1">{asset.name}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Brand</label>
                <p className="text-base text-gray-900 mt-1">{asset.brand || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Model</label>
                <p className="text-base text-gray-900 mt-1">{asset.model || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Serial Number</label>
                <p className="text-base font-mono text-gray-900 mt-1">
                  {asset.serial_number || '-'}
                </p>
              </div>

              {asset.description && (
                <div className="col-span-2">
                  <label className="text-sm font-medium text-gray-500">Deskripsi</label>
                  <p className="text-base text-gray-900 mt-1">{asset.description}</p>
                </div>
              )}
            </div>
          </Card>

          {/* Financial Information */}
          <Card title="Informasi Keuangan">
            <div className="space-y-4">
              <div className="flex justify-between items-center p-3 bg-blue-50 rounded-lg">
                <span className="text-sm font-medium text-gray-700">Nilai Perolehan</span>
                <span className="text-lg font-bold text-blue-900">
                  {formatCurrency(asset.acquisition_cost)}
                </span>
              </div>

              <div className="flex justify-between items-center p-3 bg-red-50 rounded-lg">
                <span className="text-sm font-medium text-gray-700">Akumulasi Penyusutan</span>
                <span className="text-lg font-bold text-red-900">
                  {formatCurrency(depreciation)}
                </span>
              </div>

              <div className="flex justify-between items-center p-3 bg-green-50 rounded-lg">
                <span className="text-sm font-medium text-gray-700">Nilai Buku</span>
                <span className="text-xl font-bold text-green-900">
                  {formatCurrency(asset.book_value)}
                </span>
              </div>

              <div className="pt-3 border-t">
                <div className="flex justify-between text-sm text-gray-600 mb-2">
                  <span>Penyusutan:</span>
                  <span className="font-medium">{depreciationRate.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-red-600 h-2 rounded-full"
                    style={{ width: `${depreciationRate}%` }}
                  ></div>
                </div>
              </div>
            </div>
          </Card>

          {/* Location & Usage */}
          <Card title="Lokasi & Penggunaan">
            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Lokasi</label>
                <div className="flex items-center gap-2 mt-1">
                  <MapPin className="w-4 h-4 text-gray-400" />
                  <p className="text-base text-gray-900">{asset.location || '-'}</p>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Cabang</label>
                <p className="text-base text-gray-900 mt-1">{asset.branch_name}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">PIC</label>
                <p className="text-base text-gray-900 mt-1">{asset.pic_name || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Kondisi</label>
                <p className="text-base text-gray-900 mt-1">{asset.condition || '-'}</p>
              </div>
            </div>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Status Card */}
          <Card>
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium text-gray-500">Status</label>
                <div className="mt-2">
                  <Badge variant={getStatusVariant(asset.status)}>
                    {asset.status === 'active' ? 'Aktif' :
                     asset.status === 'maintenance' ? 'Maintenance' :
                     asset.status === 'disposed' ? 'Dihapus' : 'Nonaktif'}
                  </Badge>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tanggal Perolehan</label>
                <div className="flex items-center gap-2 mt-1">
                  <Calendar className="w-4 h-4 text-gray-400" />
                  <p className="text-base text-gray-900">
                    {formatDate(asset.acquisition_date)}
                  </p>
                </div>
              </div>

              {asset.useful_life && (
                <div>
                  <label className="text-sm font-medium text-gray-500">Masa Manfaat</label>
                  <p className="text-base text-gray-900 mt-1">
                    {asset.useful_life} tahun
                  </p>
                </div>
              )}
            </div>
          </Card>

          {/* Quick Actions */}
          <Card title="Aksi Cepat">
            <div className="space-y-2">
              <Button variant="secondary" className="w-full justify-start">
                <Wrench className="w-4 h-4 mr-2" />
                Jadwalkan Maintenance
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <ArrowRightLeft className="w-4 h-4 mr-2" />
                Transfer Lokasi
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <TrendingDown className="w-4 h-4 mr-2" />
                Riwayat Penyusutan
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <DollarSign className="w-4 h-4 mr-2" />
                Riwayat Transaksi
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
