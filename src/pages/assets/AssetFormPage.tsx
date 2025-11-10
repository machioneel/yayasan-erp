import { useNavigate, useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Asset, Branch, ApiResponse } from '@/types';
import { ArrowLeft, Save, Box } from 'lucide-react';
import { formatCurrency } from '@/utils/format';

const assetSchema = z.object({
  name: z.string().min(3, 'Nama aset minimal 3 karakter'),
  category_id: z.string().min(1, 'Kategori harus dipilih'),
  branch_id: z.string().min(1, 'Cabang harus dipilih'),
  acquisition_date: z.string(),
  acquisition_cost: z.number().min(1, 'Nilai perolehan harus lebih dari 0'),
  useful_life: z.number().min(1, 'Masa manfaat harus lebih dari 0'),
  salvage_value: z.number().min(0, 'Nilai residu tidak boleh negatif'),
  brand: z.string().optional(),
  model: z.string().optional(),
  serial_number: z.string().optional(),
  location: z.string().optional(),
  condition: z.string().optional(),
  description: z.string().optional(),
});

type AssetFormData = z.infer<typeof assetSchema>;

const assetCategories = [
  { id: '1', name: 'Kendaraan' },
  { id: '2', name: 'Gedung & Bangunan' },
  { id: '3', name: 'Peralatan Kantor' },
  { id: '4', name: 'Komputer & Elektronik' },
  { id: '5', name: 'Furniture' },
  { id: '6', name: 'Peralatan Mengajar' },
];

const conditions = [
  { value: 'excellent', label: 'Sangat Baik' },
  { value: 'good', label: 'Baik' },
  { value: 'fair', label: 'Cukup' },
  { value: 'poor', label: 'Kurang Baik' },
];

