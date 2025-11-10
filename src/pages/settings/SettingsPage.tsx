import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { 
  Building2,
  Mail,
  Phone,
  MapPin,
  Globe,
  Save,
  Upload,
  Bell,
  Lock,
  Database,
  Download
} from 'lucide-react';

interface CompanySettings {
  name: string;
  email: string;
  phone: string;
  address: string;
  website: string;
  tax_id: string;
}

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<'company' | 'system' | 'backup'>('company');

  const { register, handleSubmit, formState: { errors } } = useForm<CompanySettings>({
    defaultValues: {
      name: 'Yayasan Pendidikan ABC',
      email: 'info@yayasanabc.org',
      phone: '+62 21 1234567',
      address: 'Jl. Pendidikan No. 123, Jakarta',
      website: 'www.yayasanabc.org',
      tax_id: '01.234.567.8-901.000',
    },
  });

  const onSubmit = (data: CompanySettings) => {
    console.log('Settings updated:', data);
    alert('Pengaturan berhasil disimpan!');
  };

  const tabs = [
    { id: 'company' as const, label: 'Company Info', icon: Building2 },
    { id: 'system' as const, label: 'System', icon: Lock },
    { id: 'backup' as const, label: 'Backup & Data', icon: Database },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-gray-600">Kelola pengaturan sistem dan informasi yayasan</p>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <div className="flex space-x-8">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Company Info Tab */}
      {activeTab === 'company' && (
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Card title="Informasi Yayasan">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="md:col-span-2">
                <Input
                  label="Nama Yayasan"
                  {...register('name', { required: 'Nama yayasan wajib diisi' })}
                  error={errors.name?.message}
                  required
                  icon={Building2}
                />
              </div>

              <Input
                label="Email"
                type="email"
                {...register('email', { required: 'Email wajib diisi' })}
                error={errors.email?.message}
                required
                icon={Mail}
              />

              <Input
                label="Telepon"
                {...register('phone', { required: 'Telepon wajib diisi' })}
                error={errors.phone?.message}
                required
                icon={Phone}
              />

              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1.5">
                  <MapPin className="w-4 h-4 inline mr-2" />
                  Alamat <span className="text-red-500">*</span>
                </label>
                <textarea
                  {...register('address', { required: 'Alamat wajib diisi' })}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  placeholder="Alamat lengkap yayasan"
                />
                {errors.address && (
                  <p className="mt-1.5 text-sm text-red-600">{errors.address.message}</p>
                )}
              </div>

              <Input
                label="Website"
                {...register('website')}
                icon={Globe}
                placeholder="www.example.com"
              />

              <Input
                label="NPWP"
                {...register('tax_id')}
                placeholder="00.000.000.0-000.000"
              />
            </div>
          </Card>

          <Card title="Logo Yayasan">
            <div className="flex items-center gap-6">
              <div className="w-24 h-24 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center text-white text-3xl font-bold">
                YA
              </div>
              <div className="flex-1">
                <p className="text-sm text-gray-600 mb-3">
                  Upload logo yayasan (max 2MB, format: JPG, PNG)
                </p>
                <Button variant="secondary" type="button">
                  <Upload className="w-4 h-4 mr-2" />
                  Upload Logo
                </Button>
              </div>
            </div>
          </Card>

          <div className="flex justify-end">
            <Button type="submit" variant="primary">
              <Save className="w-4 h-4 mr-2" />
              Simpan Perubahan
            </Button>
          </div>
        </form>
      )}

      {/* System Tab */}
      {activeTab === 'system' && (
        <div className="space-y-6">
          <Card title="Notifikasi">
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-gray-900">Email Notifications</div>
                  <div className="text-sm text-gray-500">
                    Terima notifikasi via email
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" defaultChecked className="sr-only peer" />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-gray-900">Low Stock Alerts</div>
                  <div className="text-sm text-gray-500">
                    Notifikasi stok rendah inventory
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" defaultChecked className="sr-only peer" />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium text-gray-900">Invoice Reminders</div>
                  <div className="text-sm text-gray-500">
                    Reminder invoice jatuh tempo
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" defaultChecked className="sr-only peer" />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
          </Card>

          <Card title="Keamanan">
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Session Timeout
                </label>
                <select className="w-full md:w-64 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
                  <option value="15">15 menit</option>
                  <option value="30" selected>30 menit</option>
                  <option value="60">1 jam</option>
                  <option value="120">2 jam</option>
                </select>
              </div>

              <div>
                <label className="flex items-center gap-2">
                  <input type="checkbox" defaultChecked className="rounded border-gray-300" />
                  <span className="text-sm font-medium text-gray-700">
                    Require password change every 90 days
                  </span>
                </label>
              </div>

              <div>
                <label className="flex items-center gap-2">
                  <input type="checkbox" className="rounded border-gray-300" />
                  <span className="text-sm font-medium text-gray-700">
                    Enable two-factor authentication (2FA)
                  </span>
                </label>
              </div>
            </div>
          </Card>

          <div className="flex justify-end">
            <Button variant="primary">
              <Save className="w-4 h-4 mr-2" />
              Simpan Pengaturan
            </Button>
          </div>
        </div>
      )}

      {/* Backup Tab */}
      {activeTab === 'backup' && (
        <div className="space-y-6">
          <Card title="Database Backup">
            <div className="space-y-4">
              <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <div className="flex items-start gap-3">
                  <Database className="w-5 h-5 text-blue-600 mt-0.5" />
                  <div className="flex-1">
                    <div className="font-medium text-blue-900 mb-1">
                      Automatic Backup Enabled
                    </div>
                    <div className="text-sm text-blue-700">
                      Backup otomatis dilakukan setiap hari pada pukul 02:00 WIB
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <div className="font-medium text-gray-900 mb-2">Last Backup</div>
                <div className="text-sm text-gray-600">
                  8 November 2025, 02:00 WIB
                </div>
              </div>

              <div>
                <Button variant="primary">
                  <Download className="w-4 h-4 mr-2" />
                  Download Backup Sekarang
                </Button>
              </div>
            </div>
          </Card>

          <Card title="Export Data">
            <div className="space-y-4">
              <p className="text-sm text-gray-600">
                Export data sistem ke format Excel untuk backup atau analisis
              </p>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Students
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Invoices
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Payments
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Employees
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Assets
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="w-4 h-4 mr-2" />
                  Inventory
                </Button>
              </div>
            </div>
          </Card>

          <Card title="Import Data">
            <div className="space-y-4">
              <p className="text-sm text-gray-600">
                Import data dari file Excel ke sistem
              </p>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Pilih Module
                </label>
                <select className="w-full md:w-64 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 mb-3">
                  <option value="">-- Pilih Module --</option>
                  <option value="students">Students</option>
                  <option value="employees">Employees</option>
                  <option value="inventory">Inventory</option>
                </select>
              </div>

              <div>
                <Button variant="secondary">
                  <Upload className="w-4 h-4 mr-2" />
                  Upload Excel File
                </Button>
              </div>

              <div className="text-xs text-gray-500">
                <p>* Download template Excel terlebih dahulu</p>
                <p>* Format file harus sesuai template</p>
                <p>* Max file size: 10MB</p>
              </div>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
}
