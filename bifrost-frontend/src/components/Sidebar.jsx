import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Menu, X, Home, Server, Settings, Monitor, Sun, Moon } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import clsx from 'clsx';

function Sidebar() {
  const location = useLocation();
  const [open, setOpen] = useState(true);
  const [darkMode, setDarkMode] = useState(false);

  const toggleSidebar = () => setOpen((prev) => !prev);
  const toggleDarkMode = () => {
    setDarkMode((prev) => !prev);
    document.documentElement.classList.toggle('dark');
  };

  const menuItems = [
    { path: '/', label: 'Home', icon: <Home size={20} /> },
    { path: '/vms', label: 'VMs', icon: <Monitor size={20} /> },
    { path: '/hosts', label: 'Hosts', icon: <Server size={20} /> },
    { path: '/settings', label: 'Settings', icon: <Settings size={20} /> },
  ];

  return (
    <div className={clsx('h-screen bg-gray-800 text-white flex flex-col transition-all', open ? 'w-64' : 'w-16')}>
      <div className="flex items-center justify-between p-4 border-b border-gray-700">
        <span className="text-xl font-bold">{open ? 'Bifrost' : 'B'}</span>
        <button onClick={toggleSidebar} className="focus:outline-none">
          {open ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>
      <nav className="flex-1 p-2 space-y-1">
        {menuItems.map(({ path, label, icon }) => (
          <Link
            key={path}
            to={path}
            className={clsx(
              'flex items-center p-2 rounded hover:bg-gray-700 transition',
              location.pathname === path && 'bg-gray-700'
            )}
          >
            {icon}
            <AnimatePresence>
              {open && (
                <motion.span
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -10 }}
                  className="ml-3"
                >
                  {label}
                </motion.span>
              )}
            </AnimatePresence>
          </Link>
        ))}
      </nav>
      <div className="p-4 border-t border-gray-700">
        <button
          onClick={toggleDarkMode}
          className="flex items-center w-full p-2 rounded hover:bg-gray-700 transition"
        >
          {darkMode ? <Sun size={20} /> : <Moon size={20} />}
          {open && <span className="ml-3">{darkMode ? 'Light Mode' : 'Dark Mode'}</span>}
        </button>
      </div>
    </div>
  );
}

export default Sidebar;
