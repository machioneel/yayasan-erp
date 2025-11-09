import { useNavigate, useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Employee, Branch, ApiResponse } from '@/types';
import { ArrowLeft, Save, User } from 'lucide-react';

const employeeSchema = z.object({
  employee_number: z.string().optional(),
  full_name: z.string().min(3, 'Nama minimal 3 karakter'),
  nik: z.string().optional(),
  npwp: z.string().optional(),
  gender: z.enum(['male', 'female']),
  birth_place: z.string().optional(),
  birth_date: z.string().optional(),
  email: z.string().email('Email tidak valid'),
  phone: z.string().min(10, 'Nomor telepon minimal 10 digit'),
  address: z.string().optional(),
  branch_id: z.string().min(1, 'Cabang harus dipilih'),
  position: z.string().min(2, 'Jabatan harus diisi'),
  department: z.string().optional(),
  employment_type: z.enum(['permanent', 'contract', 'intern', 'freelance']),
  join_date: z.string(),
  salary: z.number().min(0, 'Gaji tidak boleh negatif').optional(),
  is_teacher: z.boolean(),
  education_level: z.string().optional(),
  major: z.string().optional(),
  certification: z.string().optional(),
});

type EmployeeFormData = z.infer<typeof employeeSchema>;

export default function EmployeeFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEdit = !!id;

  // Fetch branches
  const { data: branches } = useQuery({
    queryKey: ['branches'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Branch[]>>('/branches');
      return response.data.data;
    },
  });

  // Fetch employee if editing
  const { data: employee } = useQuery({
    queryKey: ['employee', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Employee>>(`/employees/${id}`);
      return response.data.data;
    },
    enabled: isEdit,
  });

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<EmployeeFormData>({
    resolver: zodResolver(employeeSchema),
    defaultValues: employee || {
      gender: 'male',
      employment_type: 'permanent',
      join_date: new Date().toISOString().split('T')[0],
      is_teacher: false,
    },
  });

  const isTeacher = watch('is_teacher');

  const mutation = useMutation({
    mutationFn: async (data: EmployeeFormData) => {
      if (isEdit) {
        return api.put(`/employees/${id}`, data);
      }
      return api.post('/employees', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees'] });
      navigate('/employees');
    },
  });

  const onSubmit = (data: EmployeeFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/employees')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit Karyawan' : 'Tambah Karyawan Baru'}
            </h1>
            <p className="text-gray-600">
              {isEdit ? 'Update informasi karyawan' : 'Lengkapi data karyawan'}
            </p>
          </div>
        </div>
        <Button type="submit" variant="primary" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update' : 'Simpan'}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Form */}
        <div className="lg:col-span-2 space-y-6">
          {/* Personal Information */}
          <Card title="Informasi Pribadi">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="md:col-span-2">
                <Input
                  label="Nama Lengkap"
                  {...register('full_name')}
                  error={errors.full_name?.message}
                  required
                />
              </div>

              {isEdit && (
                <Input
                  label="NIP"
                  {...register('employee_number')}
                  disabled
                  helperText="NIP otomatis dibuat sistem"
                />
              )}

              <Input
                label="NIK (KTP)"
                placeholder="16 digit"
                {...register('nik')}
                error={errors.nik?.message}
              />

              <Input
                label="NPWP"
                placeholder="15 digit"
                {...register('npwp')}
                error={errors.npwp?.message}
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Jenis Kelamin <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('gender')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="male">Laki-laki</option>
                  <option value="female">Perempuan</option>
                </select>
              </div>

              <Input
                label="Tempat Lahir"
                {...register('birth_place')}
                error={errors.birth_place?.message}
              />

              <Input
                label="Tanggal Lahir"
                type="date"
                {...register('birth_date')}
                error={errors.birth_date?.message}
              />
            </div>
          </Card>

          {/* Contact Information */}
          <Card title="Informasi Kontak">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Email"
                type="email"
                {...register('email')}
                error={errors.email?.message}
                required
              />

              <Input
                label="Telepon"
                {...register('phone')}
                error={errors.phone?.message}
                required
              />

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Alamat
                </label>
                <textarea
                  {...register('address')}
                  rows={2}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Alamat lengkap"
                />
              </div>
            </div>
          </Card>

          {/* Employment Information */}
          <Card title="Informasi Kepegawaian">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Cabang <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('branch_id')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Pilih Cabang</option>
                  {branches?.map((branch) => (
                    <option key={branch.id} value={branch.id}>
                      {branch.name}
                    </option>
                  ))}
                </select>
                {errors.branch_id && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.branch_id.message}</p>
                )}
              </div>

              <Input
                label="Jabatan"
                placeholder="e.g., Guru Matematika, Staff Admin"
                {...register('position')}
                error={errors.position?.message}
                required
              />

              <Input
                label="Departemen"
                placeholder="e.g., Pendidikan, Administrasi"
                {...register('department')}
                error={errors.department?.message}
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Tipe Kepegawaian <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('employment_type')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="permanent">Tetap</option>
                  <option value="contract">Kontrak</option>
                  <option value="intern">Magang</option>
                  <option value="freelance">Freelance</option>
                </select>
              </div>

              <Input
                label="Tanggal Bergabung"
                type="date"
                {...register('join_date')}
                error={errors.join_date?.message}
                required
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  Gaji Pokok
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">Rp</span>
                  <input
                    type="number"
                    {...register('salary', { valueAsNumber: true })}
                    className="w-full pl-12 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="0"
                  />
                </div>
                {errors.salary && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.salary.message}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    {...register('is_teacher')}
                    className="rounded border-gray-300"
                  />
                  <span className="text-sm font-medium text-gray-700">
                    Karyawan ini adalah Guru/Pengajar
                  </span>
                </label>
              </div>
            </div>
          </Card>

          {/* Teacher Information - Only show if is_teacher is checked */}
          {isTeacher && (
            <Card title="Informasi Pengajar">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                  label="Pendidikan Terakhir"
                  placeholder="e.g., S1, S2, S3"
                  {...register('education_level')}
                  error={errors.education_level?.message}
                />

                <Input
                  label="Jurusan"
                  placeholder="e.g., Pendidikan Matematika"
                  {...register('major')}
                  error={errors.major?.message}
                />

                <div className="md:col-span-2">
                  <Input
                    label="Sertifikasi"
                    placeholder="e.g., Sertifikat Pendidik"
                    {...register('certification')}
                    error={errors.certification?.message}
                  />
                </div>
              </div>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-blue-100 rounded-lg">
                <User className="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">
                  {isEdit ? 'Update Data' : 'Data Baru'}
                </div>
                <div className="text-sm text-gray-500">
                  {isEdit ? 'Perbarui informasi' : 'Tambah karyawan'}
                </div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Semua field dengan tanda * wajib diisi</p>
              <p>✓ NIP akan dibuat otomatis oleh sistem</p>
              <p>✓ Data dapat diubah sewaktu-waktu</p>
              {isTeacher && (
                <p>✓ Informasi pengajar aktif</p>
              )}
            </div>
          </Card>

          <Card title="Tips">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• Pastikan email dan telepon valid</p>
              <p>• NIK dan NPWP untuk keperluan pajak</p>
              <p>• Pilih cabang sesuai penempatan</p>
              <p>• Centang "Guru" jika mengajar</p>
            </div>
          </Card>
        </div>
      </div>

      {/* Submit Button */}
      <div className="flex justify-end gap-3">
        <Button
          type="button"
          variant="secondary"
          onClick={() => navigate('/employees')}
        >
          Batal
        </Button>
        <Button type="submit" variant="primary" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update Karyawan' : 'Simpan Karyawan'}
        </Button>
      </div>
    </form>
  );
}
