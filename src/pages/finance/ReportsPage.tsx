import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  FileText, 
  Download, 
  Printer,
  Calendar,
  TrendingUp,
  TrendingDown,
  DollarSign
} from 'lucide-react';
import { formatCurrency, formatDate } from '@/utils/format';

type ReportType = 'trial-balance' | 'balance-sheet' | 'income-statement';

export default function ReportsPage() {
  const [reportType, setReportType] = useState<ReportType>('trial-balance');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

  const { data: reportData, isLoading, refetch } = useQuery({
    queryKey: ['financial-report', reportType, startDate, endDate],
    queryFn: async () => {
      const response = await api.get(`/reports/${reportType}`, {
        params: { start_date: startDate, end_date: endDate },
      });
      return response.data.data;
    },
    enabled: !!startDate && !!endDate,
  });

  const reports = [
    {
      id: 'trial-balance',
      name: 'Neraca Saldo (Trial Balance)',
      description: 'Ringkasan saldo debit dan kredit semua akun',
      icon: FileText,
      color: 'blue',
    },
    {
      id: 'balance-sheet',
      name: 'Neraca (Balance Sheet)',
      description: 'Laporan posisi keuangan (Aset, Kewajiban, Ekuitas)',
      icon: TrendingUp,
      color: 'green',
    },
    {
      id: 'income-statement',
      name: 'Laba Rugi (Income Statement)',
      description: 'Laporan pendapatan dan beban',
      icon: DollarSign,
      color: 'purple',
    },
  ];

  const handleGenerate = () => {
    refetch();
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Laporan Keuangan</h1>
        <p className="text-gray-600 mt-1">Generate dan lihat laporan keuangan</p>
      </div>

      {/* Report Selection */}
      <Card title="Pilih Jenis Laporan">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {reports.map((report) => {
            const Icon = report.icon;
            const isSelected = reportType === report.id;
            
            return (
              <button
                key={report.id}
                onClick={() => setReportType(report.id as ReportType)}
                className={`p-4 border-2 rounded-lg text-left transition-all ${
                  isSelected
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-gray-300'
                }`}
              >
                <div className="flex items-start gap-3">
                  <div className={`p-2 rounded-lg ${
                    isSelected ? 'bg-blue-100' : 'bg-gray-100'
                  }`}>
                    <Icon className={`w-6 h-6 ${
                      isSelected ? 'text-blue-600' : 'text-gray-600'
                    }`} />
                  </div>
                  <div className="flex-1">
                    <div className={`font-medium ${
                      isSelected ? 'text-blue-900' : 'text-gray-900'
                    }`}>
                      {report.name}
                    </div>
                    <div className="text-sm text-gray-500 mt-1">
                      {report.description}
                    </div>
                  </div>
                </div>
              </button>
            );
          })}
        </div>
      </Card>

      {/* Date Range & Generate */}
      <Card title="Periode Laporan">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Tanggal Mulai
            </label>
            <input
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Tanggal Akhir
            </label>
            <input
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="flex items-end">
            <Button
              variant="primary"
              onClick={handleGenerate}
              disabled={!startDate || !endDate}
              loading={isLoading}
            >
              <FileText className="w-4 h-4 mr-2" />
              Generate Laporan
            </Button>
          </div>
        </div>
      </Card>

      {/* Report Display */}
      {reportData && (
        <Card
          title="Hasil Laporan"
          action={
            <div className="flex gap-2">
              <Button variant="secondary" size="sm">
                <Printer className="w-4 h-4 mr-2" />
                Cetak
              </Button>
              <Button variant="secondary" size="sm">
                <Download className="w-4 h-4 mr-2" />
                Export Excel
              </Button>
            </div>
          }
        >
          <div className="space-y-4">
            {/* Report Header */}
            <div className="text-center pb-4 border-b">
              <h2 className="text-xl font-bold text-gray-900">
                {reports.find(r => r.id === reportType)?.name}
              </h2>
              <p className="text-gray-600 mt-1">
                Periode: {formatDate(startDate)} s/d {formatDate(endDate)}
              </p>
            </div>

            {/* Report Content - Trial Balance Example */}
            {reportType === 'trial-balance' && (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">
                        Kode
                      </th>
                      <th className="px-4 py-3 text-left text-sm font-medium text-gray-900">
                        Nama Akun
                      </th>
                      <th className="px-4 py-3 text-right text-sm font-medium text-gray-900">
                        Debit
                      </th>
                      <th className="px-4 py-3 text-right text-sm font-medium text-gray-900">
                        Kredit
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {/* Sample data - replace with actual data */}
                    <tr>
                      <td className="px-4 py-2 text-sm font-mono text-gray-600">1-1000</td>
                      <td className="px-4 py-2 text-sm text-gray-900">Kas</td>
                      <td className="px-4 py-2 text-sm text-right">{formatCurrency(50000000)}</td>
                      <td className="px-4 py-2 text-sm text-right">-</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2 text-sm font-mono text-gray-600">4-1000</td>
                      <td className="px-4 py-2 text-sm text-gray-900">Pendapatan SPP</td>
                      <td className="px-4 py-2 text-sm text-right">-</td>
                      <td className="px-4 py-2 text-sm text-right">{formatCurrency(50000000)}</td>
                    </tr>
                  </tbody>
                  <tfoot className="bg-gray-50 font-bold">
                    <tr>
                      <td colSpan={2} className="px-4 py-3 text-sm text-gray-900">
                        Total
                      </td>
                      <td className="px-4 py-3 text-sm text-right">{formatCurrency(50000000)}</td>
                      <td className="px-4 py-3 text-sm text-right">{formatCurrency(50000000)}</td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            )}

            {/* Balance Sheet Example */}
            {reportType === 'balance-sheet' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">ASET</h3>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Aset Lancar</span>
                      <span className="font-medium">{formatCurrency(100000000)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Aset Tetap</span>
                      <span className="font-medium">{formatCurrency(500000000)}</span>
                    </div>
                    <div className="flex justify-between font-bold border-t pt-2">
                      <span>Total Aset</span>
                      <span>{formatCurrency(600000000)}</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">KEWAJIBAN & EKUITAS</h3>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Kewajiban</span>
                      <span className="font-medium">{formatCurrency(100000000)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Ekuitas</span>
                      <span className="font-medium">{formatCurrency(500000000)}</span>
                    </div>
                    <div className="flex justify-between font-bold border-t pt-2">
                      <span>Total Kewajiban & Ekuitas</span>
                      <span>{formatCurrency(600000000)}</span>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Income Statement Example */}
            {reportType === 'income-statement' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">PENDAPATAN</h3>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Pendapatan SPP</span>
                      <span className="font-medium">{formatCurrency(400000000)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Pendapatan Lainnya</span>
                      <span className="font-medium">{formatCurrency(50000000)}</span>
                    </div>
                    <div className="flex justify-between font-bold border-t pt-2">
                      <span>Total Pendapatan</span>
                      <span>{formatCurrency(450000000)}</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">BEBAN</h3>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Beban Gaji</span>
                      <span className="font-medium">{formatCurrency(250000000)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Beban Operasional</span>
                      <span className="font-medium">{formatCurrency(70000000)}</span>
                    </div>
                    <div className="flex justify-between font-bold border-t pt-2">
                      <span>Total Beban</span>
                      <span>{formatCurrency(320000000)}</span>
                    </div>
                  </div>
                </div>

                <div className="bg-green-50 p-4 rounded-lg">
                  <div className="flex justify-between items-center">
                    <span className="text-lg font-bold text-green-900">LABA BERSIH</span>
                    <span className="text-2xl font-bold text-green-600">
                      {formatCurrency(130000000)}
                    </span>
                  </div>
                </div>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* Empty State */}
      {!reportData && !isLoading && (
        <Card>
          <div className="text-center py-12">
            <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">
              Pilih periode dan klik "Generate Laporan" untuk melihat hasil
            </p>
          </div>
        </Card>
      )}
    </div>
  );
}
