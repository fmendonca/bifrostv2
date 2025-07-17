import React from 'react';

function VMDetails({ vm }) {
  if (!vm) return null;

  return (
    <div className="flex-1 bg-white shadow p-4 rounded">
      <h2 className="text-xl font-bold mb-4">{vm.name}</h2>
      <p>Status: {vm.status}</p>
      <p>CPU: {vm.cpu ?? 'N/A'}</p>
      <p>Memória: {vm.memory ? `${vm.memory} MB` : 'N/A'}</p>

      <div className="mt-4 space-x-2">
        <button onClick={() => vm.onAction(vm.uuid, 'start')} className="bg-green-500 text-white px-4 py-2 rounded">
          Start
        </button>
        <button onClick={() => vm.onAction(vm.uuid, 'stop')} className="bg-red-500 text-white px-4 py-2 rounded">
          Stop
        </button>
        <button onClick={() => vm.onAction(vm.uuid, 'restart')} className="bg-yellow-500 text-white px-4 py-2 rounded">
          Restart
        </button>
      </div>
    </div>
  );
}

export default VMDetails;
