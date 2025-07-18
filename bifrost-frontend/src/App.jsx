import React, { useEffect, useState } from 'react';
import VmList from './components/VmList';
import VmDetails from './components/VmDetails';
import Sidebar from './components/Sidebar';
import Spinner from './components/Spinner';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  const [vms, setVms] = useState([]);
  const [selectedVm, setSelectedVm] = useState(null);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [apiKey, setApiKey] = useState('');

  const API_URL = process.env.REACT_APP_API_URL || '';
  const FRONTEND_SECRET = process.env.REACT_APP_FRONTEND_SECRET || 'meuSegredoForte';

  const fetchApiKey = async () => {
    try {
      const res = await fetch(`${API_URL}/api/v1/agent/frontend-key?secret=${FRONTEND_SECRET}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setApiKey(data.api_key);
      console.log('✅ API key obtida com sucesso');
    } catch (err) {
      console.error('❌ Erro ao obter API key:', err);
      toast.error('Erro ao obter API key');
    }
  };

  const fetchVMs = async () => {
    try {
      const res = await fetch(`${API_URL}/api/v1/vms`, {
        headers: { 'X-API-KEY': apiKey }
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      const sorted = data.sort((a, b) => a.name.localeCompare(b.name));
      setVms(sorted);
      if (selectedVm) {
        const updated = data.find((vm) => vm.uuid === selectedVm.uuid);
        setSelectedVm(updated || null);
      }
    } catch (err) {
      console.error('❌ Erro ao buscar VMs:', err);
      toast.error('Erro ao buscar VMs');
    } finally {
      setInitialLoading(false);
    }
  };

  useEffect(() => {
    const init = async () => {
      await fetchApiKey();
    };
    init();
  }, []);

  useEffect(() => {
    if (!apiKey) return;
    fetchVMs();
    const interval = setInterval(fetchVMs, 5000);
    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [apiKey]);

  const handleAction = async (uuid, action) => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/api/v1/vms/${uuid}/${action}`, {
        method: 'POST',
        headers: { 'X-API-KEY': apiKey }
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      toast.success(`✅ ${action.toUpperCase()} enviado para ${uuid}`);
      await fetchVMs();
    } catch (err) {
      console.error(`❌ Erro ao enviar ação ${action}:`, err);
      toast.error(`Erro ao executar ${action} na VM`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />
      <div className="flex-1 p-4">
        <ToastContainer position="top-right" autoClose={3000} />
        <h1 className="text-3xl font-bold mb-4 text-center text-bifrostBlue">
          Bifrost VM Dashboard
        </h1>

        {initialLoading ? (
          <Spinner />
        ) : (
          <>
            {loading && (
              <div className="fixed inset-0 bg-black bg-opacity-25 flex justify-center items-center z-50">
                <Spinner />
              </div>
            )}

            <div className="flex flex-col md:flex-row gap-4">
              <VmList
                vms={vms}
                onSelectVm={setSelectedVm}
                onAction={handleAction}
                loading={loading}
              />
              {selectedVm && (
                <VmDetails
                  vm={{ ...selectedVm, onAction: handleAction }}
                />
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;
