"use client";
import { useState, useEffect, useRef, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { Sidebar } from "../components/Sidebar";
import { useWebSocket } from "../context/WebSocketContext";
import { useNotificationContext } from "../context/NotificationContext";
import { fetchConversations, fetchMessages, type Conversation, type Message } from "@/src/services/chat";
import { EmojiPicker } from "../components/EmojiPicker";

function isImageOrGif(content: string) {
  if (!content) return false;
  const trimmed = content.trim();
  if (trimmed.startsWith("/uploads/")) return true;
  if (/^https?:\/\//i.test(trimmed)) {
    if (trimmed.includes("tenor.com") || trimmed.includes("giphy.com") || trimmed.includes("googleusercontent.com")) return true;
    if (/\.(png|jpe?g|gif|webp|svg)(\?.*)?$/i.test(trimmed)) return true;
  }
  return false;
}

function ChatContent() {
  const searchParams = useSearchParams();
  const queryConvId = searchParams.get("convId");

  const { isConnected, sendMessage, messages: wsMessages } = useWebSocket();
  const { refreshUnreadChat } = useNotificationContext();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConvId, setSelectedConvId] = useState<string | null>(null);
  const [chatMessages, setChatMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const [myId, setMyId] = useState<string | null>(null);

  // GIF states
  const [showGifPanel, setShowGifPanel] = useState(false);
  const [gifQuery, setGifQuery] = useState("");
  const [gifs, setGifs] = useState<string[]>([]);
  const [loadingGifs, setLoadingGifs] = useState(false);

  // New chat state
  const [showNewChat, setShowNewChat] = useState(false);
  const [userSearch, setUserSearch] = useState("");
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [searching, setSearching] = useState(false);

  // Load current user info
  useEffect(() => {
    fetch("/api/auth/me", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data && data.id) {
          setMyId(String(data.id));
        }
      })
      .catch(console.error);
  }, []);

  // Load conversations
  useEffect(() => {
    fetchConversations()
      .then((data) => {
        setConversations(data);
        if (queryConvId && data.some((c) => c.id === queryConvId)) {
          setSelectedConvId(queryConvId);
        } else if (data.length > 0) {
          setSelectedConvId(data[0].id);
        }
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [queryConvId]);

  // Load messages when conversation changes
  useEffect(() => {
    if (!selectedConvId) return;
    setLoadingMessages(true);
    fetchMessages(selectedConvId)
      .then((msgs) => {
        setChatMessages(msgs);
        refreshUnreadChat();
      })
      .catch(console.error)
      .finally(() => setLoadingMessages(false));
  }, [selectedConvId, refreshUnreadChat]);

  // Handle incoming WS messages — backend sends {type:"chat", payload:{id, from, content, timestamp, conversation_id}}
  useEffect(() => {
    if (wsMessages.length === 0) return;
    const last = wsMessages[wsMessages.length - 1];
    if (last?.type === "chat") {
      const payload = last.payload || {};
      const convId = payload.conversation_id;
      if (convId === selectedConvId) {
        const normalized: Message = {
          id: payload.id || `ws-${Date.now()}`,
          conversation_id: convId,
          sender_id: String(payload.from || ""),
          content: payload.content || "",
          created_at: payload.timestamp || Date.now() / 1000,
        };
        setChatMessages((prev) => {
          // Avoid duplicating optimistic message (same content + near same time)
          if (prev.some((m) => m.id === normalized.id)) return prev;
          // Replace matching optimistic tmp message
          const filtered = prev.filter(
            (m) =>
              !(
                m.id.startsWith("tmp-") &&
                m.content === normalized.content &&
                m.sender_id === normalized.sender_id
              )
          );
          return [...filtered, normalized];
        });
        // Silently mark as read in DB and update context count
        if (selectedConvId) {
          fetchMessages(selectedConvId)
            .then(() => refreshUnreadChat())
            .catch(console.error);
        }
      }
      // Refresh conversation list to update last_message preview
      fetchConversations().then(setConversations).catch(() => {});
    }
  }, [wsMessages, selectedConvId, refreshUnreadChat]);

  // Scroll to bottom on new messages
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chatMessages]);

  const handleUserSearch = async (q: string) => {
    setUserSearch(q);
    if (q.length < 2) { setSearchResults([]); return; }
    setSearching(true);
    try {
      const res = await fetch(`/api/users?chat=true&q=${encodeURIComponent(q)}`, { credentials: "include" });
      if (res.ok) setSearchResults(await res.json());
    } catch {}
    setSearching(false);
  };

  const startConversation = async (userId: string) => {
    try {
      const res = await fetch(`/api/chat/user/${userId}`, { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        const convId = data.conversation_id;
        const convs = await fetchConversations();
        setConversations(convs);
        setSelectedConvId(convId);
        setShowNewChat(false);
        setUserSearch("");
        setSearchResults([]);
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleSend = () => {
    if (!inputMessage.trim() || !selectedConvId || !myId) return;
    const conv = conversations.find((c) => c.id === selectedConvId);
    if (!conv) return;
    sendMessage("chat", conv.other_user_id, { content: inputMessage, conversation_id: selectedConvId });
    // Optimistic message
    const tmpId = `tmp-${Date.now()}`;
    setChatMessages((prev) => [
      ...prev,
      {
        id: tmpId,
        conversation_id: selectedConvId,
        sender_id: myId,
        content: inputMessage,
        created_at: Date.now() / 1000,
      },
    ]);
    setInputMessage("");
  };

  const sendMediaMessage = (url: string) => {
    if (!selectedConvId || !myId) return;
    const conv = conversations.find((c) => c.id === selectedConvId);
    if (!conv) return;
    sendMessage("chat", conv.other_user_id, { content: url, conversation_id: selectedConvId });
    // Optimistic message
    const tmpId = `tmp-${Date.now()}`;
    setChatMessages((prev) => [
      ...prev,
      {
        id: tmpId,
        conversation_id: selectedConvId,
        sender_id: myId,
        content: url,
        created_at: Date.now() / 1000,
      },
    ]);
  };

  // Curated GIF collection organized by category
  const gifLibrary: Record<string, string[]> = {
    trending: [
      "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
      "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
      "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
      "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
      "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
      "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
      "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
      "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
      "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
      "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
      "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
      "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
    ],
    happy: [
      "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
      "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
      "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
      "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
      "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
      "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
      "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
      "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
      "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
      "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
      "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
      "https://media1.giphy.com/media/26tknCqiJrBQG6bxC/200.gif",
    ],
    sad: [
      "https://media1.giphy.com/media/d2lcHJTG5Tscg/200.gif",
      "https://media1.giphy.com/media/OPU6wzx8JrHna/200.gif",
      "https://media1.giphy.com/media/9Y5BbDSkSTiY8/200.gif",
      "https://media1.giphy.com/media/3o6wrvdHFbwBrUFenu/200.gif",
      "https://media1.giphy.com/media/l2JhtKtDWYNKdRpoA/200.gif",
      "https://media1.giphy.com/media/3ohzdIuqJoo8QdKlnW/200.gif",
      "https://media1.giphy.com/media/26ufnwz3wDUli7GU0/200.gif",
      "https://media1.giphy.com/media/ISOckXUybVfQ4/200.gif",
      "https://media1.giphy.com/media/BEob5qwFkSJ7G/200.gif",
      "https://media1.giphy.com/media/L95W4wv8nnb9K/200.gif",
      "https://media1.giphy.com/media/ROF8OQvDmxytW/200.gif",
      "https://media1.giphy.com/media/26FPy3QZQqGtDcrja/200.gif",
    ],
    laugh: [
      "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
      "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
      "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
      "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
      "https://media1.giphy.com/media/10JhviFuU2gWD6/200.gif",
      "https://media1.giphy.com/media/ZqlvCTNHpqrio/200.gif",
      "https://media1.giphy.com/media/26BRBKqUiq586bRVm/200.gif",
      "https://media1.giphy.com/media/bC9czlgCMtw4cj8RgH/200.gif",
      "https://media1.giphy.com/media/Q7ozWVYCR0nyW2rvPW/200.gif",
      "https://media1.giphy.com/media/l3q2K5jinAlChoCLS/200.gif",
      "https://media1.giphy.com/media/xUA7aM09ByyR1w5YWc/200.gif",
      "https://media1.giphy.com/media/26tOZ42Mg6r9CUEBq/200.gif",
    ],
    love: [
      "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
      "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
      "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
      "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
      "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
      "https://media1.giphy.com/media/26FLdmIp6wJr91JAI/200.gif",
      "https://media1.giphy.com/media/3oEdv6sy3ulljPMGdy/200.gif",
      "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
      "https://media1.giphy.com/media/M90mJvfWfd5mbUuULX/200.gif",
      "https://media1.giphy.com/media/l0ExncehJzexFpRHq/200.gif",
      "https://media1.giphy.com/media/26tknCqiJrBQG6bxC/200.gif",
      "https://media1.giphy.com/media/3oEjHV0z8S7WM4MwnK/200.gif",
    ],
    cool: [
      "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
      "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
      "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
      "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
      "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
      "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
      "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
      "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
      "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
      "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
      "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
      "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
    ],
    wow: [
      "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
      "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
      "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
      "https://media1.giphy.com/media/5VKbvrjxpVJCM/200.gif",
      "https://media1.giphy.com/media/l0MYEqEzwMWFCg8rm/200.gif",
      "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
      "https://media1.giphy.com/media/26AHvEoTlR5PGWsJq/200.gif",
      "https://media1.giphy.com/media/xT0GqtpF1NWd9BcuI0/200.gif",
      "https://media1.giphy.com/media/3o7abKhOpu0NwenH3O/200.gif",
      "https://media1.giphy.com/media/l3q2K5jinAlChoCLS/200.gif",
      "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
      "https://media1.giphy.com/media/26BRBKqUiq586bRVm/200.gif",
    ],
    thumbsup: [
      "https://media1.giphy.com/media/111ebonMs90YLu/200.gif",
      "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
      "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
      "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
      "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
      "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
      "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
      "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
      "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
      "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
      "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
      "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
    ],
  };

  const gifCategories = Object.keys(gifLibrary);
  const [activeGifCategory, setActiveGifCategory] = useState("trending");

  const searchGifs = (query: string) => {
    setLoadingGifs(true);
    const term = query.trim().toLowerCase();
    if (!term) {
      setGifs(gifLibrary[activeGifCategory] || gifLibrary.trending);
    } else {
      const matchedCategory = gifCategories.find((cat) => cat.includes(term) || term.includes(cat));
      if (matchedCategory) {
        setGifs(gifLibrary[matchedCategory]);
      } else {
        const allGifs = Object.values(gifLibrary).flat();
        const unique = [...new Set(allGifs)];
        setGifs(unique.slice(0, 12));
      }
    }
    setLoadingGifs(false);
  };

  useEffect(() => {
    if (showGifPanel) {
      searchGifs("");
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [showGifPanel, activeGifCategory]);

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    try {
      const res = await fetch("/api/chat/upload", {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        const data = await res.json();
        if (data.url) {
          sendMediaMessage(data.url);
        }
      } else {
        alert("Failed to upload image");
      }
    } catch (err) {
      console.error("Image upload error:", err);
    } finally {
      e.target.value = "";
    }
  };

  const selectedConv = conversations.find((c) => c.id === selectedConvId);

  const formatTime = (ts: number) =>
    new Date(ts * 1000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content" style={{ padding: "24px" }}>
        <div style={{ maxWidth: "960px", margin: "0 auto" }}>
          <div className="page-header">
            <h1 className="page-title">Messages</h1>
            <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
              <button className="btn btn-primary btn-sm" onClick={() => setShowNewChat(!showNewChat)}>
                {showNewChat ? "Cancel" : "+ New Chat"}
              </button>
              <span
                className="dot"
                style={{ background: isConnected ? "var(--success)" : "var(--text-muted)" }}
              />
              <span style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
                {isConnected ? "Connected" : "Offline"}
              </span>
            </div>
          </div>

          {/* New Chat user search */}
          {showNewChat && (
            <div className="card fade-in" style={{ marginBottom: "16px" }}>
              <div style={{ fontWeight: 600, marginBottom: "10px", fontSize: "14px" }}>
                Start a new conversation
              </div>
              <input
                className="input"
                placeholder="Search for a user by name or username..."
                value={userSearch}
                onChange={(e) => handleUserSearch(e.target.value)}
                autoFocus
                style={{ marginBottom: "10px" }}
              />
              {searching && (
                <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>Searching...</div>
              )}
              {searchResults.length > 0 && (
                <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
                  {searchResults.map((u) => {
                    const name =
                      [u.first_name, u.last_name].filter(Boolean).join(" ") || u.username;
                    return (
                      <div
                        key={u.id}
                        style={{
                          display: "flex",
                          justifyContent: "space-between",
                          alignItems: "center",
                          padding: "8px 10px",
                          background: "var(--bg-elevated)",
                          borderRadius: "var(--radius-sm)",
                        }}
                      >
                        <div>
                          <div style={{ fontWeight: 500, fontSize: "14px" }}>{name}</div>
                          <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                            @{u.username}
                          </div>
                        </div>
                        <button
                          className="btn btn-primary btn-sm"
                          onClick={() => startConversation(u.id)}
                        >
                          Message
                        </button>
                      </div>
                    );
                  })}
                </div>
              )}
              {!searching && userSearch.length >= 2 && searchResults.length === 0 && (
                <div style={{ fontSize: "13px", color: "var(--text-muted)" }}>No users found</div>
              )}
            </div>
          )}

          <div className="chat-layout">
            {/* Conversations sidebar */}
            <div className="chat-sidebar">
              <div style={{ overflowY: "auto", flex: 1 }}>
                {loading ? (
                  <div className="loading-spinner" style={{ padding: "20px" }}>
                    <div className="spinner" />
                  </div>
                ) : conversations.length === 0 ? (
                  <div
                    style={{
                      padding: "24px",
                      textAlign: "center",
                      color: "var(--text-muted)",
                      fontSize: "13px",
                    }}
                  >
                    No conversations yet.
                    <br />
                    <button
                      className="btn btn-ghost btn-sm"
                      style={{ marginTop: "10px" }}
                      onClick={() => setShowNewChat(true)}
                    >
                      Start one
                    </button>
                  </div>
                ) : (
                  conversations.map((conv) => (
                    <div
                      key={conv.id}
                      onClick={() => setSelectedConvId(conv.id)}
                      style={{
                        display: "flex",
                        gap: "12px",
                        alignItems: "center",
                        padding: "12px 16px",
                        cursor: "pointer",
                        background:
                          selectedConvId === conv.id ? "var(--bg-hover)" : "transparent",
                        borderLeft:
                          selectedConvId === conv.id
                            ? "3px solid var(--accent)"
                            : "3px solid transparent",
                        transition: "all var(--transition)",
                      }}
                    >
                      <div
                        className="avatar avatar-md"
                        style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)" }}
                      >
                        {conv.other_username?.[0]?.toUpperCase() || "?"}
                      </div>
                      <div style={{ flex: 1, minWidth: 0 }}>
                        <div style={{ fontWeight: 500, fontSize: "14px" }}>
                          {conv.other_username}
                        </div>
                        {conv.last_message && (
                          <div
                            style={{
                              fontSize: "12px",
                              color: "var(--text-muted)",
                              overflow: "hidden",
                              textOverflow: "ellipsis",
                              whiteSpace: "nowrap",
                            }}
                          >
                            {conv.last_message}
                          </div>
                        )}
                      </div>
                      {conv.unread_count ? (
                        <span className="badge">{conv.unread_count}</span>
                      ) : null}
                    </div>
                  ))
                )}
              </div>
            </div>

            {/* Chat area */}
            <div className="chat-area">
              {selectedConv ? (
                <>
                  {/* Header */}
                  <div
                    style={{
                      display: "flex",
                      gap: "12px",
                      alignItems: "center",
                      padding: "16px 20px",
                      borderBottom: "1px solid var(--border)",
                    }}
                  >
                    <div
                      className="avatar avatar-md"
                      style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)" }}
                    >
                      {selectedConv.other_username?.[0]?.toUpperCase()}
                    </div>
                    <div>
                      <div style={{ fontWeight: 600 }}>{selectedConv.other_username}</div>
                      <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                        {isConnected ? "● Online" : "○ Offline"}
                      </div>
                    </div>
                  </div>

                  {/* Messages */}
                  <div className="chat-messages">
                    {loadingMessages ? (
                      <div className="loading-spinner">
                        <div className="spinner" />
                      </div>
                    ) : chatMessages.length === 0 ? (
                      <div
                        style={{
                          textAlign: "center",
                          color: "var(--text-muted)",
                          marginTop: "40px",
                        }}
                      >
                        Say hello to {selectedConv.other_username}! 👋
                      </div>
                    ) : (
                      chatMessages.map((msg, idx) => {
                        const isMe = myId && String(msg.sender_id) === myId;
                        const isMedia = isImageOrGif(msg.content);
                        return (
                          <div
                            key={msg.id || idx}
                            style={{
                              display: "flex",
                              justifyContent: isMe ? "flex-end" : "flex-start",
                            }}
                          >
                            <div>
                              <div
                                className={`message-bubble ${isMe ? "message-out" : "message-in"}`}
                                style={isMedia ? { padding: "8px", background: "transparent", border: "none", boxShadow: "none" } : undefined}
                              >
                                {isMedia ? (
                                  <img
                                    src={msg.content}
                                    alt="Shared media"
                                    style={{
                                      maxWidth: "240px",
                                      maxHeight: "240px",
                                      borderRadius: "12px",
                                      objectFit: "contain",
                                      display: "block",
                                      boxShadow: "var(--shadow-md)"
                                    }}
                                  />
                                ) : (
                                  msg.content
                                )}
                              </div>
                              <div className="message-time">{formatTime(msg.created_at)}</div>
                            </div>
                          </div>
                        );
                      })
                    )}
                    <div ref={messagesEndRef} />
                  </div>

                  {/* GIF Panel */}
                  {showGifPanel && (
                    <div style={{ padding: "10px 20px", borderTop: "1px solid var(--border)", background: "var(--bg-elevated)" }}>
                      <div style={{ display: "flex", gap: "6px", marginBottom: "8px", flexWrap: "wrap" }}>
                        {gifCategories.map((cat) => (
                          <button
                            key={cat}
                            type="button"
                            onClick={() => {
                              setActiveGifCategory(cat);
                              setGifQuery("");
                            }}
                            style={{
                              padding: "3px 8px",
                              fontSize: "11px",
                              borderRadius: "12px",
                              border: "1px solid var(--border)",
                              background: activeGifCategory === cat ? "var(--accent)" : "transparent",
                              color: activeGifCategory === cat ? "#fff" : "var(--text-muted)",
                              cursor: "pointer",
                              textTransform: "capitalize",
                            }}
                          >
                            {cat}
                          </button>
                        ))}
                      </div>
                      <input
                        type="text"
                        placeholder="Search GIFs..."
                        value={gifQuery}
                        onChange={(e) => {
                          setGifQuery(e.target.value);
                          searchGifs(e.target.value);
                        }}
                        className="input"
                        style={{ marginBottom: "8px", width: "100%" }}
                      />
                      {loadingGifs ? (
                        <div style={{ fontSize: "12px", color: "var(--text-muted)", padding: "10px 0" }}>Loading GIFs...</div>
                      ) : gifs.length === 0 ? (
                        <div style={{ fontSize: "12px", color: "var(--text-muted)", padding: "10px 0" }}>No GIFs found</div>
                      ) : (
                        <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: "8px", maxHeight: "150px", overflowY: "auto" }}>
                          {gifs.map((gifUrl, idx) => (
                            <img
                              key={idx}
                              src={gifUrl}
                              alt="gif"
                              style={{ width: "100%", height: "80px", objectFit: "cover", borderRadius: "4px", cursor: "pointer" }}
                              onClick={() => {
                                sendMediaMessage(gifUrl);
                                setShowGifPanel(false);
                                setGifQuery("");
                              }}
                            />
                          ))}
                        </div>
                      )}
                    </div>
                  )}

                  {/* Input controls (Image upload, Emoji, GIF) */}
                  <div style={{ display: "flex", gap: "12px", alignItems: "center", padding: "8px 20px 0 20px", borderTop: "1px solid var(--border)" }}>
                    {/* Image upload */}
                    <label style={{ cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center", padding: "6px", borderRadius: "4px", background: "var(--bg-elevated)", border: "1px solid var(--border)" }} title="Upload Image">
                      <input
                        id="chat-file-input"
                        type="file"
                        accept="image/*"
                        onChange={handleImageUpload}
                        style={{ display: "none" }}
                      />
                      📷
                    </label>

                    {/* Emoji Picker */}
                    <EmojiPicker onSelect={(emoji) => setInputMessage((p) => p + emoji)} />

                    {/* GIF button */}
                    <button
                      type="button"
                      onClick={() => setShowGifPanel(!showGifPanel)}
                      style={{
                        padding: "6px 10px",
                        fontSize: "12px",
                        fontWeight: "bold",
                        borderRadius: "4px",
                        background: showGifPanel ? "var(--accent)" : "var(--bg-elevated)",
                        color: showGifPanel ? "#fff" : "var(--text)",
                        border: "1px solid var(--border)",
                        cursor: "pointer",
                      }}
                    >
                      GIF
                    </button>
                  </div>

                  {/* Input */}
                  <div
                    style={{
                      padding: "16px 20px",
                      display: "flex",
                      gap: "10px",
                    }}
                  >
                    <input
                      type="text"
                      className="input"
                      placeholder={
                        isConnected ? "Type a message..." : "Connecting..."
                      }
                      value={inputMessage}
                      onChange={(e) => setInputMessage(e.target.value)}
                      onKeyDown={(e) => e.key === "Enter" && !e.shiftKey && handleSend()}
                      disabled={!isConnected}
                    />
                    <button
                      onClick={handleSend}
                      className="btn btn-primary"
                      disabled={!inputMessage.trim() || !isConnected}
                    >
                      Send ↗
                    </button>
                  </div>
                </>
              ) : (
                <div
                  className="empty-state"
                  style={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    height: "100%",
                  }}
                >
                  <div>
                    <div
                      style={{
                        display: "flex",
                        justifyContent: "center",
                        marginBottom: "12px",
                      }}
                    >
                      <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        style={{ width: "36px", height: "36px", color: "var(--text-muted)" }}
                      >
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                      </svg>
                    </div>
                    <p>Select a conversation or start a new one</p>
                    <div style={{ textAlign: "center", marginTop: "12px" }}>
                      <button
                        className="btn btn-primary btn-sm"
                        onClick={() => setShowNewChat(true)}
                      >
                        + New Chat
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default function Chat() {
  return (
    <Suspense fallback={<div className="loading-spinner"><div className="spinner" /></div>}>
      <ChatContent />
    </Suspense>
  );
}