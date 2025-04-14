import React, { ReactNode } from 'react';
import Sidebar from './Sidebar';
import useAppStore from '@/store/useAppStore';

interface MainLayoutProps {
    children: ReactNode;
}

const MainLayout = ({ children }: MainLayoutProps) => {
    const { activeObject, setActiveObject } = useAppStore();

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <div className="flex-1 flex overflow-hidden">
                {/* Main Panel */}
                <div className={`flex-1 overflow-auto transition-all ${activeObject ? 'pr-4' : ''}`}>
                    {children}
                </div>

                {/* Object Viewer Panel (conditionally rendered) */}
                {activeObject && (
                    <div className="w-1/3 min-w-[350px] max-w-[500px] bg-white shadow-lg overflow-auto p-4">
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-semibold">{activeObject.title}</h3>
                            <button
                                onClick={() => setActiveObject(null)}
                                className="text-gray-500 hover:text-gray-700"
                                aria-label="Close"
                            >
                                ×
                            </button>
                        </div>
                        <div>
                            {/* Render different content based on object type */}
                            {activeObject.type === 'document' && (
                                <div className="prose">
                                    <p>{JSON.stringify(activeObject.data)}</p>
                                </div>
                            )}
                            {activeObject.type === 'task' && (
                                <div>
                                    <p>Task details would go here</p>
                                </div>
                            )}
                            {activeObject.type === 'other' && (
                                <div>
                                    <p>Object details would go here</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default MainLayout;
