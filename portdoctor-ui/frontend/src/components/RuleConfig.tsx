import React, { useState, useEffect } from 'react';
import { main } from '../../wailsjs/go/models';

interface RuleConfigProps {
    port: number;
    existingRule?: main.PortRule;
    onSave: (rule: main.PortRule) => void;
    onDelete: (port: number) => void;
    onClose: () => void;
}

export const RuleConfigModal: React.FC<RuleConfigProps> = ({ port, existingRule, onSave, onDelete, onClose }) => {
    const [protectedPort, setProtectedPort] = useState(existingRule?.protected || false);
    const [allowedProcess, setAllowedProcess] = useState(existingRule?.allowedProcess || '');
    const [autoHealCmd, setAutoHealCmd] = useState(existingRule?.autoHealCmd || '');
    const [autoHealDir, setAutoHealDir] = useState(existingRule?.autoHealDir || '');

    const handleSave = () => {
        onSave({
            port: port,
            protected: protectedPort,
            allowedProcess: allowedProcess,
            autoHealCmd: autoHealCmd,
            autoHealDir: autoHealDir
        });
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 rounded-xl shadow-2xl border border-gray-700 w-full max-w-lg overflow-hidden">
                <div className="flex justify-between items-center p-5 border-b border-gray-700 bg-gray-800/80">
                    <h2 className="text-xl font-bold text-white flex items-center gap-2">
                        ⚙️ Configure Rules for Port {port}
                    </h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
                        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                
                <div className="p-6 space-y-6">
                    {/* Protection */}
                    <div className="flex items-start gap-3">
                        <input 
                            type="checkbox" 
                            id="protected"
                            checked={protectedPort}
                            onChange={(e) => setProtectedPort(e.target.checked)}
                            className="mt-1 w-4 h-4 rounded bg-gray-900 border-gray-600 text-blue-500 focus:ring-blue-500" 
                        />
                        <div>
                            <label htmlFor="protected" className="block text-sm font-medium text-gray-200">
                                Protect Port
                            </label>
                            <p className="text-xs text-gray-400 mt-1">Prevent PortDoctor from killing this port manually via UI.</p>
                        </div>
                    </div>

                    {/* Auto Kill */}
                    <div>
                        <label className="block text-sm font-medium text-gray-200 mb-1">
                            Allowed Process (Auto-Kill Strangers)
                        </label>
                        <input 
                            type="text" 
                            value={allowedProcess}
                            onChange={(e) => setAllowedProcess(e.target.value)}
                            placeholder="e.g. node.exe"
                            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-blue-500"
                        />
                        <p className="text-xs text-gray-400 mt-1">If any other process occupies this port, PortDoctor will automatically kill it.</p>
                    </div>

                    {/* Auto Heal */}
                    <div className="pt-4 border-t border-gray-700">
                        <label className="block text-sm font-medium text-gray-200 mb-1">
                            Auto-Heal Command
                        </label>
                        <input 
                            type="text" 
                            value={autoHealCmd}
                            onChange={(e) => setAutoHealCmd(e.target.value)}
                            placeholder="e.g. npm run start"
                            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-blue-500 mb-3"
                        />
                        
                        <label className="block text-sm font-medium text-gray-200 mb-1">
                            Working Directory (Optional)
                        </label>
                        <input 
                            type="text" 
                            value={autoHealDir}
                            onChange={(e) => setAutoHealDir(e.target.value)}
                            placeholder="e.g. C:\Projects\MyWeb"
                            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-blue-500"
                        />
                        <p className="text-xs text-gray-400 mt-2">If the port dies, PortDoctor will run this command to revive it.</p>
                    </div>
                </div>

                <div className="p-4 border-t border-gray-700 bg-gray-800/80 flex justify-between">
                    <button 
                        onClick={() => onDelete(port)}
                        className="px-4 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-500 rounded-lg text-sm font-medium transition-colors"
                    >
                        Delete Rule
                    </button>
                    <div className="flex gap-2">
                        <button 
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm font-medium transition-colors"
                        >
                            Cancel
                        </button>
                        <button 
                            onClick={handleSave}
                            className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-sm font-medium transition-colors"
                        >
                            Save Rule
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};
