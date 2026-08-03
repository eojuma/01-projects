import api from "./api";

export interface Conversation {
  id: string;
  other_user_id: string;
  other_username: string;
  last_message?: string;
  last_message_at?: number;
  unread_count?: number;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  sender_username?: string;
  content: string;
  created_at: number;
}

export async function fetchConversations(): Promise<Conversation[]> {
  const data = await api.get("/api/conversations");
  return data || [];
}

export async function fetchMessages(conversationId: string): Promise<Message[]> {
  const data = await api.get(`/api/conversations/${conversationId}/messages`);
  return data || [];
}

export async function fetchOrCreateConversation(userId: string): Promise<Message[]> {
  const data = await api.get(`/api/chat/user/${userId}`);
  return data || [];
}
