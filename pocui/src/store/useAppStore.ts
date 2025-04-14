import { create } from 'zustand';
import { notificationService } from '@/lib/notificationService';

// Define types for our data models
export type Dashboard = {
    id: string;
    name: string;
    widgets: Widget[];
};

export type Widget = {
    id: string;
    type: string;
    title: string;
    data: any;
    position: { x: number; y: number; w: number; h: number };
};

export type ChatMessage = {
    id: string;
    sender: 'user' | 'ai';
    content: string;
    timestamp: Date;
    attachments?: Attachment[];
};

export type Notification = {
    id: string;
    title: string;
    message: string;
    timestamp: Date;
    read: boolean;
    priority: 'high' | 'medium' | 'low';
    attachment?: Attachment;
    labels?: string[];
};

export type Attachment = {
    id: string;
    type: 'document' | 'task' | 'other';
    title: string;
    data: any;
};

export type Tab = 'dashboard' | 'chat' | 'inbox';

// Define the store state
interface AppState {
    // Navigation
    activeTab: Tab;
    setActiveTab: (tab: Tab) => void;

    // Dashboard state
    dashboards: Dashboard[];
    activeDashboardId: string | null;
    setActiveDashboardId: (id: string) => void;
    addDashboard: (dashboard: Dashboard) => void;

    // Chat state
    chatMessages: ChatMessage[];
    addChatMessage: (message: ChatMessage) => void;

    // Inbox state
    notifications: Notification[];
    markNotificationAsRead: (id: string) => void;
    fetchNotifications: () => Promise<void>;
    addNotification: (notification: Notification) => void;
    deleteNotification: (id: string) => Promise<void>;

    // Object viewer state
    activeObject: Attachment | null;
    setActiveObject: (object: Attachment | null) => void;
}

// Create the store
const useAppStore = create<AppState>((set) => ({
    // Navigation
    activeTab: 'dashboard',
    setActiveTab: (tab) => set({ activeTab: tab }),

    // Dashboard state
    dashboards: [
        {
            id: 'default',
            name: 'Main Dashboard',
            widgets: [
                {
                    id: 'summary',
                    type: 'summary',
                    title: 'Summary',
                    data: {
                        messages: 5,
                        reports: 12,
                        tasks: 8,
                    },
                    position: { x: 0, y: 0, w: 1, h: 1 },
                },
                {
                    id: 'recent-activity',
                    type: 'activity',
                    title: 'Recent Activity',
                    data: [
                        { id: '1', type: 'login', time: '09:32 AM', date: 'Today' },
                        { id: '2', type: 'file-upload', time: '11:15 AM', date: 'Today' },
                        { id: '3', type: 'message-sent', time: '02:45 PM', date: 'Yesterday' },
                    ],
                    position: { x: 1, y: 0, w: 1, h: 1 },
                },
                {
                    id: 'notifications-widget',
                    type: 'notifications',
                    title: 'Notifications',
                    data: [],
                    position: { x: 2, y: 0, w: 1, h: 1 },
                },
            ],
        },
        {
            id: 'analytics',
            name: 'Analytics Dashboard',
            widgets: [],
        },
    ],
    activeDashboardId: 'default',
    setActiveDashboardId: (id) => set({ activeDashboardId: id }),
    addDashboard: (dashboard) =>
        set((state) => ({ dashboards: [...state.dashboards, dashboard] })),

    // Chat state
    chatMessages: [],
    addChatMessage: (message) =>
        set((state) => ({ chatMessages: [...state.chatMessages, message] })),

    // Inbox state
    notifications: [],
    markNotificationAsRead: (id) =>
        set((state) => {
            // Update local state
            const updatedNotifications = state.notifications.map((notification) =>
                notification.id === id ? { ...notification, read: true } : notification
            );

            // Call API to mark as read (fire and forget)
            notificationService.markAsRead(id).catch(error => {
                console.error('Failed to mark notification as read:', error);
            });

            return { notifications: updatedNotifications };
        }),

    fetchNotifications: async () => {
        try {
            const notifications = await notificationService.getNotifications();
            set({ notifications });
        } catch (error) {
            console.error('Failed to fetch notifications:', error);
        }
    },

    addNotification: (notification) =>
        set((state) => ({
            notifications: [notification, ...state.notifications]
        })),

    deleteNotification: async (id) => {
        try {
            const success = await notificationService.deleteNotification(id);
            if (success) {
                set((state) => ({
                    notifications: state.notifications.filter(n => n.id !== id)
                }));
            }
        } catch (error) {
            console.error('Failed to delete notification:', error);
        }
    },

    // Object viewer state
    activeObject: null,
    setActiveObject: (object) => set({ activeObject: object }),
}));

export default useAppStore;
