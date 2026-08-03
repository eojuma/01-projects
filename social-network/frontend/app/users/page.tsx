"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Sidebar } from "../components/Sidebar";
import { fetchGroups } from "@/src/services/groups";

interface SearchUser {
  id: string;
  username: string;
  first_name?: string;
  last_name?: string;
  nickname?: string;
  bio?: string;
  avatar_url?: string;
  is_private: boolean;
}

interface Group {
  id: string;
  title: string;
  is_member: boolean;
}

export default function FindFriends() {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [users, setUsers] = useState<SearchUser[]>([]);
  const [loading, setLoading] = useState(false);
  const [myGroups, setMyGroups] = useState<Group[]>([]);
  const [activeInviteUserId, setActiveInviteUserId] = useState<string | null>(null);

  // Load my groups
  useEffect(() => {
    fetchGroups(1, 100)
      .then((data) => {
        // Filter groups where current user is a member
        const joined = data.filter((g: any) => g.is_member);
        setMyGroups(joined);
      })
      .catch(console.error);
  }, []);

  const handleSearch = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setLoading(true);
    try {
      const res = await fetch(`/api/users?q=${encodeURIComponent(query)}`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setUsers(data);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Run initial search
  useEffect(() => {
    handleSearch();
  }, []);

  const handleMessage = async (userId: string) => {
    try {
      const res = await fetch(`/api/chat/user/${userId}`, { credentials: "include" });
      if (res.ok) {
        const data = await res.json();
        router.push(`/chat?convId=${data.conversation_id}`);
      } else {
        alert("Failed to start conversation");
      }
    } catch {
      alert("Failed to connect to server");
    }
  };

  const handleInvite = async (userId: string, groupId: string) => {
    try {
      const res = await fetch(`/api/groups/${groupId}/invite`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ user_id: userId }),
      });
      if (res.ok) {
        alert("Invite sent successfully!");
        setActiveInviteUserId(null);
      } else {
        const txt = await res.text();
        alert(txt || "Failed to send invite");
      }
    } catch {
      alert("Failed to connect to server");
    }
  };

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div className="page-container">
          <div className="page-header">
            <h1 className="page-title">Find Friends</h1>
          </div>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="card" style={{ marginBottom: "20px", display: "flex", gap: "10px" }}>
            <input
              className="input"
              type="text"
              placeholder="Search by name, username, or nickname..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              style={{ flex: 1 }}
            />
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? "Searching..." : "Search"}
            </button>
          </form>

          {/* Users List */}
          {loading ? (
            <div className="loading-spinner"><div className="spinner" /></div>
          ) : users.length === 0 ? (
            <div className="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "12px" }}>
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
              </svg>
              <p>No users found matching your search.</p>
            </div>
          ) : (
            <div style={{ display: "grid", gap: "16px" }}>
              {users.map((user) => {
                const displayName = [user.first_name, user.last_name].filter(Boolean).join(" ") || user.username || "User";
                const isInviteOpen = activeInviteUserId === user.id;

                return (
                  <div key={user.id} className="card fade-in" style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: "16px", flexWrap: "wrap" }}>
                    <div style={{ display: "flex", gap: "12px", alignItems: "center" }}>
                      <div className="avatar avatar-md" style={{ background: "linear-gradient(135deg, #a29bfe, #fd79a8)" }}>
                        {displayName && displayName.length > 0 ? displayName[0].toUpperCase() : "?"}
                      </div>
                      <div>
                        <div style={{ fontWeight: 600, fontSize: "15px" }}>{displayName}</div>
                        <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>@{user.username}</div>
                        {user.nickname && (
                          <div style={{ fontSize: "12px", color: "var(--text-secondary)", fontStyle: "italic" }}>
                            &quot;{user.nickname}&quot;
                          </div>
                        )}
                      </div>
                    </div>

                    <div style={{ display: "flex", gap: "8px", alignItems: "center", position: "relative" }}>
                      <button className="btn btn-ghost btn-sm" onClick={() => router.push(`/users/${user.id}`)}>
                        View Profile
                      </button>
                      <button className="btn btn-ghost btn-sm" onClick={() => handleMessage(user.id)}>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "14px", height: "14px", marginRight: "4px" }}>
                          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                        </svg>
                        Message
                      </button>

                      {/* Group invite button */}
                      {myGroups.length > 0 && (
                        <div style={{ position: "relative" }}>
                          <button
                            className="btn btn-ghost btn-sm"
                            onClick={() => setActiveInviteUserId(isInviteOpen ? null : user.id)}
                          >
                            + Invite to Group
                          </button>

                          {isInviteOpen && (
                            <div className="dropdown-menu" style={{ position: "absolute", top: "100%", right: 0, background: "var(--bg-surface)", border: "1px solid var(--border)", borderRadius: "var(--radius-sm)", boxShadow: "0 4px 12px rgba(0,0,0,0.15)", zIndex: 100, minWidth: "180px", marginTop: "4px", padding: "6px" }}>
                              <div style={{ fontSize: "11px", color: "var(--text-muted)", padding: "4px 8px", borderBottom: "1px solid var(--border)", marginBottom: "4px" }}>
                                Select Group
                              </div>
                              {myGroups.map((group) => (
                                <button
                                  key={group.id}
                                  onClick={() => handleInvite(user.id, group.id)}
                                  className="dropdown-item"
                                  style={{ width: "100%", textAlign: "left", padding: "6px 8px", fontSize: "12px", background: "none", border: "none", cursor: "pointer", color: "var(--text-primary)", borderRadius: "var(--radius-xs)" }}
                                  onMouseEnter={(e) => (e.currentTarget.style.background = "var(--bg-hover)")}
                                  onMouseLeave={(e) => (e.currentTarget.style.background = "none")}
                                >
                                  {group.title}
                                </button>
                              ))}
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
