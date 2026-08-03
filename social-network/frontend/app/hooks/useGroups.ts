"use client";
import { useState, useEffect, useCallback } from "react";
import { fetchGroups, requestToJoin, acceptInvite, declineInvite } from "@/src/services/groups";

export interface Group {
  id: string;
  title: string;
  description: string;
  creator_id: string;
  privacy: "public" | "private";
  member_count?: number;
  is_member?: boolean;
  is_invited?: boolean;
  invite_id?: string;
  is_requested?: boolean;
  created_at?: number;
}

export function useGroups() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      setLoading(true);
      const data = await fetchGroups();
      setGroups(Array.isArray(data) ? data : data?.groups ?? []);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const joinGroup = async (groupId: string) => {
    try {
      await requestToJoin(groupId);
      await load();
    } catch (e: any) {
      throw new Error(e.message);
    }
  };

  const acceptGroupInvite = async (groupId: string, inviteId: string) => {
    try {
      await acceptInvite(groupId, inviteId);
      setGroups(prev => prev.map(g =>
        g.id === groupId ? { ...g, is_invited: false, is_member: true } : g
      ));
    } catch (e: any) {
      throw new Error(e.message);
    }
  };

  const declineGroupInvite = async (groupId: string, inviteId: string) => {
    try {
      await declineInvite(groupId, inviteId);
      setGroups(prev => prev.map(g =>
        g.id === groupId ? { ...g, is_invited: false } : g
      ));
    } catch (e: any) {
      throw new Error(e.message);
    }
  };

  return { groups, loading, error, joinGroup, acceptGroupInvite, declineGroupInvite, refresh: load };
}