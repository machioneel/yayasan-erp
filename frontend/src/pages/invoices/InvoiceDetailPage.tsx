import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Invoice, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Modal } from '@/components/common/Modal';
import { 
  ArrowLeft, 
  Printer, 
  Download, 
  DollarSign,
  XCircle,
  CheckCircle,
  Mail
} from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';
import { useState } from 'react';

export default function InvoiceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [showCancelModal, setShowCancelModal] = useState(false);

  const { data: invoice, isLoading } = useQuery({
    queryKey: ['invoice', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Invoice>>(`/invoices/${id}`);
      return response.data.data;
    },
  });

  const cancelMutation = useMutation({
    mutationFn: async () => {
      return api.post(`/invoices/${id}/cancel`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoice', id] });
      setShowCancelModal(false);
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!invoice) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Invoice tidak ditemukan</p>
      </div>
    );
  }

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'paid': return 'success';
      case 'partial': return 'warning';
      case 'unpaid': return 'danger';
      case 'cancelled': return 'default';
      default: return 'default';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'paid': return 'Lunas';
      case 'partial': return 'Sebagian';
      case 'unpaid': return 'Belum Bayar';
      case 'cancelled': return 'Dibatalkan';
      default: return status;
    }
  };

  const remainingAmount = invoice.total_amount - invoice.paid_amount;
  const isOverdue = new Date(invoice.due_date) < new Date() && invoice.status !== 'paid';

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            onClick={() => navigate('/invoices')}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              Invoice #{invoice.invoice_number}
            </h1>
            <p className="text-gray-600">
              {invoice.student_name}
            </p>
          </div>
        </div>

        <div className="flex gap-2">
          {invoice.status !== 'cancelled' && invoice.status !== 'paid' && (
            <Button
              variant="success"
              onClick={() => navigate(`/payments/new?invoice_id=${invoice.id}`)}
            >
              <DollarSign className="w-4 h-4 mr-2" />
              Bayar
            </Button>
          )}
          <Button variant="secondary">
            <Mail className="w-4 h-4 mr-2" />
            Kirim Email
          </Button>
          <Button 
            variant="secondary"
            onClick={() => window.print()}
          >
            <Printer className="w-4 h-4 mr-2" />
            Cetak
          </Button>
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            PDF
          </Button>
          {invoice.status !== 'cancelled' && invoice.status !== 'paid' && (
            <Button
              variant="danger"
              onClick={() => setShowCancelModal(true)}
            >
              <XCircle className="w-4 h-4 mr-2" />
              Batalkan
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Invoice Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Invoice Info */}
          <Card>
            <div className="flex items-center justify-between mb-6">
              <div>
                <div className="text-sm text-gray-500">Status</div>
                <div className="mt-1 flex items-center gap-2">
                  <Badge variant={getStatusVariant(invoice.status)}>
                    {getStatusLabel(invoice.status)}
                  </Badge>
                  {isOverdue && (
                    <Badge variant="danger">Jatuh Tempo</Badge>
                  )}
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm text-gray-500">Tanggal Invoice</div>
                <div className="text-base font-medium text-gray-900 mt-1">
                  {formatDate(invoice.invoice_date)}
                </div>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-6 mb-6">
              <div>
                <div className="text-sm text-gray-500">Jatuh Tempo</div>
                <div className="text-base font-medium text-gray-900 mt-1">
                  {formatDate(invoice.due_date)}
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-500">Cabang</div>
                <div className="text-base text-gray-900 mt-1">
                  {invoice.branch_id}
                </div>
              </div>
            </div>

            {invoice.description && (
              <div>
                <div className="text-sm text-gray-500">Keterangan</div>
                <div className="text-base text-gray-900 mt-1">
                  {invoice.description}
                </div>
              </div>
            )}
          </Card>

          {/* Invoice Items */}
          <Card title="Detail Item">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Deskripsi
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Qty
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Harga Satuan
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Jumlah
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {invoice.items?.map((item, index) => (
                    <tr key={index}>
                      <td className="px-4 py-3 text-sm text-gray-900">
                        {item.description}
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900 text-right">
                        {item.quantity}
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900 text-right">
                        {formatCurrency(item.unit_price)}
                      </td>
                      <td className="px-4 py-3 text-sm font-medium text-gray-900 text-right">
                        {formatCurrency(item.amount)}
                      </td>
                    </tr>
                  ))}
                </tbody>
                <tfoot className="border-t-2 border-gray-300">
                  <tr>
                    <td colSpan={3} className="px-4 py-3 text-sm font-medium text-gray-900 text-right">
                      Total
                    </td>
                    <td className="px-4 py-3 text-base font-bold text-gray-900 text-right">
                      {formatCurrency(invoice.total_amount)}
                    </td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Payment Summary */}
          <Card title="Ringkasan Pembayaran">
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Total Tagihan</span>
                <span className="text-base font-medium text-gray-900">
                  {formatCurrency(invoice.total_amount)}
                </span>
              </div>

              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Terbayar</span>
                <span className="text-base font-medium text-green-600">
                  {formatCurrency(invoice.paid_amount)}
                </span>
              </div>

              <div className="border-t pt-4">
                <div className="flex justify-between items-center">
                  <span className="text-base font-medium text-gray-900">Sisa Tagihan</span>
                  <span className={`text-xl font-bold ${remainingAmount > 0 ? 'text-red-600' : 'text-green-600'}`}>
                    {formatCurrency(remainingAmount)}
                  </span>
                </div>
              </div>

              {invoice.status !== 'cancelled' && remainingAmount > 0 && (
                <Button
                  variant="primary"
                  className="w-full"
                  onClick={() => navigate(`/payments/new?invoice_id=${invoice.id}`)}
                >
                  <DollarSign className="w-4 h-4 mr-2" />
                  Bayar Sekarang
                </Button>
              )}
            </div>
          </Card>

          {/* Student Info */}
          <Card title="Informasi Siswa">
            <div className="space-y-3">
              <div>
                <div className="text-sm text-gray-500">Nama</div>
                <div className="text-base font-medium text-gray-900 mt-1">
                  {invoice.student_name}
                </div>
              </div>
              <Button
                variant="secondary"
                size="sm"
                className="w-full"
                onClick={() => navigate(`/students/${invoice.student_id}`)}
              >
                Lihat Detail Siswa
              </Button>
            </div>
          </Card>

          {/* Timeline / History */}
          <Card title="Riwayat">
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-blue-600 rounded-full mt-2"></div>
                <div className="flex-1">
                  <div className="text-sm font-medium text-gray-900">
                    Invoice Dibuat
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    {formatDate(invoice.created_at)}
                  </div>
                </div>
              </div>

              {invoice.paid_amount > 0 && (
                <div className="flex items-start gap-3">
                  <div className="w-2 h-2 bg-green-600 rounded-full mt-2"></div>
                  <div className="flex-1">
                    <div className="text-sm font-medium text-gray-900">
                      Pembayaran Diterima
                    </div>
                    <div className="text-xs text-gray-500 mt-1">
                      {formatCurrency(invoice.paid_amount)}
                    </div>
                  </div>
                </div>
              )}

              {invoice.status === 'paid' && (
                <div className="flex items-start gap-3">
                  <CheckCircle className="w-5 h-5 text-green-600 mt-0.5" />
                  <div className="flex-1">
                    <div className="text-sm font-medium text-gray-900">
                      Lunas
                    </div>
                  </div>
                </div>
              )}
            </div>
          </Card>
        </div>
      </div>

      {/* Cancel Modal */}
      <Modal
        isOpen={showCancelModal}
        onClose={() => setShowCancelModal(false)}
        title="Batalkan Invoice"
        footer={
          <div className="flex justify-end gap-3">
            <Button
              variant="secondary"
              onClick={() => setShowCancelModal(false)}
            >
              Batal
            </Button>
            <Button
              variant="danger"
              onClick={() => cancelMutation.mutate()}
              loading={cancelMutation.isPending}
            >
              Ya, Batalkan
            </Button>
          </div>
        }
      >
        <p className="text-gray-700">
          Apakah Anda yakin ingin membatalkan invoice ini?
          Tindakan ini tidak dapat dibatalkan.
        </p>
      </Modal>
    </div>
  );
}
