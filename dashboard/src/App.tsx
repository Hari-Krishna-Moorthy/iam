import { useEffect, useState } from 'react';
import Login from './Login';
import Dashboard from './Dashboard';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
  const [darkMode, setDarkMode] = useState(localStorage.getItem('theme') === 'dark');

  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }, [darkMode]);

  const handleLogin = () => {
    setIsLoggedIn(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('targetTenant');
    setIsLoggedIn(false);
  };

  const toggleDarkMode = () => setDarkMode(!darkMode);

  return (
    <div className="App">
      {isLoggedIn ? (
        <Dashboard 
          onLogout={handleLogout} 
          darkMode={darkMode} 
          toggleDarkMode={toggleDarkMode} 
        />
      ) : (
        <Login onLogin={handleLogin} />
      )}
    </div>
  );
}

export default App;
