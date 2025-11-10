import { useNavigate, useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { InventoryItem, ApiResponse } from '@/types';
import { ArrowLeft, Save, Package } from 'lucide-react';
import { formatCurrency } from '@/utils/format';

const inventorySchema = z.object({
  name: z.string().min(3, 'Nama item minimal 3 karakter'),
  category: z.string().min(1, 'Kategori harus diisi'),
  unit: z.string().min(1, 'Satuan harus diisi'),
  minimum_stock: z.number().min(0, 'Stok minimum tidak boleh negatif'),
  maximum_stock: z.number().optional(),
  unit_price: z.number().min(0, 'Harga tidak boleh negatif'),
  brand: z.string().optional(),
  description: z.string().optional(),
  initial_stock: z.number().min(0, 'Stok awal tidak boleh negatif').optional(),
});

type InventoryFormData = z.infer<typeof inventorySchema>;

const categories = [
  'Alat Tulis',
  'Peralatan Kebersihan',
  'Perlengkapan Mengajar',
  'Buku & Modul',
  'Elektronik',
  'Konsumsi',
  'Obat-obatan',
  'Lainnya',
];

const units = [
  'Pcs',
  'Box',
  'Pack',
  'Lusin',
  'Kg',
  'Liter',
  'Meter',
  'Set',
];

export default function InventoryFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEdit = !!id;

  const { data: item } = useQuery({
    queryKey: ['inventory-item', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<InventoryItem>>(`/inventory/items/${id}`);
      return response.data.data;
    },
    enabled: isEdit,
  });

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<InventoryFormData>({
    resolver: zodResolver(inventorySchema),
    defaultValues: item || {
      minimum_stock: 10,
      initial_stock: 0,
    },
  });

  const unitPrice = watch('unit_price') || 0;
  const initialStock = watch('initial_stock') || 0;
  const minimumStock = watch('minimum_stock') || 0;

  const initialValue = unitPrice * initialStock;

  const mutation = useMutation({
    mutationFn: async (data: InventoryFormData) => {
      if (isEdit) {
        return api.put(`/inventory/items/${id}`, data);
      }
      return api.post('/inventory/items', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      navigate('/inventory');
    },
  });

  const onSubmit = (data: InventoryFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/inventory')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit Item Inventory' : 'Tambah Item Baru'}
            </h1>
            <p className="text-gray-600">
              {isEdit ? 'Update informasi item' : 'Lengkapi data item inventory'}
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
          <Card title="Informasi Item">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="md:col-span-2">
                <Input
                  label="Nama Item"
                  placeholder="e.g., Kertas HVS A4, Spidol Whiteboard"
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
                  {...register('category')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Kategori</option>
                  {categories.map((cat) => (
                    <option key={cat} value={cat}>
                      {cat}
                    </option>
                  ))}
                </select>
                {errors.category && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.category.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Satuan <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('unit')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Satuan</option>
                  {units.map((unit) => (
                    <option key={unit} value={unit}>
                      {unit}
                    </option>
                  ))}
                </select>
                {errors.unit && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.unit.message}</p>
                )}
              </div>

              <Input
                label="Brand/Merk"
                placeholder="e.g., Sidu, Snowman"
                {...register('brand')}
                error={errors.brand?.message}
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Harga Satuan <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">Rp</span>
                  <input
                    type="number"
                    {...register('unit_price', { valueAsNumber: true })}
                    className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="0"
                  />
                </div>
                {errors.unit_price && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.unit_price.message}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Deskripsi
                </label>
                <textarea
                  {...register('description')}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Informasi tambahan tentang item"
                />
              </div>
            </div>
          </Card>

          <Card title="Pengaturan Stok">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {!isEdit && (
                <div>
                  <Input
                    label="Stok Awal"
                    type="number"
                    {...register('initial_stock', { valueAsNumber: true })}
                    error={errors.initial_stock?.message}
                    helperText="Jumlah stok saat item pertama kali ditambahkan"
                  />
                </div>
              )}

              <div>
                <Input
                  label="Minimum Stok"
                  type="number"
                  {...register('minimum_stock', { valueAsNumber: true })}
                  error={errors.minimum_stock?.message}
                  required
                  helperText="Batas peringatan stok rendah"
                />
              </div>

              <div>
                <Input
                  label="Maximum Stok (Opsional)"
                  type="number"
                  {...register('maximum_stock', { valueAsNumber: true })}
                  error={errors.maximum_stock?.message}
                  helperText="Batas maksimal stok"
                />
              </div>
            </div>

            {!isEdit && initialStock > 0 && (
              <div className="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg">
                <div className="flex justify-between items-center">
                  <div className="text-sm text-green-700">
                    Nilai Total Stok Awal:
                  </div>
                  <div className="text-lg font-bold text-green-900">
                    {formatCurrency(initialValue)}
                  </div>
                </div>
              </div>
            )}
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-blue-100 rounded-lg">
                <Package className="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">Item Inventory</div>
                <div className="text-sm text-gray-500">
                  {isEdit ? 'Update data' : 'Data baru'}
                </div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Kode item otomatis dibuat</p>
              <p>✓ Stok akan ter-tracking otomatis</p>
              <p>✓ Alert jika stok di bawah minimum</p>
            </div>
          </Card>

          {!isEdit && (
            <Card title="Ringkasan">
              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Harga Satuan</span>
                  <span className="font-medium">{formatCurrency(unitPrice)}</span>
                </div>

                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Stok Awal</span>
                  <span className="font-medium">{initialStock}</span>
                </div>

                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Min. Stok</span>
                  <span className="font-medium">{minimumStock}</span>
                </div>

                <div className="border-t pt-3">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium">Total Nilai</span>
                    <span className="text-lg font-bold text-blue-600">
                      {formatCurrency(initialValue)}
                    </span>
                  </div>
                </div>
              </div>
            </Card>
          )}

          <Card title="Tips">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• Set minimum stok untuk alert otomatis</p>
              <p>• Harga satuan bisa diupdate sewaktu-waktu</p>
              <p>• Gunakan satuan yang konsisten</p>
              <p>• Kategori membantu dalam pelaporan</p>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
