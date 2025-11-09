import { Card } from '@/components/common/Card';
import { Badge } from '@/components/common/Badge';
import { 
  Shield, 
  Eye, 
  Edit, 
  Trash2,
  Users,
  CheckCircle,
  XCircle
} from 'lucide-react';

interface Permission {
  module: string;
  view: boolean;
  create: boolean;
  edit: boolean;
  delete: boolean;
  approve?: boolean;
}

interface Role {
  name: string;
  label: string;
  description: string;
  permissions: Permission[];
  userCount: number;
}

const roles: Role[] = [
  {
    name: 'admin',
    label: 'Administrator',
    description: 'Full system access with user management',
    userCount: 2,
    permissions: [
      { module: 'Students', view: true, create: true, edit: true, delete: true },
      { module: 'Invoices', view: true, create: true, edit: true, delete: true, approve: true },
      { module: 'Payments', view: true, create: true, edit: true, delete: true },
      { module: 'Finance', view: true, create: true, edit: true, delete: true, approve: true },
      { module: 'Employees', view: true, create: true, edit: true, delete: true },
      { module: 'Assets', view: true, create: true, edit: true, delete: true },
      { module: 'Inventory', view: true, create: true, edit: true, delete: true },
      { module: 'Users', view: true, create: true, edit: true, delete: true },
      { module: 'Settings', view: true, create: true, edit: true, delete: true },
    ],
  },
  {
    name: 'manager',
    label: 'Manager',
    description: 'Management access with approval rights',
    userCount: 5,
    permissions: [
      { module: 'Students', view: true, create: true, edit: true, delete: false },
      { module: 'Invoices', view: true, create: true, edit: true, delete: false, approve: true },
      { module: 'Payments', view: true, create: true, edit: true, delete: false },
      { module: 'Finance', view: true, create: true, edit: true, delete: false, approve: true },
      { module: 'Employees', view: true, create: false, edit: true, delete: false },
      { module: 'Assets', view: true, create: true, edit: true, delete: false },
      { module: 'Inventory', view: true, create: true, edit: true, delete: false },
      { module: 'Users', view: true, create: false, edit: false, delete: false },
      { module: 'Settings', view: true, create: false, edit: false, delete: false },
    ],
  },
  {
    name: 'staff',
    label: 'Staff',
    description: 'Basic operational access',
    userCount: 12,
    permissions: [
      { module: 'Students', view: true, create: true, edit: true, delete: false },
      { module: 'Invoices', view: true, create: true, edit: true, delete: false },
      { module: 'Payments', view: true, create: true, edit: false, delete: false },
      { module: 'Finance', view: true, create: true, edit: false, delete: false },
      { module: 'Employees', view: true, create: false, edit: false, delete: false },
      { module: 'Assets', view: true, create: false, edit: false, delete: false },
      { module: 'Inventory', view: true, create: true, edit: true, delete: false },
      { module: 'Users', view: false, create: false, edit: false, delete: false },
      { module: 'Settings', view: false, create: false, edit: false, delete: false },
    ],
  },
  {
    name: 'viewer',
    label: 'Viewer',
    description: 'Read-only access',
    userCount: 3,
    permissions: [
      { module: 'Students', view: true, create: false, edit: false, delete: false },
      { module: 'Invoices', view: true, create: false, edit: false, delete: false },
      { module: 'Payments', view: true, create: false, edit: false, delete: false },
      { module: 'Finance', view: true, create: false, edit: false, delete: false },
      { module: 'Employees', view: true, create: false, edit: false, delete: false },
      { module: 'Assets', view: true, create: false, edit: false, delete: false },
      { module: 'Inventory', view: true, create: false, edit: false, delete: false },
      { module: 'Users', view: false, create: false, edit: false, delete: false },
      { module: 'Settings', view: false, create: false, edit: false, delete: false },
    ],
  },
];

