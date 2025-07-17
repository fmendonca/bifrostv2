import React from 'react';

function VMCreate() {
  const handleSubmit = (e) => {
    e.preventDefault();
    // TODO: Chamar backend para criar VM
    alert('VM criada (mock)');
  };

  return (
    <div>
      <h1 className="text-xl font-bold mb-4">Criar Nova VM</h1>
      <form onSubmit={handleSubmit}>
        <input className="block mb-2 border p-2 w-full" placeholder="Nome da VM" required />
        <input className="block mb-2 border p-2 w-full" placeholder="CPU" required />
        <input className="block mb-2 border p-2 w-full" placeholder="Memória (MB)" required />
        <button className="bg-green-500 text-white px-4 py-2 rounded" type="submit">
          Criar
        </button>
      </form>
    </div>
  );
}

export default VMCreate;
