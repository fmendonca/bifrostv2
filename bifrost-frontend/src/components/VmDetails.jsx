import React from 'react';

function VmDetails({ vm }) {
  return (
    <div className="w-full md:w-1/2 p-4 bg-white rounded shadow">
      <h2 className="text-xl font-semibold mb-2">Detalhes da VM</h2>
      <p><span className="font-semibold">Nome:</span> {vm.name}</p>
      <p><span className="font-semibold">UUID:</span> {vm.uuid}</p>
      <p><span className="font-semibold">Estado:</span> {vm.state}</p>
      <p><span className="font-semibold">CPU:</span> {vm.cpu_allocation}</p>
      <p><span className="font-semibold">Memória (MB):</span> {vm.memory_allocation}</p>
      <div className="mt-2">
        <p className="font-semibold">Discos:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm overflow-x-auto">
          {JSON.stringify(vm.disks, null, 2)}
        </pre>
      </div>
      <div className="mt-2">
        <p className="font-semibold">Interfaces:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm overflow-x-auto">
          {JSON.stringify(vm.interfaces, null, 2)}
        </pre>
      </div>
      <div className="mt-2">
        <p className="font-semibold">Metadata:</p>
        <pre className="bg-gray-100 p-2 rounded text-sm overflow-x-auto">
          {JSON.stringify(vm.metadata, null, 2)}
        </pre>
      </div>
    </div>
  );
}

export default VmDetails;
