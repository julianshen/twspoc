import { Notification } from '@/store/useAppStore';

// Base URL for the notification service
const NOTIFICATION_API_BASE_URL = 'http://localhost:8080';

// Default user ID (in a real app, this would come from authentication)
const DEFAULT_USER_ID = 'user123';

// Enable mock mode for development without a real backend
const USE_MOCK_MODE = false;

// Mock notifications for development
const MOCK_NOTIFICATIONS: Notification[] = [
    {
        id: '1',
        title: 'Test App',
        message: 'This is a test notification created directly via Redis at 下午 6:05:14',
        timestamp: new Date(Date.now() - 27 * 60 * 1000),
        read: false,
        priority: 'high',
        labels: ['System', 'Important'],
        attachment: {
            id: 'doc1',
            type: 'document',
            title: 'Direct Redis Test',
            data: { content: 'Test document content' },
        },
    },
    {
        id: '2',
        title: 'Test App',
        message: 'This is a test notification created via debug script',
        timestamp: new Date(Date.now() - 38 * 60 * 1000),
        read: false,
        priority: 'low',
        labels: ['Debug', 'Low Priority'],
    },
];

/**
 * Notification Service client for interacting with the notification API
 */
export class NotificationService {
    private baseUrl: string;
    private userId: string;

    constructor(baseUrl = NOTIFICATION_API_BASE_URL, userId = DEFAULT_USER_ID) {
        this.baseUrl = baseUrl;
        this.userId = userId;
    }

    /**
     * Get all notifications for the current user
     */
    async getNotifications(): Promise<Notification[]> {
        if (USE_MOCK_MODE) {
            console.log('Using mock notifications data');
            return [...MOCK_NOTIFICATIONS];
        }

        try {
            const response = await fetch(`${this.baseUrl}/api/notifications?userId=${this.userId}`);

            if (!response.ok) {
                throw new Error(`Failed to fetch notifications: ${response.status}`);
            }

            const data = await response.json();
            return this.mapNotifications(data);
        } catch (error) {
            console.error('Error fetching notifications:', error);
            // Fall back to mock data if API fails
            console.log('Falling back to mock notifications data');
            return [...MOCK_NOTIFICATIONS];
        }
    }

    /**
     * Mark a notification as read
     */
    async markAsRead(notificationId: string): Promise<boolean> {
        if (USE_MOCK_MODE) {
            console.log(`Mock: Marking notification ${notificationId} as read`);
            return true;
        }

        try {
            const response = await fetch(
                `${this.baseUrl}/api/notifications/${notificationId}/read`,
                { method: 'POST' }
            );

            return response.ok;
        } catch (error) {
            console.error('Error marking notification as read:', error);
            return true; // Return success in case of error to allow UI to update
        }
    }

    /**
     * Delete a notification
     */
    async deleteNotification(notificationId: string): Promise<boolean> {
        if (USE_MOCK_MODE) {
            console.log(`Mock: Deleting notification ${notificationId}`);
            return true;
        }

        try {
            const response = await fetch(
                `${this.baseUrl}/api/notifications/${notificationId}`,
                { method: 'DELETE' }
            );

            return response.ok;
        } catch (error) {
            console.error('Error deleting notification:', error);
            return true; // Return success in case of error to allow UI to update
        }
    }

    /**
     * Subscribe to notifications using Server-Sent Events (SSE)
     * Returns an EventSource object and a cleanup function
     * In mock mode, it simulates SSE with a timer
     */
    subscribeToNotifications(
        onNotification: (notification: Notification) => void,
        onError: (error: Event) => void
    ): { eventSource: EventSource | null; cleanup: () => void } {
        if (USE_MOCK_MODE) {
            console.log('Using mock SSE subscription');

            // Set up a timer to simulate new notifications every 15 seconds
            const mockInterval = setInterval(() => {
                const mockNotification: Notification = {
                    id: `mock-${Date.now()}`,
                    title: 'New Notification',
                    message: `This is a mock notification created at ${new Date().toLocaleTimeString()}`,
                    timestamp: new Date(),
                    read: false,
                    priority: Math.random() > 0.7 ? 'high' : Math.random() > 0.4 ? 'medium' : 'low',
                    labels: ['Mock', 'Automated'],
                };

                console.log('Mock SSE: New notification', mockNotification);
                onNotification(mockNotification);
            }, 15000);

            // Return null for eventSource and a cleanup function
            return {
                eventSource: null,
                cleanup: () => {
                    clearInterval(mockInterval);
                },
            };
        }

        try {
            const eventSource = new EventSource(
                `${this.baseUrl}/api/notifications/subscribe?userId=${this.userId}`
            );

            eventSource.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    const notification = this.mapNotification(data);
                    onNotification(notification);
                } catch (error) {
                    console.error('Error parsing SSE message:', error);
                }
            };

            eventSource.onerror = onError;

            // Return the event source and a cleanup function
            return {
                eventSource,
                cleanup: () => {
                    eventSource.close();
                },
            };
        } catch (error) {
            console.error('Error setting up SSE:', error);
            return {
                eventSource: null,
                cleanup: () => { },
            };
        }
    }

    /**
     * Map notification data from the API to our frontend model
     */
    private mapNotification(data: any): Notification {
        return {
            id: data.id,
            title: data.title,
            message: data.message,
            timestamp: new Date(data.timestamp),
            read: data.read,
            priority: this.mapPriority(data.priority),
            labels: data.labels || [],
            attachment: data.attachments && data.attachments.length > 0
                ? {
                    id: data.attachments[0].id,
                    type: this.mapAttachmentType(data.attachments[0].type),
                    title: data.attachments[0].type,
                    data: { url: data.attachments[0].url },
                }
                : undefined,
        };
    }

    /**
     * Map an array of notifications
     */
    private mapNotifications(data: any[]): Notification[] {
        return data.map(item => this.mapNotification(item));
    }

    /**
     * Map priority from string to our enum
     */
    private mapPriority(priority: string): 'high' | 'medium' | 'low' {
        switch (priority?.toLowerCase()) {
            case 'high':
                return 'high';
            case 'medium':
                return 'medium';
            case 'low':
            default:
                return 'low';
        }
    }

    /**
     * Map attachment type from API to our model
     */
    private mapAttachmentType(type: string): 'document' | 'task' | 'other' {
        switch (type?.toLowerCase()) {
            case 'document':
                return 'document';
            case 'task':
                return 'task';
            default:
                return 'other';
        }
    }
}

// Export a singleton instance
export const notificationService = new NotificationService();
