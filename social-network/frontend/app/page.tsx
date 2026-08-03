"use client";

import { useState, useEffect } from "react";
import { Sidebar } from "./components/Sidebar";
import { PostCard } from "./components/PostCard";
import { usePosts } from "./hooks/usePosts";
import { EmojiPicker } from "./components/EmojiPicker";
import { GifPicker } from "./components/GifPicker";
import { togglePrivacy } from "@/src/services/profile";
import { getFollowers } from "@/src/services/follow";

export default function Home() {
  const { posts, loading, creating, error, submitPost } = usePosts();
  const [newPostContent, setNewPostContent] = useState("");
  const [newPostImages, setNewPostImages] = useState<File[]>([]);
  const [newPostImagePreviews, setNewPostImagePreviews] = useState<string[]>([]);
  const [privacy, setPrivacy] = useState("public");
  const [followers, setFollowers] = useState<any[]>([]);
  const [selectedFollowers, setSelectedFollowers] = useState<string[]>([]);
  const [loadingFollowers, setLoadingFollowers] = useState(false);

  // Account privacy state
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);
  const [isPrivateAccount, setIsPrivateAccount] = useState(false);
  const [showPrivacyModal, setShowPrivacyModal] = useState(false);
  const [togglingPrivacy, setTogglingPrivacy] = useState(false);

  useEffect(() => {
    fetch("/api/auth/me", { credentials: "include" })
      .then((res) => (res.ok ? res.json() : null))
      .then((me) => {
        if (me && me.id) {
          setCurrentUserId(String(me.id));
          setIsPrivateAccount(!!me.is_private);
          setLoadingFollowers(true);
          getFollowers(String(me.id))
            .then((data) => setFollowers(data || []))
            .catch(console.error)
            .finally(() => setLoadingFollowers(false));
        }
      })
      .catch(console.error);
  }, []);

  const handleConfirmPrivacyToggle = async () => {
    if (!currentUserId) return;
    setTogglingPrivacy(true);
    try {
      await togglePrivacy(currentUserId, !isPrivateAccount);
      setIsPrivateAccount((prev) => !prev);
    } catch (err) {
      console.error("Failed to toggle privacy:", err);
    } finally {
      setTogglingPrivacy(false);
      setShowPrivacyModal(false);
    }
  };

  const handleSelectGif = async (gifUrl: string) => {
    try {
      const response = await fetch(gifUrl);
      const blob = await response.blob();
      const file = new File([blob], "selected.gif", { type: "image/gif" });
      setNewPostImages((prev) => [...prev, file]);
      setNewPostImagePreviews((prev) => [...prev, gifUrl]);
    } catch (err) {
      console.error("Failed to select GIF:", err);
      alert("Failed to load GIF from library");
    }
  };

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPostContent.trim() && newPostImages.length === 0) return;
    try {
      await submitPost(newPostContent, privacy, newPostImages, selectedFollowers);
      setNewPostContent("");
      setNewPostImages([]);
      setNewPostImagePreviews([]);
      setSelectedFollowers([]);
    } catch (err) {
      console.error("Failed to create post:", err);
    }
  };

  const handlePostImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newFiles = Array.from(e.target.files || []);
    setNewPostImages((prev) => [...prev, ...newFiles]);
    setNewPostImagePreviews((prev) => [...prev, ...newFiles.map((f) => URL.createObjectURL(f))]);
    e.target.value = "";
  };

  const handleRemovePostImage = (index?: number) => {
    if (index === undefined) {
      setNewPostImages([]);
      setNewPostImagePreviews([]);
      return;
    }
    const imgs = [...newPostImages];
    const prevs = [...newPostImagePreviews];
    imgs.splice(index, 1);
    prevs.splice(index, 1);
    setNewPostImages(imgs);
    setNewPostImagePreviews(prevs);
  };

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <div className="page-container">

          {/* Header with privacy toggle */}
          <div className="page-header" style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
            <h1 className="page-title">Feed</h1>

            <button
              onClick={() => setShowPrivacyModal(true)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "6px 14px",
                background: "var(--bg-elevated)",
                border: "1px solid var(--border-light)",
                borderRadius: "20px",
                cursor: "pointer",
                fontSize: "13px",
                color: "var(--text-secondary)",
              }}
            >
              {isPrivateAccount ? (
                <>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ width: "14px", height: "14px" }}>
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                  </svg>
                  <span style={{ fontWeight: 500, color: "var(--text-warning, #b45309)" }}>Private</span>
                </>
              ) : (
                <>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ width: "14px", height: "14px" }}>
                    <circle cx="12" cy="12" r="10" />
                    <line x1="2" y1="12" x2="22" y2="12" />
                    <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                  </svg>
                  <span style={{ fontWeight: 500, color: "var(--text-success, #15803d)" }}>Public</span>
                </>
              )}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ width: "12px", height: "12px" }}>
                <polyline points="6 9 12 15 18 9" />
              </svg>
            </button>
          </div>

          <p style={{ fontSize: "12px", color: "var(--text-muted)", marginBottom: "16px", marginTop: "-8px" }}>
            {isPrivateAccount
              ? "Only your approved followers can see your profile and posts."
              : "Your profile and posts are visible to everyone."}
          </p>

          {/* Privacy confirmation modal */}
          {showPrivacyModal && (
            <div
              onClick={() => setShowPrivacyModal(false)}
              style={{
                position: "fixed", inset: 0,
                background: "rgba(0,0,0,0.5)",
                zIndex: 1000, display: "flex", alignItems: "center", justifyContent: "center",
              }}
            >
              <div
                onClick={(e) => e.stopPropagation()}
                style={{
                  background: "var(--bg-card, var(--bg-primary))",
                  border: "1px solid var(--border)",
                  borderRadius: "16px",
                  padding: "24px",
                  width: "320px",
                }}
              >
                <h3 style={{ fontSize: "16px", fontWeight: 600, marginBottom: "8px", color: "var(--text-primary)" }}>
                  {isPrivateAccount ? "Switch to Public?" : "Switch to Private?"}
                </h3>
                <p style={{ fontSize: "13px", color: "var(--text-secondary)", marginBottom: "20px", lineHeight: 1.5 }}>
                  {isPrivateAccount
                    ? "Everyone will be able to see your profile and follow you without approval. All pending follow requests will be accepted."
                    : "Only your approved followers will be able to see your profile and posts. New followers will need to send a request first."}
                </p>
                <div style={{ display: "flex", gap: "8px", justifyContent: "flex-end" }}>
                  <button
                    onClick={() => setShowPrivacyModal(false)}
                    className="btn btn-secondary"
                    style={{ fontSize: "13px", padding: "7px 14px" }}
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleConfirmPrivacyToggle}
                    disabled={togglingPrivacy}
                    className="btn btn-primary"
                    style={{ fontSize: "13px", padding: "7px 14px" }}
                  >
                    {togglingPrivacy ? "Saving..." : isPrivateAccount ? "Make Public" : "Make Private"}
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Create Post */}
          <form onSubmit={handleCreatePost} className="card" style={{ marginBottom: "24px" }}>
            <div style={{ display: "flex", gap: "12px" }}>
              <div className="avatar avatar-md" style={{ background: "linear-gradient(135deg, #4f8ef7, #7c3aed)", flexShrink: 0 }}>
                Me
              </div>
              <div style={{ flex: 1 }}>
                <textarea
                  className="input"
                  value={newPostContent}
                  onChange={(e) => setNewPostContent(e.target.value)}
                  placeholder="What's on your mind?"
                  rows={3}
                  style={{ marginBottom: "12px" }}
                />
                <div style={{ display: "flex", gap: "8px", marginBottom: "12px", alignItems: "center" }}>
                  <input
                    id="post-image-input"
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={handlePostImageChange}
                    disabled={creating}
                    style={{ display: "none" }}
                  />
                  <label
                    htmlFor="post-image-input"
                    style={{
                      display: "inline-flex",
                      alignItems: "center",
                      gap: "8px",
                      padding: "8px 14px",
                      backgroundColor: "var(--bg-elevated)",
                      color: "var(--accent)",
                      border: "1px solid var(--border-light)",
                      borderRadius: "var(--radius-sm)",
                      cursor: creating ? "not-allowed" : "pointer",
                      fontWeight: 500,
                      fontSize: "13px",
                      opacity: creating ? 0.6 : 1,
                    }}
                    onMouseEnter={(e) => !creating && (e.currentTarget.style.backgroundColor = "var(--bg-hover)")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "var(--bg-elevated)")}
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ width: "16px", height: "16px" }}>
                      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                      <polyline points="17 8 12 3 7 8" />
                      <line x1="12" y1="3" x2="12" y2="15" />
                    </svg>
                    <span>Add Image(s)</span>
                  </label>
                  {newPostImages.length > 0 && (
                    <span style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                      {newPostImages.length} file{newPostImages.length !== 1 ? "s" : ""} selected
                    </span>
                  )}
                </div>
                {newPostImagePreviews.length > 0 && (
                  <div style={{ display: "flex", gap: "10px", marginBottom: "12px", alignItems: "center" }}>
                    {newPostImagePreviews.map((src, idx) => (
                      <div key={idx} style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <img src={src} alt={`Preview ${idx + 1}`} style={{ width: "88px", height: "88px", borderRadius: "12px", objectFit: "cover" }} />
                        <button type="button" className="btn btn-secondary btn-sm" onClick={() => handleRemovePostImage(idx)} style={{ padding: "8px 12px" }}>
                          Remove
                        </button>
                      </div>
                    ))}
                  </div>
                )}
                {privacy === "private" && (
                  <div className="card" style={{ marginTop: "0", marginBottom: "12px", background: "var(--bg-elevated)", padding: "12px" }}>
                    <div style={{ fontWeight: 600, fontSize: "13px", marginBottom: "8px" }}>
                      Choose who can see this post:
                    </div>
                    {loadingFollowers ? (
                      <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>Loading followers...</div>
                    ) : followers.length === 0 ? (
                      <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                        You don't have any followers yet to select.
                      </div>
                    ) : (
                      <div style={{ display: "flex", flexDirection: "column", gap: "6px", maxHeight: "120px", overflowY: "auto" }}>
                        {followers.map((f: any) => (
                          <label key={f.id} style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "13px", cursor: "pointer" }}>
                            <input
                              type="checkbox"
                              checked={selectedFollowers.includes(f.id)}
                              onChange={(e) => {
                                if (e.target.checked) {
                                  setSelectedFollowers((prev: string[]) => [...prev, f.id]);
                                } else {
                                  setSelectedFollowers((prev: string[]) => prev.filter((id: string) => id !== f.id));
                                }
                              }}
                            />
                            <span>{f.first_name || f.username} (@{f.username})</span>
                          </label>
                        ))}
                      </div>
                    )}
                  </div>
                )}
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
                    <select
                      value={privacy}
                      onChange={(e) => setPrivacy(e.target.value)}
                      className="input"
                      style={{ width: "140px", padding: "7px 12px" }}
                    >
                      <option value="public">Public</option>
                      <option value="followers">Followers</option>
                      <option value="private">Private</option>
                    </select>
                    <EmojiPicker onSelect={(emoji) => setNewPostContent((p) => p + emoji)} />
                    <GifPicker onSelect={handleSelectGif} />
                  </div>
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={creating || (newPostContent.trim() === "" && newPostImages.length === 0)}
                    suppressHydrationWarning
                    style={{ marginLeft: "8px" }}
                  >
                    {creating ? "Posting..." : "Post"}
                  </button>
                </div>
              </div>
            </div>
          </form>

          {/* Feed */}
          {loading ? (
            <div className="loading-spinner"><div className="spinner" /></div>
          ) : error ? (
            <div className="card" style={{ color: "var(--danger)", textAlign: "center" }}>{error}</div>
          ) : posts.length === 0 ? (
            <div className="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "36px", height: "36px", color: "var(--text-muted)", marginBottom: "12px" }}>
                <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
                <polyline points="14 2 14 8 20 8" />
              </svg>
              <p>No posts yet. Be the first to share something!</p>
            </div>
          ) : (
            posts.map((post) => (
              <PostCard
                key={post.id}
                id={post.id}
                authorId={post.author_id}
                currentUserId={currentUserId}
                followStatus={post.follow_status}
                authorIsPrivate={post.author_is_private}
                username={post.author_username || "Unknown"}
                userInitial={(post.author_username?.[0] || "U").toUpperCase()}
                time={new Date(post.created_at * 1000).toLocaleString()}
                privacy={post.privacy || "public"}
                content={post.content}
                likes={post.likes_count || 0}
                liked={post.liked || false}
                comments={post.comment_count || 0}
                imageUrls={post.image_urls}
              />
            ))
          )}
        </div>
      </main>
    </div>
  );
}