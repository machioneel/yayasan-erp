import { useNavigate } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Account, ApiResponse } from '@/types';
import { ArrowLeft, Save, Plus, Trash2, AlertCircle, CheckCircle } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useEffect, useMemo } from 'react';

const journalItemSchema = z.object({
  account_id: z.string().min(1, 'Akun harus dipilih'),
  description: z.string().min(3, 'Deskripsi minimal 3 karakter'),
  debit: z.number().min(0, 'Debit tidak boleh negatif'),
  credit: z.number().min(0, 'Kredit tidak boleh negatif'),
});

const journalSchema = z.object({
  journal_date: z.string(),
  description: z.string().min(5, 'Deskripsi minimal 5 karakter'),
  reference_no: z.string().optional(),
  items: z.array(journalItemSchema)
    .min(2, 'Minimal 2 baris (debit dan kredit)')
    .refine((items) => {
      return items.every(item => item.debit === 0 || item.credit === 0);
    }, 'Setiap baris harus debit ATAU kredit, tidak boleh keduanya'),
});

type JournalFormData = z.infer<typeof journalSchema>;

export default function JournalFormPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: accounts } = useQuery({
    queryKey: ['accounts'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Account[]>>('/accounts');
      return response.data.data;
    },
  });

  const {
    register,
    handleSubmit,
    control,
    watch,
    formState: { errors },
  } = useForm<JournalFormData>({
    resolver: zodResolver(journalSchema),
    defaultValues: {
      journal_date: new Date().toISOString().split('T')[0],
      items: [
        { account_id: '', description: '', debit: 0, credit: 0 },
        { account_id: '', description: '', debit: 0, credit: 0 },
      ],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'items',
  });

  const watchedItems = watch('items');

  const { totalDebit, totalCredit, isBalanced } = useMemo(() => {
    const debit = watchedItems.reduce((sum, item) => sum + (item.debit || 0), 0);
    const credit = watchedItems.reduce((sum, item) => sum + (item.credit || 0), 0);
    return {
      totalDebit: debit,
      totalCredit: credit,
      isBalanced: debit === credit && debit > 0,
    };
  }, [watchedItems]);

  const mutation = useMutation({
    mutationFn: async (data: JournalFormData) => {
      return api.post('/journals', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['journals'] });
      navigate('/finance/journals');
    },
  });

  const onSubmit = (data: JournalFormData) => {
    if (!isBalanced) {
      alert('Debit dan Kredit harus balance!');
      return;
    }
    mutation.mutate(data);
  };

  // Filter out header accounts
  const selectableAccounts = accounts?.filter(acc => !acc.is_header) || [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/finance/journals')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Buat Jurnal Baru</h1>
            <p className="text-gray-600">Input transaksi jurnal umum</p>
          </div>
        </div>
        <Button 
          type="submit" 
          variant="primary" 
          loading={mutation.isPending}
          disabled={!isBalanced}
        >
          <Save className="w-4 h-4 mr-2" />
          Simpan Jurnal
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          {/* Header Info */}
          <Card title="Informasi Jurnal">
            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Tanggal Jurnal"
                type="date"
                {...register('journal_date')}
                error={errors.journal_date?.message}
                required
              />

              <Input
                label="No. Referensi"
                placeholder="e.g., INV-001, PV-123"
                {...register('reference_no')}
                helperText="Opsional - untuk cross reference"
              />

              <div className="col-span-2">
                <Input
                  label="Deskripsi Transaksi"
                  placeholder="e.g., Pembayaran gaji bulan November 2024"
                  {...register('description')}
                  error={errors.description?.message}
                  required
                />
              </div>
            </div>
          </Card>

          {/* Journal Items */}
          <Card
            title="Baris Jurnal"
            action={
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => append({ account_id: '', description: '', debit: 0, credit: 0 })}
              >
                <Plus className="w-4 h-4 mr-2" />
                Tambah Baris
              </Button>
            }
          >
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b">
                  <tr>
                    <th className="px-3 py-2 text-left text-xs font-medium text-gray-700">
                      Akun <span className="text-red-500">*</span>
                    </th>
                    <th className="px-3 py-2 text-left text-xs font-medium text-gray-700">
                      Deskripsi <span className="text-red-500">*</span>
                    </th>
                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-700">
                      Debit
                    </th>
                    <th className="px-3 py-2 text-right text-xs font-medium text-gray-700">
                      Kredit
                    </th>
                    <th className="px-3 py-2 text-center text-xs font-medium text-gray-700">
                      Aksi
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {fields.map((field, index) => (
                    <tr key={field.id} className="hover:bg-gray-50">
                      <td className="px-3 py-2">
                        <select
                          {...register(`items.${index}.account_id`)}
                          className="w-40 px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                        >
                          <option value="">Pilih Akun</option>
                          {selectableAccounts.map((acc) => (
                            <option key={acc.id} value={acc.id}>
                              {acc.code} - {acc.name}
                            </option>
                          ))}
                        </select>
                        {errors.items?.[index]?.account_id && (
                          <p className="text-xs text-red-600 mt-1">
                            {errors.items[index]?.account_id?.message}
                          </p>
                        )}
                      </td>
                      <td className="px-3 py-2">
                        <input
                          {...register(`items.${index}.description`)}
                          className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                          placeholder="Deskripsi baris"
                        />
                        {errors.items?.[index]?.description && (
                          <p className="text-xs text-red-600 mt-1">
                            {errors.items[index]?.description?.message}
                          </p>
                        )}
                      </td>
                      <td className="px-3 py-2">
                        <input
                          type="number"
                          {...register(`items.${index}.debit`, { valueAsNumber: true })}
                          className="w-32 px-2 py-1 text-sm text-right border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                          placeholder="0"
                        />
                      </td>
                      <td className="px-3 py-2">
                        <input
                          type="number"
                          {...register(`items.${index}.credit`, { valueAsNumber: true })}
                          className="w-32 px-2 py-1 text-sm text-right border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                          placeholder="0"
                        />
                      </td>
                      <td className="px-3 py-2 text-center">
                        {fields.length > 2 && (
                          <Button
                            type="button"
                            variant="danger"
                            size="sm"
                            onClick={() => remove(index)}
                          >
                            <Trash2 className="w-3 h-3" />
                          </Button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
                <tfoot className="bg-gray-100 font-bold border-t-2 border-gray-300">
                  <tr>
                    <td colSpan={2} className="px-3 py-3 text-right text-sm">
                      Total:
                    </td>
                    <td className="px-3 py-3 text-right text-base">
                      {formatCurrency(totalDebit)}
                    </td>
                    <td className="px-3 py-3 text-right text-base">
                      {formatCurrency(totalCredit)}
                    </td>
                    <td></td>
                  </tr>
                </tfoot>
              </table>
            </div>

            {errors.items && typeof errors.items === 'object' && 'message' in errors.items && (
              <p className="mt-2 text-sm text-red-600">{errors.items.message as string}</p>
            )}
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Balance Check */}
          <Card title="Status Balance">
            <div className="space-y-4">
              <div className="p-4 bg-blue-50 rounded-lg">
                <div className="text-sm text-gray-600 mb-1">Total Debit</div>
                <div className="text-2xl font-bold text-blue-900">
                  {formatCurrency(totalDebit)}
                </div>
              </div>

              <div className="p-4 bg-green-50 rounded-lg">
                <div className="text-sm text-gray-600 mb-1">Total Kredit</div>
                <div className="text-2xl font-bold text-green-900">
                  {formatCurrency(totalCredit)}
                </div>
              </div>

              <div className={`p-4 rounded-lg border-2 ${
                isBalanced 
                  ? 'bg-green-50 border-green-300' 
                  : 'bg-red-50 border-red-300'
              }`}>
                <div className="flex items-center gap-2 mb-2">
                  {isBalanced ? (
                    <CheckCircle className="w-5 h-5 text-green-600" />
                  ) : (
                    <AlertCircle className="w-5 h-5 text-red-600" />
                  )}
                  <div className={`font-medium ${
                    isBalanced ? 'text-green-900' : 'text-red-900'
                  }`}>
                    {isBalanced ? 'Balance ✓' : 'Tidak Balance'}
                  </div>
                </div>
                <div className={`text-sm ${
                  isBalanced ? 'text-green-700' : 'text-red-700'
                }`}>
                  {isBalanced 
                    ? 'Debit dan Kredit sudah balance' 
                    : `Selisih: ${formatCurrency(Math.abs(totalDebit - totalCredit))}`}
                </div>
              </div>
            </div>
          </Card>

          {/* Tips */}
          <Card title="Panduan Jurnal">
            <div className="text-sm text-gray-600 space-y-3">
              <div>
                <p className="font-medium text-gray-900 mb-1">Aturan Debit-Kredit:</p>
                <ul className="list-disc list-inside space-y-1 text-xs">
                  <li>Aset: Debit (+) Kredit (-)</li>
                  <li>Kewajiban: Debit (-) Kredit (+)</li>
                  <li>Ekuitas: Debit (-) Kredit (+)</li>
                  <li>Pendapatan: Debit (-) Kredit (+)</li>
                  <li>Beban: Debit (+) Kredit (-)</li>
                </ul>
              </div>

              <div>
                <p className="font-medium text-gray-900 mb-1">Tips:</p>
                <ul className="space-y-1 text-xs">
                  <li>• Setiap baris hanya debit ATAU kredit</li>
                  <li>• Minimal 2 baris (1 debit, 1 kredit)</li>
                  <li>• Total debit = total kredit</li>
                  <li>• Deskripsi harus jelas</li>
                </ul>
              </div>

              <div className="pt-3 border-t">
                <p className="font-medium text-gray-900 mb-1">Contoh:</p>
                <div className="text-xs bg-gray-50 p-2 rounded">
                  <p>Bayar gaji karyawan Rp 5.000.000:</p>
                  <p className="mt-1">D: Beban Gaji 5.000.000</p>
                  <p>K: Kas 5.000.000</p>
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
