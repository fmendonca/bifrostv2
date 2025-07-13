import React from 'react';

function VmList({ vms, onSelectVm }) {
  return (
    <div className="w-full md:w-1/2 bg-white rounded-lg shadow p-4">
      <h2 className="text-xl font-semibold mb-2">Lista de VMs</h2>
      <ul className="divide-y">
        {vms.map((vm) => (
          <li
            key={vm.uuid}
            className="py-2 cursor-pointer hover:bg-gray-100"
            onClick={() => onSelectVm(vm)}
          >
            <span className="font-medium">{vm.name}</span> — {vm.state}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VmList;
