import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Employee, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { 
  ArrowLeft, 
  Edit, 
  Mail, 
  Phone, 
  MapPin, 
  Calendar,
  Briefcase,
  DollarSign,
  FileText
} from 'lucide-react';
import { formatDate, formatCurrency } from '@/utils/format';

export default function EmployeeDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: employee, isLoading } = useQuery({
    queryKey: ['employee', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Employee>>(`/employees/${id}`);
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

  if (!employee) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Karyawan tidak ditemukan</p>
      </div>
    );
  }

  const getEmploymentTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      permanent: 'Tetap',
      contract: 'Kontrak',
      intern: 'Magang',
      freelance: 'Freelance',
    };
    return labels[type] || type;
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/employees')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{employee.full_name}</h1>
            <p className="text-gray-600">{employee.employee_number}</p>
          </div>
        </div>
        <Button variant="primary" onClick={() => navigate(`/employees/${id}/edit`)}>
          <Edit className="w-4 h-4 mr-2" />
          Edit
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Personal Information */}
          <Card title="Informasi Pribadi">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Nama Lengkap</label>
                <p className="text-base font-medium text-gray-900 mt-1">{employee.full_name}</p>
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-500">NIP</label>
                <p className="text-base font-mono text-gray-900 mt-1">{employee.employee_number}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">NIK</label>
                <p className="text-base text-gray-900 mt-1">{employee.nik || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">NPWP</label>
                <p className="text-base text-gray-900 mt-1">{employee.npwp || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Jenis Kelamin</label>
                <p className="text-base text-gray-900 mt-1">
                  {employee.gender === 'male' ? 'Laki-laki' : 'Perempuan'}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tempat, Tanggal Lahir</label>
                <p className="text-base text-gray-900 mt-1">
                  {employee.birth_place && employee.birth_date
                    ? `${employee.birth_place}, ${formatDate(employee.birth_date)}`
                    : '-'}
                </p>
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
                  <p className="text-base text-gray-900">{employee.email}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <Phone className="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <label className="text-sm font-medium text-gray-500">Telepon</label>
                  <p className="text-base text-gray-900">{employee.phone}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-gray-400 mt-0.5" />
                <div>
                  <label className="text-sm font-medium text-gray-500">Alamat</label>
                  <p className="text-base text-gray-900">{employee.address || '-'}</p>
                </div>
              </div>
            </div>
          </Card>

          {/* Employment Information */}
          <Card title="Informasi Kepegawaian">
            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-500">Jabatan</label>
                <p className="text-base font-medium text-gray-900 mt-1">{employee.position}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Departemen</label>
                <p className="text-base text-gray-900 mt-1">{employee.department || '-'}</p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tipe Kepegawaian</label>
                <p className="text-base text-gray-900 mt-1">
                  {getEmploymentTypeLabel(employee.employment_type)}
                </p>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Tanggal Bergabung</label>
                <p className="text-base text-gray-900 mt-1">
                  {formatDate(employee.join_date)}
                </p>
              </div>
            </div>
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
                  <Badge variant={employee.status === 'active' ? 'success' : 'default'}>
                    {employee.status === 'active' ? 'Aktif' : 'Tidak Aktif'}
                  </Badge>
                </div>
              </div>

              <div>
                <label className="text-sm font-medium text-gray-500">Cabang</label>
                <p className="text-base font-medium text-gray-900 mt-1">
                  {employee.branch_name}
                </p>
              </div>

              {employee.is_teacher && (
                <div>
                  <Badge variant="info">Guru/Pengajar</Badge>
                </div>
              )}

              <div>
                <label className="text-sm font-medium text-gray-500">Bergabung Sejak</label>
                <div className="flex items-center gap-2 mt-1">
                  <Calendar className="w-4 h-4 text-gray-400" />
                  <p className="text-base text-gray-900">
                    {formatDate(employee.created_at)}
                  </p>
                </div>
              </div>
            </div>
          </Card>

          {/* Quick Actions */}
          <Card title="Aksi Cepat">
            <div className="space-y-2">
              <Button variant="secondary" className="w-full justify-start">
                <DollarSign className="w-4 h-4 mr-2" />
                Riwayat Gaji
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <Calendar className="w-4 h-4 mr-2" />
                Absensi
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <FileText className="w-4 h-4 mr-2" />
                Kontrak
              </Button>
              <Button variant="secondary" className="w-full justify-start">
                <Briefcase className="w-4 h-4 mr-2" />
                Cuti
              </Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
