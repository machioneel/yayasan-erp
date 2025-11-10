import { useNavigate, useParams } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Student, ApiResponse } from '@/types';
import { ArrowLeft, Save, Plus, Trash2, Calculator } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useEffect } from 'react';

const invoiceItemSchema = z.object({
  description: z.string().min(3, 'Deskripsi minimal 3 karakter'),
  quantity: z.number().min(1, 'Jumlah minimal 1'),
  unit_price: z.number().min(1, 'Harga satuan harus lebih dari 0'),
  amount: z.number(),
});

const invoiceSchema = z.object({
  student_id: z.string().min(1, 'Siswa harus dipilih'),
  invoice_date: z.string(),
  due_date: z.string(),
  description: z.string().optional(),
  items: z.array(invoiceItemSchema).min(1, 'Minimal 1 item harus diisi'),
});

type InvoiceFormData = z.infer<typeof invoiceSchema>;

export default function InvoiceFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEdit = !!id;

  const {
    register,
    handleSubmit,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useForm<InvoiceFormData>({
    resolver: zodResolver(invoiceSchema),
    defaultValues: {
      invoice_date: new Date().toISOString().split('T')[0],
      due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      items: [{ description: '', quantity: 1, unit_price: 0, amount: 0 }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'items',
  });

  // Fetch students
  const { data: students } = useQuery({
    queryKey: ['students-active'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Student[]>>('/students', {
        params: { status: 'active' },
      });
      return response.data.data;
    },
  });

  const watchedItems = watch('items');

  // Auto-calculate amount for each item
  useEffect(() => {
    watchedItems.forEach((item, index) => {
      const amount = item.quantity * item.unit_price;
      if (amount !== item.amount) {
        setValue(`items.${index}.amount`, amount);
      }
    });
  }, [watchedItems, setValue]);

  const totalAmount = watchedItems.reduce((sum, item) => sum + (item.amount || 0), 0);

  const mutation = useMutation({
    mutationFn: async (data: InvoiceFormData) => {
      if (isEdit) {
        return api.put(`/invoices/${id}`, data);
      }
      return api.post('/invoices', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      navigate('/invoices');
    },
  });

  const onSubmit = (data: InvoiceFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/invoices')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit Invoice' : 'Buat Invoice Baru'}
            </h1>
            <p className="text-gray-600">Buat tagihan untuk siswa</p>
          </div>
        </div>
        <Button type="submit" variant="primary" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update' : 'Simpan'}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Invoice Form */}
        <div className="lg:col-span-2 space-y-6">
          <Card title="Informasi Invoice">
            <div className="grid grid-cols-2 gap-4">
              <div className="col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Siswa <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('student_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Siswa</option>
                  {students?.map((student) => (
                    <option key={student.id} value={student.id}>
                      {student.full_name} - {student.class_name}
                    </option>
                  ))}
                </select>
                {errors.student_id && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.student_id.message}</p>
                )}
              </div>

              <Input
                label="Tanggal Invoice"
                type="date"
                {...register('invoice_date')}
                error={errors.invoice_date?.message}
                required
              />

              <Input
                label="Jatuh Tempo"
                type="date"
                {...register('due_date')}
                error={errors.due_date?.message}
                required
              />

              <div className="col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Keterangan
                </label>
                <textarea
                  {...register('description')}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Keterangan tambahan (opsional)"
                />
              </div>
            </div>
          </Card>

          <Card
            title="Item Invoice"
            action={
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => append({ description: '', quantity: 1, unit_price: 0, amount: 0 })}
              >
                <Plus className="w-4 h-4 mr-2" />
                Tambah Item
              </Button>
            }
          >
            <div className="space-y-4">
              {fields.map((field, index) => (
                <div key={field.id} className="p-4 border border-gray-200 rounded-lg">
                  <div className="flex items-start justify-between mb-3">
                    <h4 className="font-medium text-gray-900">Item #{index + 1}</h4>
                    {fields.length > 1 && (
                      <Button
                        type="button"
                        variant="danger"
                        size="sm"
                        onClick={() => remove(index)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    )}
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-12 gap-3">
                    <div className="md:col-span-6">
                      <Input
                        label="Deskripsi"
                        placeholder="e.g., SPP Bulan Januari"
                        {...register(`items.${index}.description`)}
                        error={errors.items?.[index]?.description?.message}
                        required
                      />
                    </div>

                    <div className="md:col-span-2">
                      <Input
                        label="Jumlah"
                        type="number"
                        {...register(`items.${index}.quantity`, { valueAsNumber: true })}
                        error={errors.items?.[index]?.quantity?.message}
                        required
                      />
                    </div>

                    <div className="md:col-span-4">
                      <label className="block text-sm font-medium text-gray-700 mb-1.5">
                        Harga Satuan <span className="text-red-500">*</span>
                      </label>
                      <div className="relative">
                        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">Rp</span>
                        <input
                          type="number"
                          {...register(`items.${index}.unit_price`, { valueAsNumber: true })}
                          className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                        />
                      </div>
                      {errors.items?.[index]?.unit_price && (
                        <p className="mt-1.5 text-sm text-red-600">
                          {errors.items[index]?.unit_price?.message}
                        </p>
                      )}
                    </div>

                    <div className="md:col-span-12">
                      <div className="flex justify-end items-center gap-2 text-sm">
                        <span className="text-gray-600">Subtotal:</span>
                        <span className="font-bold text-gray-900">
                          {formatCurrency(watchedItems[index]?.amount || 0)}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}

              {errors.items && typeof errors.items === 'object' && 'message' in errors.items && (
                <p className="text-sm text-red-600">{errors.items.message as string}</p>
              )}
            </div>
          </Card>
        </div>

        {/* Summary Sidebar */}
        <div>
          <Card title="Ringkasan">
            <div className="space-y-4">
              <div className="flex items-center justify-between text-gray-700">
                <span>Jumlah Item</span>
                <span className="font-medium">{fields.length}</span>
              </div>

              <div className="border-t pt-4">
                <div className="flex items-center justify-between">
                  <span className="text-lg font-medium text-gray-900">Total</span>
                  <span className="text-2xl font-bold text-blue-600">
                    {formatCurrency(totalAmount)}
                  </span>
                </div>
              </div>

              <Button type="submit" variant="primary" className="w-full" loading={mutation.isPending}>
                <Save className="w-4 h-4 mr-2" />
                {isEdit ? 'Update Invoice' : 'Buat Invoice'}
              </Button>
            </div>
          </Card>

          <Card className="mt-4">
            <div className="flex items-start gap-3">
              <Calculator className="w-5 h-5 text-blue-600 mt-0.5" />
              <div className="text-sm text-gray-600">
                <p className="font-medium text-gray-900 mb-1">Tips:</p>
                <ul className="space-y-1 list-disc list-inside">
                  <li>Total akan dihitung otomatis</li>
                  <li>Jatuh tempo default 30 hari</li>
                  <li>Invoice bisa diedit sebelum dibayar</li>
                </ul>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
