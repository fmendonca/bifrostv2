import React, { useEffect, useState } from 'react';
import VmList from './components/VmList';
import VmDetails from './components/VmDetails';

function App() {
  const [vms, setVms] = useState([]);
  const [selectedVm, setSelectedVm] = useState(null);

  // Lê a variável de ambiente REACT_APP_API_URL (setada no build)
  const API_URL = process.env.REACT_APP_API_URL || '';

  useEffect(() => {
    fetch(`${API_URL}/api/v1/vms`)
      .then((res) => {
        if (!res.ok) {
          throw new Error(`Erro HTTP ${res.status}`);
        }
        return res.json();
      })
      .then((data) => setVms(data))
      .catch((err) => console.error('Erro ao buscar VMs:', err));
  }, [API_URL]);

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-3xl font-bold mb-4 text-center text-blue-700">Bifrost VM Dashboard</h1>
      <div className="flex flex-col md:flex-row gap-4">
        <VmList vms={vms} onSelectVm={setSelectedVm} />
        {selectedVm && <VmDetails vm={selectedVm} />}
      </div>
    </div>
  );
}

export default App;
