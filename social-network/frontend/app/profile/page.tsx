"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Sidebar } from "../components/Sidebar";
import { PostCard } from "../components/PostCard";
import { getUserProfile, updateProfile, togglePrivacy, getUserPosts } from "@/src/services/profile";
import { getFollowers, getFollowing } from "@/src/services/follow";

const PRESELECTED_AVATARS = [
  "https://api.dicebear.com/7.x/bottts/svg?seed=Felix",
  "https://api.dicebear.com/7.x/bottts/svg?seed=Aneka",
  "https://api.dicebear.com/7.x/bottts/svg?seed=Jack",
  "https://api.dicebear.com/7.x/bottts/svg?seed=Boots",
  "https://api.dicebear.com/7.x/bottts/svg?seed=Coco",
  "https://api.dicebear.com/7.x/bottts/svg?seed=Zoe"
];

export default function Profile() {
  const router = useRouter();
  const [profile, setProfile] = useState<any>(null);
  const [posts, setPosts] = useState<any[]>([]);
  const [followersCount, setFollowersCount] = useState(0);
  const [followingCount, setFollowingCount] = useState(0);
  const [followers, setFollowers] = useState<any[]>([]);
  const [following, setFollowing] = useState<any[]>([]);
  const [profileTab, setProfileTab] = useState<"posts" | "followers" | "following">("posts");
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    bio: "",
    nickname: "",
    avatar_url: "",
  });

  useEffect(() => {
    fetchCurrentUser();
  }, []);

  const fetchCurrentUser = async () => {
    try {
      const res = await fetch("/api/auth/me", { credentials: "include" });
      if (!res.ok) {
        if (res.status === 401) { router.push("/login"); return; }
        throw new Error("Failed to get current user");
      }
      const user = await res.json();
      await loadProfile(user.username, user.id);
    } catch (err) {
      setError("Please log in to view your profile");
      setLoading(false);
    }
  };

  const loadProfile = async (username: string, userId: string) => {
    try {
      const data = await getUserProfile(username);
      setProfile({ ...data, id: userId });
      setFormData({
        first_name: data.first_name || "",
        last_name: data.last_name || "",
        bio: data.bio || "",
        nickname: data.nickname || "",
        avatar_url: data.avatar_url || "",
      });
      const userPosts = await getUserPosts(userId);
      setPosts(userPosts);
      // Load follower/following counts and lists
      try {
        const fwrs = await getFollowers(userId);
        setFollowers(fwrs || []);
        setFollowersCount(Array.isArray(fwrs) ? fwrs.length : 0);
      } catch { /* ignore */ }
      try {
        const fwing = await getFollowing(userId);
        setFollowing(fwing || []);
        setFollowingCount(Array.isArray(fwing) ? fwing.length : 0);
      } catch { /* ignore */ }
    } catch (err) {
      setError("Could not load profile");
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const uploadData = new FormData();
    uploadData.append("image", file);

    try {
      const res = await fetch("/api/chat/upload", {
        method: "POST",
        body: uploadData,
      });
      if (res.ok) {
        const data = await res.json();
        setFormData(prev => ({ ...prev, avatar_url: data.url }));
      } else {
        alert("Failed to upload avatar");
      }
    } catch (err) {
      console.error("Avatar upload failed:", err);
      alert("Error uploading avatar");
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!profile?.id) return;
    setSaving(true);
    try {
      const updated = await updateProfile(profile.id, formData);
      setProfile((prev: any) => ({ ...prev, ...updated }));
      setEditing(false);
    } catch (err: any) {
      alert(err.message || "Failed to save profile");
    } finally {
      setSaving(false);
    }
  };

  const handleTogglePrivacy = async () => {
    if (!profile?.id) return;
    try {
      await togglePrivacy(profile.id, !profile.is_private);
      setProfile((prev: any) => ({ ...prev, is_private: !prev.is_private }));
    } catch (err: any) {
      alert(err.message || "Failed to update privacy");
    }
  };

  const displayName = [profile?.first_name, profile?.last_name].filter(Boolean).join(" ") || profile?.username;

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div className="page-container">
          {loading ? (
            <div className="loading-spinner"><div className="spinner" /></div>
          ) : error ? (
            <div className="card" style={{ textAlign: "center" }}>
              <p style={{ color: "var(--danger)", marginBottom: "16px" }}>{error}</p>
              <button className="btn btn-primary" onClick={() => router.push("/login")}>Go to Login</button>
            </div>
          ) : !profile ? null : (
            <>
              {/* Profile Header */}
              <div className="card" style={{ marginBottom: "20px" }}>
                <div style={{ display: "flex", gap: "20px", alignItems: "flex-start" }}>
                  {/* Avatar */}
                  {profile.avatar_url ? (
                    <img
                      src={profile.avatar_url}
                      alt={displayName}
                      className="avatar avatar-xl"
                      style={{ objectFit: "cover" }}
                    />
                  ) : (
                    <div
                      className="avatar avatar-xl"
                      style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)" }}
                    >
                      {displayName[0]?.toUpperCase() || "?"}
                    </div>
                  )}

                  {/* Info */}
                  <div style={{ flex: 1 }}>
                    <div style={{ display: "flex", alignItems: "center", gap: "10px", flexWrap: "wrap", marginBottom: "4px" }}>
                      <h1 style={{ fontSize: "22px", fontWeight: 700 }}>{displayName}</h1>
                      <span className={`tag ${profile.is_private ? "tag-muted" : "tag-blue"}`}>
                        {profile.is_private ? "Private" : "Public"}
                      </span>
                    </div>
                    <p style={{ color: "var(--text-muted)", fontSize: "14px" }}>@{profile.username}</p>
                    {profile.email && (
                      <p style={{ color: "var(--text-secondary)", fontSize: "14px", marginTop: "2px" }}>
                        📧 {profile.email}
                      </p>
                    )}
                    {profile.nickname && (
                      <p style={{ color: "var(--text-secondary)", fontSize: "13px", marginTop: "2px" }}>
                        "{profile.nickname}"
                      </p>
                    )}
                    {profile.bio && (
                      <p style={{ color: "var(--text-secondary)", fontSize: "14px", marginTop: "10px", lineHeight: "1.6" }}>
                        {profile.bio}
                      </p>
                    )}
                    <p style={{ color: "var(--text-muted)", fontSize: "12px", marginTop: "10px" }}>
                      Joined {new Date(profile.created_at).toLocaleDateString("en-US", { year: "numeric", month: "long" })}
                    </p>
                    {/* Stats */}
                   <div style={{ display: "flex", gap: "4px", marginTop: "14px" }}>
                    <div className="stat-chip" onClick={() => setProfileTab("followers")} style={{ cursor: "pointer" }}>
                      <span className="stat-chip-value">{followersCount}</span>
                       <span className="stat-chip-label">Followers</span>
                       </div>
                       <div className="stat-chip" onClick={() => setProfileTab("following")} style={{ cursor: "pointer" }}>
                        <span className="stat-chip-value">{followingCount}</span>
                        <span className="stat-chip-label">Following</span>
                        </div>
                        <div className="stat-chip" onClick={() => setProfileTab("posts")} style={{ cursor: "pointer" }}>
                          <span className="stat-chip-value">{posts.length}</span>
                          <span className="stat-chip-label">Posts</span>
                          </div>
                          </div>
                          </div>

                  {/* Actions */}
                  <button
                    className="btn btn-ghost btn-sm"
                    onClick={() => setEditing(!editing)}
                  >
                    {editing ? "Cancel" : "Edit Profile"}
                  </button>
                </div>
              </div>

              {/* Edit Form */}
              {editing && (
                <div className="card fade-in" style={{ marginBottom: "20px" }}>
                  <h3 style={{ fontWeight: 600, marginBottom: "16px" }}>Edit Profile</h3>
                  <form onSubmit={handleUpdate}>
                    <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "12px" }}>
                      <div>
                        <label style={{ display: "block", fontSize: "12px", color: "var(--text-muted)", marginBottom: "6px" }}>
                          FIRST NAME
                        </label>
                        <input
                          className="input"
                          value={formData.first_name}
                          onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                          placeholder="First name"
                        />
                      </div>
                      <div>
                        <label style={{ display: "block", fontSize: "12px", color: "var(--text-muted)", marginBottom: "6px" }}>
                          LAST NAME
                        </label>
                        <input
                          className="input"
                          value={formData.last_name}
                          onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                          placeholder="Last name"
                        />
                      </div>
                    </div>
                    <div style={{ marginBottom: "12px" }}>
                      <label style={{ display: "block", fontSize: "12px", color: "var(--text-muted)", marginBottom: "6px" }}>
                        NICKNAME
                      </label>
                      <input
                        className="input"
                        value={formData.nickname}
                        onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
                        placeholder="Nickname"
                      />
                    </div>
                    <div style={{ marginBottom: "16px" }}>
                      <label style={{ display: "block", fontSize: "12px", color: "var(--text-muted)", marginBottom: "6px" }}>
                        BIO
                      </label>
                      <textarea
                        className="input"
                        rows={3}
                        value={formData.bio}
                        onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                        placeholder="Tell people about yourself..."
                      />
                    </div>
                    {/* Avatar Selection */}
                    <div style={{ marginBottom: "16px" }}>
                      <label style={{ display: "block", fontSize: "12px", color: "var(--text-muted)", marginBottom: "8px" }}>
                        AVATAR
                      </label>
                      {formData.avatar_url && (
                        <div style={{ display: "flex", alignItems: "center", gap: "12px", marginBottom: "12px" }}>
                          <img
                            src={formData.avatar_url}
                            alt="Current avatar"
                            style={{ width: "56px", height: "56px", borderRadius: "50%", objectFit: "cover", border: "2px solid var(--accent)" }}
                          />
                          <button
                            type="button"
                            className="btn btn-ghost btn-sm"
                            onClick={() => setFormData({ ...formData, avatar_url: "" })}
                            style={{ color: "var(--danger)", fontSize: "12px" }}
                          >
                            Remove Avatar
                          </button>
                        </div>
                      )}
                      <div style={{ display: "grid", gridTemplateColumns: "repeat(6, 1fr)", gap: "8px", marginBottom: "10px" }}>
                        {PRESELECTED_AVATARS.map((url) => (
                          <img
                            key={url}
                            src={url}
                            alt="Avatar option"
                            onClick={() => setFormData({ ...formData, avatar_url: url })}
                            style={{
                              width: "100%",
                              aspectRatio: "1",
                              borderRadius: "50%",
                              cursor: "pointer",
                              border: formData.avatar_url === url ? "3px solid var(--accent)" : "2px solid var(--border)",
                              objectFit: "cover",
                              transition: "border-color 0.2s, transform 0.2s",
                              transform: formData.avatar_url === url ? "scale(1.1)" : "scale(1)",
                            }}
                          />
                        ))}
                      </div>
                      <label
                        style={{
                          display: "inline-flex",
                          alignItems: "center",
                          gap: "6px",
                          padding: "6px 12px",
                          background: "var(--bg-elevated)",
                          border: "1px solid var(--border)",
                          borderRadius: "var(--radius-sm)",
                          cursor: "pointer",
                          fontSize: "12px",
                          color: "var(--text-secondary)",
                        }}
                      >
                        <input
                          type="file"
                          accept="image/*"
                          onChange={handleAvatarUpload}
                          style={{ display: "none" }}
                        />
                        📷 Upload Custom Avatar
                      </label>
                    </div>
                    <div style={{ display: "flex", gap: "10px", justifyContent: "space-between", alignItems: "center" }}>
                      <button
                        type="button"
                        className="btn btn-ghost btn-sm"
                        onClick={handleTogglePrivacy}
                      >
                        {profile.is_private ? "Make Public" : "Make Private"}
                      </button>
                      <div style={{ display: "flex", gap: "8px" }}>
                        <button type="button" className="btn btn-ghost btn-sm" onClick={() => setEditing(false)}>
                          Cancel
                        </button>
                        <button type="submit" className="btn btn-primary btn-sm" disabled={saving}>
                          {saving ? "Saving..." : "Save Changes"}
                        </button>
                      </div>
                    </div>
                  </form>
                </div>
              )}

              {/* Profile Tabs */}
              <div style={{ display: "flex", gap: "4px", marginBottom: "20px", background: "var(--bg-surface)", borderRadius: "var(--radius-sm)", padding: "4px" }}>
                {(["posts", "followers", "following"] as const).map((t) => (
                  <button
                    key={t}
                    className={`btn btn-sm ${profileTab === t ? "btn-primary" : "btn-ghost"}`}
                    style={{ flex: 1, justifyContent: "center", textTransform: "capitalize" }}
                    onClick={() => setProfileTab(t)}
                  >
                    {t}
                  </button>
                ))}
              </div>

              {/* Tab Content */}
              {profileTab === "posts" && (
                posts.length === 0 ? (
                  <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "8px" }}>
                      <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
                      <polyline points="14 2 14 8 20 8" />
                    </svg>
                    <p>No posts yet</p>
                  </div>
                ) : (
                  posts.map((post: any) => (
                    <PostCard
                      key={post.id}
                      id={post.id}
                      username={post.author_username || profile.username}
                      userInitial={(post.author_username?.[0] || profile.username?.[0] || "U").toUpperCase()}
                      time={new Date(post.created_at * 1000).toLocaleString()}
                      privacy={post.privacy || "public"}
                      content={post.content}
                      likes={post.likes_count || 0}
                      liked={post.liked || false}
                      comments={post.comment_count || 0}
                    />
                  ))
                )
              )}

              {profileTab === "followers" && (
                followers.length === 0 ? (
                  <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "8px" }}>
                      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
                      <circle cx="12" cy="7" r="4" />
                    </svg>
                    <p>No followers yet</p>
                  </div>
                ) : (
                  <div className="card">
                    {followers.map((user) => (
                      <div
                        key={user.id}
                        className="user-item"
                        onClick={() => router.push(`/users/${user.id}`)}
                      >
                        <div
                          className="avatar avatar-md"
                          style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)", flexShrink: 0 }}
                        >
                          {(user.username?.[0] || "U").toUpperCase()}
                        </div>
                        <div>
                          <p style={{ fontWeight: 600, fontSize: "14px" }}>
                            {[user.first_name, user.last_name].filter(Boolean).join(" ") || user.username}
                          </p>
                          <p style={{ color: "var(--text-muted)", fontSize: "12px" }}>@{user.username}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                )
              )}

              {profileTab === "following" && (
                following.length === 0 ? (
                  <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "8px" }}>
                      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                      <circle cx="9" cy="7" r="4" />
                      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                    </svg>
                    <p>Not following anyone yet</p>
                  </div>
                ) : (
                  <div className="card">
                    {following.map((user) => (
                      <div
                        key={user.id}
                        className="user-item"
                        onClick={() => router.push(`/users/${user.id}`)}
                      >
                        <div
                          className="avatar avatar-md"
                          style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)", flexShrink: 0 }}
                        >
                          {(user.username?.[0] || "U").toUpperCase()}
                        </div>
                        <div>
                          <p style={{ fontWeight: 600, fontSize: "14px" }}>
                            {[user.first_name, user.last_name].filter(Boolean).join(" ") || user.username}
                          </p>
                          <p style={{ color: "var(--text-muted)", fontSize: "12px" }}>@{user.username}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                )
              )}
            </>
          )}
        </div>
      </main>
    </div>
  );
}