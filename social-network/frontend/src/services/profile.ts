const API_URL = "";

export interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  nickname?: string;
  bio: string;
  avatar_url: string;
  is_private: boolean;
  created_at: string;
  updated_at: string;
}

export async function getUserProfile(username: string): Promise<User> {
  const res = await fetch(`${API_URL}/api/users/${username}`, {
    method: "GET",
    credentials: "include",
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to fetch user profile");
  }

  return res.json();
}

export async function updateProfile(
  userId: string,
  data: {
    first_name?: string;
    last_name?: string;
    bio?: string;
    nickname?: string;
    avatar_url?: string;
  }
): Promise<User> {
  const res = await fetch(`${API_URL}/api/users/${userId}/profile`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to update profile");
  }

  return res.json();
}

export async function togglePrivacy(userId: string, isPrivate: boolean): Promise<void> {
  const res = await fetch(`${API_URL}/api/users/${userId}/privacy`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ is_private: isPrivate }),
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to update privacy");
  }
}

export async function getUserPosts(userId: string) {
  const res = await fetch(`${API_URL}/api/users/${userId}/posts`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.posts ?? data ?? [];
}