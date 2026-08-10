import React, { useState } from 'react';

interface HelpGuideProps {
    onClose: () => void;
}

export const HelpGuideModal: React.FC<HelpGuideProps> = ({ onClose }) => {
    const [activeTab, setActiveTab] = useState<string>('share');

    const tabs = [
        { id: 'share', name: '🌍 Share' },
        { id: 'inspect', name: '🔎 Inspect' },
        { id: 'details', name: '📄 Details' },
        { id: 'rules', name: '🛡️ Rules' }
    ];

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 rounded-xl shadow-2xl border border-gray-700 w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">
                <div className="flex justify-between items-center p-6 border-b border-gray-700 bg-gray-900">
                    <div>
                        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                            ℹ️ PortDoctor Guide
                        </h2>
                        <p className="text-gray-400 text-sm mt-1">Learn how to use advanced features</p>
                    </div>
                    <button onClick={onClose} className="text-gray-400 hover:text-red-400 transition-colors">
                        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                
                <div className="flex flex-1 overflow-hidden">
                    {/* Sidebar Tabs */}
                    <div className="w-48 bg-gray-900 border-r border-gray-700 flex flex-col">
                        {tabs.map(tab => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`text-left px-6 py-4 font-medium transition-colors border-l-4 ${
                                    activeTab === tab.id 
                                    ? 'bg-gray-800 text-blue-400 border-blue-500' 
                                    : 'text-gray-400 border-transparent hover:bg-gray-800/50 hover:text-gray-200'
                                }`}
                            >
                                {tab.name}
                            </button>
                        ))}
                    </div>

                    {/* Content Area */}
                    <div className="flex-1 p-8 overflow-y-auto bg-gray-800 custom-scrollbar text-gray-300 space-y-6">
                        
                        {activeTab === 'share' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-blue-400 mb-4">🌍 Đưa Localhost lên Internet (Share)</h3>
                                <p className="mb-4 text-gray-300">Tính năng này giúp bạn chia sẻ trang web/ứng dụng đang chạy trên máy tính của bạn (localhost) cho bất kỳ ai trên thế giới truy cập qua một đường link Public.</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">Cách sử dụng:</h4>
                                <ol className="list-decimal pl-5 space-y-2 text-gray-400">
                                    <li>Đảm bảo ứng dụng của bạn (Web, API...) đang chạy và chiếm một Port (VD: Port 3000).</li>
                                    <li>Bấm nút <span className="bg-blue-500/20 text-blue-400 px-2 py-0.5 rounded text-xs border border-blue-500/50">Share</span> ở cột Action.</li>
                                    <li>Chờ vài giây, một đường link màu xanh có dạng <code className="text-emerald-400 bg-gray-900 px-1 py-0.5 rounded">https://xyz.lhr.life</code> sẽ xuất hiện.</li>
                                    <li>Copy đường link đó và gửi cho bạn bè hoặc khách hàng. Họ có thể xem trực tiếp sản phẩm của bạn!</li>
                                    <li>Bấm <span className="bg-orange-500/20 text-orange-400 px-2 py-0.5 rounded text-xs border border-orange-500/50">Stop Share</span> khi không muốn chia sẻ nữa.</li>
                                </ol>
                                
                                <div className="mt-6 p-4 bg-blue-900/20 border border-blue-800 rounded-lg text-sm">
                                    <strong className="text-blue-400">Lưu ý:</strong> PortDoctor sử dụng kết nối mã hoá SSH qua dịch vụ `localhost.run`. Việc chia sẻ này là an toàn và tự động ngắt kết nối khi bạn tắt app.
                                </div>
                            </div>
                        )}

                        {activeTab === 'inspect' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-emerald-400 mb-4">🔎 Soi Gói Tin API (HTTP Traffic Inspector)</h3>
                                <p className="mb-4 text-gray-300">Tính năng này biến PortDoctor thành một công cụ như Wireshark hay Chrome DevTools Network Tab. Giúp bạn theo dõi toàn bộ request HTTP đi qua một Port.</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">Cách sử dụng:</h4>
                                <ol className="list-decimal pl-5 space-y-2 text-gray-400">
                                    <li>Bấm nút <span className="bg-emerald-500/20 text-emerald-400 px-2 py-0.5 rounded text-xs border border-emerald-500/50">Inspect</span> ở dòng Port bạn muốn soi (VD: Port 8080).</li>
                                    <li>Một cửa sổ Proxy sẽ hiện ra. PortDoctor sẽ tạo một <strong className="text-white">Proxy Port</strong> ảo (VD: 18080).</li>
                                    <li>Thay vì truy cập vào <code className="text-gray-400">http://localhost:8080</code>, bạn hãy đổi trình duyệt hoặc Postman sang truy cập vào <code className="text-emerald-400 bg-gray-900 px-1 py-0.5 rounded">http://localhost:18080</code>.</li>
                                    <li>Dữ liệu vẫn được trả về bình thường, nhưng toàn bộ lịch sử gửi/nhận (Headers, Body) sẽ được log trực tiếp trên màn hình PortDoctor!</li>
                                </ol>

                                <div className="mt-6 p-4 bg-emerald-900/20 border border-emerald-800 rounded-lg text-sm">
                                    <strong className="text-emerald-400">Mẹo Debug:</strong> Rất hiệu quả khi bạn gọi API từ App Mobile xuống Backend Localhost và cần xem Body JSON hoặc Token có truyền đúng hay không.
                                </div>
                            </div>
                        )}

                        {activeTab === 'details' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-indigo-400 mb-4">📄 Soi Ruột Tiến Trình (Process Details)</h3>
                                <p className="mb-4 text-gray-300">Hiển thị các thông tin "tuyệt mật" mà Task Manager thông thường không cho bạn xem chi tiết.</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">Các thông tin cung cấp:</h4>
                                <ul className="list-disc pl-5 space-y-3 text-gray-400">
                                    <li>
                                        <strong className="text-gray-200">Command Line:</strong> Câu lệnh chính xác đã được dùng để gọi ứng dụng lên. Rất hữu ích để xem ứng dụng có được chạy kèm cờ (flag) nào không (VD: <code className="text-gray-500 bg-gray-900 px-1 py-0.5 rounded">node server.js --env=prod</code>).
                                    </li>
                                    <li>
                                        <strong className="text-gray-200">Working Directory:</strong> Thư mục gốc nơi ứng dụng đang được thực thi.
                                    </li>
                                    <li>
                                        <strong className="text-gray-200">Environment Variables (Biến môi trường):</strong> Toàn bộ cấu hình ngầm của máy tính và ứng dụng. Rất quan trọng khi code Node.js, Python, hay Docker để kiểm tra xem file <code className="text-gray-500 bg-gray-900 px-1 py-0.5 rounded">.env</code> có được nạp đúng chưa.
                                    </li>
                                </ul>
                            </div>
                        )}

                        {activeTab === 'rules' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-purple-400 mb-4">🛡️ Tự động Bảo vệ & Hồi sinh (Rule Engine)</h3>
                                <p className="mb-4 text-gray-300">Cho phép PortDoctor trở thành một người gác cổng (Watchdog) giám sát 24/7 các Port của bạn.</p>
                                
                                <div className="space-y-6 mt-6">
                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">1. Protect Port (Bảo vệ Port)</h4>
                                        <p className="text-sm text-gray-400">Ngăn chặn việc lỡ tay ấn nhầm nút Kill trên PortDoctor. Nút Kill sẽ bị khoá (Disabled).</p>
                                    </div>

                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">2. Allowed Process (Tự động diệt kẻ lạ)</h4>
                                        <p className="text-sm text-gray-400">Nhập vào tên tiến trình (VD: <code className="text-gray-500">nginx.exe</code>). Nếu một ứng dụng bất kỳ khác (không phải nginx) dám bật lên và chiếm dụng Port này, PortDoctor sẽ âm thầm tự động "Bóp cổ" (Kill) ứng dụng đó ngay lập tức.</p>
                                    </div>

                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">3. Auto-Heal (Tự động Hồi sinh)</h4>
                                        <p className="text-sm text-gray-400">Nhập vào câu lệnh khởi chạy (VD: <code className="text-gray-500">npm run dev</code>) và thư mục dự án. Nếu PortDoctor phát hiện Port bị sập (ứng dụng crash), nó sẽ tự động chạy câu lệnh này trong chế độ chạy ngầm (Background) để gọi ứng dụng sống dậy.</p>
                                    </div>
                                </div>
                            </div>
                        )}

                    </div>
                </div>
            </div>
        </div>
    );
};
