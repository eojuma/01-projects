"use client";

import { Sidebar } from "../components/Sidebar";
import { useNotificationContext } from "../context/NotificationContext";

const typeIcon: Record<string, React.ReactNode> = {
  follow_request: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <line x1="19" y1="8" x2="19" y2="14" />
      <line x1="22" y1="11" x2="16" y2="11" />
    </svg>
  ),
  group_invite: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  ),
  group_join_request: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M22 17a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V9.5C2 7 4 5 6.5 5H18c2.2 0 4 1.8 4 4v8Z" />
      <polyline points="22,9 12,14 2,9" />
    </svg>
  ),
  event: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
      <line x1="16" y1="2" x2="16" y2="6" />
      <line x1="8" y1="2" x2="8" y2="6" />
      <line x1="3" y1="10" x2="21" y2="10" />
    </svg>
  ),
  follow_accepted: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px", color: "var(--success)" }}>
      <polyline points="20 6 9 17 4 12" />
    </svg>
  ),
  like: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px", color: "var(--danger)" }}>
      <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
    </svg>
  ),
  default: (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
      <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
    </svg>
  )
};

export default function Notifications() {
  const {
    notifications,
    loading,
    error,
    unreadCount,
    markAllRead,
    markAsRead,
    acceptAction,
    declineAction,
    refresh,
  } = useNotificationContext();

  const formatTime = (ts: string) => {
    try {
      const d = new Date(ts);
      const diff = Date.now() - d.getTime();
      const mins = Math.floor(diff / 60000);
      if (mins < 1) return "just now";
      if (mins < 60) return `${mins}m ago`;
      const hrs = Math.floor(mins / 60);
      if (hrs < 24) return `${hrs}h ago`;
      return `${Math.floor(hrs / 24)}d ago`;
    } catch {
      return ts;
    }
  };

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div className="page-container">
          <div className="page-header">
            <h1 className="page-title">
              Notifications
              {unreadCount > 0 && (
                <span
                  className="tag tag-blue"
                  style={{ marginLeft: "10px", fontSize: "13px", verticalAlign: "middle" }}
                >
                  {unreadCount} new
                </span>
              )}
            </h1>
            <div style={{ display: "flex", gap: "8px" }}>
              {unreadCount > 0 && (
                <button className="btn btn-ghost btn-sm" onClick={markAllRead}>
                  Mark all read
                </button>
              )}
              <button
                className="btn btn-ghost btn-sm"
                onClick={refresh}
                style={{ display: "flex", alignItems: "center", gap: "6px" }}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "14px", height: "14px" }}>
                  <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67" />
                </svg>
                <span>Refresh</span>
              </button>
            </div>
          </div>

          {loading ? (
            <div className="loading-spinner"><div className="spinner" /></div>
          ) : error ? (
            <div className="card" style={{ color: "var(--danger)", textAlign: "center" }}>{error}</div>
          ) : notifications.length === 0 ? (
            <div className="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "12px" }}>
                <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
                <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
              </svg>
              <p>You're all caught up!</p>
            </div>
          ) : (
            notifications.map((notif) => (
              <div
                key={notif.id}
                className="card fade-in"
                onClick={() => { if (!notif.read) markAsRead(notif.id); }}
                style={{
                  display: "flex",
                  gap: "14px",
                  alignItems: "flex-start",
                  borderLeft: !notif.read ? "3px solid var(--accent)" : "3px solid transparent",
                  background: !notif.read ? "var(--bg-elevated)" : "var(--bg-surface)",
                  cursor: !notif.read ? "pointer" : "default",
                  transition: "background var(--transition)",
                }}
              >
                <div
                  className="avatar avatar-md"
                  style={{ background: "var(--bg-hover)", fontSize: "18px" }}
                >
                  {typeIcon[notif.type] || typeIcon.default}
                </div>
                <div style={{ flex: 1 }}>
                  <p style={{ fontSize: "14px", color: "var(--text-primary)", lineHeight: "1.5" }}>
                    {notif.message}
                  </p>
                  <p style={{ fontSize: "12px", color: "var(--text-muted)", marginTop: "4px" }}>
                    {formatTime(notif.timestamp)}
                  </p>
                  {notif.actionable && (
                    <div style={{ display: "flex", gap: "8px", marginTop: "12px" }}>
                      <button
                        className="btn btn-primary btn-sm"
                        onClick={(e) => { e.stopPropagation(); acceptAction(notif); }}
                      >
                        Accept
                      </button>
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={(e) => { e.stopPropagation(); declineAction(notif); }}
                      >
                        Decline
                      </button>
                    </div>
                  )}
                </div>
                {!notif.read && (
                  <span className="dot" style={{ marginTop: "6px", flexShrink: 0 }} />
                )}
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}