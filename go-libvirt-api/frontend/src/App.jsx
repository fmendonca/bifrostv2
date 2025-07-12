import React, { useState } from 'react';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';

function App() {
    const [token, setToken] = useState(localStorage.getItem('token') || '');

    if (!token) {
        return <LoginPage setToken={setToken} />;
    }
    return <DashboardPage token={token} />;
}

export default App;