export default function AssetFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEdit = !!id;

  const { data: branches } = useQuery({
    queryKey: ['branches'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Branch[]>>('/branches');
      return response.data.data;
    },
  });

  const { data: asset } = useQuery({
    queryKey: ['asset', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Asset>>(`/assets/${id}`);
      return response.data.data;
    },
    enabled: isEdit,
  });

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<AssetFormData>({
    resolver: zodResolver(assetSchema),
    defaultValues: asset || {
      acquisition_date: new Date().toISOString().split('T')[0],
      salvage_value: 0,
    },
  });

  const acquisitionCost = watch('acquisition_cost') || 0;
  const usefulLife = watch('useful_life') || 1;
  const salvageValue = watch('salvage_value') || 0;

  const annualDepreciation = (acquisitionCost - salvageValue) / usefulLife;
  const monthlyDepreciation = annualDepreciation / 12;

  const mutation = useMutation({
    mutationFn: async (data: AssetFormData) => {
      if (isEdit) {
        return api.put(`/assets/${id}`, data);
      }
      return api.post('/assets', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['assets'] });
      navigate('/assets');
    },
  });

  const onSubmit = (data: AssetFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/assets')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit Aset' : 'Tambah Aset Baru'}
            </h1>
            <p className="text-gray-600">
              {isEdit ? 'Update informasi aset' : 'Lengkapi data aset tetap'}
            </p>
          </div>
        </div>
        <Button type="submit" variant="primary" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update' : 'Simpan'}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card title="Informasi Aset">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="md:col-span-2">
                <Input
                  label="Nama Aset"
                  placeholder="e.g., Laptop Dell Latitude 5420"
                  {...register('name')}
                  error={errors.name?.message}
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Kategori <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('category_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Kategori</option>
                  {assetCategories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
                {errors.category_id && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.category_id.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Cabang <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('branch_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Cabang</option>
                  {branches?.map((branch) => (
                    <option key={branch.id} value={branch.id}>
                      {branch.name}
                    </option>
                  ))}
                </select>
                {errors.branch_id && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.branch_id.message}</p>
                )}
              </div>

              <Input
                label="Brand/Merk"
                placeholder="e.g., Dell, Toyota"
                {...register('brand')}
                error={errors.brand?.message}
              />

              <Input
                label="Model/Tipe"
                placeholder="e.g., Latitude 5420, Avanza 1.5"
                {...register('model')}
                error={errors.model?.message}
              />

              <Input
                label="Serial Number"
                placeholder="e.g., ABC123XYZ"
                {...register('serial_number')}
                error={errors.serial_number?.message}
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Kondisi
                </label>
                <select
                  {...register('condition')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Kondisi</option>
                  {conditions.map((cond) => (
                    <option key={cond.value} value={cond.value}>
                      {cond.label}
                    </option>
                  ))}
                </select>
              </div>

              <div className="md:col-span-2">
                <Input
                  label="Lokasi"
                  placeholder="e.g., Ruang Guru, Kantor Pusat"
                  {...register('location')}
                  error={errors.location?.message}
                />
              </div>

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Deskripsi
                </label>
                <textarea
                  {...register('description')}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Informasi tambahan tentang aset"
                />
              </div>
            </div>
          </Card>

          <Card title="Informasi Keuangan">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Tanggal Perolehan"
                type="date"
                {...register('acquisition_date')}
                error={errors.acquisition_date?.message}
                required
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Nilai Perolehan <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">Rp</span>
                  <input
                    type="number"
                    {...register('acquisition_cost', { valueAsNumber: true })}
                    className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="0"
                  />
                </div>
                {errors.acquisition_cost && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.acquisition_cost.message}</p>
                )}
              </div>

              <div>
                <Input
                  label="Masa Manfaat (Tahun)"
                  type="number"
                  {...register('useful_life', { valueAsNumber: true })}
                  error={errors.useful_life?.message}
                  required
                  helperText="Berapa tahun aset akan digunakan"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Nilai Residu
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">Rp</span>
                  <input
                    type="number"
                    {...register('salvage_value', { valueAsNumber: true })}
                    className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="0"
                  />
                </div>
                <p className="mt-1.5 text-xs text-gray-500">
                  Nilai estimasi saat aset tidak digunakan lagi
                </p>
              </div>
            </div>
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-blue-100 rounded-lg">
                <Box className="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">Aset Tetap</div>
                <div className="text-sm text-gray-500">
                  {isEdit ? 'Update data' : 'Data baru'}
                </div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Kode aset otomatis dibuat sistem</p>
              <p>✓ Penyusutan dihitung otomatis</p>
              <p>✓ Data dapat diubah sewaktu-waktu</p>
            </div>
          </Card>

          <Card title="Kalkulasi Penyusutan">
            <div className="space-y-4">
              <div className="p-3 bg-gray-50 rounded-lg">
                <div className="text-sm text-gray-600">Nilai Perolehan</div>
                <div className="text-lg font-bold text-gray-900 mt-1">
                  {formatCurrency(acquisitionCost)}
                </div>
              </div>

              <div className="p-3 bg-gray-50 rounded-lg">
                <div className="text-sm text-gray-600">Nilai Residu</div>
                <div className="text-lg font-bold text-gray-900 mt-1">
                  {formatCurrency(salvageValue)}
                </div>
              </div>

              <div className="p-3 bg-blue-50 rounded-lg">
                <div className="text-sm text-gray-600">Penyusutan/Tahun</div>
                <div className="text-lg font-bold text-blue-900 mt-1">
                  {formatCurrency(annualDepreciation)}
                </div>
              </div>

              <div className="p-3 bg-green-50 rounded-lg">
                <div className="text-sm text-gray-600">Penyusutan/Bulan</div>
                <div className="text-lg font-bold text-green-900 mt-1">
                  {formatCurrency(monthlyDepreciation)}
                </div>
              </div>
            </div>
          </Card>

          <Card title="Tips">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• Pastikan nilai perolehan sesuai bukti pembelian</p>
              <p>• Masa manfaat sesuai peraturan perpajakan</p>
              <p>• Nilai residu bisa 0 jika tidak ada nilai jual</p>
              <p>• Serial number untuk tracking yang lebih baik</p>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
