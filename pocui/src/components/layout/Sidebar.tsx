import React from 'react';
import { FiGrid, FiMessageSquare, FiInbox } from 'react-icons/fi';
import { cn } from '@/lib/utils';
import useAppStore from '@/store/useAppStore';

const Sidebar = () => {
    const { activeTab, setActiveTab } = useAppStore();

    const navItems = [
        {
            icon: FiGrid,
            label: 'Dashboard',
            value: 'dashboard' as const,
        },
        {
            icon: FiMessageSquare,
            label: 'AI Chat',
            value: 'chat' as const,
        },
        {
            icon: FiInbox,
            label: 'Inbox',
            value: 'inbox' as const,
            badge: useAppStore().notifications.filter(n => !n.read).length,
        },
    ];

    return (
        <div className="h-full w-16 bg-[#1a1d2d] flex flex-col items-center py-6">
            <div className="flex flex-col space-y-8 items-center">
                {navItems.map((item) => (
                    <button
                        key={item.value}
                        onClick={() => setActiveTab(item.value)}
                        className={cn(
                            "relative w-10 h-10 flex items-center justify-center rounded-lg transition-colors",
                            activeTab === item.value
                                ? "bg-blue-600 text-white"
                                : "text-gray-400 hover:text-white hover:bg-[#2a2d3d]"
                        )}
                        aria-label={item.label}
                    >
                        <item.icon size={20} />

                        {/* Badge for notifications */}
                        {item.badge ? (
                            <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs w-5 h-5 flex items-center justify-center rounded-full">
                                {item.badge}
                            </span>
                        ) : null}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default Sidebar;
