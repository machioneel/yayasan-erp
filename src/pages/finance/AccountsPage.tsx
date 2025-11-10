import { useQuery } from '@tanstack/react-query';
import { api } from '@/services/api';
import { Account, ApiResponse } from '@/types';
import { Card } from '@/components/common/Card';
import { Badge } from '@/components/common/Badge';
import { ChevronRight, ChevronDown } from 'lucide-react';
import { formatCurrency } from '@/utils/format';
import { useState } from 'react';

export default function AccountsPage() {
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  const { data: accounts, isLoading } = useQuery({
    queryKey: ['accounts'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<Account[]>>('/accounts');
      return response.data.data;
    },
  });

  const toggleExpand = (id: string) => {
    const newExpanded = new Set(expanded);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpanded(newExpanded);
  };

  const renderAccount = (account: Account, children: Account[], level: number = 0) => {
    const hasChildren = children.length > 0;
    const isExpanded = expanded.has(account.id);

    return (
      <div key={account.id}>
        <div
          className={`flex items-center justify-between p-3 hover:bg-gray-50 border-b cursor-pointer`}
          style={{ paddingLeft: `${level * 24 + 12}px` }}
          onClick={() => hasChildren && toggleExpand(account.id)}
        >
          <div className="flex items-center gap-3 flex-1">
            {hasChildren && (
              <button className="text-gray-400">
                {isExpanded ? (
                  <ChevronDown className="w-4 h-4" />
                ) : (
                  <ChevronRight className="w-4 h-4" />
                )}
              </button>
            )}
            {!hasChildren && <div className="w-4" />}
            
            <div className="flex items-center gap-3">
              <span className="font-mono text-sm text-gray-600">{account.code}</span>
              <span className={account.is_header ? 'font-semibold text-gray-900' : 'text-gray-700'}>
                {account.name}
              </span>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <Badge variant={account.normal_balance === 'debit' ? 'info' : 'success'}>
              {account.normal_balance === 'debit' ? 'Debit' : 'Kredit'}
            </Badge>
            
            {!account.is_header && account.balance !== undefined && (
              <span className="font-medium text-gray-900 w-40 text-right">
                {formatCurrency(account.balance)}
              </span>
            )}

            <Badge variant={account.is_active ? 'success' : 'default'}>
              {account.is_active ? 'Aktif' : 'Nonaktif'}
            </Badge>
          </div>
        </div>

        {hasChildren && isExpanded && (
          <div>
            {children.map(child => {
              const childChildren = (accounts || []).filter(a => a.parent_id === child.id);
              return renderAccount(child, childChildren, level + 1);
            })}
          </div>
        )}
      </div>
    );
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const rootAccounts = (accounts || []).filter(a => !a.parent_id);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Chart of Accounts</h1>
        <p className="text-gray-600 mt-1">Bagan akun keuangan</p>
      </div>

      <Card>
        <div className="overflow-x-auto">
          <div className="min-w-full">
            {rootAccounts.map(account => {
              const children = (accounts || []).filter(a => a.parent_id === account.id);
              return renderAccount(account, children);
            })}
          </div>
        </div>
      </Card>
    </div>
  );
}
