"use client";
import React, { useState } from "react";
import styles from "./SimpleUILayout.module.css";

const TABS = [
  { key: "dashboard", label: "Dashboard", icon: "🏠" },
  { key: "aichat", label: "AI Chat", icon: "🤖" },
  { key: "inbox", label: "Inbox", icon: "📥" },
];

export default function SimpleUILayout({ children }: { children: React.ReactNode }) {
  const [activeTab, setActiveTab] = useState("dashboard");
  const [showObjectViewer, setShowObjectViewer] = useState(false);

  // For demo: show object viewer when a button is clicked in AI Chat or Inbox
  const handleShowObjectViewer = () => setShowObjectViewer(true);
  const handleHideObjectViewer = () => setShowObjectViewer(false);

  return (
    <div className={styles.root}>
      <aside className={styles.sidebar}>
        <div className={styles.avatarContainer}>
          <img
            src="https://i.pravatar.cc/48"
            alt="User Avatar"
            className={styles.avatar}
          />
        </div>
        <nav className={styles.tabNav}>
          {TABS.map((tab) => (
            <button
              key={tab.key}
              className={`${styles.tabButton} ${activeTab === tab.key ? styles.active : ""}`}
              onClick={() => {
                setActiveTab(tab.key);
                setShowObjectViewer(false);
              }}
              title={tab.label}
            >
              <span className={styles.tabIcon}>{tab.icon}</span>
            </button>
          ))}
        </nav>
      </aside>
      <main className={styles.main}>
        {activeTab === "dashboard" && (
          <section className={styles.panel}>
            <h1>Dashboard</h1>
            <p>Welcome to the dashboard.</p>
          </section>
        )}
        {(activeTab === "aichat" || activeTab === "inbox") && (
          <section className={styles.twoPanel}>
            <div className={styles.leftPanel}>
              <h1>{activeTab === "aichat" ? "AI Chat" : "Inbox"}</h1>
              <p>
                {activeTab === "aichat"
                  ? "This is the AI chat panel."
                  : "This is your inbox."}
              </p>
              {activeTab === "aichat" && (
                <div style={{ margin: "24px 0" }}>
                  {children}
                </div>
              )}
              {!showObjectViewer && (
                <button onClick={handleShowObjectViewer} className={styles.showObjBtn}>
                  Show Object Viewer
                </button>
              )}
            </div>
            {showObjectViewer && (
              <div className={styles.rightPanel}>
                <div className={styles.objectViewerHeader}>
                  <span>Object Viewer</span>
                  <button onClick={handleHideObjectViewer} className={styles.closeBtn}>
                    ✕
                  </button>
                </div>
                <div className={styles.objectViewerContent}>
                  <p>Object details go here.</p>
                </div>
              </div>
            )}
          </section>
        )}
      </main>
    </div>
  );
}
