import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Home, Server, Cpu, Sun, Moon, Menu, X } from 'lucide-react';
import { useDarkMode } from './DarkModeContext';

function Sidebar() {
  const location = useLocation();
  const { dark, toggleDark } = useDarkMode();
  const [open, setOpen] = useState(true);

  const menuItems = [
    { path: '/', label: 'Home', icon: <Home size={18} /> },
    { path: '/vms', label: 'VMs', icon: <Cpu size={18} /> },
    { path: '/hosts', label: 'Hosts', icon: <Server size={18} /> },
  ];

  return (
    <div className={`h-screen ${open ? 'w-64' : 'w-16'} bg-gray-900 text-gray-200 flex flex-col transition-all duration-300 shadow-lg`}>
      <div className="flex justify-between items-center p-4 border-b border-gray-700">
        <span className="text-xl font-extrabold text-bifrostBlue">
          {open ? '⚡ Bifrost' : '⚡'}
        </span>
        <button onClick={() => setOpen(!open)} className="text-gray-400 hover:text-white">
          {open ? <X size={20} /> : <Menu size={20} />}
        </button>
      </div>

      <nav className="flex-1 p-2 space-y-1">
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
              <span>{item.icon}</span>
              {open && <span className="ml-3 text-sm font-medium">{item.label}</span>}
            </Link>
          );
        })}
      </nav>

      <div className="p-4 flex justify-center border-t border-gray-700">
        <button onClick={toggleDark} className="text-gray-400 hover:text-white">
          {dark ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </div>
    </div>
  );
}

export default Sidebar;
