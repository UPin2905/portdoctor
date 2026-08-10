import { useState, useEffect } from 'react';
import { ScanPorts, KillPort, SharePort, StopSharePort } from '../wailsjs/go/main/App';
import { main } from '../wailsjs/go/models';

function App() {
  const [ports, setPorts] = useState<main.UIPortInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');
  const [sharing, setSharing] = useState<Record<number, boolean>>({});

  const loadPorts = async () => {
    setLoading(true);
    setError('');
    try {
      const result = await ScanPorts();
      setPorts(result.sort((a, b) => a.port - b.port));
    } catch (err: any) {
      setError(err.toString());
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPorts();
  }, []);

  const handleKill = async (p: number) => {
    if (!window.confirm(`Are you sure you want to kill the process on port ${p}?`)) return;
    try {
      await KillPort(p);
      await loadPorts();
    } catch (err: any) {
      alert("Error killing port: " + err);
    }
  };

  const handleShare = async (p: number) => {
    setSharing(prev => ({...prev, [p]: true}));
    try {
      await SharePort(p);
      await loadPorts();
    } catch (err: any) {
      alert("Error sharing port: " + err);
    } finally {
      setSharing(prev => ({...prev, [p]: false}));
    }
  };

  const handleStopShare = async (p: number) => {
    try {
      await StopSharePort(p);
      await loadPorts();
    } catch (err: any) {
      alert("Error stopping share: " + err);
    }
  };

  const formatMem = (bytes: number) => {
    if (!bytes) return '-';
    return (bytes / 1024 / 1024).toFixed(1) + ' MB';
  };

  const formatCPU = (pct: number) => {
    if (!pct) return '-';
    return pct.toFixed(1) + '%';
  };

  const filteredPorts = ports.filter(p => {
    const s = search.toLowerCase();
    return (
      p.port.toString().includes(s) ||
      (p.status || '').toLowerCase().includes(s) ||
      (p.processName || '').toLowerCase().includes(s) ||
      (p.project || '').toLowerCase().includes(s)
    );
  });

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-4 md:p-8">
      <div className="w-full h-full mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-emerald-400">
              🩺 PortDoctor
            </h1>
            <p className="text-gray-400 mt-2">Visual Port Management Interface</p>
          </div>
          <button
            onClick={loadPorts}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors flex items-center shadow-lg shadow-blue-500/30"
          >
            {loading ? 'Scanning...' : '🔄 Refresh'}
          </button>
        </div>

        {error && (
          <div className="bg-red-500/20 text-red-400 p-4 rounded-lg mb-6 border border-red-500/50">
            {error}
          </div>
        )}

        <div className="mb-6 relative">
          <input
            type="text"
            placeholder="Live search by port, status, or process name..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 text-white px-4 py-3 pr-12 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all shadow-inner"
          />
          {search && (
            <button
              onClick={() => setSearch('')}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white bg-gray-700 hover:bg-gray-600 rounded-full p-1 transition-colors flex items-center justify-center"
              title="Clear search"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>

        <div className="bg-gray-800 rounded-xl shadow-xl overflow-x-auto border border-gray-700">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-800/50 text-gray-400 uppercase text-xs tracking-wider border-b border-gray-700">
                <th className="px-6 py-4 font-medium">Port</th>
                <th className="px-6 py-4 font-medium">Status</th>
                <th className="px-6 py-4 font-medium">Process</th>
                <th className="px-6 py-4 font-medium">Project</th>
                <th className="px-6 py-4 font-medium">Resources</th>
                <th className="px-6 py-4 font-medium text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {filteredPorts.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                    No ports found
                  </td>
                </tr>
              ) : (
                filteredPorts.map((p) => (
                  <tr key={p.port} className="hover:bg-gray-700/30 transition-colors">
                    <td className="px-6 py-4">
                      <span className="font-mono text-blue-400 font-bold text-lg">{p.port}</span>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-semibold ${
                        p.status === 'OCCUPIED' ? 'bg-amber-500/20 text-amber-400 border border-amber-500/30' :
                        p.status === 'FREE' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' :
                        'bg-gray-500/20 text-gray-400 border border-gray-500/30'
                      }`}>
                        {p.status || 'UNKNOWN'}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-col">
                        <span className="font-medium text-gray-200">{p.processName || '-'}</span>
                        <span className="text-xs text-gray-500">PID: {p.pid > 0 ? p.pid : '-'}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      {p.project ? (
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 font-medium text-sm">
                          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                            <path d="M2 6a2 2 0 012-2h5.586a1 1 0 01.707.293l2.828 2.828a1 1 0 00.707.293H16a2 2 0 012 2v5a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                          </svg>
                          {p.project}
                        </span>
                      ) : (
                        <span className="text-gray-600">-</span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      {p.pid > 0 ? (
                        <div className="flex flex-col gap-1 text-xs">
                          <div className="flex items-center gap-2">
                            <span className="text-gray-500 w-8">CPU</span>
                            <div className="w-16 h-1.5 bg-gray-700 rounded-full overflow-hidden">
                              <div className="h-full bg-blue-500 rounded-full" style={{ width: `${Math.min(p.cpu || 0, 100)}%` }} />
                            </div>
                            <span className="text-gray-400">{formatCPU(p.cpu)}</span>
                          </div>
                          <div className="flex items-center gap-2">
                            <span className="text-gray-500 w-8">RAM</span>
                            <div className="w-16 h-1.5 bg-gray-700 rounded-full overflow-hidden">
                              <div className="h-full bg-emerald-500 rounded-full" style={{ width: `${Math.min((p.ram || 0) / 1024 / 1024 / 16, 100)}%` }} />
                            </div>
                            <span className="text-gray-400">{formatMem(p.ram)}</span>
                          </div>
                        </div>
                      ) : (
                        <span className="text-gray-600">-</span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right">
                      {p.status === 'OCCUPIED' && (
                        <div className="flex items-center justify-end gap-2 flex-wrap max-w-[250px] ml-auto">
                          {p.sharedUrl ? (
                            <div className="flex flex-col items-end gap-1">
                              <a href={p.sharedUrl} target="_blank" rel="noreferrer" className="text-xs text-blue-400 hover:text-blue-300 underline underline-offset-2">
                                {p.sharedUrl.replace('https://', '')}
                              </a>
                              <button
                                onClick={() => handleStopShare(p.port)}
                                className="px-3 py-1.5 bg-orange-500/10 hover:bg-orange-500 text-orange-500 hover:text-white rounded-md transition-all border border-orange-500/50 hover:shadow-[0_0_15px_rgba(249,115,22,0.5)] text-xs font-medium"
                              >
                                Stop Share
                              </button>
                            </div>
                          ) : (
                            <button
                              onClick={() => handleShare(p.port)}
                              disabled={sharing[p.port]}
                              className="px-3 py-1.5 bg-blue-500/10 hover:bg-blue-500 text-blue-400 hover:text-white rounded-md transition-all border border-blue-500/50 hover:shadow-[0_0_15px_rgba(59,130,246,0.5)] text-xs font-medium disabled:opacity-50 flex items-center gap-1.5"
                            >
                              {sharing[p.port] ? (
                                <svg className="animate-spin h-3.5 w-3.5" viewBox="0 0 24 24">
                                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none"></circle>
                                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                              ) : (
                                <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                                  <path d="M11 3a1 1 0 100 2h2.586l-6.293 6.293a1 1 0 101.414 1.414L15 6.414V9a1 1 0 102 0V4a1 1 0 00-1-1h-5z" />
                                  <path d="M5 5a2 2 0 00-2 2v8a2 2 0 002 2h8a2 2 0 002-2v-3a1 1 0 10-2 0v3H5V7h3a1 1 0 000-2H5z" />
                                </svg>
                              )}
                              Share
                            </button>
                          )}
                          <button
                            onClick={() => handleKill(p.port)}
                            className="px-3 py-1.5 bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white rounded-md transition-all border border-red-500/50 hover:shadow-[0_0_15px_rgba(239,68,68,0.5)] text-xs font-medium"
                          >
                            Kill
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default App;
