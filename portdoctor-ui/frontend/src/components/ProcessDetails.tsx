import React from 'react';
import { main } from '../../wailsjs/go/models';

interface ProcessDetailsProps {
    details: main.ProcessDetails | null;
    onClose: () => void;
}

export const ProcessDetailsModal: React.FC<ProcessDetailsProps> = ({ details, onClose }) => {
    if (!details) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 rounded-xl shadow-2xl border border-gray-700 w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">
                <div className="flex justify-between items-center p-6 border-b border-gray-700 bg-gray-800/80">
                    <div>
                        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                            📄 Process Details
                        </h2>
                        <p className="text-gray-400 mt-1">
                            {details.name} (PID: {details.pid})
                        </p>
                    </div>
                    <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
                        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                
                <div className="p-6 overflow-y-auto flex-1 space-y-6 custom-scrollbar">
                    {/* Basic Info */}
                    <div className="bg-gray-900 rounded-lg p-4 border border-gray-700">
                        <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">Runtime Info</h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <span className="text-gray-500 text-xs">Working Directory</span>
                                <p className="text-gray-200 font-mono text-sm break-all">{details.cwd || 'N/A'}</p>
                            </div>
                            <div>
                                <span className="text-gray-500 text-xs">Username</span>
                                <p className="text-gray-200 font-mono text-sm">{details.username || 'N/A'}</p>
                            </div>
                        </div>
                    </div>

                    {/* Command Line */}
                    <div className="bg-gray-900 rounded-lg p-4 border border-gray-700">
                        <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">Command Line</h3>
                        <div className="bg-black/50 p-3 rounded border border-gray-800 font-mono text-sm text-green-400 break-words">
                            {details.cmdline ? details.cmdline.join(' ') : 'N/A'}
                        </div>
                    </div>

                    {/* Environment Variables */}
                    <div className="bg-gray-900 rounded-lg p-4 border border-gray-700">
                        <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3 flex justify-between">
                            <span>Environment Variables</span>
                            <span className="text-xs bg-gray-800 px-2 py-1 rounded text-gray-300">
                                {Object.keys(details.envVars || {}).length} vars
                            </span>
                        </h3>
                        <div className="overflow-x-auto">
                            <table className="w-full text-left text-sm border-collapse">
                                <thead>
                                    <tr className="border-b border-gray-700 text-gray-500">
                                        <th className="py-2 font-medium w-1/3">Key</th>
                                        <th className="py-2 font-medium">Value</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-800">
                                    {Object.entries(details.envVars || {}).sort().map(([key, value]) => (
                                        <tr key={key} className="hover:bg-gray-800/50">
                                            <td className="py-2 pr-4 font-mono text-blue-400 break-all">{key}</td>
                                            <td className="py-2 font-mono text-gray-300 break-all">{value}</td>
                                        </tr>
                                    ))}
                                    {Object.keys(details.envVars || {}).length === 0 && (
                                        <tr>
                                            <td colSpan={2} className="py-4 text-center text-gray-500">No environment variables found or access denied.</td>
                                        </tr>
                                    )}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};
