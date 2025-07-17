import React from 'react';

function VmList({ vms, onSelectVm, loading }) {
  return (
    <div className="w-full md:w-1/2 p-4 bg-white rounded shadow">
      <h2 className="text-xl font-semibold mb-4 border-b pb-2">Lista de VMs</h2>
      {vms.length === 0 ? (
        <p className="text-gray-500">Nenhuma VM encontrada.</p>
      ) : (
        <ul className="space-y-2">
          {vms.map((vm) => (
            <li
              key={vm.uuid}
              className={`p-2 border rounded cursor-pointer ${
                loading ? 'opacity-50 pointer-events-none' : 'hover:bg-gray-100'
              }`}
              onClick={() => onSelectVm(vm)} // 👉 passa o OBJETO COMPLETO
            >
              <div className="flex justify-between items-center">
                <span>{vm.name}</span>
                <span
                  className={`px-2 py-0.5 text-white text-xs rounded ${
                    vm.state.includes('running')
                      ? 'bg-bifrostGreen'
                      : vm.state.includes('shut')
                      ? 'bg-bifrostRed'
                      : vm.state.includes('paused')
                      ? 'bg-bifrostYellow'
                      : 'bg-gray-500'
                  }`}
                >
                  {vm.state}
                </span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default VmList;
