"use client";

import { createContext, useContext, useEffect, useState, ReactNode, useRef, useCallback } from "react";
import { useAuth } from "./AuthContext";

interface WebSocketContextType {
  isConnected: boolean;
  sendMessage: (type: string, to: string, payload: any) => void;
  messages: any[];
  connect: () => void;
  disconnect: () => void;
}

const WebSocketContext = createContext<WebSocketContextType | null>(null);

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<any[]>([]);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const socketRef = useRef<WebSocket | null>(null);
  // Stop retrying after 3 attempts — backend may not be running
  const maxReconnectAttempts = 1000;
  const isConnectingRef = useRef(false);
  const { isAuthenticated } = useAuth();

  const disconnect = useCallback(() => {
    isConnectingRef.current = false;
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (socketRef.current) {
      socketRef.current.close(1000, "Normal closure");
      socketRef.current = null;
    }
    setSocket(null);
    setIsConnected(false);
  }, []);

  const connect = useCallback(() => {
    // Don't connect if already open or connecting
    if (
      isConnectingRef.current ||
      (socketRef.current && socketRef.current.readyState === WebSocket.OPEN)
    ) {
      return;
    }

    isConnectingRef.current = true;
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsHost = window.location.port
      ? `${window.location.hostname}:8080`
      : window.location.host;
    const ws = new WebSocket(`${protocol}//${wsHost}/ws`);
    socketRef.current = ws;

    ws.onopen = () => {
      isConnectingRef.current = false;
      setIsConnected(true);
      reconnectAttempts.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setMessages((prev) => [...prev, data]);
      } catch {
        // ignore parse errors
      }
    };

    ws.onclose = (event) => {
      isConnectingRef.current = false;
      setIsConnected(false);
      socketRef.current = null;
      // Only retry on unexpected closure (not code 1000 = normal, 1001 = going away)
      if (event.code !== 1000 && event.code !== 1001) {
        if (reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++;
          // Exponential backoff capped at 10 seconds
          const delay = Math.min(10000, 3000 * Math.pow(2, reconnectAttempts.current - 1));
          reconnectTimeoutRef.current = setTimeout(connect, delay);
        }
      }
    };

    ws.onerror = () => {
      // Error is always followed by onclose — nothing to do here
    };

    setSocket(ws);
  }, []);

  useEffect(() => {
    if (isAuthenticated) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  const sendMessage = (type: string, to: string, payload: any) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify({ type, to, payload }));
    }
  };

  return (
    <WebSocketContext.Provider value={{ isConnected, sendMessage, messages, connect, disconnect }}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
}