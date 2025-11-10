import { useNavigate, useSearchParams } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { InventoryItem, ApiResponse } from '@/types';
import { ArrowLeft, Save, Plus, Trash2, TrendingDown, AlertCircle } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useMemo } from 'react';

const stockOutItemSchema = z.object({
  item_id: z.string().min(1, 'Item harus dipilih'),
  quantity: z.number().min(1, 'Quantity minimal 1'),
  notes: z.string().optional(),
});

const stockOutSchema = z.object({
  transaction_date: z.string(),
  transaction_type: z.literal('out'),
  destination: z.string().min(3, 'Tujuan harus diisi'),
  reference_no: z.string().optional(),
  notes: z.string().optional(),
  items: z.array(stockOutItemSchema).min(1, 'Minimal 1 item'),
});

type StockOutFormData = z.infer<typeof stockOutSchema>;

export default function StockOutPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [searchParams] = useSearchParams();
  const preselectedItemId = searchParams.get('item_id');

  const { data: items } = useQuery({
    queryKey: ['inventory-items'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<InventoryItem[]>>('/inventory/items');
      return response.data.data;
    },
  });

  const {
    register,
    handleSubmit,
    control,
    watch,
    formState: { errors },
  } = useForm<StockOutFormData>({
    resolver: zodResolver(stockOutSchema),
    defaultValues: {
      transaction_date: new Date().toISOString().split('T')[0],
      transaction_type: 'out',
      items: preselectedItemId 
        ? [{ item_id: preselectedItemId, quantity: 1, notes: '' }]
        : [{ item_id: '', quantity: 1, notes: '' }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'items',
  });

  const watchedItems = watch('items');

  const totalValue = useMemo(() => {
    return watchedItems.reduce((sum, item) => {
      const selectedItem = items?.find(i => i.id === item.item_id);
      if (!selectedItem) return sum;
      return sum + ((item.quantity || 0) * selectedItem.unit_price);
    }, 0);
  }, [watchedItems, items]);

  const mutation = useMutation({
    mutationFn: async (data: StockOutFormData) => {
      return api.post('/inventory/transactions', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      queryClient.invalidateQueries({ queryKey: ['inventory-transactions'] });
      navigate('/inventory');
    },
  });

  const onSubmit = (data: StockOutFormData) => {
    // Validate stock availability
    const errors: string[] = [];
    data.items.forEach((item, index) => {
      const inventoryItem = items?.find(i => i.id === item.item_id);
      if (inventoryItem && item.quantity > inventoryItem.current_stock) {
        errors.push(`Item #${index + 1}: Stok tidak cukup (tersedia: ${inventoryItem.current_stock})`);
      }
    });

    if (errors.length > 0) {
      alert('Stok tidak mencukupi:\n\n' + errors.join('\n'));
      return;
    }

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
            <h1 className="text-2xl font-bold text-gray-900">Stok Keluar</h1>
            <p className="text-gray-600">Input stok keluar (penggunaan/pemakaian)</p>
          </div>
        </div>
        <Button type="submit" variant="danger" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          Simpan Transaksi
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          {/* Transaction Info */}
          <Card title="Informasi Transaksi">
            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Tanggal Transaksi"
                type="date"
                {...register('transaction_date')}
                error={errors.transaction_date?.message}
                required
              />

              <Input
                label="Tujuan"
                placeholder="e.g., Kelas 1A, Ruang Guru, Kegiatan"
                {...register('destination')}
                error={errors.destination?.message}
                required
                helperText="Untuk apa/kemana stok digunakan"
              />

              <Input
                label="No. Referensi"
                placeholder="e.g., REQ-001, ACT-123"
                {...register('reference_no')}
                helperText="Opsional - nomor permintaan atau referensi"
              />

              <div className="col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Catatan
                </label>
                <textarea
                  {...register('notes')}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Catatan tambahan (opsional)"
                />
              </div>
            </div>
          </Card>

          {/* Items */}
          <Card
            title="Daftar Item"
            action={
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => append({ item_id: '', quantity: 1, notes: '' })}
              >
                <Plus className="w-4 h-4 mr-2" />
                Tambah Item
              </Button>
            }
          >
            <div className="space-y-4">
              {fields.map((field, index) => {
                const selectedItem = items?.find(
                  item => item.id === watchedItems[index]?.item_id
                );
                const requestedQty = watchedItems[index]?.quantity || 0;
                const isOverStock = selectedItem && requestedQty > selectedItem.current_stock;

                return (
                  <div key={field.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-12 gap-4">
                      {/* Item Selection */}
                      <div className="col-span-12 md:col-span-7">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Item <span className="text-red-500">*</span>
                        </label>
                        <select
                          {...register(`items.${index}.item_id`)}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                        >
                          <option value="">Pilih Item</option>
                          {items?.map((item) => (
                            <option key={item.id} value={item.id}>
                              {item.name} (Stok: {item.current_stock} {item.unit})
                            </option>
                          ))}
                        </select>
                        {errors.items?.[index]?.item_id && (
                          <p className="mt-1.5 text-sm text-red-600">
                            {errors.items[index]?.item_id?.message}
                          </p>
                        )}
                        {selectedItem && (
                          <div className="mt-1 flex items-center gap-2">
                            <p className="text-xs text-gray-500">
                              Stok tersedia: {selectedItem.current_stock} {selectedItem.unit}
                            </p>
                            {selectedItem.current_stock <= selectedItem.minimum_stock && (
                              <span className="text-xs text-yellow-600 font-medium">
                                (Stok Rendah!)
                              </span>
                            )}
                          </div>
                        )}
                      </div>

                      {/* Quantity */}
                      <div className="col-span-8 md:col-span-3">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Quantity <span className="text-red-500">*</span>
                        </label>
                        <input
                          type="number"
                          {...register(`items.${index}.quantity`, { valueAsNumber: true })}
                          className={`w-full px-3 py-2 border rounded-lg focus:ring-2 ${
                            isOverStock 
                              ? 'border-red-300 focus:ring-red-500' 
                              : 'border-gray-300 focus:ring-blue-500'
                          }`}
                          placeholder="0"
                        />
                        {errors.items?.[index]?.quantity && (
                          <p className="mt-1.5 text-sm text-red-600">
                            {errors.items[index]?.quantity?.message}
                          </p>
                        )}
                        {isOverStock && (
                          <p className="mt-1.5 text-sm text-red-600 flex items-center gap-1">
                            <AlertCircle className="w-3 h-3" />
                            Stok tidak cukup!
                          </p>
                        )}
                      </div>

                      {/* Delete Button */}
                      <div className="col-span-4 md:col-span-2 flex items-end">
                        {fields.length > 1 && (
                          <Button
                            type="button"
                            variant="danger"
                            className="w-full"
                            onClick={() => remove(index)}
                          >
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        )}
                      </div>

                      {/* Value */}
                      {selectedItem && watchedItems[index]?.quantity && (
                        <div className="col-span-12 pt-2 border-t">
                          <div className="flex justify-between items-center text-sm">
                            <span className="text-gray-600">
                              Nilai @ {formatCurrency(selectedItem.unit_price)}:
                            </span>
                            <span className="font-bold text-red-600">
                              {formatCurrency(
                                (watchedItems[index]?.quantity || 0) * selectedItem.unit_price
                              )}
                            </span>
                          </div>
                        </div>
                      )}

                      {/* Notes */}
                      <div className="col-span-12">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Catatan Item (Opsional)
                        </label>
                        <input
                          {...register(`items.${index}.notes`)}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                          placeholder="Catatan untuk item ini"
                        />
                      </div>
                    </div>
                  </div>
                );
              })}

              {errors.items && typeof errors.items === 'object' && 'message' in errors.items && (
                <p className="text-sm text-red-600">{errors.items.message as string}</p>
              )}
            </div>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-red-100 rounded-lg">
                <TrendingDown className="w-6 h-6 text-red-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">Stok Keluar</div>
                <div className="text-sm text-gray-500">Penggunaan Barang</div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Stok akan berkurang otomatis</p>
              <p>✓ Validasi ketersediaan stok</p>
              <p>✓ History tercatat lengkap</p>
            </div>
          </Card>

          <Card title="Ringkasan Transaksi">
            <div className="space-y-4">
              <div className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                <span className="text-sm text-gray-600">Total Item</span>
                <span className="text-lg font-bold text-gray-900">
                  {fields.length}
                </span>
              </div>

              <div className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                <span className="text-sm text-gray-600">Total Quantity</span>
                <span className="text-lg font-bold text-gray-900">
                  {watchedItems.reduce((sum, item) => sum + (item.quantity || 0), 0)}
                </span>
              </div>

              <div className="p-4 bg-gradient-to-r from-red-50 to-rose-50 rounded-lg border border-red-200">
                <div className="text-sm text-red-700 mb-1">Total Nilai</div>
                <div className="text-2xl font-bold text-red-900">
                  {formatCurrency(totalValue)}
                </div>
              </div>
            </div>
          </Card>

          <Card title="Tips">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• Pastikan stok mencukupi</p>
              <p>• System akan cek ketersediaan</p>
              <p>• Tidak bisa input melebihi stok</p>
              <p>• Alert jika stok menjadi rendah</p>
              <p>• Catatan untuk tracking</p>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
