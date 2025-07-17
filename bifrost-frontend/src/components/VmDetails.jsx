import React from 'react';

function VmDetails({ vm }) {
  const safeJson = (data) => {
    try {
      return JSON.stringify(data, null, 2);
    } catch {
      return 'Formato inválido';
    }
  };

  const getStatusColor = (state) => {
    if (state.includes('running')) return 'bg-bifrostGreen';
    if (state.includes('shut')) return 'bg-bifrostRed';
    if (state.includes('paused')) return 'bg-bifrostYellow';
    return 'bg-gray-500';
  };

  return (
    <div className="w-full md:w-1/2 p-4 bg-white rounded shadow">
      <h2 className="text-xl font-semibold mb-4 border-b pb-2">Detalhes da VM</h2>
      <div className="space-y-1">
        <p>
          <span className="font-semibold">Nome:</span> {vm.name}
        </p>
        <p>
          <span className="font-semibold">UUID:</span> {vm.uuid}
        </p>
        <p className="flex items-center space-x-2">
          <span className="font-semibold">Estado:</span>
          <span className={`px-2 py-0.5 text-white text-xs rounded ${getStatusColor(vm.state)}`}>
            {vm.state}
          </span>
        </p>
        <p>
          <span className="font-semibold">CPU:</span> {vm.cpu_allocation}
        </p>
        <p>
          <span className="font-semibold">Memória (MB):</span> {vm.memory_allocation}
        </p>
      </div>

      <div className="mt-4">
        <p className="font-semibold mb-1">Discos:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm max-h-40 overflow-y-auto">
          {safeJson(vm.disks)}
        </pre>
      </div>

      <div className="mt-4">
        <p className="font-semibold mb-1">Interfaces:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm max-h-40 overflow-y-auto">
          {safeJson(vm.interfaces)}
        </pre>
      </div>

      <div className="mt-4">
        <p className="font-semibold mb-1">Metadata:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm max-h-40 overflow-y-auto">
          {safeJson(vm.metadata)}
        </pre>
      </div>
    </div>
  );
}

export default VmDetails;
