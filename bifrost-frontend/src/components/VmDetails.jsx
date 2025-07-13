import React from 'react';

function VmDetails({ vm }) {
  return (
    <div className="w-full md:w-1/2 bg-white rounded-lg shadow p-4">
      <h2 className="text-xl font-semibold mb-2">Detalhes da VM: {vm.name}</h2>
      <p><strong>UUID:</strong> {vm.uuid}</p>
      <p><strong>Estado:</strong> {vm.state}</p>
      <p><strong>CPU:</strong> {vm.cpu_allocation}</p>
      <p><strong>Memória:</strong> {vm.memory_allocation} MB</p>
      <div>
        <strong>Discos:</strong>
        <ul className="list-disc ml-6">
          {vm.disks.map((d, i) => (
            <li key={i}>{d.device}: {d.path}</li>
          ))}
        </ul>
      </div>
      <div>
        <strong>Interfaces:</strong>
        <ul className="list-disc ml-6">
          {vm.interfaces.map((iface, i) => (
            <li key={i}>{iface.name} ({iface.mac}) — {iface.addrs.join(', ')}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default VmDetails;
