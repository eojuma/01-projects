"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";
import { Sidebar } from "../../components/Sidebar";
import { PostCard } from "../../components/PostCard";
import { EmojiPicker } from "../../components/EmojiPicker";
import { useWebSocket } from "../../context/WebSocketContext";
import {
  fetchGroupPosts,
  createGroupPost,
  fetchGroupEvents,
  createGroupEvent,
  respondToEvent,
  requestToJoin,
  approveRequest,
  declineRequest,
  inviteUser
} from "@/src/services/groups";

interface GroupDetails {
  id: string;
  title: string;
  description: string;
  creator_id: string;
  created_at: number;
  is_member: boolean;
  is_creator: boolean;
  privacy?: string;
  is_requested?: boolean;
}

interface GroupPost {
  id: string;
  group_id: string;
  author_id: string;
  author_username: string;
  content: string;
  created_at: number;
  comment_count: number;
}

interface GroupEvent {
  id: string;
  group_id: string;
  creator_id: string;
  title: string;
  description: string;
  event_date: number;
  created_at: number;
  going: number;
  not_going: number;
  user_response: string;
}

interface JoinRequest {
  id: string;
  group_id: string;
  user_id: string;
  status: string;
  created_at: number;
  username: string;
  first_name: string;
  last_name: string;
}

interface GroupMember {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url: string;
  joined_at: number;
}

