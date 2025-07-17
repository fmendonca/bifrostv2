import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../../services/api';
import { toast } from 'react-toastify';

function VMDetails() {
  const { id } = useParams();
  const [vm, setVm] = useState(null);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);

  const loadVM = () => {
    api.get(`/vms/${id}`)
      .then((response) => setVm(response.data))
      .catch(() => toast.error('Erro ao carregar detalhes da VM'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    loadVM();
  }, [id]);

  const handleAction = (action) => {
    setActionLoading(true);
    api.post(`/vms/${id}/action`, { action })
      .then(() => {
        toast.success(`Ação '${action}' executada`);
        loadVM();
      })
      .catch(() => toast.error(`Erro ao executar '${action}'`))
      .finally(() => setActionLoading(false));
  };

  if (loading) return <p>Carregando detalhes...</p>;
  if (!vm) return <p>VM não encontrada.</p>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">{vm.name}</h1>
      <div className="mb-4">
        <p>Status: {vm.status}</p>
        <p>CPU: {vm.cpu}</p>
        <p>Memória: {vm.memory} MB</p>
      </div>

      <div className="flex gap-2 mb-4">
        <button
          onClick={() => handleAction('start')}
          disabled={actionLoading}
          className="bg-green-500 text-white px-3 py-1 rounded disabled:opacity-50"
        >
          Start
        </button>
        <button
          onClick={() => handleAction('stop')}
          disabled={actionLoading}
          className="bg-red-500 text-white px-3 py-1 rounded disabled:opacity-50"
        >
          Stop
        </button>
        <button
          onClick={() => handleAction('restart')}
          disabled={actionLoading}
          className="bg-yellow-500 text-white px-3 py-1 rounded disabled:opacity-50"
        >
          Restart
        </button>
        <Link
          to={`/compute/${id}/console`}
          className="bg-purple-500 text-white px-3 py-1 rounded"
        >
          Console
        </Link>
      </div>
    </div>
  );
}

export default VMDetails;
