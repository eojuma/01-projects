"use client";
 
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Sidebar } from "../components/Sidebar";
import { useGroups } from "../hooks/useGroups";
import { createGroup, requestToJoin, inviteUser } from "@/src/services/groups";
 
export default function Groups() {
  const router = useRouter();
  const { groups, loading, error, acceptGroupInvite, declineGroupInvite, refresh } = useGroups();
  const [showCreate, setShowCreate] = useState(false);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({ title: "", description: "", privacy: "public" });
  const [joinLoading, setJoinLoading] = useState<string | null>(null);

  const [allUsers, setAllUsers] = useState<any[]>([]);
  const [selectedUserIds, setSelectedUserIds] = useState<string[]>([]);
  const [loadingUsers, setLoadingUsers] = useState(false);

  useEffect(() => {
    if (showCreate) {
      setLoadingUsers(true);
      fetch(`/api/users?q=`, { credentials: "include" })
        .then((res) => (res.ok ? res.json() : []))
        .then((data) => setAllUsers(data))
        .catch(console.error)
        .finally(() => setLoadingUsers(false));
    }
  }, [showCreate]);
 
  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title.trim()) return;
    setCreating(true);
    try {
      const newGroup = await createGroup(form.title, form.description);
      // Invite selected users if any
      if (selectedUserIds.length > 0) {
        await Promise.all(
          selectedUserIds.map((userId) =>
            inviteUser(newGroup.id, userId).catch((err) =>
              console.error(`Failed to invite user ${userId}:`, err)
            )
          )
        );
      }
      setForm({ title: "", description: "", privacy: "public" });
      setSelectedUserIds([]);
      setShowCreate(false);
      await refresh();
    } catch (err: any) {
      alert(err.message);
    } finally {
      setCreating(false);
    }
  };
 
  const handleJoin = async (groupId: string) => {
    setJoinLoading(groupId);
    try {
      await requestToJoin(groupId);
      await refresh();
    } catch (err: any) {
      alert(err.message || "Failed to join group");
    } finally {
      setJoinLoading(null);
    }
  };
 
  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div className="page-container">
          <div className="page-header">
            <h1 className="page-title">Groups</h1>
            <button
              className="btn btn-primary"
              onClick={() => setShowCreate(!showCreate)}
            >
              {showCreate ? "Cancel" : "+ Create Group"}
            </button>
          </div>
 
          {/* Create Group Form */}
          {showCreate && (
            <form onSubmit={handleCreate} className="card fade-in" style={{ marginBottom: "24px" }}>
              <h3 style={{ fontWeight: 600, marginBottom: "16px" }}>New Group</h3>
              <input
                className="input"
                placeholder="Group name *"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                style={{ marginBottom: "12px" }}
                required
              />
              <textarea
                className="input"
                placeholder="Description (optional)"
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                rows={3}
                style={{ marginBottom: "12px" }}
              />

              {/* Select Members to Invite */}
              <div style={{ marginBottom: "16px" }}>
                <label className="form-label" style={{ marginBottom: "8px", display: "block" }}>Invite Members</label>
                {loadingUsers ? (
                  <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>Loading users...</div>
                ) : allUsers.length === 0 ? (
                  <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>No other users to invite.</div>
                ) : (
                  <div style={{ maxHeight: "150px", overflowY: "auto", border: "1px solid var(--border)", borderRadius: "var(--radius-sm)", padding: "8px", background: "var(--bg-elevated)" }}>
                    {allUsers.map((u) => {
                      const isSelected = selectedUserIds.includes(u.id);
                      const displayName = [u.first_name, u.last_name].filter(Boolean).join(" ") || u.username || "User";
                      return (
                        <label key={u.id} style={{ display: "flex", alignItems: "center", gap: "10px", padding: "6px 8px", cursor: "pointer", borderRadius: "var(--radius-xs)", transition: "background 0.15s" }}
                          onMouseEnter={(e) => (e.currentTarget.style.background = "var(--bg-hover)")}
                          onMouseLeave={(e) => (e.currentTarget.style.background = "none")}
                        >
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => {
                              setSelectedUserIds((prev) =>
                                isSelected ? prev.filter((id) => id !== u.id) : [...prev, u.id]
                              );
                            }}
                            style={{ cursor: "pointer" }}
                          />
                          <span style={{ fontSize: "13px" }}>{displayName} (@{u.username})</span>
                        </label>
                      );
                    })}
                  </div>
                )}
              </div>

              <div style={{ display: "flex", gap: "10px", justifyContent: "flex-end" }}>
                <select
                  className="input"
                  style={{ width: "140px" }}
                  value={form.privacy}
                  onChange={(e) => setForm({ ...form, privacy: e.target.value })}
                >
                  <option value="public">Public</option>
                  <option value="private">Private</option>
                </select>
                <button type="submit" className="btn btn-primary" disabled={creating}>
                  {creating ? "Creating..." : "Create"}
                </button>
              </div>
            </form>
          )}

          {/* Groups List */}
          {loading ? (
            <div className="loading-spinner"><div className="spinner" /></div>
          ) : error ? (
            <div className="card" style={{ color: "var(--danger)", textAlign: "center" }}>{error}</div>
          ) : groups.length === 0 ? (
            <div className="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "12px" }}>
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                <path d="M16 3.13a4 4 0 0 1 0 7.75" />
              </svg>
              <p>No groups yet. Create one to get started!</p>
            </div>
          ) : (
            groups.map((group) => (
              <div key={group.id} className="card fade-in">
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px" }}>
                  <div style={{ flex: 1 }}>
                    <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "4px" }}>
                      <h3
                        style={{ fontWeight: 600, fontSize: "16px", cursor: "pointer" }}
                        onClick={() => router.push(`/groups/${group.id}`)}
                        onMouseEnter={(e) => (e.currentTarget.style.color = "var(--accent)")}
                        onMouseLeave={(e) => (e.currentTarget.style.color = "inherit")}
                      >
                        {group.title || (group as any).name}
                      </h3>
                      <span className={`tag ${group.privacy === "private" ? "tag-muted" : "tag-blue"}`}>
                        {group.privacy === "private" ? "Private" : "Public"}
                      </span>
                    </div>
                    {group.member_count !== undefined && (
                      <p style={{ fontSize: "12px", color: "var(--text-muted)", marginBottom: "8px" }}>
                        {group.member_count.toLocaleString()} members
                      </p>
                    )}
                    {group.description && (
                      <p style={{ fontSize: "14px", color: "var(--text-secondary)", lineHeight: "1.5" }}>
                        {group.description}
                      </p>
                    )}
                  </div>
                  <div style={{ display: "flex", gap: "8px", flexShrink: 0 }}>
                    {group.is_member ? (
                      <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
                        <span className="tag tag-green">✓ Member</span>
                        <button
                          className="btn btn-primary btn-sm"
                          onClick={() => router.push(`/groups/${group.id}`)}
                        >
                          View Group
                        </button>
                      </div>
                    ) : group.is_invited ? (
                      <div style={{ display: "flex", gap: "6px" }}>
                        <button
                          className="btn btn-primary btn-sm"
                          onClick={() => group.invite_id && acceptGroupInvite(group.id, group.invite_id)}
                        >
                          Accept Invite
                        </button>
                        <button
                          className="btn btn-ghost btn-sm"
                          onClick={() => group.invite_id && declineGroupInvite(group.id, group.invite_id)}
                        >
                          Decline
                        </button>
                      </div>
                    ) : group.is_requested ? (
                      <button className="btn btn-ghost btn-sm" disabled style={{ opacity: 0.6, cursor: "not-allowed" }}>
                        Requested
                      </button>
                    ) : (
                      <button
                        className="btn btn-ghost btn-sm"
                        disabled={joinLoading === group.id}
                        onClick={() => handleJoin(group.id)}
                      >
                        {joinLoading === group.id ? "..." : "Join"}
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}