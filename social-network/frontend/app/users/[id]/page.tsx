"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Sidebar } from "../../components/Sidebar";
import { PostCard } from "../../components/PostCard";
import {
  followUser,
  unfollowUser,
  getFollowers,
  getFollowing,
  type FollowUser,
} from "@/src/services/follow";

interface UserProfile {
  id: string;
  email?: string;
  username: string;
  first_name?: string;
  last_name?: string;
  nickname?: string;
  bio?: string;
  avatar_url?: string;
  is_private: boolean;
  created_at: string;
  follow_status?: string; // "following" | "pending" | "none"
}

interface Post {
  id: string;
  author_id: string;
  author_username?: string;
  author_is_private: boolean;
  follow_status?: string;
  created_at: number;
  privacy?: string;
  content: string;
  likes_count?: number;
  liked?: boolean;
  comment_count?: number;
  image_urls?: string[];
}

type Tab = "posts" | "followers" | "following";

export default function UserProfilePage() {
  const params = useParams();
  const router = useRouter();
  const userId = params.id as string;

  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [followers, setFollowers] = useState<FollowUser[]>([]);
  const [following, setFollowing] = useState<FollowUser[]>([]);
  const [tab, setTab] = useState<Tab>("posts");
  const [loading, setLoading] = useState(true);
  const [followLoading, setFollowLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!userId) return;
    loadProfile();
  }, [userId]);

  const loadProfile = async () => {
    setLoading(true);
    setError("");
    try {
      // Fetch profile via /api/users/{id} (by ID - using the existing route)
      const res = await fetch(`/api/users/${userId}`, {
        credentials: "include",
      });
      if (!res.ok) {
        if (res.status === 401) { router.push("/login"); return; }
        throw new Error("User not found");
      }
      const data = await res.json();
      setProfile(data);

      // Fetch posts
      const postsRes = await fetch(`/api/users/${userId}/posts`, {
        credentials: "include",
      });
      if (postsRes.ok) {
        const postsData = await postsRes.json();
        // Ensure image URLs are absolute (start with '/') so they don't resolve
        // relative to the current route (e.g. /users/:id/uploads/...)
        const mapped = (postsData.posts ?? postsData ?? []).map((p: any) => {
          const raw: string[] = p.image_urls || [];
          const image_urls = raw.map((u: string) => (u && u.startsWith("/") ? u : `/${u}`));
          let image_url: string | undefined = undefined;
          if (image_urls.length > 0) image_url = image_urls[0];
          return { ...p, image_url, image_urls };
        });
        setPosts(mapped);
      }

      // Fetch followers/following counts
      try {
        const fwrs = await getFollowers(userId);
        setFollowers(fwrs);
      } catch { /* private profile */ }

      try {
        const fwing = await getFollowing(userId);
        setFollowing(fwing);
      } catch { /* private profile */ }
    } catch (err: any) {
      setError(err.message || "Failed to load profile");
    } finally {
      setLoading(false);
    }
  };

  const handleFollow = async () => {
    if (!profile) return;
    setFollowLoading(true);
    try {
      if (profile.follow_status === "following") {
        await unfollowUser(userId);
        await loadProfile();
      } else {
        const result = await followUser(userId);
        if (result.status === "accepted") {
          await loadProfile();
        } else {
          setProfile((p) => p ? { ...p, follow_status: "pending" } : p);
        }
      }
    } catch (err: any) {
      alert(err.message || "Action failed");
    } finally {
      setFollowLoading(false);
    }
  };

  const handleMessage = async () => {
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

  const displayName =
    [profile?.first_name, profile?.last_name].filter(Boolean).join(" ") ||
    profile?.username ||
    "Unknown";

  const followBtnLabel = () => {
    if (followLoading) return "…";
    if (profile?.follow_status === "following") return "Following";
    if (profile?.follow_status === "pending") return "Requested";
    return "Follow";
  };

  const followBtnClass = () => {
    if (profile?.follow_status === "following") return "btn btn-ghost btn-sm";
    return "btn btn-primary btn-sm";
  };

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
              <button className="btn btn-ghost btn-sm" onClick={() => router.back()}>← Go Back</button>
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
                      style={{ objectFit: "cover", flexShrink: 0 }}
                    />
                  ) : (
                    <div
                      className="avatar avatar-xl"
                      style={{ background: "linear-gradient(135deg, #7c3aed, #db2777)", flexShrink: 0 }}
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
                        &quot;{profile.nickname}&quot;
                      </p>
                    )}
                    {profile.bio && (
                      <p style={{ color: "var(--text-secondary)", fontSize: "14px", marginTop: "10px", lineHeight: "1.6" }}>
                        {profile.bio}
                      </p>
                    )}

                    {/* Stats */}
                    <div style={{ display: "flex", gap: "4px", marginTop: "14px" }}>
                      <button
                        className="stat-chip"
                        onClick={() => setTab("followers")}
                      >
                        <span className="stat-chip-value">{followers.length}</span>
                        <span className="stat-chip-label">Followers</span>
                      </button>
                      <button
                        className="stat-chip"
                        onClick={() => setTab("following")}
                      >
                        <span className="stat-chip-value">{following.length}</span>
                        <span className="stat-chip-label">Following</span>
                      </button>
                      <button
                        className="stat-chip"
                        onClick={() => setTab("posts")}
                      >
                        <span className="stat-chip-value">{posts.length}</span>
                        <span className="stat-chip-label">Posts</span>
                      </button>
                    </div>
                  </div>

                  {/* Actions */}
                  <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
                    <button
                      className={followBtnClass()}
                      onClick={handleFollow}
                      disabled={followLoading || profile.follow_status === "pending"}
                    >
                      {followBtnLabel()}
                    </button>
                    <button
                      className="btn btn-ghost btn-sm"
                      onClick={handleMessage}
                      style={{ display: "flex", alignItems: "center", gap: "6px" }}
                    >
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "16px", height: "16px" }}>
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                      </svg>
                      Message
                    </button>
                  </div>
                </div>
              </div>

              {/* Tabs */}
              <div style={{ display: "flex", gap: "4px", marginBottom: "20px", background: "var(--bg-surface)", borderRadius: "var(--radius-sm)", padding: "4px" }}>
                {(["posts", "followers", "following"] as Tab[]).map((t) => (
                  <button
                    key={t}
                    className={`btn btn-sm ${tab === t ? "btn-primary" : "btn-ghost"}`}
                    style={{ flex: 1, justifyContent: "center", textTransform: "capitalize" }}
                    onClick={() => setTab(t)}
                  >
                    {t}
                  </button>
                ))}
              </div>

              {/* Tab Content */}
              {tab === "posts" && (
                posts.length === 0 ? (
                  <div className="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "8px" }}>
                      <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
                      <polyline points="14 2 14 8 20 8" />
                    </svg>
                    <p>No posts yet</p>
                  </div>
                ) : (
                  posts.map((post) => (
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
                      imageUrls={post.image_urls}
                    />
                  ))
                )
              )}

              {tab === "followers" && (
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

              {tab === "following" && (
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
