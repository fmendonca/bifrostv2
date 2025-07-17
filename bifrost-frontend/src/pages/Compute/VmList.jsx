import React from 'react';

function VMList({ vms, onSelectVm, loading }) {
  return (
    <div className="flex-1 bg-white shadow p-4 rounded">
      <h2 className="text-xl font-bold mb-4">Máquinas Virtuais</h2>
      <ul className="space-y-2">
        {vms.map((vm) => (
          <li
            key={vm.uuid}
            className={`p-2 border rounded cursor-pointer ${
              loading ? 'opacity-50 pointer-events-none' : 'hover:bg-gray-100'
            }`}
            onClick={() => onSelectVm(vm)}
          >
            <div className="flex justify-between">
              <span>{vm.name}</span>
              <span
                className={`text-sm ${
                  vm.status === 'running'
                    ? 'text-green-600'
                    : vm.status === 'stopped'
                    ? 'text-red-600'
                    : 'text-gray-600'
                }`}
              >
                {vm.status}
              </span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VMList;
