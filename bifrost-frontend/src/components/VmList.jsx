import React, { useState } from 'react';
import { format } from 'date-fns';

function VmList({ vms, onSelectVm, onAction, loading }) {
  const [sortBy, setSortBy] = useState('name');
  const [sortAsc, setSortAsc] = useState(true);

  const getStatusColor = (state) => {
    if (state.includes('running')) return 'bg-bifrostGreen';
    if (state.includes('shut')) return 'bg-bifrostRed';
    if (state.includes('paused')) return 'bg-bifrostYellow';
    return 'bg-gray-500';
  };

  const sortedVms = [...vms].sort((a, b) => {
    let res = 0;
    if (sortBy === 'name') res = a.name.localeCompare(b.name);
    else if (sortBy === 'timestamp') res = new Date(a.timestamp) - new Date(b.timestamp);
    return sortAsc ? res : -res;
  });

  return (
    <div className="w-full md:w-1/2">
      <div className="flex justify-between items-center mb-2">
        <h2 className="text-xl font-semibold">Máquinas Virtuais</h2>
        <div className="flex items-center gap-2">
          <select value={sortBy} onChange={(e) => setSortBy(e.target.value)} className="border rounded p-1 text-sm">
            <option value="name">Nome</option>
            <option value="timestamp">Data</option>
          </select>
          <button onClick={() => setSortAsc(!sortAsc)} className="text-gray-600 hover:text-gray-900 text-sm px-2 py-1 border rounded">
            {sortAsc ? '↑' : '↓'}
          </button>
        </div>
      </div>

      <ul className="bg-white rounded shadow p-4 max-h-[70vh] overflow-y-auto">
        {sortedVms.map((vm) => (
          <li key={vm.uuid} className="flex justify-between items-center mb-2 p-2 border-b hover:bg-gray-50 transition">
            <div className="flex flex-col cursor-pointer w-2/3" onClick={() => onSelectVm(vm)}>
              <div className="flex items-center space-x-2">
                <span className={`w-3 h-3 rounded-full ${getStatusColor(vm.state)}`}></span>
                <span className="font-semibold truncate">{vm.name}</span>
              </div>
              <span className="text-xs text-gray-500">
                {vm.timestamp ? format(new Date(vm.timestamp), 'dd/MM/yyyy HH:mm:ss') : 'Sem data'}
              </span>
            </div>
            <div className="flex space-x-1">
              <button
                className="px-2 py-1 bg-bifrostGreen text-white rounded text-sm disabled:opacity-50 hover:bg-green-700"
                disabled={loading || vm.state.includes('running')}
                onClick={() => onAction(vm.uuid, 'start')}
              >
                Start
              </button>
              <button
                className="px-2 py-1 bg-bifrostRed text-white rounded text-sm disabled:opacity-50 hover:bg-red-700"
                disabled={loading || vm.state.includes('shut')}
                onClick={() => onAction(vm.uuid, 'stop')}
              >
                Stop
              </button>
            </div>
          </li>
        ))}
        {sortedVms.length === 0 && (
          <li className="text-center text-gray-500 py-4">Nenhuma VM encontrada</li>
        )}
      </ul>
    </div>
  );
}

export default VmList;
