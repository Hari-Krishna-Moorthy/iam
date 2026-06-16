import { useEffect, useState } from 'react';
import api from './api';
import { Users, LogOut, Shield, ArrowRightLeft } from 'lucide-react';

const Dashboard: React.FC<{ onLogout: () => void }> = ({ onLogout }) => {
  const [users, setUsers] = useState<any[]>([]);
  const [tenants, setTenants] = useState<any[]>([]);
  const [isSuperAdmin, setIsSuperAdmin] = useState(false);
  const [targetTenant, setTargetTenant] = useState(localStorage.getItem('targetTenant') || '');

  useEffect(() => {
    fetchUsers();
    checkSuperAdmin();
  }, [targetTenant]);

  const checkSuperAdmin = async () => {
    try {
      // In a real app, the /me endpoint would return tenant info
      await api.get('/me');
      // If we can fetch tenants, we are probably a super admin
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
      setUsers(resp.data);
    } catch (err) {
      console.error('Failed to fetch users', err);
    }
  };

  const handleBulkCreate = async () => {
    const usernames = prompt("Enter usernames separated by comma:");
    if (!usernames) return;
    
    const users = usernames.split(',').map(u => ({
      username: u.trim(),
      email: `${u.trim()}@example.com`,
      password: 'P@ssword123',
      role_id: 'default-role-id' // Placeholder
    }));

    try {
      const resp = await api.post('/bulk/users/create', { users });
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
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <div className="w-64 bg-indigo-900 text-white flex flex-col">
        <div className="p-6 text-xl font-bold flex items-center gap-2">
          <Shield size={24} />
          <span>IAM Admin</span>
        </div>
        <nav className="flex-1 px-4 py-4 space-y-2">
          <a href="#" className="flex items-center gap-3 p-2 bg-indigo-800 rounded">
            <Users size={20} />
            <span>Users</span>
          </a>
        </nav>
        <button 
          onClick={onLogout}
          className="p-4 border-t border-indigo-800 flex items-center gap-3 hover:bg-indigo-800"
        >
          <LogOut size={20} />
          <span>Logout</span>
        </button>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        <header className="bg-white shadow-sm p-4 flex justify-between items-center">
          <div className="flex items-center gap-4">
            <h2 className="text-xl font-semibold text-gray-800">User Management</h2>
            <button 
              onClick={handleBulkCreate}
              className="bg-indigo-600 text-white px-3 py-1 rounded text-sm hover:bg-indigo-700"
            >
              Bulk Create
            </button>
          </div>
          
          {isSuperAdmin && (
            <div className="flex items-center gap-2">
              <ArrowRightLeft size={18} className="text-gray-500" />
              <select 
                className="border rounded p-1 text-sm bg-gray-50"
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
        </header>

        <main className="p-6">
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <table className="w-full text-left">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-6 py-3 text-sm font-medium text-gray-500 uppercase">Username</th>
                  <th className="px-6 py-3 text-sm font-medium text-gray-500 uppercase">Email</th>
                  <th className="px-6 py-3 text-sm font-medium text-gray-500 uppercase">Role</th>
                  <th className="px-6 py-3 text-sm font-medium text-gray-500 uppercase text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {users.map(u => (
                  <tr key={u.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">{u.username}</td>
                    <td className="px-6 py-4">{u.email}</td>
                    <td className="px-6 py-4">
                      <span className="px-2 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium">
                        {u.role}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button className="text-indigo-600 hover:text-indigo-900 text-sm font-medium">Edit</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </main>
      </div>
    </div>
  );
};

export default Dashboard;
