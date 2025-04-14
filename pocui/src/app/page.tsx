'use client';

import React from 'react';
import MainLayout from '@/components/layout/MainLayout';
import Dashboard from '@/components/tabs/Dashboard';
import AiChat from '@/components/tabs/AiChat';
import Inbox from '@/components/tabs/Inbox';
import useAppStore from '@/store/useAppStore';

export default function Home() {
  const { activeTab } = useAppStore();

  // Render the appropriate tab based on the active tab
  const renderActiveTab = () => {
    switch (activeTab) {
      case 'dashboard':
        return <Dashboard />;
      case 'chat':
        return <AiChat />;
      case 'inbox':
        return <Inbox />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <MainLayout>
      {renderActiveTab()}
    </MainLayout>
  );
}
