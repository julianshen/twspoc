import React, { useState, useEffect } from 'react';
import { FiClock, FiCheck, FiFile, FiTag, FiFlag, FiTrash2 } from 'react-icons/fi';
import useAppStore from '@/store/useAppStore';
import { notificationService } from '@/lib/notificationService';

type FilterTab = 'all' | 'unread' | 'read';

const Inbox = () => {
    const {
        notifications,
        markNotificationAsRead,
        setActiveObject,
        fetchNotifications,
        addNotification,
        deleteNotification
    } = useAppStore();
    const [activeFilter, setActiveFilter] = useState<FilterTab>('all');

    // Fetch notifications on component mount
    useEffect(() => {
        // Initial fetch
        fetchNotifications();

        // Subscribe to SSE for real-time updates
        const { cleanup } = notificationService.subscribeToNotifications(
            (notification) => {
                console.log('Received notification via SSE:', notification);
                addNotification(notification);
            },
            (error) => {
                console.error('SSE error:', error);
            }
        );

        // Cleanup on unmount
        return cleanup;
    }, [fetchNotifications, addNotification]);

    // Handle clicking on a notification
    const handleNotificationClick = (id: string) => {
        markNotificationAsRead(id);
    };

    // Handle deleting a notification
    const handleDeleteNotification = (id: string, e: React.MouseEvent) => {
        e.stopPropagation(); // Prevent triggering the notification click
        deleteNotification(id);
    };

    // Handle clicking on an attachment
    const handleAttachmentClick = (notificationId: string, e: React.MouseEvent) => {
        e.stopPropagation(); // Prevent triggering the notification click

        const notification = notifications.find(n => n.id === notificationId);
        if (notification?.attachment) {
            markNotificationAsRead(notificationId);
            setActiveObject(notification.attachment);
        }
    };

    // Format timestamp to relative time
    const formatRelativeTime = (date: Date) => {
        const now = new Date();
        const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60));

        if (diffInMinutes < 1) return 'Just now';
        if (diffInMinutes < 60) return `${diffInMinutes} minute${diffInMinutes === 1 ? '' : 's'} ago`;

        const diffInHours = Math.floor(diffInMinutes / 60);
        if (diffInHours < 24) return `${diffInHours} hour${diffInHours === 1 ? '' : 's'} ago`;

        const diffInDays = Math.floor(diffInHours / 24);
        return `${diffInDays} day${diffInDays === 1 ? '' : 's'} ago`;
    };

    // Filter notifications based on active tab
    const filteredNotifications = notifications.filter(notification => {
        if (activeFilter === 'all') return true;
        if (activeFilter === 'unread') return !notification.read;
        if (activeFilter === 'read') return notification.read;
        return true;
    });

    return (
        <div className="flex flex-col h-full">
            <div className="p-6 border-b">
                <h1 className="text-2xl font-bold">Inbox</h1>
                <p className="text-gray-500 mt-1">
                    {notifications.filter(n => !n.read).length} unread notifications
                </p>
            </div>

            {/* Filter buttons */}
            <div className="flex p-4 gap-2">
                <button
                    onClick={() => setActiveFilter('all')}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${activeFilter === 'all'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                >
                    All
                </button>
                <button
                    onClick={() => setActiveFilter('unread')}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${activeFilter === 'unread'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                >
                    Unread
                </button>
                <button
                    onClick={() => setActiveFilter('read')}
                    className={`px-4 py-2 rounded-md font-medium transition-colors ${activeFilter === 'read'
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                >
                    Read
                </button>
            </div>

            <div className="flex-1 overflow-y-auto">
                {filteredNotifications.length === 0 ? (
                    <div className="flex flex-col items-center justify-center h-full text-gray-500">
                        <p className="text-lg mb-2">No notifications</p>
                        <p className="text-sm">You're all caught up!</p>
                    </div>
                ) : (
                    <div className="divide-y">
                        {filteredNotifications.map((notification) => (
                            <div
                                key={notification.id}
                                onClick={() => handleNotificationClick(notification.id)}
                                className={`p-6 hover:bg-gray-50 cursor-pointer transition-colors ${notification.read ? 'bg-white' : 'bg-blue-50'
                                    }`}
                            >
                                <div className="flex items-start">
                                    <div
                                        className={`w-2 h-2 mt-2 rounded-full flex-shrink-0 ${notification.read ? 'bg-gray-300' : 'bg-blue-600'
                                            }`}
                                    />

                                    <div className="ml-4 flex-1">
                                        <div className="flex justify-between">
                                            <div className="flex items-center">
                                                <h3 className="font-medium">{notification.title}</h3>
                                                <span
                                                    className={`ml-2 px-1.5 py-0.5 text-xs rounded-full flex items-center ${notification.priority === 'high'
                                                        ? 'bg-red-100 text-red-800'
                                                        : notification.priority === 'medium'
                                                            ? 'bg-yellow-100 text-yellow-800'
                                                            : 'bg-green-100 text-green-800'
                                                        }`}
                                                >
                                                    <FiFlag className="mr-1" size={10} />
                                                    {notification.priority}
                                                </span>
                                            </div>
                                            <span className="text-sm text-gray-500 flex items-center">
                                                <FiClock className="mr-1" size={14} />
                                                {formatRelativeTime(notification.timestamp)}
                                            </span>
                                        </div>

                                        <p className="text-gray-600 mt-1">{notification.message}</p>

                                        {/* Labels */}
                                        {notification.labels && notification.labels.length > 0 && (
                                            <div className="mt-2 flex flex-wrap gap-2">
                                                {notification.labels.map((label, index) => (
                                                    <span
                                                        key={index}
                                                        className="inline-flex items-center px-2 py-1 bg-gray-100 rounded-md text-xs text-gray-700"
                                                    >
                                                        <FiTag className="mr-1" size={12} />
                                                        {label}
                                                    </span>
                                                ))}
                                            </div>
                                        )}

                                        {notification.attachment && (
                                            <div
                                                onClick={(e) => handleAttachmentClick(notification.id, e)}
                                                className="mt-3 inline-flex items-center px-3 py-1.5 bg-gray-100 hover:bg-gray-200 rounded-md text-sm text-gray-700 transition-colors"
                                            >
                                                <FiFile className="mr-2" size={14} />
                                                {notification.attachment.title}
                                            </div>
                                        )}

                                        <div className="mt-2 flex items-center justify-between">
                                            {notification.read && (
                                                <div className="flex items-center text-sm text-gray-500">
                                                    <FiCheck className="mr-1" size={14} />
                                                    Read
                                                </div>
                                            )}
                                            <button
                                                onClick={(e) => handleDeleteNotification(notification.id, e)}
                                                className="ml-auto p-1 text-gray-400 hover:text-red-500 transition-colors"
                                                aria-label="Delete notification"
                                            >
                                                <FiTrash2 size={16} />
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Inbox;
