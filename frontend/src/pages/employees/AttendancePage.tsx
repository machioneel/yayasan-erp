import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Card } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Badge } from '@/components/common/Badge';
import { Employee, ApiResponse } from '@/types';
import { 
  Calendar,
  CheckCircle,
  XCircle,
  Clock,
  AlertCircle,
  Download,
  ChevronLeft,
  ChevronRight,
  Users
} from 'lucide-react';
import { formatDate } from '@/utils/format';

type AttendanceStatus = 'present' | 'absent' | 'late' | 'sick' | 'permission' | 'leave';

interface AttendanceRecord {
  employee_id: string;
  date: string;
  status: AttendanceStatus;
  check_in?: string;
  check_out?: string;
  notes?: string;
}

export default function AttendancePage() {
  const queryClient = useQueryClient();
  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
  const [attendanceData, setAttendanceData] = useState<Map<string, AttendanceStatus>>(new Map());

  const { data: employees } = useQuery({
    queryKey: ['employees'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Employee[]>>('/employees');
      return response.data.data?.filter(emp => emp.employment_type !== 'freelance') || [];
    },
  });

  const saveMutation = useMutation({
    mutationFn: async (data: AttendanceRecord[]) => {
      return api.post('/attendance/bulk', { date: selectedDate, records: data });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendance'] });
      alert('Absensi berhasil disimpan!');
    },
  });

  const handleStatusChange = (employeeId: string, status: AttendanceStatus) => {
    const newData = new Map(attendanceData);
    newData.set(employeeId, status);
    setAttendanceData(newData);
  };

  const handleSave = () => {
    const records: AttendanceRecord[] = Array.from(attendanceData.entries()).map(
      ([employee_id, status]) => ({
        employee_id,
        date: selectedDate,
        status,
      })
    );

    if (records.length === 0) {
      alert('Tidak ada data absensi untuk disimpan');
      return;
    }

    if (confirm(`Simpan absensi untuk ${records.length} karyawan?`)) {
      saveMutation.mutate(records);
    }
  };

  const handleBulkMark = (status: AttendanceStatus) => {
    if (!employees) return;
    
    if (confirm(`Tandai semua karyawan sebagai "${getStatusLabel(status)}"?`)) {
      const newData = new Map<string, AttendanceStatus>();
      employees.forEach(emp => {
        newData.set(emp.id, status);
      });
      setAttendanceData(newData);
    }
  };

  const getStatusLabel = (status: AttendanceStatus) => {
    switch (status) {
      case 'present': return 'Hadir';
      case 'absent': return 'Alpa';
      case 'late': return 'Terlambat';
      case 'sick': return 'Sakit';
      case 'permission': return 'Izin';
      case 'leave': return 'Cuti';
      default: return status;
    }
  };

  const getStatusVariant = (status: AttendanceStatus) => {
    switch (status) {
      case 'present': return 'success';
      case 'late': return 'warning';
      case 'sick': return 'info';
      case 'permission': return 'info';
      case 'leave': return 'default';
      case 'absent': return 'danger';
      default: return 'default';
    }
  };

  const getStatusIcon = (status: AttendanceStatus) => {
    switch (status) {
      case 'present': return CheckCircle;
      case 'late': return Clock;
      case 'absent': return XCircle;
      case 'sick': return AlertCircle;
      case 'permission': return AlertCircle;
      case 'leave': return Calendar;
      default: return CheckCircle;
    }
  };

  const statusOptions: { value: AttendanceStatus; label: string; color: string }[] = [
    { value: 'present', label: 'Hadir', color: 'bg-green-100 hover:bg-green-200' },
    { value: 'late', label: 'Terlambat', color: 'bg-yellow-100 hover:bg-yellow-200' },
    { value: 'sick', label: 'Sakit', color: 'bg-blue-100 hover:bg-blue-200' },
    { value: 'permission', label: 'Izin', color: 'bg-purple-100 hover:bg-purple-200' },
    { value: 'leave', label: 'Cuti', color: 'bg-gray-100 hover:bg-gray-200' },
    { value: 'absent', label: 'Alpa', color: 'bg-red-100 hover:bg-red-200' },
  ];

  const summary = {
    total: employees?.length || 0,
    present: Array.from(attendanceData.values()).filter(s => s === 'present').length,
    late: Array.from(attendanceData.values()).filter(s => s === 'late').length,
    sick: Array.from(attendanceData.values()).filter(s => s === 'sick').length,
    permission: Array.from(attendanceData.values()).filter(s => s === 'permission').length,
    leave: Array.from(attendanceData.values()).filter(s => s === 'leave').length,
    absent: Array.from(attendanceData.values()).filter(s => s === 'absent').length,
  };

  const changeDate = (days: number) => {
    const current = new Date(selectedDate);
    current.setDate(current.getDate() + days);
    setSelectedDate(current.toISOString().split('T')[0]);
    setAttendanceData(new Map()); // Reset attendance data when date changes
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Attendance Management</h1>
          <p className="text-gray-600">Kelola absensi karyawan harian</p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary">
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button 
            variant="primary"
            onClick={handleSave}
            disabled={attendanceData.size === 0 || saveMutation.isPending}
            loading={saveMutation.isPending}
          >
            <CheckCircle className="w-4 h-4 mr-2" />
            Simpan Absensi
          </Button>
        </div>
      </div>

      {/* Date Navigator */}
      <Card>
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            onClick={() => changeDate(-1)}
          >
            <ChevronLeft className="w-4 h-4 mr-2" />
            Kemarin
          </Button>

          <div className="flex items-center gap-4">
            <Calendar className="w-5 h-5 text-gray-400" />
            <input
              type="date"
              value={selectedDate}
              onChange={(e) => {
                setSelectedDate(e.target.value);
                setAttendanceData(new Map());
              }}
              className="px-4 py-2 text-lg font-medium border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <Button
            variant="ghost"
            onClick={() => changeDate(1)}
          >
            Besok
            <ChevronRight className="w-4 h-4 ml-2" />
          </Button>
        </div>
      </Card>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-4">
        <Card>
          <div className="flex items-center gap-2 mb-2">
            <Users className="w-4 h-4 text-gray-500" />
            <span className="text-xs text-gray-600">Total</span>
          </div>
          <div className="text-2xl font-bold text-gray-900">{summary.total}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <CheckCircle className="w-4 h-4 text-green-500" />
            <span className="text-xs text-gray-600">Hadir</span>
          </div>
          <div className="text-2xl font-bold text-green-600">{summary.present}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <Clock className="w-4 h-4 text-yellow-500" />
            <span className="text-xs text-gray-600">Terlambat</span>
          </div>
          <div className="text-2xl font-bold text-yellow-600">{summary.late}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <AlertCircle className="w-4 h-4 text-blue-500" />
            <span className="text-xs text-gray-600">Sakit</span>
          </div>
          <div className="text-2xl font-bold text-blue-600">{summary.sick}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <AlertCircle className="w-4 h-4 text-purple-500" />
            <span className="text-xs text-gray-600">Izin</span>
          </div>
          <div className="text-2xl font-bold text-purple-600">{summary.permission}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <Calendar className="w-4 h-4 text-gray-500" />
            <span className="text-xs text-gray-600">Cuti</span>
          </div>
          <div className="text-2xl font-bold text-gray-600">{summary.leave}</div>
        </Card>

        <Card>
          <div className="flex items-center gap-2 mb-2">
            <XCircle className="w-4 h-4 text-red-500" />
            <span className="text-xs text-gray-600">Alpa</span>
          </div>
          <div className="text-2xl font-bold text-red-600">{summary.absent}</div>
        </Card>
      </div>

      {/* Bulk Actions */}
      <Card title="Aksi Massal">
        <div className="flex flex-wrap gap-2">
          {statusOptions.map((option) => (
            <Button
              key={option.value}
              variant="secondary"
              size="sm"
              onClick={() => handleBulkMark(option.value)}
              className={option.color}
            >
              Tandai Semua: {option.label}
            </Button>
          ))}
        </div>
      </Card>

      {/* Attendance List */}
      <Card title="Daftar Karyawan">
        <div className="space-y-2">
          {employees?.map((employee) => {
            const status = attendanceData.get(employee.id);
            const StatusIcon = status ? getStatusIcon(status) : Users;

            return (
              <div
                key={employee.id}
                className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
              >
                <div className="flex items-center gap-4 flex-1">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold">
                    {employee.full_name.charAt(0)}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">{employee.full_name}</div>
                    <div className="text-sm text-gray-500">
                      {employee.position} • {employee.department || 'N/A'}
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {status && (
                    <Badge variant={getStatusVariant(status)}>
                      <StatusIcon className="w-3 h-3 mr-1" />
                      {getStatusLabel(status)}
                    </Badge>
                  )}

                  <div className="flex gap-1">
                    {statusOptions.map((option) => {
                      const Icon = getStatusIcon(option.value);
                      return (
                        <button
                          key={option.value}
                          onClick={() => handleStatusChange(employee.id, option.value)}
                          className={`p-2 rounded transition-colors ${
                            status === option.value
                              ? option.color.replace('hover:', '')
                              : 'bg-white hover:bg-gray-100 border border-gray-300'
                          }`}
                          title={option.label}
                        >
                          <Icon className={`w-4 h-4 ${
                            status === option.value ? 'opacity-100' : 'opacity-50'
                          }`} />
                        </button>
                      );
                    })}
                  </div>
                </div>
              </div>
            );
          })}

          {(!employees || employees.length === 0) && (
            <div className="text-center py-12">
              <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500">Tidak ada data karyawan</p>
            </div>
          )}
        </div>
      </Card>

      {/* Legend */}
      <Card title="Keterangan">
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
          {statusOptions.map((option) => {
            const Icon = getStatusIcon(option.value);
            return (
              <div key={option.value} className="flex items-center gap-2">
                <Icon className="w-4 h-4" />
                <span>{option.label}</span>
              </div>
            );
          })}
        </div>
      </Card>
    </div>
  );
}
