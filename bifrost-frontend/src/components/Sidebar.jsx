import React from 'react';
import { Link, useLocation } from 'react-router-dom';

function Sidebar() {
  const location = useLocation();

  const menuItem = (path, label) => (
    <Link
      to={path}
      className={`block p-2 rounded hover:bg-gray-700 ${
        location.pathname === path ? 'bg-gray-700' : ''
      }`}
    >
      {label}
    </Link>
  );

  return (
    <div className="h-full w-64 bg-gray-800 text-white flex flex-col">
      <div className="p-4 text-2xl font-bold border-b border-gray-700 text-center">
        Bifrost
      </div>
      <nav className="flex-1 p-2 space-y-1">
        {menuItem('/', '🏠 Home')}
        {menuItem('/vms', '🖥️ VMs')}
        {menuItem('/hosts', '🛠️ Hosts')}
        {/* Você pode adicionar mais rotas aqui */}
      </nav>
    </div>
  );
}

export default Sidebar;
