import React from 'react';
import { Link, useLocation } from 'react-router-dom';

function Sidebar() {
  const location = useLocation();

  const menuItems = [
    { path: '/', label: '🏠 Home' },
    { path: '/vms', label: '🖥️ VMs' },
    { path: '/hosts', label: '🛠️ Hosts' },
    // { path: '/storage', label: '💾 Storage' },
    // { path: '/network', label: '🌐 Network' },
    // { path: '/logs', label: '📜 Logs' },
  ];

  return (
    <div className="h-full w-64 bg-gray-900 text-gray-200 flex flex-col shadow-lg">
      <div className="p-6 text-2xl font-extrabold text-center text-bifrostBlue border-b border-gray-700">
        ⚡ Bifrost
      </div>

      <nav className="flex-1 p-4 space-y-1">
        {menuItems.map((item) => {
          const active = location.pathname === item.path;
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center p-2 rounded transition duration-200 ${
                active
                  ? 'bg-bifrostBlue text-white'
                  : 'hover:bg-gray-700 hover:text-white'
              }`}
            >
              <span className="text-lg mr-2">{item.label.split(' ')[0]}</span>
              <span className="text-sm font-medium">{item.label.split(' ')[1]}</span>
            </Link>
          );
        })}
      </nav>

      <div className="p-4 text-center text-xs text-gray-500 border-t border-gray-700">
        v1.0.0
      </div>
    </div>
  );
}

export default Sidebar;
