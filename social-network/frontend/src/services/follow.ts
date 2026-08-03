const API_URL = "";

export interface FollowUser {
  id: string;
  username: string;
  first_name?: string;
  last_name?: string;
  avatar_url?: string;
  follow_status?: string; // "following", "pending", "none"
}

export interface FollowResult {
  status: string; // "following", "pending"
  message?: string;
}

export async function followUser(userId: string): Promise<FollowResult> {
  const res = await fetch(`${API_URL}/api/users/${userId}/follow`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function unfollowUser(userId: string): Promise<void> {
  const res = await fetch(`${API_URL}/api/users/${userId}/unfollow`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error(await res.text());
}

export async function respondToFollowRequest(
  requestId: string,
  action: "accept" | "decline"
): Promise<void> {
  const res = await fetch(
    `${API_URL}/api/follow-requests/${requestId}/respond?action=${action}`,
    { method: "POST", credentials: "include" }
  );
  if (!res.ok) throw new Error(await res.text());
}

export async function getFollowers(userId: string): Promise<FollowUser[]> {
  const res = await fetch(`${API_URL}/api/users/${userId}/followers`, {
    credentials: "include",
  });
  if (!res.ok) {
    if (res.status === 403) throw new Error("private");
    throw new Error(await res.text());
  }
  const data = await res.json();
  // backend returns a paginated object { users: [...] } — accept either shape
  return data.users ?? data.followers ?? data ?? [];
}

export async function getFollowing(userId: string): Promise<FollowUser[]> {
  const res = await fetch(`${API_URL}/api/users/${userId}/following`, {
    credentials: "include",
  });
  if (!res.ok) {
    if (res.status === 403) throw new Error("private");
    throw new Error(await res.text());
  }
  const data = await res.json();
  // backend returns a paginated object { users: [...] } — accept either shape
  return data.users ?? data.following ?? data ?? [];
}
