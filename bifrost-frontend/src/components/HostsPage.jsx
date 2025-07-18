import React, { useState } from 'react';
import { toast } from 'react-toastify';

function HostsPage({ apiUrl }) {
  const [hostName, setHostName] = useState('');
  const [registering, setRegistering] = useState(false);

  const handleRegister = async (e) => {
    e.preventDefault();
    if (!hostName.trim()) {
      toast.error('Nome do host não pode ser vazio');
      return;
    }
    setRegistering(true);
    try {
      const res = await fetch(`${apiUrl}/api/v1/agent/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: hostName.trim() }),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      toast.success(`✅ Host registrado: ${data.name}`);
      toast.info(`📦 API Key: ${data.api_key}`);
      toast.info(`📢 Redis Channel: ${data.redis_channel}`);
    } catch (err) {
      console.error(err);
      toast.error('❌ Falha ao registrar host');
    } finally {
      setRegistering(false);
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white p-4 rounded shadow">
      <h2 className="text-xl font-semibold mb-4">Registrar Novo Host</h2>
      <form onSubmit={handleRegister} className="space-y-4">
        <input
          type="text"
          placeholder="Nome do Host"
          className="w-full border rounded p-2"
          value={hostName}
          onChange={(e) => setHostName(e.target.value)}
          disabled={registering}
        />
        <button
          type="submit"
          className="w-full bg-bifrostBlue text-white py-2 rounded hover:bg-blue-800 disabled:opacity-50"
          disabled={registering}
        >
          {registering ? 'Registrando...' : 'Registrar'}
        </button>
      </form>
      <p className="mt-4 text-sm text-gray-600">
        Após o registro, copie e configure a API Key no agente correspondente.
      </p>
    </div>
  );
}

export default HostsPage;
