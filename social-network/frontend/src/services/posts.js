// frontend/src/services/posts.js
const API_URL = "";

export async function createPost(content, privacy = "public", imageFiles = [], allowedUsers = []) {
  const formData = new FormData();
  formData.append("content", content);
  // Map "private" to "selected" for backend compatibility
  const backendPrivacy = privacy === "private" ? "selected" : privacy;
  formData.append("privacy", backendPrivacy);
  if (imageFiles && imageFiles.length) {
    for (const f of imageFiles) {
      formData.append("image", f);
    }
  }
  if (allowedUsers && allowedUsers.length) {
    for (const u of allowedUsers) {
      formData.append("allowed_users[]", u);
    }
  }

  const res = await fetch(`${API_URL}/api/posts`, {
    method: "POST",
    credentials: "include",
    body: formData,
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to create post");
  }

  return res.json();
}

export async function fetchFeed(page = 1, perPage = 20) {
  const params = new URLSearchParams({ page: String(page), per_page: String(perPage) });
  const res = await fetch(`${API_URL}/api/posts/feed?${params}`, {
    method: "GET",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to fetch feed");
  }

  return res.json();
}

export async function fetchUserPosts(userId, page = 1, perPage = 20) {
  const params = new URLSearchParams({ page: String(page), per_page: String(perPage) });
  const res = await fetch(`${API_URL}/api/users/${userId}/posts?${params}`, {
    method: "GET",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to fetch user posts");
  }

  return res.json();
}

export async function likePost(postId) {
  const res = await fetch(`${API_URL}/api/posts/${postId}/like`, {
    method: "POST",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to like post");
  }

  return res.json();
}

export async function unlikePost(postId) {
  const res = await fetch(`${API_URL}/api/posts/${postId}/like`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to unlike post");
  }

  return res.json();
}

/**
 * @param {string} postId
 * @param {string} content
 * @param {File | null | undefined} [imageFile]
 */
export async function createComment(postId, content, imageFile = undefined) {
  const formData = new FormData();
  formData.append("content", content);
  if (imageFile) {
    formData.append("image", imageFile);
  }

  const res = await fetch(`${API_URL}/api/posts/${postId}/comments`, {
    method: "POST",
    credentials: "include",
    body: formData,
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to create comment");
  }

  return res.json();
}

export async function fetchComments(postId, page = 1, perPage = 20) {
  const params = new URLSearchParams({ page: String(page), per_page: String(perPage) });
  const res = await fetch(`${API_URL}/api/posts/${postId}/comments?${params}`, {
    method: "GET",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to fetch comments");
  }

  return res.json();
}