export default function GroupDetailsPage({ params }: { params: Promise<{ id: string }> }) {
  const router = useRouter();
  const { id } = use(params);

  const [group, setGroup] = useState<GroupDetails | null>(null);
  const [posts, setPosts] = useState<GroupPost[]>([]);
  const [events, setEvents] = useState<GroupEvent[]>([]);
  const [requests, setRequests] = useState<JoinRequest[]>([]);
  const [members, setMembers] = useState<GroupMember[]>([]);

  const [loading, setLoading] = useState(true);
  const [loadingFeed, setLoadingFeed] = useState(false);
  const [submittingPost, setSubmittingPost] = useState(false);
  const [newPostContent, setNewPostContent] = useState("");

  // Event form
  const [showEventForm, setShowEventForm] = useState(false);
  const [eventTitle, setEventTitle] = useState("");
  const [eventDesc, setEventDesc] = useState("");
  const [eventDate, setEventDate] = useState("");
  const [creatingEvent, setCreatingEvent] = useState(false);

  // Invite section
  const [allUsers, setAllUsers] = useState<any[]>([]);
  const [inviteUserId, setInviteUserId] = useState("");
  const [inviting, setInviting] = useState(false);

  // Group editing
  const [showEditForm, setShowEditForm] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editDesc, setEditDesc] = useState("");
  const [saving, setSaving] = useState(false);

  // Tab for right sidebar
  const [rightTab, setRightTab] = useState<"events" | "members" | "chat" | "manage">("events");

  // Group chat states
  const [groupMessages, setGroupMessages] = useState<any[]>([]);
  const [groupInput, setGroupInput] = useState("");
  const [loadingGroupMessages, setLoadingGroupMessages] = useState(false);
  const [myId, setMyId] = useState<string | null>(null);
  const [myUsername, setMyUsername] = useState("");
  const { messages: wsMessages, sendMessage } = useWebSocket();

  const loadGroupDetails = async () => {
    try {
      const res = await fetch(`/api/groups/${id}`, { credentials: "include" });
      if (!res.ok) {
        if (res.status === 404) {
          alert("Group not found");
          router.push("/groups");
          return;
        }
        throw new Error("Failed to load group");
      }
      const data = await res.json();
      setGroup(data);
      setEditTitle(data.title);
      setEditDesc(data.description || "");

      // Load members for all visitors so count is correct
      try {
        const memRes = await fetch(`/api/groups/${id}/members`, { credentials: "include" });
        if (memRes.ok) setMembers(await memRes.json());
      } catch (e) {
        console.error("Failed to load members", e);
      }

      if (data.is_member) {
        setLoadingFeed(true);
        const postData = await fetchGroupPosts(id, 1, 100);
        setPosts(postData || []);

        const eventData = await fetchGroupEvents(id);
        setEvents(eventData || []);

        if (data.is_creator) {
          const reqRes = await fetch(`/api/groups/${id}/requests`, { credentials: "include" });
          if (reqRes.ok) setRequests(await reqRes.json() || []);

          const usersRes = await fetch("/api/users?q=", { credentials: "include" });
          if (usersRes.ok) setAllUsers(await usersRes.json() || []);
        }
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
      setLoadingFeed(false);
    }
  };

  useEffect(() => {
    loadGroupDetails();
  }, [id]);

  useEffect(() => {
    fetch("/api/auth/me", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data && data.id) {
          setMyId(String(data.id));
          setMyUsername(data.username);
        }
      })
      .catch(console.error);
  }, []);

  useEffect(() => {
    if (rightTab === "chat" && group?.id) {
      setLoadingGroupMessages(true);
      fetch(`/api/groups/${group.id}/messages`, { credentials: "include" })
        .then((res) => (res.ok ? res.json() : []))
        .then((data) => {
          setGroupMessages(data);
        })
        .catch(console.error)
        .finally(() => setLoadingGroupMessages(false));
    }
  }, [rightTab, group?.id]);

  useEffect(() => {
    if (wsMessages.length === 0 || !group?.id) return;
    const last = wsMessages[wsMessages.length - 1];
    if (last?.type === "group_chat") {
      const payload = last.payload || {};
      if (payload.group_id === group.id) {
        setGroupMessages((prev) => {
          if (prev.some((m) => m.id === payload.id)) return prev;
          const filtered = prev.filter(
            (m) =>
              !(
                m.id.startsWith("tmp-group-") &&
                m.content === payload.content &&
                m.sender_id === String(payload.from || "")
              )
          );
          return [...filtered, {
            id: payload.id,
            group_id: payload.group_id,
            sender_id: String(payload.from || ""),
            username: payload.username || "Other",
            content: payload.content || "",
            created_at: payload.timestamp || Date.now() / 1000,
          }];
        });
      }
    }
  }, [wsMessages, group?.id]);

  const handleSendGroupMessage = () => {
    if (!groupInput.trim() || !group?.id || !myId) return;
    sendMessage("group_chat", group.id, { content: groupInput });
    
    const tmpId = `tmp-group-${Date.now()}`;
    const newMsg = {
      id: tmpId,
      group_id: group.id,
      sender_id: myId,
      username: myUsername || "Me",
      content: groupInput,
      created_at: Date.now() / 1000,
    };
    setGroupMessages((prev) => [...prev, newMsg]);
    setGroupInput("");
  };

  const handleJoin = async () => {
    if (!group) return;
    try {
      await requestToJoin(group.id);
      await loadGroupDetails();
    } catch (err: any) {
      alert(err.message || "Failed to join group");
    }
  };

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPostContent.trim() || !group) return;
    setSubmittingPost(true);
    try {
      const newPost = await createGroupPost(group.id, newPostContent);
      const enrichedPost: GroupPost = {
        ...newPost,
        author_username: newPost.author_username || "Me"
      };
      setPosts((prev) => [enrichedPost, ...prev]);
      setNewPostContent("");
    } catch (err: any) {
      alert(err.message || "Failed to create post");
    } finally {
      setSubmittingPost(false);
    }
  };

  const handleCreateEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!eventTitle.trim() || !eventDate || !group) return;
    setCreatingEvent(true);
    try {
      const unixTime = new Date(eventDate).getTime() / 1000;
      const newEvent = await createGroupEvent(group.id, eventTitle, eventDesc, unixTime);
      setEvents((prev) => [newEvent, ...prev]);
      setEventTitle("");
      setEventDesc("");
      setEventDate("");
      setShowEventForm(false);
    } catch (err: any) {
      alert(err.message || "Failed to create event");
    } finally {
      setCreatingEvent(false);
    }
  };

  const handleRespondEvent = async (eventId: string, response: "going" | "not_going") => {
    try {
      await respondToEvent(eventId, response);
      // Update local state
      setEvents((prev) =>
        prev.map((ev) =>
          ev.id === eventId ? { ...ev, user_response: response } : ev
        )
      );
    } catch (err: any) {
      alert(err.message || "Failed to respond to event");
    }
  };

  const handleApproveRequest = async (reqId: string) => {
    if (!group) return;
    try {
      await approveRequest(group.id, reqId);
      setRequests((prev) => prev.filter((r) => r.id !== reqId));
      await loadGroupDetails();
    } catch (err: any) {
      alert(err.message || "Failed to approve request");
    }
  };

  const handleDeclineRequest = async (reqId: string) => {
    if (!group) return;
    try {
      await declineRequest(group.id, reqId);
      setRequests((prev) => prev.filter((r) => r.id !== reqId));
    } catch (err: any) {
      alert(err.message || "Failed to decline request");
    }
  };

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inviteUserId || !group) return;
    setInviting(true);
    try {
      await inviteUser(group.id, inviteUserId);
      alert("User invited successfully!");
      setInviteUserId("");
    } catch (err: any) {
      alert(err.message || "Failed to invite user");
    } finally {
      setInviting(false);
    }
  };

  const handleUpdateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editTitle.trim() || !group) return;
    setSaving(true);
    try {
      const res = await fetch(`/api/groups/${group.id}`, {
        method: "PATCH",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: editTitle, description: editDesc }),
      });
      if (!res.ok) {
        const err = await res.text();
        throw new Error(err);
      }
      const updated = await res.json();
      setGroup((prev) => prev ? { ...prev, title: updated.title, description: updated.description } : prev);
      setShowEditForm(false);
      alert("Group updated!");
    } catch (err: any) {
      alert(err.message || "Failed to update group");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="layout">
        <Sidebar />
        <main className="main-content">
          <div className="loading-spinner"><div className="spinner" /></div>
        </main>
      </div>
    );
  }

  if (!group) {
    return (
      <div className="layout">
        <Sidebar />
        <main className="main-content">
          <div className="card" style={{ color: "var(--danger)", textAlign: "center" }}>Group not found</div>
        </main>
      </div>
    );
  }

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div style={{ maxWidth: "960px", margin: "0 auto" }}>
          {/* Header */}
          <div className="card fade-in" style={{ marginBottom: "24px" }}>
            {showEditForm && group.is_creator ? (
              <form onSubmit={handleUpdateGroup} style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                <h3 style={{ fontWeight: 600, marginBottom: "4px" }}>Edit Group</h3>
                <input
                  className="input"
                  value={editTitle}
                  onChange={(e) => setEditTitle(e.target.value)}
                  placeholder="Group name *"
                  required
                />
                <textarea
                  className="input"
                  value={editDesc}
                  onChange={(e) => setEditDesc(e.target.value)}
                  placeholder="Description (optional)"
                  rows={3}
                />
                <div style={{ display: "flex", gap: "8px", justifyContent: "flex-end" }}>
                  <button type="button" className="btn btn-ghost btn-sm" onClick={() => setShowEditForm(false)}>Cancel</button>
                  <button type="submit" className="btn btn-primary btn-sm" disabled={saving}>
                    {saving ? "Saving..." : "Save Changes"}
                  </button>
                </div>
              </form>
            ) : (
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px" }}>
                <div>
                  <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "4px" }}>
                    <h1 className="page-title" style={{ fontSize: "24px" }}>{group.title}</h1>
                    <span className={`tag ${group.privacy === "private" ? "tag-muted" : "tag-blue"}`}>
                      {group.privacy === "private" ? "Private" : "Public"}
                    </span>
                  </div>
                  <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>{group.description || "No description provided."}</p>
                  <p style={{ fontSize: "12px", color: "var(--text-muted)", marginTop: "6px" }}>
                    {members.length} members
                  </p>
                </div>
                <div style={{ display: "flex", gap: "8px", flexShrink: 0 }}>
                  {group.is_creator && (
                    <button className="btn btn-ghost btn-sm" onClick={() => setShowEditForm(true)}>
                      ✎ Edit
                    </button>
                  )}
                  {!group.is_member && (
                    group.is_requested ? (
                      <button className="btn btn-ghost" disabled style={{ opacity: 0.6, cursor: "not-allowed" }}>
                        Request Pending
                      </button>
                    ) : (
                      <button className="btn btn-primary" onClick={handleJoin}>
                        Join Group
                      </button>
                    )
                  )}
                  {group.is_member && !group.is_creator && (
                    <span className="tag tag-green">✓ Member</span>
                  )}
                </div>
              </div>
            )}
          </div>

          {/* Group Content */}
          {!group.is_member ? (
            <div className="card" style={{ textAlign: "center", padding: "40px", color: "var(--text-secondary)" }}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "48px", height: "48px", color: "var(--text-muted)", marginBottom: "16px" }}>
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
              <h3>This group is private</h3>
              <p style={{ fontSize: "13px", marginTop: "8px" }}>Join this group to view posts and events.</p>
            </div>
          ) : (
            <div style={{ display: "grid", gridTemplateColumns: "1fr 300px", gap: "24px" }}>
              {/* Main Column - Feed */}
              <div>
                {/* Create post */}
                <form onSubmit={handleCreatePost} className="card" style={{ marginBottom: "24px" }}>
                  <div style={{ display: "flex", gap: "12px" }}>
                    <div className="avatar avatar-md" style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)" }}>Me</div>
                    <div style={{ flex: 1 }}>
                      <textarea
                        className="input"
                        placeholder={`Post to ${group.title}...`}
                        value={newPostContent}
                        onChange={(e) => setNewPostContent(e.target.value)}
                        rows={3}
                        style={{ marginBottom: "12px" }}
                        required
                      />
                      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <EmojiPicker onSelect={(emoji) => setNewPostContent((p) => p + emoji)} />
                        <button type="submit" className="btn btn-primary" disabled={submittingPost || !newPostContent.trim()}>
                          {submittingPost ? "Posting..." : "Post"}
                        </button>
                      </div>
                    </div>
                  </div>
                </form>

                {/* Posts List */}
                {loadingFeed ? (
                  <div className="loading-spinner"><div className="spinner" /></div>
                ) : posts.length === 0 ? (
                  <div className="card" style={{ textAlign: "center", padding: "40px", color: "var(--text-muted)" }}>
                    No posts yet. Be the first to post!
                  </div>
                ) : (
                  posts.map((post) => (
                    <PostCard
                      key={post.id}
                      id={post.id}
                      username={post.author_username || "Unknown"}
                      userInitial={(post.author_username?.[0] || "U").toUpperCase()}
                      time={new Date(post.created_at * 1000).toLocaleString()}
                      privacy="group"
                      content={post.content}
                      comments={post.comment_count}
                    />
                  ))
                )}
              </div>

              {/* Sidebar */}
              <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                {/* Tab bar */}
                <div style={{ display: "flex", gap: "4px", background: "var(--bg-elevated)", borderRadius: "var(--radius-sm)", padding: "4px" }}>
                  {(["events", "members", "chat", ...(group.is_creator ? ["manage"] : [])] as typeof rightTab[]).map((tab) => (
                    <button
                      key={tab}
                      onClick={() => setRightTab(tab as typeof rightTab)}
                      style={{
                        flex: 1,
                        padding: "6px",
                        fontSize: "12px",
                        fontWeight: rightTab === tab ? 600 : 400,
                        background: rightTab === tab ? "var(--bg-surface)" : "transparent",
                        border: "none",
                        borderRadius: "var(--radius-xs)",
                        cursor: "pointer",
                        color: rightTab === tab ? "var(--text-primary)" : "var(--text-muted)",
                        textTransform: "capitalize",
                      }}
                    >
                      {tab}
                    </button>
                  ))}
                </div>

                {/* Events Tab */}
                {rightTab === "events" && (
                  <div className="card">
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "14px" }}>
                      <h3 style={{ fontWeight: 600, fontSize: "15px" }}>Events</h3>
                      <button className="btn btn-ghost btn-sm" onClick={() => setShowEventForm(!showEventForm)}>
                        {showEventForm ? "Cancel" : "+ Event"}
                      </button>
                    </div>

                    {showEventForm && (
                      <form onSubmit={handleCreateEvent} style={{ borderBottom: "1px solid var(--border)", paddingBottom: "14px", marginBottom: "14px", display: "flex", flexDirection: "column", gap: "10px" }}>
                        <input className="input" placeholder="Event title *" value={eventTitle} onChange={(e) => setEventTitle(e.target.value)} required />
                        <textarea className="input" placeholder="Event details" value={eventDesc} onChange={(e) => setEventDesc(e.target.value)} rows={2} />
                        <input className="input" type="datetime-local" value={eventDate} onChange={(e) => setEventDate(e.target.value)} required />
                        <button type="submit" className="btn btn-primary btn-sm" disabled={creatingEvent}>
                          {creatingEvent ? "Creating..." : "Create"}
                        </button>
                      </form>
                    )}

                    {events.length === 0 ? (
                      <p style={{ fontSize: "13px", color: "var(--text-muted)", textAlign: "center" }}>No upcoming events.</p>
                    ) : (
                      <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                        {events.map((ev) => (
                          <div key={ev.id} style={{ border: "1px solid var(--border)", borderRadius: "var(--radius-sm)", padding: "10px", background: "var(--bg-elevated)" }}>
                            <div style={{ fontWeight: 600, fontSize: "13px", marginBottom: "2px" }}>{ev.title}</div>
                            <p style={{ fontSize: "11px", color: "var(--text-secondary)", marginBottom: "4px" }}>
                              {new Date(ev.event_date * 1000).toLocaleString()}
                            </p>
                            {ev.description && <p style={{ fontSize: "12px", color: "var(--text-secondary)", marginBottom: "8px" }}>{ev.description}</p>}
                            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", fontSize: "11px", color: "var(--text-muted)", marginBottom: "8px" }}>
                              <span>Going: {ev.going} · Not: {ev.not_going}</span>
                              {ev.user_response && <span className="tag tag-green" style={{ fontSize: "9px" }}>{ev.user_response}</span>}
                            </div>
                            <div style={{ display: "flex", gap: "6px" }}>
                              <button
                                className={`btn btn-sm ${ev.user_response === "going" ? "btn-primary" : "btn-ghost"}`}
                                style={{ flex: 1, padding: "3px" }}
                                onClick={() => handleRespondEvent(ev.id, "going")}
                              >Going</button>
                              <button
                                className={`btn btn-sm ${ev.user_response === "not_going" ? "btn-primary" : "btn-ghost"}`}
                                style={{ flex: 1, padding: "3px" }}
                                onClick={() => handleRespondEvent(ev.id, "not_going")}
                              >Not Going</button>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}

                {/* Members Tab */}
                {rightTab === "members" && (
                  <div className="card">
                    <h3 style={{ fontWeight: 600, fontSize: "15px", marginBottom: "12px" }}>Members ({members.length})</h3>
                    {members.length === 0 ? (
                      <p style={{ fontSize: "13px", color: "var(--text-muted)", textAlign: "center" }}>No members found.</p>
                    ) : (
                      <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                        {members.map((m) => {
                          const name = [m.first_name, m.last_name].filter(Boolean).join(" ") || m.username;
                          return (
                            <div key={m.id} style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                              <div className="avatar avatar-sm" style={{ background: "linear-gradient(135deg, #a29bfe, #fd79a8)", fontSize: "11px" }}>
                                {name[0]?.toUpperCase() || "?"}
                              </div>
                              <div>
                                <div style={{ fontWeight: 500, fontSize: "13px" }}>{name}</div>
                                <div style={{ fontSize: "11px", color: "var(--text-muted)" }}>@{m.username}</div>
                              </div>
                              {m.id === group.creator_id && (
                                <span className="tag tag-blue" style={{ marginLeft: "auto", fontSize: "9px" }}>Creator</span>
                              )}
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </div>
                )}

                {/* Chat Tab */}
                {rightTab === "chat" && (
                  <div className="card" style={{ display: "flex", flexDirection: "column", height: "450px" }}>
                    <h3 style={{ fontWeight: 600, fontSize: "15px", marginBottom: "14px" }}>Group Chat Room</h3>

                    {/* Messages list */}
                    <div style={{ flex: 1, overflowY: "auto", display: "flex", flexDirection: "column", gap: "8px", marginBottom: "10px", paddingRight: "4px" }}>
                      {loadingGroupMessages ? (
                        <div style={{ textAlign: "center", color: "var(--text-muted)", fontSize: "12px" }}>Loading messages...</div>
                      ) : groupMessages.length === 0 ? (
                        <div style={{ textAlign: "center", color: "var(--text-muted)", fontSize: "12px", marginTop: "20px" }}>No messages yet. Say hi! 👋</div>
                      ) : (
                        groupMessages.map((msg, idx) => {
                          const isMe = myId && String(msg.sender_id) === myId;
                          return (
                            <div key={msg.id || idx} style={{ display: "flex", flexDirection: "column", alignItems: isMe ? "flex-end" : "flex-start" }}>
                              {!isMe && <span style={{ fontSize: "10px", color: "var(--text-muted)", marginLeft: "4px", marginBottom: "2px" }}>@{msg.username}</span>}
                              <div
                                style={{
                                  background: isMe ? "var(--accent)" : "var(--bg-elevated)",
                                  color: isMe ? "#fff" : "var(--text-primary)",
                                  padding: "8px 12px",
                                  borderRadius: "12px",
                                  maxWidth: "85%",
                                  fontSize: "13px",
                                  wordBreak: "break-word",
                                  textAlign: "left",
                                }}
                              >
                                {msg.content}
                              </div>
                            </div>
                          );
                        })
                      )}
                    </div>

                    {/* Input box */}
                    <div style={{ display: "flex", gap: "6px" }}>
                      <input
                        type="text"
                        className="input"
                        placeholder="Type a message..."
                        value={groupInput}
                        onChange={(e) => setGroupInput(e.target.value)}
                        onKeyDown={(e) => e.key === "Enter" && handleSendGroupMessage()}
                        style={{ fontSize: "12px", padding: "6px 10px", flex: 1 }}
                      />
                      <button className="btn btn-primary btn-sm" onClick={handleSendGroupMessage} disabled={!groupInput.trim()}>
                        Send
                      </button>
                    </div>
                  </div>
                )}

                {/* Manage Tab (creator only) */}
                {rightTab === "manage" && group.is_creator && (
                  <>
                    {/* Invite users */}
                    <div className="card">
                      <h3 style={{ fontWeight: 600, fontSize: "15px", marginBottom: "12px" }}>Invite Members</h3>
                      <form onSubmit={handleInvite} style={{ display: "flex", gap: "8px" }}>
                        <select
                          className="input"
                          value={inviteUserId}
                          onChange={(e) => setInviteUserId(e.target.value)}
                          required
                          style={{ flex: 1 }}
                        >
                          <option value="">Select user...</option>
                          {allUsers.filter(u => !members.some(m => m.id === u.id)).map((u) => (
                            <option key={u.id} value={u.id}>{u.first_name || u.username} (@{u.username})</option>
                          ))}
                        </select>
                        <button type="submit" className="btn btn-primary btn-sm" disabled={inviting || !inviteUserId}>
                          Invite
                        </button>
                      </form>
                    </div>

                    {/* Pending join requests */}
                    <div className="card">
                      <h3 style={{ fontWeight: 600, fontSize: "15px", marginBottom: "12px" }}>
                        Join Requests {requests.length > 0 && <span className="badge">{requests.length}</span>}
                      </h3>
                      {requests.length === 0 ? (
                        <p style={{ fontSize: "13px", color: "var(--text-muted)", textAlign: "center" }}>No pending requests.</p>
                      ) : (
                        <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                          {requests.map((r) => {
                            const name = [r.first_name, r.last_name].filter(Boolean).join(" ") || r.username;
                            return (
                              <div key={r.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: "8px" }}>
                                <div style={{ minWidth: 0 }}>
                                  <div style={{ fontWeight: 500, fontSize: "13px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{name}</div>
                                  <div style={{ fontSize: "11px", color: "var(--text-muted)" }}>@{r.username}</div>
                                </div>
                                <div style={{ display: "flex", gap: "4px", flexShrink: 0 }}>
                                  <button className="btn btn-primary btn-sm" style={{ padding: "3px 8px" }} onClick={() => handleApproveRequest(r.id)}>✓</button>
                                  <button className="btn btn-ghost btn-sm" style={{ padding: "3px 8px" }} onClick={() => handleDeclineRequest(r.id)}>✗</button>
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </div>
                  </>
                )}
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
