import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api'; // usa o client axios/fetch já existente
import { toast } from 'react-toastify';

function VMList() {
  const [vms, setVms] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/vms')
      .then((response) => setVms(response.data))
      .catch((error) => {
        console.error('Erro ao carregar VMs:', error);
        toast.error('Erro ao carregar VMs');
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Máquinas Virtuais</h1>
      <Link to="/compute/create" className="bg-blue-600 text-white px-4 py-2 rounded">
        + Nova VM
      </Link>
      {loading ? (
        <p className="mt-4">Carregando...</p>
      ) : (
        <ul className="mt-4 divide-y divide-gray-200">
          {vms.map((vm) => (
            <li key={vm.id} className="py-2 flex justify-between items-center">
              <Link to={`/compute/${vm.id}`} className="text-blue-600 hover:underline">
                {vm.name} ({vm.status})
              </Link>
              <span className={`px-2 py-1 text-sm rounded ${vm.status === 'running' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                {vm.status}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default VMList;
