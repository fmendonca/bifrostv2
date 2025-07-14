import React, { useEffect, useState } from 'react';
import VmList from './components/VmList';
import VmDetails from './components/VmDetails';

function App() {
  const [vms, setVms] = useState([]);
  const [selectedVm, setSelectedVm] = useState(null);

  const API_URL = process.env.REACT_APP_API_URL || '';

  const fetchVms = () => {
    fetch(`${API_URL}/api/v1/vms`)
      .then((res) => {
        if (!res.ok) {
          throw new Error(`Erro HTTP ${res.status}`);
        }
        return res.json();
      })
      .then((data) => setVms(data))
      .catch((err) => console.error('Erro ao buscar VMs:', err));
  };

  useEffect(() => {
    fetchVms();
  }, [API_URL]);

  const handleAction = (uuid, action) => {
    fetch(`${API_URL}/api/v1/vms/${uuid}/${action}`, { method: 'POST' })
      .then((res) => {
        if (!res.ok) {
          throw new Error(`Erro ao ${action} VM: HTTP ${res.status}`);
        }
        return res.text();
      })
      .then((message) => {
        console.log(message);
        fetchVms();
      })
      .catch((err) => console.error(`Erro ao ${action} VM:`, err));
  };

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-3xl font-bold mb-4 text-center text-blue-700">Bifrost VM Dashboard</h1>
      <div className="flex flex-col md:flex-row gap-4">
        <VmList vms={vms} onSelectVm={setSelectedVm} onAction={handleAction} />
        {selectedVm && <VmDetails vm={selectedVm} />}
      </div>
    </div>
  );
}

export default App;
