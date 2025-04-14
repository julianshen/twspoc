import React from 'react';
import useAppStore from '@/store/useAppStore';
import { FiMessageSquare, FiFileText, FiCheckSquare } from 'react-icons/fi';

const Dashboard = () => {
    const { dashboards, activeDashboardId, setActiveDashboardId } = useAppStore();

    const activeDashboard = dashboards.find(d => d.id === activeDashboardId) || dashboards[0];

    return (
        <div className="p-6 h-full">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Dashboard</h1>

                {/* Dashboard Selector */}
                <div className="relative">
                    <select
                        value={activeDashboardId || ''}
                        onChange={(e) => setActiveDashboardId(e.target.value)}
                        className="bg-white border border-gray-300 rounded-md py-2 px-4 pr-8 appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        {dashboards.map(dashboard => (
                            <option key={dashboard.id} value={dashboard.id}>
                                {dashboard.name}
                            </option>
                        ))}
                    </select>
                    <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none">
                        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                    </div>
                </div>
            </div>

            {/* Dashboard Widgets */}
            <div className="grid grid-cols-3 gap-6">
                {/* Summary Widget */}
                <div className="bg-white rounded-lg shadow p-6">
                    <h2 className="text-lg font-semibold mb-4">Summary</h2>
                    <div className="space-y-4">
                        <div className="flex items-center">
                            <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center mr-3">
                                <FiMessageSquare className="text-blue-600" size={20} />
                            </div>
                            <div>
                                <p className="text-sm text-gray-500">Messages</p>
                                <p className="text-xl font-semibold">5</p>
                            </div>
                        </div>
                        <div className="flex items-center">
                            <div className="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center mr-3">
                                <FiFileText className="text-purple-600" size={20} />
                            </div>
                            <div>
                                <p className="text-sm text-gray-500">Reports</p>
                                <p className="text-xl font-semibold">12</p>
                            </div>
                        </div>
                        <div className="flex items-center">
                            <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center mr-3">
                                <FiCheckSquare className="text-green-600" size={20} />
                            </div>
                            <div>
                                <p className="text-sm text-gray-500">Tasks</p>
                                <p className="text-xl font-semibold">8</p>
                            </div>
                        </div>
                        <div className="pt-2">
                            <p className="text-sm text-gray-500 mb-1">Task Completion Rate</p>
                            <div className="w-full bg-gray-200 rounded-full h-2.5">
                                <div className="bg-blue-600 h-2.5 rounded-full" style={{ width: '75%' }}></div>
                            </div>
                            <p className="text-right text-sm text-gray-500 mt-1">75%</p>
                        </div>
                    </div>
                </div>

                {/* Recent Activity Widget */}
                <div className="bg-white rounded-lg shadow p-6">
                    <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
                    <div className="space-y-4">
                        {activeDashboard.widgets
                            .find(w => w.id === 'recent-activity')?.data
                            .map((activity: any) => (
                                <div key={activity.id} className="flex items-center">
                                    <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center mr-3">
                                        <span className="text-white text-xs">
                                            {activity.type === 'login' && 'L'}
                                            {activity.type === 'file-upload' && 'F'}
                                            {activity.type === 'message-sent' && 'M'}
                                        </span>
                                    </div>
                                    <div className="flex-1">
                                        <p className="text-sm font-medium">
                                            {activity.type === 'login' && 'Login'}
                                            {activity.type === 'file-upload' && 'File upload'}
                                            {activity.type === 'message-sent' && 'Message sent'}
                                        </p>
                                        <div className="flex justify-between">
                                            <p className="text-xs text-gray-500">{activity.time}</p>
                                            <p className="text-xs text-gray-500">{activity.date}</p>
                                        </div>
                                    </div>
                                </div>
                            )) || []}
                    </div>
                </div>

                {/* Notifications Widget */}
                <div className="bg-white rounded-lg shadow p-6">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-lg font-semibold">Notifications</h2>
                        <a href="#" className="text-xs text-blue-600 hover:underline">View all (1 unread)</a>
                    </div>
                    <div className="space-y-4">
                        {useAppStore().notifications.slice(0, 3).map(notification => (
                            <div key={notification.id} className={`p-3 rounded-lg ${notification.read ? 'bg-gray-50' : 'bg-blue-50'}`}>
                                <div className="flex justify-between items-start">
                                    <p className="text-sm font-medium">{notification.title}</p>
                                    <p className="text-xs text-gray-500">
                                        {notification.timestamp.getMinutes()} minutes ago
                                    </p>
                                </div>
                                <p className="text-sm text-gray-600 mt-1">{notification.message}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
