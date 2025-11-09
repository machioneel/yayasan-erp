import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Payment, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  ArrowLeft, 
  Printer, 
  Download,
  CheckCircle,
  DollarSign,
  Calendar,
  CreditCard,
  FileText,
  User
} from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

export default function PaymentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: payment, isLoading } = useQuery({
    queryKey: ['payment', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Payment>>(`/payments/${id}`);
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

  if (!payment) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Pembayaran tidak ditemukan</p>
      </div>
    );
  }

  const getMethodBadge = (method: string) => {
    switch (method) {
      case 'cash': return { variant: 'success' as const, label: 'Tunai', icon: DollarSign };
      case 'transfer': return { variant: 'info' as const, label: 'Transfer Bank', icon: CreditCard };
      case 'card': return { variant: 'warning' as const, label: 'Kartu Debit/Kredit', icon: CreditCard };
      case 'ewallet': return { variant: 'info' as const, label: 'E-Wallet', icon: DollarSign };
      default: return { variant: 'default' as const, label: method, icon: DollarSign };
    }
  };

  const methodInfo = getMethodBadge(payment.payment_method);
  const MethodIcon = methodInfo.icon;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/payments')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              Pembayaran #{payment.payment_number}
            </h1>
            <p className="text-gray-600">Detail transaksi pembayaran</p>
          </div>
        </div>

        <div className="flex gap-2">
          <Button 
            variant="secondary"
            onClick={() => window.print()}
          >
            <Printer className="w-4 h-4 mr-2" />
            Cetak
          </Button>
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            Download PDF
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Payment Status */}
          <Card>
            <div className="flex items-center justify-center py-8">
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-20 h-20 bg-green-100 rounded-full mb-4">
                  <CheckCircle className="w-12 h-12 text-green-600" />
                </div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">
                  Pembayaran Berhasil
                </h2>
                <p className="text-gray-600">
                  Transaksi telah diverifikasi dan tercatat
                </p>
              </div>
            </div>
          </Card>

          {/* Payment Information */}
          <Card title="Informasi Pembayaran">
            <div className="space-y-6">
              {/* Amount */}
              <div className="p-6 bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg border border-green-200">
                <div className="text-sm text-green-700 mb-2">Jumlah Pembayaran</div>
                <div className="text-4xl font-bold text-green-900">
                  {formatCurrency(payment.amount)}
                </div>
              </div>

              {/* Details Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="text-sm font-medium text-gray-500">No. Pembayaran</label>
                  <p className="text-base font-mono font-medium text-gray-900 mt-1">
                    {payment.payment_number}
                  </p>
                </div>

                <div>
                  <label className="text-sm font-medium text-gray-500">Tanggal Pembayaran</label>
                  <div className="flex items-center gap-2 mt-1">
                    <Calendar className="w-4 h-4 text-gray-400" />
                    <p className="text-base text-gray-900">
                      {formatDate(payment.payment_date)}
                    </p>
                  </div>
                </div>

                <div>
                  <label className="text-sm font-medium text-gray-500">Metode Pembayaran</label>
                  <div className="mt-2">
                    <Badge variant={methodInfo.variant}>
                      <MethodIcon className="w-3 h-3 mr-1" />
                      {methodInfo.label}
                    </Badge>
                  </div>
                </div>

                {payment.reference_no && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">No. Referensi</label>
                    <p className="text-base font-mono text-gray-900 mt-1">
                      {payment.reference_no}
                    </p>
                  </div>
                )}

                {payment.notes && (
                  <div className="md:col-span-2">
                    <label className="text-sm font-medium text-gray-500">Catatan</label>
                    <p className="text-base text-gray-900 mt-1">
                      {payment.notes}
                    </p>
                  </div>
                )}
              </div>
            </div>
          </Card>

          {/* Invoice Information */}
          <Card title="Informasi Invoice">
            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">No. Invoice</label>
                <div className="flex items-center gap-2 mt-1">
                  <FileText className="w-4 h-4 text-gray-400" />
                  <button
                    onClick={() => navigate(`/invoices/${payment.invoice_id}`)}
                    className="text-base font-medium text-blue-600 hover:text-blue-700"
                  >
                    {payment.invoice_number}
                  </button>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Nama Siswa</label>
                <div className="flex items-center gap-2 mt-1">
                  <User className="w-4 h-4 text-gray-400" />
                  <button
                    onClick={() => navigate(`/students/${payment.student_id}`)}
                    className="text-base font-medium text-blue-600 hover:text-blue-700"
                  >
                    {payment.student_name}
                  </button>
                </div>
              </div>
            </div>
          </Card>

          {/* Payment Proof */}
          {payment.attachment_url && (
            <Card title="Bukti Pembayaran">
              <div className="p-4 border border-gray-200 rounded-lg">
                <img 
                  src={payment.attachment_url} 
                  alt="Bukti Pembayaran"
                  className="w-full rounded"
                />
              </div>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Receipt Card */}
          <Card title="Kwitansi">
            <div className="space-y-4">
              <div className="p-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
                <div className="text-center">
                  <div className="text-sm text-gray-600 mb-2">Kwitansi Resmi</div>
                  <div className="text-lg font-bold text-gray-900 mb-4">
                    {payment.payment_number}
                  </div>
                  <Button
                    variant="primary"
                    className="w-full"
                    onClick={() => window.print()}
                  >
                    <Printer className="w-4 h-4 mr-2" />
                    Cetak Kwitansi
                  </Button>
                </div>
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Diterima dari:</span>
                  <span className="font-medium text-gray-900">{payment.student_name}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Untuk pembayaran:</span>
                  <span className="font-medium text-gray-900">Invoice</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Sejumlah:</span>
                  <span className="font-medium text-gray-900">{formatCurrency(payment.amount)}</span>
                </div>
              </div>
            </div>
          </Card>

          {/* Transaction Info */}
          <Card title="Info Transaksi">
            <div className="space-y-3 text-sm">
              <div>
                <div className="text-gray-500">Diinput oleh</div>
                <div className="font-medium text-gray-900 mt-1">
                  {payment.created_by_name || 'System'}
                </div>
              </div>

              <div>
                <div className="text-gray-500">Waktu Input</div>
                <div className="font-medium text-gray-900 mt-1">
                  {formatDate(payment.created_at)}
                </div>
              </div>

              {payment.approved_by && (
                <>
                  <div className="border-t pt-3">
                    <div className="text-gray-500">Disetujui oleh</div>
                    <div className="font-medium text-gray-900 mt-1">
                      {payment.approved_by_name}
                    </div>
                  </div>

                  <div>
                    <div className="text-gray-500">Waktu Approval</div>
                    <div className="font-medium text-gray-900 mt-1">
                      {formatDate(payment.approved_at)}
                    </div>
                  </div>
                </>
              )}
            </div>
          </Card>

          {/* Quick Actions */}
          <Card title="Aksi Lainnya">
            <div className="space-y-2">
              <Button
                variant="secondary"
                className="w-full justify-start"
                onClick={() => navigate(`/invoices/${payment.invoice_id}`)}
              >
                <FileText className="w-4 h-4 mr-2" />
                Lihat Invoice
              </Button>
              <Button
                variant="secondary"
                className="w-full justify-start"
                onClick={() => navigate(`/students/${payment.student_id}`)}
              >
                <User className="w-4 h-4 mr-2" />
                Lihat Data Siswa
              </Button>
              <Button
                variant="secondary"
                className="w-full justify-start"
              >
                <Download className="w-4 h-4 mr-2" />
                Download Kwitansi
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
