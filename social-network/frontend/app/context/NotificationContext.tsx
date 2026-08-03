"use client";

import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import { useWebSocket } from "./WebSocketContext";
import { useAuth } from "./AuthContext";
import {
  fetchNotifications,
  markAllNotificationsRead,
  respondToFollowRequest,
  respondToGroupInvite,
  respondToGroupJoinRequest,
  type Notification
} from "@/src/services/notifications";

interface NotificationContextType {
  notifications: Notification[];
  loading: boolean;
  error: string | null;
  unreadCount: number;
  unreadChatCount: number;
  markAllRead: () => Promise<void>;
  markAsRead: (id: string) => Promise<void>;
  acceptAction: (notif: Notification) => Promise<void>;
  declineAction: (notif: Notification) => Promise<void>;
  refresh: () => Promise<void>;
  refreshUnreadChat: () => Promise<void>;
}

const NotificationContext = createContext<NotificationContextType | null>(null);

export function NotificationProvider({ children }: { children: React.ReactNode }) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadChatCount, setUnreadChatCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { messages } = useWebSocket();
  const { isAuthenticated } = useAuth();

  const refreshUnreadChat = useCallback(async () => {
    if (!isAuthenticated) {
      setUnreadChatCount(0);
      return;
    }
    try {
      const chatRes = await fetch("/api/chat/conversations", { credentials: "include" });
      if (chatRes.ok) {
        const chats = await chatRes.json();
        const totalUnread = chats.reduce((sum: number, c: any) => sum + (c.unread_count || 0), 0);
        setUnreadChatCount(totalUnread);
      }
    } catch (e) {
      console.error("Failed to fetch chat unread counts:", e);
    }
  }, [isAuthenticated]);

  const load = useCallback(async () => {
    if (!isAuthenticated) {
      setNotifications([]);
      setUnreadChatCount(0);
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const data = await fetchNotifications();
      setNotifications(data);

      const chatRes = await fetch("/api/chat/conversations", { credentials: "include" });
      if (chatRes.ok) {
        const chats = await chatRes.json();
        const totalUnread = chats.reduce((sum: number, c: any) => sum + (c.unread_count || 0), 0);
        setUnreadChatCount(totalUnread);
      }
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    load();
    // Poll every 30 seconds for new notifications as fallback
    const interval = setInterval(load, 30000);
    return () => clearInterval(interval);
  }, [load]);

  // Listen to incoming WebSocket messages
  useEffect(() => {
    if (messages.length === 0) return;
    const lastMsg = messages[messages.length - 1];
    if (lastMsg && lastMsg.type === "notification") {
      const notif = lastMsg.payload as Notification;
      if (notif && notif.id) {
        setNotifications((prev) => {
          if (prev.some((n) => n.id === notif.id)) return prev;
          return [notif, ...prev];
        });
      }
    } else if (lastMsg && lastMsg.type === "chat") {
      refreshUnreadChat();
    }
  }, [messages, refreshUnreadChat]);

  const markAllRead = async () => {
    try {
      await markAllNotificationsRead();
      setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    } catch (e) {
      console.error("Failed to mark notifications read:", e);
    }
  };

  const markAsRead = async (notifId: string) => {
    try {
      await fetch(`/api/notifications/read?id=${notifId}`, {
        method: "POST",
        credentials: "include",
      });
      setNotifications((prev) =>
        prev.map((n) => (n.id === notifId ? { ...n, read: true } : n))
      );
    } catch (e) {
      console.error("Failed to mark notification read:", e);
    }
  };

  const acceptAction = async (notif: Notification) => {
    setNotifications((prev) => prev.filter((n) => n.id !== notif.id));
    try {
      if (notif.type === "follow_request" && notif.ref_id) {
        await respondToFollowRequest(notif.ref_id, "accept");
      } else if (notif.type === "group_invite" && notif.ref_id) {
        await respondToGroupInvite(notif.group_id || "any", notif.ref_id, "accept");
      } else if (notif.type === "group_join_request" && notif.ref_id) {
        await respondToGroupJoinRequest(notif.group_id || "any", notif.ref_id, "approve");
      }
    } catch (e) {
      console.error("Action failed:", e);
      load();
    }
  };

  const declineAction = async (notif: Notification) => {
    setNotifications((prev) => prev.filter((n) => n.id !== notif.id));
    try {
      if (notif.type === "follow_request" && notif.ref_id) {
        await respondToFollowRequest(notif.ref_id, "decline");
      } else if (notif.type === "group_invite" && notif.ref_id) {
        await respondToGroupInvite(notif.group_id || "any", notif.ref_id, "decline");
      } else if (notif.type === "group_join_request" && notif.ref_id) {
        await respondToGroupJoinRequest(notif.group_id || "any", notif.ref_id, "decline");
      }
    } catch (e) {
      console.error("Action failed:", e);
      load();
    }
  };

  const unreadCount = notifications.filter((n) => !n.read).length;

  return (
    <NotificationContext.Provider
      value={{
        notifications,
        loading,
        error,
        unreadCount,
        unreadChatCount,
        markAllRead,
        markAsRead,
        acceptAction,
        declineAction,
        refresh: load,
        refreshUnreadChat,
      }}
    >
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotificationContext() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error("useNotificationContext must be used within a NotificationProvider");
  }
  return context;
}
