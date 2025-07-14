import React from 'react';

function VmList({ vms, onSelectVm, onAction }) {
  return (
    <div className="w-full md:w-1/2">
      <ul>
        {vms.map((vm) => (
          <li key={vm.uuid} className="p-4 mb-2 bg-white rounded shadow">
            <div className="flex justify-between items-center">
              <div onClick={() => onSelectVm(vm)} className="cursor-pointer">
                <p className="font-semibold">{vm.name}</p>
                <p className="text-sm text-gray-600">Estado: {vm.state}</p>
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => onAction(vm.uuid, 'start')}
                  className="bg-green-500 hover:bg-green-600 text-white px-2 py-1 rounded"
                >
                  Start
                </button>
                <button
                  onClick={() => onAction(vm.uuid, 'stop')}
                  className="bg-red-500 hover:bg-red-600 text-white px-2 py-1 rounded"
                >
                  Stop
                </button>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VmList;
