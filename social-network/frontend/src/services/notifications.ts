import api from "./api";

export interface Notification {
  id: string;
  type: "follow_request" | "group_invite" | "group_join_request" | "event" | "follow_accepted" | "like";
  message: string;
  timestamp: string;
  actionable: boolean;
  read: boolean;
  ref_id?: string;
  group_id?: string;
}

export async function fetchNotifications(): Promise<Notification[]> {
  const data = await api.get("/api/notifications");
  return data || [];
}

export async function markAllNotificationsRead(): Promise<void> {
  await api.post("/api/notifications/read");
}

export async function respondToFollowRequest(requestID: string, action: "accept" | "decline"): Promise<void> {
  await api.post(`/api/follow-requests/${requestID}/respond?action=${action}`);
}

export async function respondToGroupInvite(groupID: string, inviteID: string, action: "accept" | "decline"): Promise<void> {
  await api.post(`/api/groups/${groupID}/invites/${inviteID}/${action}`);
}

export async function respondToGroupJoinRequest(groupID: string, requestID: string, action: "approve" | "decline"): Promise<void> {
  await api.post(`/api/groups/${groupID}/requests/${requestID}/${action}`);
}

