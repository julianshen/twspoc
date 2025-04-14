import React, { useState, useRef, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { FiSend, FiPaperclip } from 'react-icons/fi';
import useAppStore, { ChatMessage, Attachment } from '@/store/useAppStore';

const AiChat = () => {
    const { chatMessages, addChatMessage, setActiveObject } = useAppStore();
    const [message, setMessage] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);

    // Auto-scroll to bottom when messages change
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [chatMessages]);

    // Handle sending a new message
    const handleSendMessage = () => {
        if (!message.trim()) return;

        // Add user message
        const userMessage: ChatMessage = {
            id: Date.now().toString(),
            sender: 'user',
            content: message,
            timestamp: new Date(),
        };
        addChatMessage(userMessage);
        setMessage('');

        // Simulate AI response after a short delay
        setTimeout(() => {
            const aiMessage: ChatMessage = {
                id: (Date.now() + 1).toString(),
                sender: 'ai',
                content: generateAiResponse(message),
                timestamp: new Date(),
                attachments: message.toLowerCase().includes('document') ? [
                    {
                        id: 'doc-' + Date.now(),
                        type: 'document',
                        title: 'Sample Document',
                        data: { content: 'This is a sample document attached to the AI response.' }
                    }
                ] : undefined
            };
            addChatMessage(aiMessage);
        }, 1000);
    };

    // Generate a simple AI response based on user input
    const generateAiResponse = (userMessage: string): string => {
        const lowerMessage = userMessage.toLowerCase();

        if (lowerMessage.includes('hello') || lowerMessage.includes('hi')) {
            return "Hello! How can I assist you today?";
        } else if (lowerMessage.includes('help')) {
            return "I'm here to help! You can ask me questions, request information, or discuss any topic you'd like.";
        } else if (lowerMessage.includes('document')) {
            return "I've attached a sample document for you. You can click on it to view the details.\n\n```json\n{\n  \"type\": \"document\",\n  \"title\": \"Sample Document\",\n  \"content\": \"This is a sample document\"\n}\n```";
        } else if (lowerMessage.includes('markdown')) {
            return "# Markdown Support\n\nThis chat supports markdown formatting:\n\n* **Bold text** for emphasis\n* *Italic text* for slight emphasis\n* `code blocks` for code\n* > Blockquotes for quotations\n\n```javascript\n// Code blocks with syntax highlighting\nfunction hello() {\n  console.log('Hello, world!');\n}\n```";
        } else {
            return "I understand you're interested in this topic. Can you tell me more about what you're looking for?";
        }
    };

    // Handle clicking on an attachment
    const handleAttachmentClick = (attachment: Attachment) => {
        setActiveObject(attachment);
    };

    return (
        <div className="flex flex-col h-full">
            <div className="p-6 border-b">
                <h1 className="text-2xl font-bold">AI Chat</h1>
            </div>

            {/* Messages Container */}
            <div className="flex-1 overflow-y-auto p-6 space-y-4">
                {chatMessages.length === 0 ? (
                    <div className="flex flex-col items-center justify-center h-full text-gray-500">
                        <p className="text-lg mb-2">No messages yet</p>
                        <p className="text-sm">Start a conversation with the AI assistant</p>
                    </div>
                ) : (
                    chatMessages.map((msg) => (
                        <div
                            key={msg.id}
                            className={`flex ${msg.sender === 'user' ? 'justify-end' : 'justify-start'}`}
                        >
                            <div
                                className={`max-w-[70%] rounded-lg p-4 ${msg.sender === 'user'
                                        ? 'bg-blue-600 text-white'
                                        : 'bg-gray-200 text-gray-800'
                                    }`}
                            >
                                <div className="prose max-w-none dark:prose-invert">
                                    <ReactMarkdown>
                                        {msg.content}
                                    </ReactMarkdown>
                                </div>

                                {/* Attachments */}
                                {msg.attachments && msg.attachments.length > 0 && (
                                    <div className="mt-3 pt-3 border-t border-gray-300 dark:border-gray-700">
                                        <p className="text-sm mb-2">Attachments:</p>
                                        <div className="flex flex-wrap gap-2">
                                            {msg.attachments.map(attachment => (
                                                <button
                                                    key={attachment.id}
                                                    onClick={() => handleAttachmentClick(attachment)}
                                                    className="flex items-center gap-1 px-3 py-1 bg-white bg-opacity-20 rounded text-sm hover:bg-opacity-30 transition-colors"
                                                >
                                                    <FiPaperclip size={14} />
                                                    <span>{attachment.title}</span>
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                <div className="text-xs mt-2 opacity-70 text-right">
                                    {new Intl.DateTimeFormat('en-US', {
                                        hour: '2-digit',
                                        minute: '2-digit'
                                    }).format(msg.timestamp)}
                                </div>
                            </div>
                        </div>
                    ))
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* Message Input */}
            <div className="p-4 border-t">
                <div className="flex items-center bg-white rounded-lg border">
                    <input
                        type="text"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && handleSendMessage()}
                        placeholder="Type a message..."
                        className="flex-1 px-4 py-2 bg-transparent outline-none"
                    />
                    <button
                        onClick={handleSendMessage}
                        disabled={!message.trim()}
                        className="p-2 rounded-r-lg text-white bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed"
                    >
                        <FiSend size={20} />
                    </button>
                </div>
                <p className="text-xs text-gray-500 mt-2">
                    Tip: Try typing "markdown" to see formatting options
                </p>
            </div>
        </div>
    );
};

export default AiChat;
