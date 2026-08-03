"use client";

import { useState, useEffect } from "react";
import { fetchComments, createComment } from "@/src/services/posts";
import { postStore } from "@/src/store/postStore";
import { EmojiPicker } from "./EmojiPicker";
import { GifPicker } from "./GifPicker";
import { followUser, unfollowUser } from "@/src/services/follow";

interface CommentAuthor {
  id: string;
  username: string;
  avatar_url?: string;
}

interface Comment {
  id: string;
  post_id: string;
  author_id: string;
  content: string;
  image_path?: string;
  created_at: number;
  author: CommentAuthor;
}

interface PostCardProps {
  id: string;
  authorId?: string;
  currentUserId?: string | null;
  followStatus?: string;
  authorIsPrivate?: boolean;
  username: string;
  userInitial: string;
  time: string;
  privacy: string;
  content: string;
  likes?: number;
  liked?: boolean;
  comments?: number;
  imageUrls?: string[];
  onComment?: (postId: string) => void;
}

export function PostCard({
  id,
  authorId,
  currentUserId,
  followStatus,
  authorIsPrivate = false,
  username,
  userInitial,
  time,
  privacy,
  content,
  likes = 0,
  liked = false,
  comments = 0,
  imageUrls = [],
}: PostCardProps) {
  const [likedState, setLikedState] = useState(() => postStore.getLiked(id, liked));
  const [likeCount, setLikeCount] = useState(() => postStore.getLikeCount(id, likes));
  const [commentCount, setCommentCount] = useState(() => postStore.getCommentCount(id, comments));

  const [showComments, setShowComments] = useState(false);
  const [commentsList, setCommentsList] = useState<Comment[]>([]);
  const [loadingComments, setLoadingComments] = useState(false);
  const [newCommentText, setNewCommentText] = useState("");
  const [submittingComment, setSubmittingComment] = useState(false);
  const [commentImage, setCommentImage] = useState<File | null>(null);
  const [commentImagePreview, setCommentImagePreview] = useState<string | null>(null);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  type FollowState = "self" | "following" | "pending" | "none";
  const [followState, setFollowState] = useState<FollowState>(
    (followStatus as FollowState) || "none"
  );
  const [followLoading, setFollowLoading] = useState(false);
  const [hoveringFollow, setHoveringFollow] = useState(false);

  const isOwnPost = followStatus === "self" || (currentUserId && authorId && currentUserId === authorId);

  const handleSelectCommentGif = async (gifUrl: string) => {
    try {
      const response = await fetch(gifUrl);
      const blob = await response.blob();
      const file = new File([blob], "comment.gif", { type: "image/gif" });
      setCommentImage(file);
      setCommentImagePreview(gifUrl);
    } catch (err) {
      console.error("Failed to select GIF:", err);
      alert("Failed to load GIF");
    }
  };

  useEffect(() => {
    postStore.setLiked(id, liked);
    postStore.setLikeCount(id, likes);
    postStore.setCommentCount(id, comments);
    const unsubscribe = postStore.subscribe(id, (data) => {
      setLikedState(data.liked);
      setLikeCount(data.likeCount);
      setCommentCount(data.commentCount);
    });
    return unsubscribe;
  }, [id, likes, liked, comments]);

  useEffect(() => {
    if (followStatus) setFollowState(followStatus as FollowState);
  }, [followStatus]);

  const handleFollowToggle = async () => {
    if (!authorId || followLoading) return;
    setFollowLoading(true);
    try {
      if (followState === "following") {
        await unfollowUser(authorId);
        setFollowState("none");
      } else if (followState === "none") {
        const result = await followUser(authorId);
        setFollowState(result.status as FollowState);
      }
    } catch (err) {
      console.error("Follow action failed:", err);
    } finally {
      setFollowLoading(false);
    }
  };

  const followButtonLabel = () => {
    if (followLoading) return "...";
    if (followState === "following") return hoveringFollow ? "Unfollow" : "Following ✓";
    if (followState === "pending") return "Requested";
    return "+ Follow";
  };

  const followButtonStyle: React.CSSProperties = {
    fontSize: "12px",
    fontWeight: 500,
    padding: "4px 12px",
    borderRadius: "20px",
    border: "1px solid",
    cursor: followLoading || followState === "pending" ? "not-allowed" : "pointer",
    transition: "all 0.15s",
    ...(followState === "following"
      ? hoveringFollow
        ? { background: "#fef2f2", color: "#b91c1c", borderColor: "#fca5a5" }
        : { background: "#eff6ff", color: "#1d4ed8", borderColor: "#bfdbfe" }
      : followState === "pending"
      ? { background: "var(--bg-elevated)", color: "var(--text-muted)", borderColor: "var(--border-light)", opacity: 0.7 }
      : { background: "var(--bg-elevated)", color: "var(--text-secondary)", borderColor: "var(--border-light)" }),
  };

  const handleLike = async () => {
    try {
      const { likePost } = await import("@/src/services/posts");
      const res = await likePost(id);
      postStore.setLiked(id, res.liked);
      postStore.setLikeCount(id, res.likes_count);
    } catch (err) {
      console.error("Failed to like post:", err);
    }
  };

  const toggleComments = async () => {
    const nextShow = !showComments;
    setShowComments(nextShow);
    if (nextShow && commentsList.length === 0) {
      setLoadingComments(true);
      try {
        const data = await fetchComments(id, 1, 100);
        setCommentsList(data || []);
      } catch (err) {
        console.error("Failed to fetch comments:", err);
      } finally {
        setLoadingComments(false);
      }
    }
  };

  const handleShare = async () => {
    const shareUrl = `${window.location.origin}${window.location.pathname}#post-${id}`;
    if (navigator.share) {
      try { await navigator.share({ title: `Post by ${username}`, text: content.slice(0, 120), url: shareUrl }); return; } catch {}
    }
    try { await navigator.clipboard.writeText(shareUrl); alert("Post link copied to clipboard"); }
    catch { prompt("Copy this link", shareUrl); }
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    setCommentImage(file);
    setCommentImagePreview(file ? URL.createObjectURL(file) : null);
  };

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCommentText.trim() && !commentImage) return;
    setSubmittingComment(true);
    try {
      const addedComment = await createComment(id, newCommentText.trim(), commentImage ?? undefined);
      setCommentsList((prev) => [...prev, addedComment]);
      postStore.setCommentCount(id, commentCount + 1);
      setNewCommentText("");
      setCommentImage(null);
      setCommentImagePreview(null);
    } catch (err) {
      console.error("Failed to add comment:", err);
      alert("Failed to submit comment");
    } finally {
      setSubmittingComment(false);
    }
  };

  const privacyLabel =
    privacy === "public" ? "Public" :
    privacy === "followers" ? "Followers" :
    privacy === "selected" ? "Private" : privacy;

  const gradientColors = [
    "linear-gradient(135deg, #4f8ef7, #7c3aed)",
    "linear-gradient(135deg, #43d97b, #0ea5e9)",
    "linear-gradient(135deg, #f6b93b, #e55039)",
    "linear-gradient(135deg, #a29bfe, #fd79a8)",
  ];
  const colorIdx = username.charCodeAt(0) % gradientColors.length;

  return (
    <article id={`post-${id}`} className="card fade-in">
      {/* Author row */}
      <div style={{ display: "flex", gap: "12px", alignItems: "center" }}>
        <div className="avatar avatar-md" style={{ background: gradientColors[colorIdx] }}>{userInitial}</div>
        <div style={{ flex: 1 }}>
          <div style={{ fontWeight: 600, fontSize: "14px" }}>{username}</div>
          <div style={{ fontSize: "12px", color: "var(--text-muted)", marginTop: "1px" }}>
            {time} · <span style={{ color: "var(--text-secondary)" }}>{privacyLabel}</span>
            {authorIsPrivate && (
              <span style={{ marginLeft: "6px", fontSize: "11px" }}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ width: "11px", height: "11px", verticalAlign: "middle" }}>
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
              </span>
            )}
          </div>
        </div>

        {!isOwnPost && authorId && (
          <button
            onClick={handleFollowToggle}
            onMouseEnter={() => followState === "following" && setHoveringFollow(true)}
            onMouseLeave={() => setHoveringFollow(false)}
            disabled={followLoading || followState === "pending"}
            style={followButtonStyle}
          >
            {followButtonLabel()}
          </button>
        )}
      </div>

      {/* Content */}
      <p style={{ marginTop: "14px", fontSize: "14px", color: "var(--text-primary)", lineHeight: "1.65" }}>{content}</p>

      {/* Images carousel */}
      {imageUrls && imageUrls.length > 0 && (
        <div style={{ marginTop: "12px", position: "relative" }}>
          <img src={imageUrls[currentImageIndex]} alt={`Post image ${currentImageIndex + 1}`} style={{ width: "100%", borderRadius: "8px", maxHeight: "360px", objectFit: "cover" }} />
          {imageUrls.length > 1 && (
            <>
              <button onClick={() => setCurrentImageIndex((i) => (i - 1 + imageUrls.length) % imageUrls.length)} aria-label="Previous image" style={{ position: "absolute", left: 8, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,0.4)", border: "none", color: "white", padding: "6px 8px", borderRadius: "6px", cursor: "pointer" }}>‹</button>
              <button onClick={() => setCurrentImageIndex((i) => (i + 1) % imageUrls.length)} aria-label="Next image" style={{ position: "absolute", right: 8, top: "50%", transform: "translateY(-50%)", background: "rgba(0,0,0,0.4)", border: "none", color: "white", padding: "6px 8px", borderRadius: "6px", cursor: "pointer" }}>›</button>
              <div style={{ position: "absolute", left: "50%", transform: "translateX(-50%)", bottom: 8, display: "flex", gap: "6px" }}>
                {imageUrls.map((_, idx) => (
                  <span key={idx} onClick={() => setCurrentImageIndex(idx)} style={{ width: 8, height: 8, borderRadius: 8, background: idx === currentImageIndex ? "white" : "rgba(255,255,255,0.6)", display: "inline-block", cursor: "pointer" }} />
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="post-actions">
        <button onClick={handleLike} className={`post-action-btn ${likedState ? "liked" : ""}`} style={{ display: "flex", alignItems: "center", gap: "4px" }}>
          <svg viewBox="0 0 24 24" fill={likedState ? "var(--danger)" : "none"} stroke={likedState ? "var(--danger)" : "currentColor"} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "16px", height: "16px" }}>
            <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
          </svg>
          <span>{likeCount > 0 ? likeCount : "Like"}</span>
        </button>
        <button className="post-action-btn" onClick={toggleComments} style={{ display: "flex", alignItems: "center", gap: "4px" }}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "16px", height: "16px" }}>
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
          <span>{commentCount > 0 ? commentCount : "Comment"}</span>
        </button>
        <button className="post-action-btn" onClick={handleShare} style={{ display: "flex", alignItems: "center", gap: "4px" }}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "16px", height: "16px" }}>
            <circle cx="18" cy="5" r="3" /><circle cx="6" cy="12" r="3" /><circle cx="18" cy="19" r="3" />
            <line x1="8.59" y1="13.51" x2="15.42" y2="17.49" /><line x1="15.41" y1="6.51" x2="8.59" y2="10.49" />
          </svg>
          <span>Share</span>
        </button>
      </div>

      {/* Comments Section */}
      {showComments && (
        <div style={{ marginTop: "16px", borderTop: "1px solid var(--border)", paddingTop: "16px" }}>
          {loadingComments ? (
            <div style={{ display: "flex", justifyContent: "center", padding: "12px 0" }}>
              <div className="spinner" style={{ width: "20px", height: "20px", borderWidth: "2px" }} />
            </div>
          ) : (
            <div style={{ display: "flex", flexDirection: "column", gap: "12px", marginBottom: "16px" }}>
              {commentsList.length === 0 ? (
                <div style={{ color: "var(--text-muted)", fontSize: "13px", textAlign: "center", padding: "8px 0" }}>No comments yet.</div>
              ) : (
                commentsList.map((c) => {
                  const initial = (c.author?.username?.[0] || "?").toUpperCase();
                  const cGradient = (c.author?.username || "").charCodeAt(0) % gradientColors.length || 0;
                  return (
                    <div key={c.id} style={{ display: "flex", gap: "10px", alignItems: "flex-start" }}>
                      <div className="avatar avatar-sm" style={{ background: gradientColors[cGradient] }}>{initial}</div>
                      <div style={{ flex: 1, background: "var(--bg-elevated)", padding: "10px 12px", borderRadius: "12px" }}>
                        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "2px" }}>
                          <span style={{ fontWeight: 600, fontSize: "13px" }}>{c.author?.username || "Unknown"}</span>
                          <span style={{ fontSize: "10px", color: "var(--text-muted)" }}>
                            {new Date(c.created_at * 1000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                          </span>
                        </div>
                        <p style={{ fontSize: "13px", color: "var(--text-primary)", lineHeight: "1.4" }}>{c.content}</p>
                        {c.image_path && <img src={c.image_path} alt="Comment attachment" style={{ marginTop: "10px", width: "100%", borderRadius: "10px", maxHeight: "300px", objectFit: "cover" }} />}
                      </div>
                    </div>
                  );
                })
              )}
            </div>
          )}
          <form onSubmit={handleSubmitComment} style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
            <input className="input" type="text" placeholder="Write a comment..." value={newCommentText} onChange={(e) => setNewCommentText(e.target.value)} style={{ width: "100%", fontSize: "13px", padding: "8px 12px" }} disabled={submittingComment} />
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <div style={{ display: "flex", gap: "12px", alignItems: "center" }}>
                <label style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "13px", color: "var(--text-muted)", cursor: "pointer" }}>
                  <input type="file" accept="image/*" onChange={handleImageChange} disabled={submittingComment} style={{ display: "none" }} />
                  <span style={{ textDecoration: "underline" }}>Attach image/GIF</span>
                </label>
                <EmojiPicker onSelect={(emoji) => setNewCommentText((p) => p + emoji)} />
                <GifPicker onSelect={handleSelectCommentGif} />
              </div>
              <button type="submit" className="btn btn-primary btn-sm" disabled={submittingComment || (!newCommentText.trim() && !commentImage)} style={{ padding: "8px 16px" }}>
                {submittingComment ? "..." : "Send"}
              </button>
            </div>
            {commentImagePreview && (
              <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                <img src={commentImagePreview} alt="Comment preview" style={{ width: "72px", height: "72px", borderRadius: "12px", objectFit: "cover" }} />
                <button type="button" className="btn btn-secondary btn-sm" onClick={() => { setCommentImage(null); setCommentImagePreview(null); }} style={{ padding: "8px 12px" }}>Remove attachment</button>
              </div>
            )}
          </form>
        </div>
      )}
    </article>
  );
}