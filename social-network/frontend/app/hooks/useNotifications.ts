"use client";
import { useState, useEffect, useCallback } from "react";
import {
  fetchNotifications,
  markAllNotificationsRead,
  respondToFollowRequest,
  respondToGroupInvite,
  respondToGroupJoinRequest,
  type Notification
} from "@/src/services/notifications";

export function useNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      setLoading(true);
      const data = await fetchNotifications();
      setNotifications(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
    // Poll every 30 seconds for new notifications
    const interval = setInterval(load, 30000);
    return () => clearInterval(interval);
  }, [load]);

  const markAllRead = async () => {
    try {
      await markAllNotificationsRead();
      setNotifications(prev => prev.map(n => ({ ...n, read: true })));
    } catch (e) {
      console.error("Failed to mark notifications read:", e);
    }
  };

  const acceptAction = async (notif: Notification) => {
    setNotifications(prev => prev.filter(n => n.id !== notif.id));
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
    setNotifications(prev => prev.filter(n => n.id !== notif.id));
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

  const markAsRead = async (notifId: string) => {
    try {
      await fetch(`/api/notifications/read?id=${notifId}`, {
        method: "POST",
        credentials: "include",
      });
      setNotifications(prev =>
        prev.map(n => (n.id === notifId ? { ...n, read: true } : n))
      );
    } catch (e) {
      console.error("Failed to mark notification read:", e);
    }
  };

  const unreadCount = notifications.filter(n => !n.read).length;

  return {
    notifications,
    loading,
    error,
    unreadCount,
    markAllRead,
    markAsRead,
    acceptAction,
    declineAction,
    refresh: load
  };
}