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
import { ArrowLeft, Save, Plus, Trash2, TrendingUp } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useMemo } from 'react';

const stockInItemSchema = z.object({
  item_id: z.string().min(1, 'Item harus dipilih'),
  quantity: z.number().min(1, 'Quantity minimal 1'),
  unit_price: z.number().min(0, 'Harga tidak boleh negatif'),
  notes: z.string().optional(),
});

const stockInSchema = z.object({
  transaction_date: z.string(),
  transaction_type: z.literal('in'),
  source: z.string().min(3, 'Sumber harus diisi'),
  reference_no: z.string().optional(),
  notes: z.string().optional(),
  items: z.array(stockInItemSchema).min(1, 'Minimal 1 item'),
});

type StockInFormData = z.infer<typeof stockInSchema>;

export default function StockInPage() {
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
    setValue,
    formState: { errors },
  } = useForm<StockInFormData>({
    resolver: zodResolver(stockInSchema),
    defaultValues: {
      transaction_date: new Date().toISOString().split('T')[0],
      transaction_type: 'in',
      items: preselectedItemId 
        ? [{ item_id: preselectedItemId, quantity: 1, unit_price: 0, notes: '' }]
        : [{ item_id: '', quantity: 1, unit_price: 0, notes: '' }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'items',
  });

  const watchedItems = watch('items');

  const totalValue = useMemo(() => {
    return watchedItems.reduce((sum, item) => {
      return sum + ((item.quantity || 0) * (item.unit_price || 0));
    }, 0);
  }, [watchedItems]);

  const mutation = useMutation({
    mutationFn: async (data: StockInFormData) => {
      return api.post('/inventory/transactions', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory'] });
      queryClient.invalidateQueries({ queryKey: ['inventory-transactions'] });
      navigate('/inventory');
    },
  });

  const onSubmit = (data: StockInFormData) => {
    mutation.mutate(data);
  };

  const handleItemSelect = (index: number, itemId: string) => {
    const selectedItem = items?.find(item => item.id === itemId);
    if (selectedItem) {
      setValue(`items.${index}.unit_price`, selectedItem.unit_price);
    }
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
            <h1 className="text-2xl font-bold text-gray-900">Stok Masuk</h1>
            <p className="text-gray-600">Input stok masuk (pembelian/penerimaan)</p>
          </div>
        </div>
        <Button type="submit" variant="success" loading={mutation.isPending}>
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
                label="Sumber"
                placeholder="e.g., Supplier, Donatur, Transfer Cabang"
                {...register('source')}
                error={errors.source?.message}
                required
                helperText="Dari mana stok diperoleh"
              />

              <Input
                label="No. Referensi"
                placeholder="e.g., PO-001, DO-123"
                {...register('reference_no')}
                helperText="Opsional - nomor PO, DO, atau referensi lain"
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
                onClick={() => append({ item_id: '', quantity: 1, unit_price: 0, notes: '' })}
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

                return (
                  <div key={field.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-12 gap-4">
                      {/* Item Selection */}
                      <div className="col-span-12 md:col-span-5">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Item <span className="text-red-500">*</span>
                        </label>
                        <select
                          {...register(`items.${index}.item_id`)}
                          onChange={(e) => handleItemSelect(index, e.target.value)}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                        >
                          <option value="">Pilih Item</option>
                          {items?.map((item) => (
                            <option key={item.id} value={item.id}>
                              {item.name} ({item.current_stock} {item.unit})
                            </option>
                          ))}
                        </select>
                        {errors.items?.[index]?.item_id && (
                          <p className="mt-1.5 text-sm text-red-600">
                            {errors.items[index]?.item_id?.message}
                          </p>
                        )}
                        {selectedItem && (
                          <p className="mt-1 text-xs text-gray-500">
                            Stok saat ini: {selectedItem.current_stock} {selectedItem.unit}
                          </p>
                        )}
                      </div>

                      {/* Quantity */}
                      <div className="col-span-6 md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Quantity <span className="text-red-500">*</span>
                        </label>
                        <input
                          type="number"
                          {...register(`items.${index}.quantity`, { valueAsNumber: true })}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                          placeholder="0"
                        />
                        {errors.items?.[index]?.quantity && (
                          <p className="mt-1.5 text-sm text-red-600">
                            {errors.items[index]?.quantity?.message}
                          </p>
                        )}
                      </div>

                      {/* Unit Price */}
                      <div className="col-span-6 md:col-span-3">
                        <label className="block text-sm font-medium text-gray-700 mb-1.5">
                          Harga Satuan
                        </label>
                        <input
                          type="number"
                          {...register(`items.${index}.unit_price`, { valueAsNumber: true })}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                          placeholder="0"
                        />
                      </div>

                      {/* Delete Button */}
                      <div className="col-span-12 md:col-span-2 flex items-end">
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

                      {/* Subtotal */}
                      {watchedItems[index]?.quantity && watchedItems[index]?.unit_price && (
                        <div className="col-span-12 pt-2 border-t">
                          <div className="flex justify-between items-center text-sm">
                            <span className="text-gray-600">Subtotal:</span>
                            <span className="font-bold text-green-600">
                              {formatCurrency(
                                (watchedItems[index]?.quantity || 0) * 
                                (watchedItems[index]?.unit_price || 0)
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
              <div className="p-3 bg-green-100 rounded-lg">
                <TrendingUp className="w-6 h-6 text-green-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">Stok Masuk</div>
                <div className="text-sm text-gray-500">Penerimaan Barang</div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Stok akan bertambah otomatis</p>
              <p>✓ Nilai inventory terupdate</p>
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

              <div className="p-4 bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg border border-green-200">
                <div className="text-sm text-green-700 mb-1">Total Nilai</div>
                <div className="text-2xl font-bold text-green-900">
                  {formatCurrency(totalValue)}
                </div>
              </div>
            </div>
          </Card>

          <Card title="Tips">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• Pilih item yang akan ditambah stoknya</p>
              <p>• Input quantity sesuai jumlah yang diterima</p>
              <p>• Harga satuan otomatis dari data item</p>
              <p>• Bisa edit harga jika berbeda</p>
              <p>• Tambahkan catatan jika diperlukan</p>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
