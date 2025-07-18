import React, { createContext, useContext, useState } from 'react';

const DarkModeContext = createContext();

export function DarkModeProvider({ children }) {
  const [dark, setDark] = useState(false);
  const toggleDark = () => setDark((prev) => !prev);

  return (
    <DarkModeContext.Provider value={{ dark, toggleDark }}>
      <div className={dark ? 'dark' : ''}>
        {children}
      </div>
    </DarkModeContext.Provider>
  );
}

export function useDarkMode() {
  return useContext(DarkModeContext);
}
