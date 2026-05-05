import { useEffect, useState, useCallback } from 'react';

type ToastInfo = { id: number; message: string; type: 'success' | 'error' };

function App() {
  const [isShieldActive, setIsShieldActive] = useState(false);
  const [toasts, setToasts] = useState<ToastInfo[]>([]);
  const [relationalGlitch, setRelationalGlitch] = useState(false);
  const [goGlitch, setGoGlitch] = useState(false);

  // Mock data to simulate tables
  const initialRelational = [
    { id: 1, username: 'admin', data: 'Super Secret Setup Data' },
    { id: 2, username: 'user1', data: 'Hello World! Welcome' },
  ];
  
  const initialGoData = [
    { id: 3, username: 'go_admin', data: 'System logs initialized' },
    { id: 4, username: 'go_user', data: 'Cache warming complete' },
  ];

  const [relationalData] = useState(initialRelational);
  const [goData] = useState(initialGoData);

  const addToast = useCallback((message: string, type: 'success' | 'error') => {
    const id = Date.now();
    setToasts((prev) => [...prev, { id, message, type }]);
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 3000);
  }, []);

  // The Invincible Switch - Hidden keyboard shortcut Alt + X
  useEffect(() => {
    const handleKeyDown = async (e: KeyboardEvent) => {
      if (e.altKey && e.key.toLowerCase() === 'x') {
        const newShieldState = !isShieldActive;
        setIsShieldActive(newShieldState);
        
        // Optimistically set UI, then fire API calls in background
        try {
          await Promise.allSettled([
            fetch('http://localhost:8080/api/toggle-shield', { method: 'POST' }), // Java Spring
            fetch('http://localhost:8081/api/shield/toggle', { method: 'POST' })  // Go Gin
          ]);
        } catch (err) {
          console.error("Failed to toggle shield on backend nodes", err);
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isShieldActive]);

  const triggerAttack = (target: 'relational' | 'go', attackName: string) => {
    if (isShieldActive) {
      addToast(`Threat Blocked: ${attackName} intercepted by Shield.`, 'success');
      return;
    }

    // Shield is OFF - Vulnerable!
    addToast(`System Compromised: ${attackName} successful!`, 'error');
    
    if (target === 'relational') {
      setRelationalGlitch(true);
      setTimeout(() => setRelationalGlitch(false), 2000);
    } else {
      setGoGlitch(true);
      setTimeout(() => setGoGlitch(false), 2000);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 font-sans selection:bg-indigo-100">
      {/* Navigation Bar */}
      <nav className="bg-white shadow-sm border-b border-slate-200 px-6 py-4 flex justify-between items-center sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <img src="/favicon.png" alt="Shield Logo" className="w-8 h-8 rounded-lg shadow-sm" />
          <h1 className="text-xl font-bold tracking-tight text-slate-800">Threat Visualizer</h1>
        </div>
        <div className="flex items-center gap-4 text-sm font-medium text-slate-500">
          <span>Status:</span>
          <span className={`px-4 py-1.5 rounded-full text-xs font-bold tracking-wider uppercase shadow-sm transition-all duration-300 ${isShieldActive ? 'bg-emerald-100 text-emerald-700 border border-emerald-200' : 'bg-rose-100 text-rose-700 border border-rose-200'}`}>
            {isShieldActive ? 'Protected' : 'Vulnerable'}
          </span>
        </div>
      </nav>

      {/* Main Content: Split Screen */}
      <main className="max-w-7xl mx-auto p-6 grid grid-cols-1 lg:grid-cols-2 gap-8 mt-4">
        
        {/* Left Column: Relational Data (MySQL / Java) */}
        <section className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden flex flex-col transition-shadow hover:shadow-md">
          <div className="px-6 py-5 border-b border-slate-100 bg-slate-50/50">
            <h2 className="text-lg font-semibold text-slate-800">Relational Data (MySQL)</h2>
            <p className="text-sm text-slate-500 mt-1">Managed by Java Spring Boot</p>
          </div>
          
          <div className="p-6 flex-grow">
            <div className="overflow-x-auto rounded-lg border border-slate-200 shadow-sm">
              <table className="w-full text-sm text-left">
                <thead className="bg-slate-50 text-slate-600 font-medium border-b border-slate-200">
                  <tr>
                    <th className="px-4 py-3">ID</th>
                    <th className="px-4 py-3">Username</th>
                    <th className="px-4 py-3">Data</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {relationalData.map((row) => (
                    <tr key={row.id} className={`vulnerable-row ${relationalGlitch ? 'glitching' : ''}`}>
                      <td className="px-4 py-3 font-medium text-slate-500">{row.id}</td>
                      <td className="px-4 py-3 font-medium text-indigo-900">{row.username}</td>
                      <td className="px-4 py-3 text-slate-600">{row.data}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div className="px-6 py-5 bg-slate-50 border-t border-slate-100 flex flex-wrap gap-3">
            <button 
              onClick={() => triggerAttack('relational', 'Wipeout (TRUNCATE)')}
              className="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white text-sm font-medium rounded-lg shadow-sm transition-all hover:shadow focus:ring-2 focus:ring-rose-500 focus:ring-offset-2"
            >
              Launch Wipeout
            </button>
            <button 
              onClick={() => triggerAttack('relational', 'SQL Injection')}
              className="px-4 py-2 bg-rose-600 hover:bg-rose-700 text-white text-sm font-medium rounded-lg shadow-sm transition-all hover:shadow focus:ring-2 focus:ring-rose-500 focus:ring-offset-2"
            >
              Launch SQLi
            </button>
          </div>
        </section>

        {/* Right Column: Relational Data (MySQL / Go) */}
        <section className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden flex flex-col transition-shadow hover:shadow-md">
          <div className="px-6 py-5 border-b border-slate-100 bg-slate-50/50">
            <h2 className="text-lg font-semibold text-slate-800">Relational Data (MySQL)</h2>
            <p className="text-sm text-slate-500 mt-1">Managed by Go Gin</p>
          </div>
          
          <div className="p-6 flex-grow">
            <div className="overflow-x-auto rounded-lg border border-slate-200 shadow-sm">
              <table className="w-full text-sm text-left">
                <thead className="bg-slate-50 text-slate-600 font-medium border-b border-slate-200">
                  <tr>
                    <th className="px-4 py-3">ID</th>
                    <th className="px-4 py-3">Username</th>
                    <th className="px-4 py-3">Data</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {goData.map((row) => (
                    <tr key={row.id} className={`vulnerable-row ${goGlitch ? 'shaking' : ''}`}>
                      <td className="px-4 py-3 font-medium text-slate-500">{row.id}</td>
                      <td className="px-4 py-3 font-medium text-indigo-900">{row.username}</td>
                      <td className="px-4 py-3 text-slate-600">{row.data}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          <div className="px-6 py-5 bg-slate-50 border-t border-slate-100 flex flex-wrap gap-3">
            <button 
              onClick={() => triggerAttack('go', 'Rapid Duplication (DoS)')}
              className="px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white text-sm font-medium rounded-lg shadow-sm transition-all hover:shadow focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
            >
              Launch Duplication
            </button>
            <button 
              onClick={() => triggerAttack('go', 'IDOR Alteration')}
              className="px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white text-sm font-medium rounded-lg shadow-sm transition-all hover:shadow focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"
            >
              Launch Alteration
            </button>
          </div>
        </section>

      </main>

      {/* Toast Notifications */}
      <div className="fixed bottom-6 right-6 flex flex-col gap-3 z-50">
        {toasts.map(toast => (
          <div 
            key={toast.id} 
            className={`px-6 py-4 rounded-xl shadow-lg border text-sm font-medium flex items-center gap-3 transform transition-all duration-300 translate-y-0 opacity-100
              ${toast.type === 'success' 
                ? 'bg-emerald-50 border-emerald-200 text-emerald-800' 
                : 'bg-rose-50 border-rose-200 text-rose-800'}`}
          >
            {toast.type === 'success' ? (
              <svg className="w-5 h-5 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            ) : (
              <svg className="w-5 h-5 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            )}
            {toast.message}
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
