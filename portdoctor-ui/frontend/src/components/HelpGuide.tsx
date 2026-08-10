import React, { useState } from 'react';

interface HelpGuideProps {
    onClose: () => void;
}

export const HelpGuideModal: React.FC<HelpGuideProps> = ({ onClose }) => {
    const [activeTab, setActiveTab] = useState<string>('share');
    const [lang, setLang] = useState<'vi' | 'en'>('vi');

    const tabs = [
        { id: 'share', name: '🌍 Share' },
        { id: 'inspect', name: '🔎 Inspect' },
        { id: 'details', name: '📄 Details' },
        { id: 'rules', name: '🛡️ Rules' }
    ];

    const t = {
        vi: {
            title: 'ℹ️ Hướng dẫn PortDoctor',
            subtitle: 'Tìm hiểu cách sử dụng các tính năng nâng cao',
            share: {
                title: '🌍 Đưa Localhost lên Internet (Share)',
                desc: 'Tính năng này giúp bạn chia sẻ trang web/ứng dụng đang chạy trên máy tính của bạn (localhost) cho bất kỳ ai trên thế giới truy cập qua một đường link Public.',
                usage: 'Cách sử dụng:',
                step1: 'Đảm bảo ứng dụng của bạn (Web, API...) đang chạy và chiếm một Port (VD: Port 3000).',
                step2: 'Bấm nút',
                step2Btn: 'Share',
                step2Text: 'ở cột Action.',
                step3: 'Chờ vài giây, một đường link màu xanh có dạng',
                step3Url: 'https://xyz.lhr.life',
                step3Text: 'sẽ xuất hiện.',
                step4: 'Copy đường link đó và gửi cho bạn bè hoặc khách hàng. Họ có thể xem trực tiếp sản phẩm của bạn!',
                step5: 'Bấm',
                step5Btn: 'Stop Share',
                step5Text: 'khi không muốn chia sẻ nữa.',
                noteTitle: 'Lưu ý:',
                noteDesc: 'PortDoctor sử dụng kết nối mã hoá SSH qua dịch vụ `localhost.run`. Việc chia sẻ này là an toàn và tự động ngắt kết nối khi bạn tắt app.'
            },
            inspect: {
                title: '🔎 Soi Gói Tin API (HTTP Traffic Inspector)',
                desc: 'Tính năng này biến PortDoctor thành một công cụ như Wireshark hay Chrome DevTools Network Tab. Giúp bạn theo dõi toàn bộ request HTTP đi qua một Port.',
                usage: 'Cách sử dụng:',
                step1: 'Bấm nút',
                step1Btn: 'Inspect',
                step1Text: 'ở dòng Port bạn muốn soi (VD: Port 8080).',
                step2: 'Một cửa sổ Proxy sẽ hiện ra. PortDoctor sẽ tạo một Proxy Port ảo (VD: 18080).',
                step3: 'Thay vì truy cập vào',
                step3Url1: 'http://localhost:8080',
                step3Text: ', bạn hãy đổi trình duyệt hoặc Postman sang truy cập vào',
                step3Url2: 'http://localhost:18080',
                step4: 'Dữ liệu vẫn được trả về bình thường, nhưng toàn bộ lịch sử gửi/nhận (Headers, Body) sẽ được log trực tiếp trên màn hình PortDoctor!',
                tipTitle: 'Mẹo Debug:',
                tipDesc: 'Rất hiệu quả khi bạn gọi API từ App Mobile xuống Backend Localhost và cần xem Body JSON hoặc Token có truyền đúng hay không.'
            },
            details: {
                title: '📄 Soi Ruột Tiến Trình (Process Details)',
                desc: 'Hiển thị các thông tin "tuyệt mật" mà Task Manager thông thường không cho bạn xem chi tiết.',
                infoList: 'Các thông tin cung cấp:',
                cmdline: 'Command Line:',
                cmdlineDesc: 'Câu lệnh chính xác đã được dùng để gọi ứng dụng lên. Rất hữu ích để xem ứng dụng có được chạy kèm cờ (flag) nào không (VD: node server.js --env=prod).',
                cwd: 'Working Directory:',
                cwdDesc: 'Thư mục gốc nơi ứng dụng đang được thực thi.',
                env: 'Environment Variables (Biến môi trường):',
                envDesc: 'Toàn bộ cấu hình ngầm của máy tính và ứng dụng. Rất quan trọng khi code Node.js, Python, hay Docker để kiểm tra xem file .env có được nạp đúng chưa.'
            },
            rules: {
                title: '🛡️ Tự động Bảo vệ & Hồi sinh (Rule Engine)',
                desc: 'Cho phép PortDoctor trở thành một người gác cổng (Watchdog) giám sát 24/7 các Port của bạn.',
                r1Title: '1. Protect Port (Bảo vệ Port)',
                r1Desc: 'Ngăn chặn việc lỡ tay ấn nhầm nút Kill trên PortDoctor. Nút Kill sẽ bị khoá (Disabled).',
                r2Title: '2. Allowed Process (Tự động diệt kẻ lạ)',
                r2Desc: 'Nhập vào tên tiến trình (VD: nginx.exe). Nếu một ứng dụng bất kỳ khác (không phải nginx) dám bật lên và chiếm dụng Port này, PortDoctor sẽ âm thầm tự động "Bóp cổ" (Kill) ứng dụng đó ngay lập tức.',
                r3Title: '3. Auto-Heal (Tự động Hồi sinh)',
                r3Desc: 'Nhập vào câu lệnh khởi chạy (VD: npm run dev) và thư mục dự án. Nếu PortDoctor phát hiện Port bị sập (ứng dụng crash), nó sẽ tự động chạy câu lệnh này trong chế độ chạy ngầm (Background) để gọi ứng dụng sống dậy.'
            }
        },
        en: {
            title: 'ℹ️ PortDoctor Guide',
            subtitle: 'Learn how to use advanced features',
            share: {
                title: '🌍 Expose Localhost to the Internet (Share)',
                desc: 'This feature allows you to share your website/app running on localhost with anyone in the world via a public link.',
                usage: 'How to use:',
                step1: 'Ensure your app (Web, API...) is running on a Port (e.g. Port 3000).',
                step2: 'Click the',
                step2Btn: 'Share',
                step2Text: 'button in the Action column.',
                step3: 'Wait a few seconds, a blue link like',
                step3Url: 'https://xyz.lhr.life',
                step3Text: 'will appear.',
                step4: 'Copy that link and send it to friends or clients. They can instantly view your product!',
                step5: 'Click',
                step5Btn: 'Stop Share',
                step5Text: 'when you want to stop sharing.',
                noteTitle: 'Note:',
                noteDesc: 'PortDoctor uses a secure SSH tunnel via the `localhost.run` service. Sharing is secure and automatically disconnects when you close the app.'
            },
            inspect: {
                title: '🔎 Intercept API Traffic (HTTP Traffic Inspector)',
                desc: 'This turns PortDoctor into a tool like Wireshark or Chrome DevTools Network Tab. It helps you monitor all HTTP requests going through a port.',
                usage: 'How to use:',
                step1: 'Click the',
                step1Btn: 'Inspect',
                step1Text: 'button on the Port you want to inspect (e.g. Port 8080).',
                step2: 'A Proxy window will appear. PortDoctor creates a virtual Proxy Port (e.g. 18080).',
                step3: 'Instead of navigating to',
                step3Url1: 'http://localhost:8080',
                step3Text: ', point your browser or Postman to',
                step3Url2: 'http://localhost:18080',
                step4: 'Data is returned normally, but the entire send/receive history (Headers, Body) will be logged live on the PortDoctor screen!',
                tipTitle: 'Debugging Tip:',
                tipDesc: 'Very effective when calling APIs from a Mobile App to a Localhost Backend and you need to see if the JSON Body or Token is being passed correctly.'
            },
            details: {
                title: '📄 Inspect Process Internals (Process Details)',
                desc: 'Displays "top secret" information that standard Task Manager hides from you.',
                infoList: 'Provided information:',
                cmdline: 'Command Line:',
                cmdlineDesc: 'The exact command used to launch the application. Very useful to see if the app was run with any flags (e.g. node server.js --env=prod).',
                cwd: 'Working Directory:',
                cwdDesc: 'The root directory where the application is being executed.',
                env: 'Environment Variables:',
                envDesc: 'The entire hidden configuration of the computer and application. Crucial when coding Node.js, Python, or Docker to check if the .env file is loaded correctly.'
            },
            rules: {
                title: '🛡️ Auto Protect & Revive (Rule Engine)',
                desc: 'Allows PortDoctor to act as a 24/7 Watchdog monitoring your Ports.',
                r1Title: '1. Protect Port',
                r1Desc: 'Prevents accidental clicks on the Kill button in PortDoctor. The Kill button will be Disabled.',
                r2Title: '2. Allowed Process (Auto-kill intruders)',
                r2Desc: 'Enter a process name (e.g. nginx.exe). If any other application (not nginx) dares to start and occupy this Port, PortDoctor will silently and automatically Kill it instantly.',
                r3Title: '3. Auto-Heal (Auto-Revive)',
                r3Desc: 'Enter a startup command (e.g. npm run dev) and the project directory. If PortDoctor detects the Port has crashed, it will automatically run this command in the Background to resurrect the application.'
            }
        }
    };

    const text = t[lang];

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-gray-800 rounded-xl shadow-2xl border border-gray-700 w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">
                <div className="flex justify-between items-center p-6 border-b border-gray-700 bg-gray-900">
                    <div>
                        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
                            {text.title}
                        </h2>
                        <p className="text-gray-400 text-sm mt-1">{text.subtitle}</p>
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="flex bg-gray-800 p-1 rounded-lg border border-gray-700">
                            <button 
                                onClick={() => setLang('vi')} 
                                className={`px-3 py-1 rounded-md text-sm font-medium transition-colors ${lang === 'vi' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'}`}
                            >
                                VI
                            </button>
                            <button 
                                onClick={() => setLang('en')} 
                                className={`px-3 py-1 rounded-md text-sm font-medium transition-colors ${lang === 'en' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'}`}
                            >
                                EN
                            </button>
                        </div>
                        <button onClick={onClose} className="text-gray-400 hover:text-red-400 transition-colors">
                            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
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
                                <h3 className="text-xl font-bold text-blue-400 mb-4">{text.share.title}</h3>
                                <p className="mb-4 text-gray-300">{text.share.desc}</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">{text.share.usage}</h4>
                                <ol className="list-decimal pl-5 space-y-2 text-gray-400">
                                    <li>{text.share.step1}</li>
                                    <li>{text.share.step2} <span className="bg-blue-500/20 text-blue-400 px-2 py-0.5 rounded text-xs border border-blue-500/50">{text.share.step2Btn}</span> {text.share.step2Text}</li>
                                    <li>{text.share.step3} <code className="text-emerald-400 bg-gray-900 px-1 py-0.5 rounded">{text.share.step3Url}</code> {text.share.step3Text}</li>
                                    <li>{text.share.step4}</li>
                                    <li>{text.share.step5} <span className="bg-orange-500/20 text-orange-400 px-2 py-0.5 rounded text-xs border border-orange-500/50">{text.share.step5Btn}</span> {text.share.step5Text}</li>
                                </ol>
                                
                                <div className="mt-6 p-4 bg-blue-900/20 border border-blue-800 rounded-lg text-sm">
                                    <strong className="text-blue-400">{text.share.noteTitle}</strong> {text.share.noteDesc}
                                </div>
                            </div>
                        )}

                        {activeTab === 'inspect' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-emerald-400 mb-4">{text.inspect.title}</h3>
                                <p className="mb-4 text-gray-300">{text.inspect.desc}</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">{text.inspect.usage}</h4>
                                <ol className="list-decimal pl-5 space-y-2 text-gray-400">
                                    <li>{text.inspect.step1} <span className="bg-emerald-500/20 text-emerald-400 px-2 py-0.5 rounded text-xs border border-emerald-500/50">{text.inspect.step1Btn}</span> {text.inspect.step1Text}</li>
                                    <li>{text.inspect.step2}</li>
                                    <li>{text.inspect.step3} <code className="text-gray-400">{text.inspect.step3Url1}</code>{text.inspect.step3Text} <code className="text-emerald-400 bg-gray-900 px-1 py-0.5 rounded">{text.inspect.step3Url2}</code>.</li>
                                    <li>{text.inspect.step4}</li>
                                </ol>

                                <div className="mt-6 p-4 bg-emerald-900/20 border border-emerald-800 rounded-lg text-sm">
                                    <strong className="text-emerald-400">{text.inspect.tipTitle}</strong> {text.inspect.tipDesc}
                                </div>
                            </div>
                        )}

                        {activeTab === 'details' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-indigo-400 mb-4">{text.details.title}</h3>
                                <p className="mb-4 text-gray-300">{text.details.desc}</p>
                                
                                <h4 className="font-semibold text-gray-200 mt-6 mb-2">{text.details.infoList}</h4>
                                <ul className="list-disc pl-5 space-y-3 text-gray-400">
                                    <li>
                                        <strong className="text-gray-200">{text.details.cmdline}</strong> {text.details.cmdlineDesc}
                                    </li>
                                    <li>
                                        <strong className="text-gray-200">{text.details.cwd}</strong> {text.details.cwdDesc}
                                    </li>
                                    <li>
                                        <strong className="text-gray-200">{text.details.env}</strong> {text.details.envDesc}
                                    </li>
                                </ul>
                            </div>
                        )}

                        {activeTab === 'rules' && (
                            <div className="animate-fade-in">
                                <h3 className="text-xl font-bold text-purple-400 mb-4">{text.rules.title}</h3>
                                <p className="mb-4 text-gray-300">{text.rules.desc}</p>
                                
                                <div className="space-y-6 mt-6">
                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">{text.rules.r1Title}</h4>
                                        <p className="text-sm text-gray-400">{text.rules.r1Desc}</p>
                                    </div>

                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">{text.rules.r2Title}</h4>
                                        <p className="text-sm text-gray-400">{text.rules.r2Desc}</p>
                                    </div>

                                    <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-700">
                                        <h4 className="font-bold text-gray-200 mb-2">{text.rules.r3Title}</h4>
                                        <p className="text-sm text-gray-400">{text.rules.r3Desc}</p>
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
