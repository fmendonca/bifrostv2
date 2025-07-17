import React, { useEffect, useState } from 'react';
import VmList from './components/VmList';
import VmDetails from './components/VmDetails';
import Spinner from './components/Spinner';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  const [vms, setVms] = useState([]);
  const [selectedVm, setSelectedVm] = useState(null);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const API_URL = process.env.REACT_APP_API_URL || '';

  const fetchVMs = async () => {
    try {
      const res = await fetch(`${API_URL}/api/v1/vms`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      const sorted = data.sort((a, b) => a.name.localeCompare(b.name));
      setVms(sorted);
    } catch (err) {
      console.error('❌ Erro ao buscar VMs:', err);
      toast.error('Erro ao buscar VMs');
    } finally {
      setInitialLoading(false);
    }
  };

  useEffect(() => {
    fetchVMs();
    const interval = setInterval(fetchVMs, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleAction = async (uuid, action) => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/api/v1/vms/${uuid}/${action}`, { method: 'POST' });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      toast.success(`✅ ${action.toUpperCase()} enviado para ${uuid}`);
      await fetchVMs();
    } catch (err) {
      console.error(`❌ Erro ao enviar ação ${action}:`, err);
      toast.error(`❌ Erro ao executar ${action} na VM`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <ToastContainer position="top-right" autoClose={3000} />
      <h1 className="text-3xl font-bold mb-4 text-center text-bifrostBlue">Bifrost VM Dashboard</h1>

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
            <VmList vms={vms} onSelectVm={setSelectedVm} onAction={handleAction} loading={loading} />
            {selectedVm && <VmDetails vm={selectedVm} />}
          </div>
        </>
      )}
    </div>
  );
}

export default App;
