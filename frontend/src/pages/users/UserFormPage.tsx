import { useNavigate, useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { ApiResponse } from '@/types';
import { ArrowLeft, Save, User as UserIcon, Shield, Mail, Lock } from 'lucide-react';

const userSchema = z.object({
  username: z.string().min(3, 'Username minimal 3 karakter'),
  email: z.string().email('Email tidak valid'),
  full_name: z.string().min(3, 'Nama lengkap minimal 3 karakter'),
  password: z.string().min(6, 'Password minimal 6 karakter').optional(),
  confirm_password: z.string().optional(),
  role: z.enum(['admin', 'manager', 'staff', 'viewer']),
  is_active: z.boolean(),
}).refine((data) => {
  if (data.password && data.password !== data.confirm_password) {
    return false;
  }
  return true;
}, {
  message: "Password tidak cocok",
  path: ["confirm_password"],
});

type UserFormData = z.infer<typeof userSchema>;

interface User extends UserFormData {
  id: string;
  created_at: string;
}

export default function UserFormPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEdit = !!id;

  const { data: user } = useQuery({
    queryKey: ['user', id],
    queryFn: async () => {
      const response = await api.get<ApiResponse<User>>(`/users/${id}`);
      return response.data.data;
    },
    enabled: isEdit,
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    defaultValues: user || {
      role: 'staff',
      is_active: true,
    },
  });

  const mutation = useMutation({
    mutationFn: async (data: UserFormData) => {
      const { confirm_password, ...submitData } = data;
      if (isEdit) {
        // Don't send password if empty
        if (!submitData.password) {
          delete submitData.password;
        }
        return api.put(`/users/${id}`, submitData);
      }
      return api.post('/users', submitData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      navigate('/users');
    },
  });

  const onSubmit = (data: UserFormData) => {
    mutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button type="button" variant="ghost" onClick={() => navigate('/users')}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {isEdit ? 'Edit User' : 'Tambah User Baru'}
            </h1>
            <p className="text-gray-600">
              {isEdit ? 'Update informasi user' : 'Buat akun user baru'}
            </p>
          </div>
        </div>
        <Button type="submit" variant="primary" loading={mutation.isPending}>
          <Save className="w-4 h-4 mr-2" />
          {isEdit ? 'Update' : 'Simpan'}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          {/* Account Info */}
          <Card title="Informasi Akun">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Username"
                {...register('username')}
                error={errors.username?.message}
                required
                icon={UserIcon}
                placeholder="username123"
              />

              <Input
                label="Email"
                type="email"
                {...register('email')}
                error={errors.email?.message}
                required
                icon={Mail}
                placeholder="user@example.com"
              />

              <div className="md:col-span-2">
                <Input
                  label="Nama Lengkap"
                  {...register('full_name')}
                  error={errors.full_name?.message}
                  required
                  placeholder="John Doe"
                />
              </div>
            </div>
          </Card>

          {/* Password */}
          <Card title={isEdit ? 'Ubah Password (Opsional)' : 'Password'}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Password"
                type="password"
                {...register('password')}
                error={errors.password?.message}
                required={!isEdit}
                icon={Lock}
                placeholder="••••••••"
                helperText={isEdit ? 'Kosongkan jika tidak ingin mengubah password' : undefined}
              />

              <Input
                label="Konfirmasi Password"
                type="password"
                {...register('confirm_password')}
                error={errors.confirm_password?.message}
                required={!isEdit}
                icon={Lock}
                placeholder="••••••••"
              />
            </div>
          </Card>

          {/* Permissions */}
          <Card title="Role & Permissions">
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <Shield className="w-4 h-4 inline mr-2" />
                  Role <span className="text-red-500">*</span>
                </label>
                <select
                  {...register('role')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  <option value="viewer">Viewer - View only access</option>
                  <option value="staff">Staff - Basic operations</option>
                  <option value="manager">Manager - Management access</option>
                  <option value="admin">Admin - Full access</option>
                </select>
                {errors.role && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.role.message}</p>
                )}
              </div>

              <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <div className="text-sm text-blue-900 font-medium mb-2">Role Permissions:</div>
                <div className="text-sm text-blue-700 space-y-1">
                  <p>• <strong>Viewer:</strong> View only, no edit/delete</p>
                  <p>• <strong>Staff:</strong> Create, edit own records</p>
                  <p>• <strong>Manager:</strong> Full CRUD, approve workflows</p>
                  <p>• <strong>Admin:</strong> All access + user management</p>
                </div>
              </div>

              <div>
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    {...register('is_active')}
                    className="rounded border-gray-300"
                  />
                  <span className="text-sm font-medium text-gray-700">
                    User Active (dapat login)
                  </span>
                </label>
              </div>
            </div>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-blue-100 rounded-lg">
                <UserIcon className="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <div className="font-medium text-gray-900">
                  {isEdit ? 'Update User' : 'New User'}
                </div>
                <div className="text-sm text-gray-500">
                  {isEdit ? 'Modify account' : 'Create account'}
                </div>
              </div>
            </div>

            <div className="space-y-3 text-sm text-gray-600">
              <p>✓ Username harus unique</p>
              <p>✓ Email untuk notifikasi</p>
              <p>✓ Password minimal 6 karakter</p>
              <p>✓ Role menentukan hak akses</p>
            </div>
          </Card>

          <Card title="Security">
            <div className="text-sm text-gray-600 space-y-2">
              <p className="font-medium text-gray-900">Password Guidelines:</p>
              <ul className="list-disc list-inside space-y-1">
                <li>Minimal 6 karakter</li>
                <li>Kombinasi huruf & angka</li>
                <li>Hindari password umum</li>
                <li>Jangan share password</li>
              </ul>
            </div>
          </Card>

          <Card title="Account Status">
            <div className="text-sm text-gray-600 space-y-2">
              <p>• <strong>Active:</strong> User dapat login</p>
              <p>• <strong>Inactive:</strong> Login diblokir</p>
              <p className="mt-3 text-xs text-gray-500">
                Non-aktifkan user untuk suspend akses tanpa menghapus data
              </p>
            </div>
          </Card>
        </div>
      </div>
    </form>
  );
}
