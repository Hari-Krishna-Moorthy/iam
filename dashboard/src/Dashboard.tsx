import React, { useEffect, useState } from 'react';
import api from './api';
import { Users, LogOut, Shield, ArrowRightLeft, Activity, Moon, Sun } from 'lucide-react';

interface DashboardProps {
  onLogout: () => void;
  darkMode: boolean;
  toggleDarkMode: () => void;
}

const Dashboard: React.FC<DashboardProps> = ({ onLogout, darkMode, toggleDarkMode }) => {
  const [users, setUsers] = useState<any[]>([]);
  const [tenants, setTenants] = useState<any[]>([]);
  const [auditLogs, setAuditLogs] = useState<any[]>([]);
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [targetTenant, setTargetTenant] = useState(localStorage.getItem('targetTenant') || '');
  const [activeTab, setActiveTab] = useState<'users' | 'audit'>('users');

  useEffect(() => {
    fetchUsers();
    fetchAuditLogs();
    checkSuperAdmin();
  }, [targetTenant]);

  const checkSuperAdmin = async () => {
    try {
      await api.get('/me');
      const tenantsResp = await api.get('/tenants');
      setTenants(tenantsResp.data);
      setIsSuperAdmin(true);
    } catch (e) {
      setIsSuperAdmin(false);
    }
  };

  const fetchUsers = async () => {
    try {
      const resp = await api.get('/users');
      setUsers(resp.data || []);
    } catch (err) {
      console.error('Failed to fetch users', err);
    }
  };

  const fetchAuditLogs = async () => {
    try {
      const resp = await api.get('/audit');
      setAuditLogs(resp.data || []);
    } catch (err) {
      console.error('Failed to fetch audit logs', err);
    }
  };

  const handleBulkCreate = async () => {
    const usernames = prompt("Enter usernames separated by comma:");
    if (!usernames) return;
    
    const newUsers = usernames.split(',').map(u => ({
      username: u.trim(),
      email: `${u.trim()}@example.com`,
      password: 'P@ssword123!', // Ensure it passes policy
      role_id: 'default-role-id' // Placeholder
    }));

    try {
      const resp = await api.post('/bulk/users/create', { users: newUsers });
      alert(`Bulk job submitted! Job ID: ${resp.data.job_id}`);
    } catch (err) {
      alert('Failed to submit bulk job');
    }
  };

  const handleTenantSwitch = (tid: string) => {
    setTargetTenant(tid);
    if (tid) {
      localStorage.setItem('targetTenant', tid);
    } else {
      localStorage.removeItem('targetTenant');
    }
  };

  return (
    <div className="flex h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
      {/* Sidebar */}
      <div className="w-64 bg-indigo-900 dark:bg-indigo-950 text-white flex flex-col transition-colors duration-200">
        <div className="p-6 text-xl font-bold flex items-center gap-2">
          <Shield size={24} />
          <span>IAM Admin</span>
        </div>
        <nav className="flex-1 px-4 py-4 space-y-2">
          <button 
            onClick={() => setActiveTab('users')}
            className={`w-full flex items-center gap-3 p-2 rounded transition-colors ${activeTab === 'users' ? 'bg-indigo-800' : 'hover:bg-indigo-800/50'}`}
          >
            <Users size={20} />
            <span>Users</span>
          </button>
          <button 
            onClick={() => setActiveTab('audit')}
            className={`w-full flex items-center gap-3 p-2 rounded transition-colors ${activeTab === 'audit' ? 'bg-indigo-800' : 'hover:bg-indigo-800/50'}`}
          >
            <Activity size={20} />
            <span>Audit Logs</span>
          </button>
        </nav>
        <button 
          onClick={onLogout}
          className="p-4 border-t border-indigo-800 flex items-center gap-3 hover:bg-indigo-800 transition-colors"
        >
          <LogOut size={20} />
          <span>Logout</span>
        </button>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="bg-white dark:bg-gray-800 shadow-sm p-4 flex justify-between items-center transition-colors duration-200">
          <div className="flex items-center gap-4">
            <h2 className="text-xl font-semibold text-gray-800 dark:text-gray-100">
              {activeTab === 'users' ? 'User Management' : 'Audit Logs'}
            </h2>
            {activeTab === 'users' && (
              <button 
                onClick={handleBulkCreate}
                className="bg-indigo-600 text-white px-3 py-1 rounded text-sm hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600 transition-colors"
              >
                Bulk Create
              </button>
            )}
          </div>
          
          <div className="flex items-center gap-4">
            {isSuperAdmin && (
              <div className="flex items-center gap-2">
                <ArrowRightLeft size={18} className="text-gray-500 dark:text-gray-400" />
                <select 
                  className="border rounded p-1 text-sm bg-gray-50 dark:bg-gray-700 dark:border-gray-600 dark:text-gray-200 transition-colors"
                  value={targetTenant}
                  onChange={(e) => handleTenantSwitch(e.target.value)}
                >
                  <option value="">Default Tenant</option>
                  {tenants.map(t => (
                    <option key={t.id} value={t.id}>{t.name}</option>
                  ))}
                </select>
              </div>
            )}
            <button 
              onClick={toggleDarkMode}
              className="p-2 text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-400 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              {darkMode ? <Sun size={20} /> : <Moon size={20} />}
            </button>
          </div>
        </header>

        <main className="p-6 overflow-y-auto flex-1">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden transition-colors duration-200">
            
            {activeTab === 'users' ? (
              <table className="w-full text-left">
                <thead className="bg-gray-50 dark:bg-gray-900/50 border-b dark:border-gray-700">
                  <tr>
                    <th className="px-6 py-3 text-sm font-medium text-gray-500 dark:text-gray-400 uppercase">Username</th>
                    <th className="px-6 py-3 text-sm font-medium text-gray-500 dark:text-gray-400 uppercase">Email</th>
                    <th className="px-6 py-3 text-sm font-medium text-gray-500 dark:text-gray-400 uppercase">Role ID</th>
                  </tr>
                </thead>
                <tbody className="divide-y dark:divide-gray-700">
                  {users.length === 0 ? (
                    <tr><td colSpan={3} className="px-6 py-4 text-center text-gray-500 dark:text-gray-400">No users found</td></tr>
                  ) : users.map(u => (
                    <tr key={u.ID} className="hover:bg-gray-50 dark:hover:bg-gray-750">
                      <td className="px-6 py-4 dark:text-gray-200">{u.Username}</td>
                      <td className="px-6 py-4 dark:text-gray-200">{u.Email}</td>
                      <td className="px-6 py-4">
                        <span className="px-2 py-1 bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 rounded-full text-xs font-medium">
                          {u.RoleID}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <table className="w-full text-left text-sm">
                <thead className="bg-gray-50 dark:bg-gray-900/50 border-b dark:border-gray-700">
                  <tr>
                    <th className="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 uppercase">Time</th>
                    <th className="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 uppercase">Action</th>
                    <th className="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 uppercase">Resource</th>
                    <th className="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 uppercase">User ID</th>
                  </tr>
                </thead>
                <tbody className="divide-y dark:divide-gray-700">
                  {auditLogs.length === 0 ? (
                    <tr><td colSpan={4} className="px-6 py-4 text-center text-gray-500 dark:text-gray-400">No audit logs found</td></tr>
                  ) : auditLogs.map(log => (
                    <tr key={log.ID} className="hover:bg-gray-50 dark:hover:bg-gray-750">
                      <td className="px-6 py-4 whitespace-nowrap text-gray-500 dark:text-gray-400">
                        {new Date(log.CreatedAt).toLocaleString()}
                      </td>
                      <td className="px-6 py-4 font-medium dark:text-gray-200">{log.Action}</td>
                      <td className="px-6 py-4 dark:text-gray-300">{log.Resource}</td>
                      <td className="px-6 py-4 text-gray-500 dark:text-gray-400">{log.UserID || 'N/A'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

          </div>
        </main>
      </div>
    </div>
  );
};

export default Dashboard;
