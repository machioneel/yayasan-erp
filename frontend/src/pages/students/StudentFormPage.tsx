import { useNavigate, useParams } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { ArrowLeft, Save, Plus, Trash2 } from 'lucide-react';
import { Student, ApiResponse, Branch } from '@/types';

const studentSchema = z.object({
  full_name: z.string().min(3, 'Nama minimal 3 karakter'),
  nick_name: z.string().optional(),
  nisn: z.string().optional(),
  gender: z.enum(['male', 'female']),
  birth_place: z.string().optional(),
  birth_date: z.string().optional(),
  religion: z.string().optional(),
  email: z.string().email('Email tidak valid').optional().or(z.literal('')),
  phone: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  branch_id: z.string().min(1, 'Cabang harus dipilih'),
  registration_date: z.string(),
  parents: z.array(z.object({
    full_name: z.string().min(3, 'Nama orang tua minimal 3 karakter'),
    relationship: z.string().min(1, 'Hubungan harus diisi'),
    phone: z.string().min(10, 'Nomor telepon minimal 10 digit'),
    email: z.string().email('Email tidak valid').optional().or(z.literal('')),
    occupation: z.string().optional(),
    is_primary: z.boolean(),
    is_financial: z.boolean(),
  })).optional(),
});

type StudentFormData = z.infer<typeof studentSchema>;

export default function StudentFormPage() {
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

  // Fetch student data if editing
  const { data: student } = useQuery({
    queryKey: ['student', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Student>>(`/students/${id}`);
      return response.data.data;
    },
    enabled: isEdit,
  });

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<StudentFormData>({
    resolver: zodResolver(studentSchema),
    defaultValues: student || {
      gender: 'male',
      registration_date: new Date().toISOString().split('T')[0],
      parents: [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'parents',
  });

  const mutation = useMutation({
    mutationFn: async (data: StudentFormData) => {
      if (isEdit) {
        return api.put(`/students/${id}`, data);
      }
      return api.post('/students', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['students'] });
      navigate('/students');
    },
  });

  const onSubmit = (data: StudentFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            type="button"
            variant="ghost"
            onClick={() => navigate('/students')}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit Siswa' : 'Tambah Siswa Baru'}
            </h1>
            <p className="text-gray-600">
              {isEdit ? 'Update informasi siswa' : 'Lengkapi formulir pendaftaran siswa'}
            </p>
          </div>
        </div>
        <Button
          type="submit"
          variant="primary"
          loading={mutation.isPending}
        >
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update' : 'Simpan'}
        </Button>
      </div>

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

          <Input
            label="Nama Panggilan"
            {...register('nick_name')}
            error={errors.nick_name?.message}
          />

          <Input
            label="NISN"
            {...register('nisn')}
            error={errors.nisn?.message}
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

          <Input
            label="Agama"
            {...register('religion')}
            error={errors.religion?.message}
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
          />

          <Input
            label="Telepon"
            {...register('phone')}
            error={errors.phone?.message}
          />

          <div className="md:col-span-2">
            <Input
              label="Alamat"
              {...register('address')}
              error={errors.address?.message}
            />
          </div>

          <Input
            label="Kota"
            {...register('city')}
            error={errors.city?.message}
          />
        </div>
      </Card>

      {/* School Information */}
      <Card title="Informasi Sekolah">
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
            label="Tanggal Pendaftaran"
            type="date"
            {...register('registration_date')}
            error={errors.registration_date?.message}
            required
          />
        </div>
      </Card>

      {/* Parents Information */}
      <Card
        title="Informasi Orang Tua / Wali"
        action={
          <Button
            type="button"
            variant="secondary"
            size="sm"
            onClick={() =>
              append({
                full_name: '',
                relationship: '',
                phone: '',
                email: '',
                occupation: '',
                is_primary: fields.length === 0,
                is_financial: fields.length === 0,
              })
            }
          >
            <Plus className="w-4 h-4 mr-2" />
            Tambah Orang Tua
          </Button>
        }
      >
        {fields.length === 0 ? (
          <p className="text-gray-500 text-center py-8">
            Belum ada data orang tua. Click "Tambah Orang Tua" untuk menambahkan.
          </p>
        ) : (
          <div className="space-y-6">
            {fields.map((field, index) => (
              <div key={field.id} className="p-4 border border-gray-200 rounded-lg">
                <div className="flex items-center justify-between mb-4">
                  <h4 className="font-medium text-gray-900">Orang Tua #{index + 1}</h4>
                  <Button
                    type="button"
                    variant="danger"
                    size="sm"
                    onClick={() => remove(index)}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Input
                    label="Nama Lengkap"
                    {...register(`parents.${index}.full_name`)}
                    error={errors.parents?.[index]?.full_name?.message}
                    required
                  />

                  <Input
                    label="Hubungan"
                    placeholder="e.g., Ayah, Ibu, Wali"
                    {...register(`parents.${index}.relationship`)}
                    error={errors.parents?.[index]?.relationship?.message}
                    required
                  />

                  <Input
                    label="Telepon"
                    {...register(`parents.${index}.phone`)}
                    error={errors.parents?.[index]?.phone?.message}
                    required
                  />

                  <Input
                    label="Email"
                    type="email"
                    {...register(`parents.${index}.email`)}
                    error={errors.parents?.[index]?.email?.message}
                  />

                  <Input
                    label="Pekerjaan"
                    {...register(`parents.${index}.occupation`)}
                  />

                  <div className="flex gap-4 items-center">
                    <label className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        {...register(`parents.${index}.is_primary`)}
                        className="rounded border-gray-300"
                      />
                      <span className="text-sm text-gray-700">Kontak Utama</span>
                    </label>

                    <label className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        {...register(`parents.${index}.is_financial`)}
                        className="rounded border-gray-300"
                      />
                      <span className="text-sm text-gray-700">Penanggung Biaya</span>
                    </label>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Submit Buttons */}
      <div className="flex justify-end gap-3">
        <Button
          type="button"
          variant="secondary"
          onClick={() => navigate('/students')}
        >
          Batal
        </Button>
        <Button
          type="submit"
          variant="primary"
          loading={mutation.isPending}
        >
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update Siswa' : 'Simpan Siswa'}
        </Button>
      </div>
    </form>
  );
}
