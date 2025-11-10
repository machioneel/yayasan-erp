import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Student, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  ArrowLeft, 
  Edit, 
  Trash2, 
  User, 
  Phone, 
  Mail, 
  MapPin, 
  Calendar,
  Users,
  FileText,
  DollarSign
} from 'lucide-react';
import { formatDate } from '@/utils/format';

export default function StudentDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: student, isLoading } = useQuery({
    queryKey: ['student', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Student>>(`/students/${id}`);
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

  if (!student) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Siswa tidak ditemukan</p>
      </div>
    );
  }

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'inactive': return 'default';
      case 'graduated': return 'info';
      case 'dropped': return 'danger';
      default: return 'default';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            onClick={() => navigate('/students')}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{student.full_name}</h1>
            <p className="text-gray-600">No. Registrasi: {student.registration_number}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button
            variant="primary"
            onClick={() => navigate(`/students/${id}/edit`)}
          >
            <Edit className="w-4 h-4 mr-2" />
            Edit
          </Button>
          <Button variant="danger">
            <Trash2 className="w-4 h-4 mr-2" />
            Hapus
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Personal Information */}
          <Card title="Informasi Pribadi">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Nama Lengkap</label>
                <p className="text-base font-medium text-gray-900 mt-1">{student.full_name}</p>
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-500">Nama Panggilan</label>
                <p className="text-base text-gray-900 mt-1">{student.nick_name || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">NISN</label>
                <p className="text-base text-gray-900 mt-1">{student.nisn || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Jenis Kelamin</label>
                <p className="text-base text-gray-900 mt-1">
                  {student.gender === 'male' ? 'Laki-laki' : 'Perempuan'}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tempat, Tanggal Lahir</label>
                <p className="text-base text-gray-900 mt-1">
                  {student.birth_place && student.birth_date
                    ? `${student.birth_place}, ${formatDate(student.birth_date)}`
                    : '-'}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Agama</label>
                <p className="text-base text-gray-900 mt-1">{student.religion || '-'}</p>
              </div>
            </div>
          </Card>

          {/* Contact Information */}
          <Card title="Informasi Kontak">
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <Mail className="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <label className="text-sm font-medium text-gray-500">Email</label>
                  <p className="text-base text-gray-900">{student.email || '-'}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <Phone className="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <label className="text-sm font-medium text-gray-500">Telepon</label>
                  <p className="text-base text-gray-900">{student.phone || '-'}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <label className="text-sm font-medium text-gray-500">Alamat</label>
                  <p className="text-base text-gray-900">
                    {student.address || '-'}
                    {student.city && `, ${student.city}`}
                  </p>
                </div>
              </div>
            </div>
          </Card>

          {/* Parents Information */}
          <Card title="Informasi Orang Tua / Wali">
            {student.parents && student.parents.length > 0 ? (
              <div className="space-y-4">
                {student.parents.map((parent, index) => (
                  <div key={index} className="p-4 bg-gray-50 rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-2">
                        <Users className="w-5 h-5 text-gray-600" />
                        <span className="font-medium text-gray-900">{parent.full_name}</span>
                      </div>
                      <div className="flex gap-2">
                        {parent.is_primary && (
                          <Badge variant="info">Utama</Badge>
                        )}
                        {parent.is_financial && (
                          <Badge variant="success">Penanggung Biaya</Badge>
                        )}
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-3 text-sm">
                      <div>
                        <span className="text-gray-500">Hubungan:</span>
                        <span className="ml-2 text-gray-900">{parent.relationship}</span>
                      </div>
                      <div>
                        <span className="text-gray-500">Telepon:</span>
                        <span className="ml-2 text-gray-900">{parent.phone}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-500 text-center py-4">Belum ada data orang tua</p>
            )}
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Status Card */}
          <Card>
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium text-gray-500">Status</label>
                <div className="mt-2">
                  <Badge variant={getStatusVariant(student.status)}>
                    {student.status_name || student.status}
                  </Badge>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Cabang</label>
                <p className="text-base font-medium text-gray-900 mt-1">
                  {student.branch_name}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Kelas</label>
                <p className="text-base text-gray-900 mt-1">
                  {student.class_name || '-'}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tanggal Daftar</label>
                <div className="flex items-center gap-2 mt-1">
                  <Calendar className="w-4 h-4 text-gray-400" />
                  <p className="text-base text-gray-900">
                    {formatDate(student.registration_date)}
                  </p>
                </div>
              </div>
            </div>
          </Card>

          {/* Quick Actions */}
          <Card title="Aksi Cepat">
            <div className="space-y-2">
              <Button
                variant="secondary"
                className="w-full justify-start"
                onClick={() => navigate(`/students/${id}/invoices`)}
              >
                <FileText className="w-4 h-4 mr-2" />
                Lihat Tagihan
              </Button>
              <Button
                variant="secondary"
                className="w-full justify-start"
                onClick={() => navigate(`/students/${id}/payments`)}
              >
                <DollarSign className="w-4 h-4 mr-2" />
                Riwayat Pembayaran
              </Button>
              <Button
                variant="secondary"
                className="w-full justify-start"
              >
                <User className="w-4 h-4 mr-2" />
                Data Akademik
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
