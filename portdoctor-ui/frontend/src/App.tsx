import { useState, useEffect } from 'react';
import { ScanPorts, KillPort } from '../wailsjs/go/main/App';
import { main } from '../wailsjs/go/models';

function App() {
  const [ports, setPorts] = useState<main.UIPortInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [error, setError] = useState('');

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

  const filteredPorts = ports.filter(p => {
    const s = search.toLowerCase();
    return (
      p.port.toString().includes(s) ||
      (p.status || '').toLowerCase().includes(s) ||
      (p.processName || '').toLowerCase().includes(s)
    );
  });

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 p-8">
      <div className="max-w-5xl mx-auto">
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
            className="w-full bg-gray-800 border border-gray-700 text-white px-4 py-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all shadow-inner"
          />
        </div>

        <div className="bg-gray-800 rounded-xl shadow-xl overflow-hidden border border-gray-700">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-800/50 text-gray-400 uppercase text-xs tracking-wider border-b border-gray-700">
                <th className="px-6 py-4 font-medium">Port</th>
                <th className="px-6 py-4 font-medium">Status</th>
                <th className="px-6 py-4 font-medium">Process</th>
                <th className="px-6 py-4 font-medium">PID</th>
                <th className="px-6 py-4 font-medium text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {filteredPorts.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
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
                      <span className="font-medium text-gray-200">{p.processName || '-'}</span>
                    </td>
                    <td className="px-6 py-4">
                      <span className="font-mono text-gray-400">{p.pid > 0 ? p.pid : '-'}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      {p.status === 'OCCUPIED' && (
                        <button
                          onClick={() => handleKill(p.port)}
                          className="px-3 py-1.5 bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white rounded-md transition-all border border-red-500/50 hover:shadow-[0_0_15px_rgba(239,68,68,0.5)]"
                        >
                          Kill Process
                        </button>
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
