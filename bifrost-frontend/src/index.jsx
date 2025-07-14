import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

// Pega o elemento root único no HTML
const container = document.getElementById('root');

// Cria o root React (só uma vez!)
const root = ReactDOM.createRoot(container);

// Renderiza o App
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
