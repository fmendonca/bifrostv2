import React, { useEffect, useState } from 'react';
import VMList from './pages/Compute/VmList';
import VMDetails from './pages/Compute/VmDetails';
import VMConsole from './pages/Compute/VMConsole';
import Sidebar from './components/Sidebar';
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

      // Se tinha uma VM selecionada, atualiza ela
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
    fetchVMs();
    const interval = setInterval(fetchVMs, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleAction = async (uuid, action) => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/api/v1/vms/${uuid}/${action}`, {
        method: 'POST',
      });
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

  const handleSelectVm = async (vm) => {
    try {
      const res = await fetch(`${API_URL}/api/v1/vms/${vm.uuid}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const detailedVm = await res.json();
      setSelectedVm({ ...detailedVm, onAction: handleAction });
    } catch (err) {
      console.error('❌ Erro ao buscar detalhes da VM:', err);
      toast.error('Erro ao buscar detalhes da VM');
    }
  };
  
  return (
    <div className="min-h-screen bg-gray-100 flex">
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
              <VMList
                vms={vms}
                onSelectVm={setSelectedVm}
                onAction={handleAction}
                loading={loading}
              />
              {selectedVm && (
                <VMDetails
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
