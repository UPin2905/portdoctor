import React, { useState, useEffect, useRef } from 'react';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

interface TrafficInspectorProps {
    port: number;
    proxyPort: number;
    onClose: () => void;
}

interface HttpLog {
    id: string;
    method?: string;
    url?: string;
    status?: number;
    reqHeaders?: Record<string, string>;
    resHeaders?: Record<string, string>;
    reqBody?: string;
    resBody?: string;
    timestamp: Date;
}

export const TrafficInspectorModal: React.FC<TrafficInspectorProps> = ({ port, proxyPort, onClose }) => {
    const [logs, setLogs] = useState<HttpLog[]>([]);
    const [selectedLogId, setSelectedLogId] = useState<string | null>(null);
    const logsEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const reqEventName = `http-req-${port}`;
        const resEventName = `http-res-${port}`;

        EventsOn(reqEventName, (data: any) => {
            setLogs(prev => {
                const newLog: HttpLog = {
                    id: data.id,
                    method: data.method,
                    url: data.url,
                    reqHeaders: data.headers,
                    reqBody: data.body,
                    timestamp: new Date()
                };
                return [...prev, newLog];
            });
        });

        EventsOn(resEventName, (data: any) => {
            setLogs(prev => prev.map(log => {
                if (log.id === data.id) {
                    return {
                        ...log,
                        status: data.status,
                        resHeaders: data.headers,
                        resBody: data.body,
                    };
                }
                return log;
            }));
        });

        return () => {
            EventsOff(reqEventName);
            EventsOff(resEventName);
        };
    }, [port]);

    useEffect(() => {
        // Auto-scroll to bottom on new logs
        logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [logs.length]);

    const selectedLog = logs.find(l => l.id === selectedLogId);

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 rounded-xl shadow-2xl border border-gray-700 w-full max-w-6xl h-[90vh] flex flex-col overflow-hidden">
                <div className="flex justify-between items-center p-4 border-b border-gray-700 bg-gray-900">
                    <div>
                        <h2 className="text-xl font-bold text-emerald-400 flex items-center gap-2">
                            🔎 HTTP Traffic Inspector
                        </h2>
                        <p className="text-gray-400 text-sm mt-1">
                            Listening on proxy port <span className="text-white font-mono">{proxyPort}</span> ➔ Forwarding to <span className="text-white font-mono">{port}</span>
                        </p>
                    </div>
                    <div className="flex gap-4">
                        <button onClick={() => setLogs([])} className="px-3 py-1 bg-gray-700 hover:bg-gray-600 rounded text-sm text-white">Clear</button>
                        <button onClick={onClose} className="text-gray-400 hover:text-red-400 transition-colors">
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>
                
                <div className="flex flex-1 overflow-hidden">
                    {/* Log List */}
                    <div className="w-1/3 border-r border-gray-700 flex flex-col bg-gray-900/50 overflow-y-auto">
                        {logs.length === 0 ? (
                            <div className="flex-1 flex items-center justify-center text-gray-500 p-4 text-center text-sm">
                                No requests yet.<br/>Send HTTP traffic to http://127.0.0.1:{proxyPort}
                            </div>
                        ) : (
                            <div className="divide-y divide-gray-800">
                                {logs.map(log => (
                                    <div 
                                        key={log.id} 
                                        onClick={() => setSelectedLogId(log.id)}
                                        className={`p-3 cursor-pointer hover:bg-gray-700 transition-colors ${selectedLogId === log.id ? 'bg-gray-700 border-l-4 border-emerald-500' : 'border-l-4 border-transparent'}`}
                                    >
                                        <div className="flex justify-between items-center mb-1">
                                            <span className={`font-bold text-sm ${
                                                log.method === 'GET' ? 'text-blue-400' : 
                                                log.method === 'POST' ? 'text-green-400' : 
                                                log.method === 'DELETE' ? 'text-red-400' : 'text-yellow-400'
                                            }`}>{log.method}</span>
                                            <span className={`text-xs px-2 py-0.5 rounded ${
                                                !log.status ? 'bg-gray-600' :
                                                log.status < 300 ? 'bg-green-900 text-green-300' :
                                                log.status < 400 ? 'bg-yellow-900 text-yellow-300' :
                                                'bg-red-900 text-red-300'
                                            }`}>
                                                {log.status || 'Pending'}
                                            </span>
                                        </div>
                                        <div className="text-gray-300 text-sm truncate" title={log.url}>{log.url}</div>
                                        <div className="text-gray-500 text-xs mt-1 text-right">{log.timestamp.toLocaleTimeString()}</div>
                                    </div>
                                ))}
                                <div ref={logsEndRef} />
                            </div>
                        )}
                    </div>

                    {/* Log Details */}
                    <div className="w-2/3 flex flex-col bg-gray-800 overflow-y-auto p-4 custom-scrollbar">
                        {selectedLog ? (
                            <div className="space-y-6">
                                {/* Request */}
                                <div>
                                    <h3 className="text-sm font-bold text-gray-400 uppercase tracking-wider mb-2 border-b border-gray-700 pb-1">Request</h3>
                                    <div className="bg-gray-900 rounded p-3 text-sm font-mono text-gray-300">
                                        <span className="text-blue-400 font-bold">{selectedLog.method}</span> {selectedLog.url}
                                    </div>
                                    <h4 className="text-xs text-gray-500 mt-3 mb-1">Headers</h4>
                                    <div className="bg-gray-900 rounded p-3 text-xs font-mono text-gray-400">
                                        {Object.entries(selectedLog.reqHeaders || {}).map(([k, v]) => (
                                            <div key={k}><span className="text-purple-400">{k}:</span> {v}</div>
                                        ))}
                                    </div>
                                    {selectedLog.reqBody && (
                                        <>
                                            <h4 className="text-xs text-gray-500 mt-3 mb-1">Body</h4>
                                            <pre className="bg-gray-900 rounded p-3 text-xs font-mono text-gray-300 overflow-x-auto whitespace-pre-wrap">
                                                {selectedLog.reqBody}
                                            </pre>
                                        </>
                                    )}
                                </div>

                                {/* Response */}
                                <div>
                                    <h3 className="text-sm font-bold text-gray-400 uppercase tracking-wider mb-2 border-b border-gray-700 pb-1">Response</h3>
                                    {!selectedLog.status ? (
                                        <div className="text-gray-500 text-sm italic">Waiting for response...</div>
                                    ) : (
                                        <>
                                            <div className="bg-gray-900 rounded p-3 text-sm font-mono text-gray-300 flex gap-2">
                                                Status: <span className={selectedLog.status < 400 ? 'text-green-400' : 'text-red-400'}>{selectedLog.status}</span>
                                            </div>
                                            <h4 className="text-xs text-gray-500 mt-3 mb-1">Headers</h4>
                                            <div className="bg-gray-900 rounded p-3 text-xs font-mono text-gray-400">
                                                {Object.entries(selectedLog.resHeaders || {}).map(([k, v]) => (
                                                    <div key={k}><span className="text-purple-400">{k}:</span> {v}</div>
                                                ))}
                                            </div>
                                            {selectedLog.resBody && (
                                                <>
                                                    <h4 className="text-xs text-gray-500 mt-3 mb-1">Body</h4>
                                                    <pre className="bg-gray-900 rounded p-3 text-xs font-mono text-gray-300 overflow-x-auto whitespace-pre-wrap">
                                                        {selectedLog.resBody}
                                                    </pre>
                                                </>
                                            )}
                                        </>
                                    )}
                                </div>
                            </div>
                        ) : (
                            <div className="h-full flex items-center justify-center text-gray-600">
                                Select a request to view details
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};