export default function RolesPage() {
  const getRoleBadge = (name: string) => {
    switch (name) {
      case 'admin': return 'danger';
      case 'manager': return 'warning';
      case 'staff': return 'info';
      case 'viewer': return 'default';
      default: return 'default';
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Roles & Permissions</h1>
        <p className="text-gray-600">View role-based access control matrix</p>
      </div>

      {/* Roles Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {roles.map((role) => (
          <Card key={role.name}>
            <div className="flex items-center justify-between mb-3">
              <Badge variant={getRoleBadge(role.name) as any}>
                <Shield className="w-3 h-3 mr-1" />
                {role.label}
              </Badge>
              <div className="flex items-center gap-1 text-sm text-gray-600">
                <Users className="w-4 h-4" />
                {role.userCount}
              </div>
            </div>
            <p className="text-sm text-gray-600">{role.description}</p>
          </Card>
        ))}
      </div>

      {/* Permission Matrix */}
      {roles.map((role) => (
        <Card key={role.name} title={role.label}>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Module</th>
                  <th className="px-4 py-2 text-center font-medium text-gray-700">
                    <Eye className="w-4 h-4 mx-auto" title="View" />
                  </th>
                  <th className="px-4 py-2 text-center font-medium text-gray-700">
                    <Edit className="w-4 h-4 mx-auto" title="Create" />
                  </th>
                  <th className="px-4 py-2 text-center font-medium text-gray-700">
                    <Edit className="w-4 h-4 mx-auto" title="Edit" />
                  </th>
                  <th className="px-4 py-2 text-center font-medium text-gray-700">
                    <Trash2 className="w-4 h-4 mx-auto" title="Delete" />
                  </th>
                  <th className="px-4 py-2 text-center font-medium text-gray-700">
                    <CheckCircle className="w-4 h-4 mx-auto" title="Approve" />
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {role.permissions.map((perm) => (
                  <tr key={perm.module} className="hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-900">{perm.module}</td>
                    <td className="px-4 py-3 text-center">
                      {perm.view ? (
                        <CheckCircle className="w-4 h-4 text-green-600 mx-auto" />
                      ) : (
                        <XCircle className="w-4 h-4 text-gray-300 mx-auto" />
                      )}
                    </td>
                    <td className="px-4 py-3 text-center">
                      {perm.create ? (
                        <CheckCircle className="w-4 h-4 text-green-600 mx-auto" />
                      ) : (
                        <XCircle className="w-4 h-4 text-gray-300 mx-auto" />
                      )}
                    </td>
                    <td className="px-4 py-3 text-center">
                      {perm.edit ? (
                        <CheckCircle className="w-4 h-4 text-green-600 mx-auto" />
                      ) : (
                        <XCircle className="w-4 h-4 text-gray-300 mx-auto" />
                      )}
                    </td>
                    <td className="px-4 py-3 text-center">
                      {perm.delete ? (
                        <CheckCircle className="w-4 h-4 text-green-600 mx-auto" />
                      ) : (
                        <XCircle className="w-4 h-4 text-gray-300 mx-auto" />
                      )}
                    </td>
                    <td className="px-4 py-3 text-center">
                      {perm.approve !== undefined ? (
                        perm.approve ? (
                          <CheckCircle className="w-4 h-4 text-green-600 mx-auto" />
                        ) : (
                          <XCircle className="w-4 h-4 text-gray-300 mx-auto" />
                        )
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      ))}

      {/* Legend */}
      <Card title="Permission Legend">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm">
          <div className="flex items-center gap-2">
            <Eye className="w-4 h-4 text-gray-600" />
            <span>View - Can view data</span>
          </div>
          <div className="flex items-center gap-2">
            <Edit className="w-4 h-4 text-gray-600" />
            <span>Create - Can create new</span>
          </div>
          <div className="flex items-center gap-2">
            <Edit className="w-4 h-4 text-gray-600" />
            <span>Edit - Can modify</span>
          </div>
          <div className="flex items-center gap-2">
            <Trash2 className="w-4 h-4 text-gray-600" />
            <span>Delete - Can remove</span>
          </div>
          <div className="flex items-center gap-2">
            <CheckCircle className="w-4 h-4 text-gray-600" />
            <span>Approve - Can approve</span>
          </div>
        </div>
      </Card>
    </div>
  );
}
