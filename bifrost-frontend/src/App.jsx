import React, { useEffect, useState } from 'react';
import VmList from './components/VmList';
import VmDetails from './components/VmDetails';

function App() {
  const [vms, setVms] = useState([]);
  const [selectedVm, setSelectedVm] = useState(null);

  useEffect(() => {
    fetch('/api/v1/vms')
      .then((res) => res.json())
      .then((data) => setVms(data))
      .catch((err) => console.error('Erro ao buscar VMs:', err));
  }, []);

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-3xl font-bold mb-4 text-center">Bifrost VM Dashboard</h1>
      <div className="flex flex-col md:flex-row gap-4">
        <VmList vms={vms} onSelectVm={setSelectedVm} />
        {selectedVm && <VmDetails vm={selectedVm} />}
      </div>
    </div>
  );
}

export default App;
