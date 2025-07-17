import React from 'react';

function VmList({ vms, onSelectVm, onAction, loading }) {
  const getStatusColor = (state) => {
    if (state.includes('running')) return 'bg-green-500';
    if (state.includes('shut')) return 'bg-red-500';
    if (state.includes('paused')) return 'bg-yellow-500';
    return 'bg-gray-500';
  };

  return (
    <div className="w-full md:w-1/2">
      <ul className="bg-white rounded shadow p-4">
        {vms.map((vm) => (
          <li key={vm.uuid} className="flex justify-between items-center mb-2 p-2 border-b">
            <div className="flex items-center space-x-2 cursor-pointer" onClick={() => onSelectVm(vm)}>
              <span className={`w-3 h-3 rounded-full ${getStatusColor(vm.state)}`}></span>
              <span className="font-semibold">{vm.name}</span>
            </div>
            <div className="flex space-x-1">
              <button
                className="px-2 py-1 bg-green-600 text-white rounded text-sm disabled:opacity-50"
                disabled={loading || vm.state.includes('running')}
                onClick={() => onAction(vm.uuid, 'start')}
              >
                Start
              </button>
              <button
                className="px-2 py-1 bg-red-600 text-white rounded text-sm disabled:opacity-50"
                disabled={loading || vm.state.includes('shut')}
                onClick={() => onAction(vm.uuid, 'stop')}
              >
                Stop
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VmList;
