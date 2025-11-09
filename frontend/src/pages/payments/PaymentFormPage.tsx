import { useNavigate, useSearchParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Invoice, ApiResponse } from '@/types';
import { ArrowLeft, Save, DollarSign } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useEffect } from 'react';

const paymentSchema = z.object({
  invoice_id: z.string().min(1, 'Invoice harus dipilih'),
  payment_date: z.string().min(1, 'Tanggal pembayaran harus diisi'),
  amount: z.number().min(1, 'Jumlah pembayaran harus lebih dari 0'),
  payment_method: z.enum(['cash', 'transfer', 'card', 'ewallet']),
  reference_no: z.string().optional(),
  notes: z.string().optional(),
});

type PaymentFormData = z.infer<typeof paymentSchema>;

const paymentMethods = [
  { value: 'cash', label: 'Tunai' },
  { value: 'transfer', label: 'Transfer Bank' },
  { value: 'card', label: 'Kartu Debit/Kredit' },
  { value: 'ewallet', label: 'E-Wallet (OVO, GoPay, Dana)' },
];

export default function PaymentFormPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [searchParams] = useSearchParams();
  const invoiceIdFromUrl = searchParams.get('invoice_id');

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<PaymentFormData>({
    resolver: zodResolver(paymentSchema),
    defaultValues: {
      invoice_id: invoiceIdFromUrl || '',
      payment_date: new Date().toISOString().split('T')[0],
      payment_method: 'cash',
      amount: 0,
    },
  });

  const selectedInvoiceId = watch('invoice_id');
  const paymentMethod = watch('payment_method');

  // Fetch invoices for selection
  const { data: invoices } = useQuery({
    queryKey: ['invoices-unpaid'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Invoice[]>>('/invoices', {
        params: { status: 'unpaid,partial' },
      });
      return response.data.data;
    },
  });

  // Fetch selected invoice detail
  const { data: selectedInvoice } = useQuery({
    queryKey: ['invoice', selectedInvoiceId],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Invoice>>(`/invoices/${selectedInvoiceId}`);
      return response.data.data;
    },
    enabled: !!selectedInvoiceId,
  });

  // Auto-fill amount when invoice selected
  useEffect(() => {
    if (selectedInvoice) {
      const remaining = selectedInvoice.total_amount - selectedInvoice.paid_amount;
      setValue('amount', remaining);
    }
  }, [selectedInvoice, setValue]);

  const mutation = useMutation({
    mutationFn: async (data: PaymentFormData) => {
      return api.post('/payments', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      navigate('/payments');
    },
  });

  const onSubmit = (data: PaymentFormData) => {
    mutation.mutate(data);
  };

  const remainingAmount = selectedInvoice
    ? selectedInvoice.total_amount - selectedInvoice.paid_amount
    : 0;

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            type="button"
            variant="ghost"
            onClick={() => navigate('/payments')}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Input Pembayaran</h1>
            <p className="text-gray-600">Catat pembayaran dari siswa</p>
          </div>
        </div>
        <Button
          type="submit"
          variant="primary"
          loading={mutation.isPending}
        >
          <Save className="w-4 h-4 mr-2" />
          Simpan Pembayaran
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Payment Form */}
        <div className="lg:col-span-2">
          <Card title="Informasi Pembayaran">
            <div className="space-y-4">
              {/* Invoice Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Invoice <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('invoice_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  disabled={!!invoiceIdFromUrl}
                >
                  <option value="">Pilih Invoice</option>
                  {invoices?.map((invoice) => (
                    <option key={invoice.id} value={invoice.id}>
                      {invoice.invoice_number} - {invoice.student_name} - {formatCurrency(invoice.total_amount - invoice.paid_amount)}
                    </option>
                  ))}
                </select>
                {errors.invoice_id && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.invoice_id.message}</p>
                )}
              </div>

              {/* Payment Date */}
              <Input
                label="Tanggal Pembayaran"
                type="date"
                {...register('payment_date')}
                error={errors.payment_date?.message}
                required
              />

              {/* Amount */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Jumlah Pembayaran <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <span className="text-gray-500">Rp</span>
                  </div>
                  <input
                    type="number"
                    {...register('amount', { valueAsNumber: true })}
                    className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="0"
                  />
                </div>
                {errors.amount && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.amount.message}</p>
                )}
                {selectedInvoice && (
                  <p className="mt-1.5 text-sm text-gray-500">
                    Sisa tagihan: {formatCurrency(remainingAmount)}
                  </p>
                )}
              </div>

              {/* Payment Method */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Metode Pembayaran <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('payment_method')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  {paymentMethods.map((method) => (
                    <option key={method.value} value={method.value}>
                      {method.label}
                    </option>
                  ))}
                </select>
              </div>

              {/* Reference Number */}
              {paymentMethod !== 'cash' && (
                <Input
                  label="Nomor Referensi"
                  placeholder="No. transfer, no. kartu, dll"
                  {...register('reference_no')}
                  helperText="Opsional - untuk rekonsiliasi"
                />
              )}

              {/* Notes */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Catatan
                </label>
                <textarea
                  {...register('notes')}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Catatan tambahan (opsional)"
                />
              </div>
            </div>
          </Card>
        </div>

        {/* Sidebar - Invoice Summary */}
        <div>
          {selectedInvoice ? (
            <Card title="Ringkasan Invoice">
              <div className="space-y-4">
                <div>
                  <div className="text-sm text-gray-500">No. Invoice</div>
                  <div className="text-base font-medium text-gray-900 mt-1">
                    {selectedInvoice.invoice_number}
                  </div>
                </div>

                <div>
                  <div className="text-sm text-gray-500">Nama Siswa</div>
                  <div className="text-base text-gray-900 mt-1">
                    {selectedInvoice.student_name}
                  </div>
                </div>

                <div className="border-t pt-4">
                  <div className="flex justify-between mb-2">
                    <span className="text-sm text-gray-600">Total Tagihan</span>
                    <span className="text-sm font-medium">
                      {formatCurrency(selectedInvoice.total_amount)}
                    </span>
                  </div>
                  <div className="flex justify-between mb-2">
                    <span className="text-sm text-gray-600">Terbayar</span>
                    <span className="text-sm font-medium text-green-600">
                      {formatCurrency(selectedInvoice.paid_amount)}
                    </span>
                  </div>
                  <div className="flex justify-between pt-2 border-t">
                    <span className="text-base font-medium text-gray-900">Sisa</span>
                    <span className="text-lg font-bold text-red-600">
                      {formatCurrency(remainingAmount)}
                    </span>
                  </div>
                </div>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  className="w-full"
                  onClick={() => navigate(`/invoices/${selectedInvoice.id}`)}
                >
                  Lihat Detail Invoice
                </Button>
              </div>
            </Card>
          ) : (
            <Card>
              <div className="text-center py-8">
                <DollarSign className="w-12 h-12 text-gray-300 mx-auto mb-3" />
                <p className="text-gray-500">Pilih invoice untuk melihat detail</p>
              </div>
            </Card>
          )}
        </div>
      </div>
    </form>
  );
}
