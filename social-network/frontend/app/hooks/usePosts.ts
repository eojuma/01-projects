"use client";
import { useState, useEffect, useCallback } from "react";
import { fetchFeed, createPost } from "@/src/services/posts";

export interface Post {
  id: string;
  author_id?: string;
  author_username: string;
  author_is_private: boolean; 
  content: string;
  privacy: string;
  image_url?: string;
  image_urls?: string[];
  created_at: number;
  comment_count?: number;
  likes_count?: number;
  liked?: boolean;
  follow_status?: string;
}

export function usePosts() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [creating, setCreating] = useState(false);

  const loadFeed = useCallback(async () => {
    try {
      setLoading(true);
      const data = await fetchFeed(1, 20);
      const mapped = (data?.posts ?? []).map((p: any) => {
        const raw: string[] = p.image_urls || [];
        const image_urls = raw.map((u: string) => (u.startsWith("/") ? u : `/${u}`));
        let image_url: string | undefined = undefined;
        if (image_urls.length > 0) image_url = image_urls[0];
        return { ...p, image_url, image_urls };
      });
      setPosts(mapped);
    } catch (e: any) {
      setError(e.message);
      setPosts([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadFeed();
  }, [loadFeed]);

  const submitPost = async (
    content: string,
    privacy: string = "public",
    imageFiles: File[] | null = null,
    allowedUsers: string[] = []
  ) => {
    setCreating(true);
    try {
      const newPost = await createPost(content, privacy, imageFiles || [], allowedUsers);
      setPosts(prev => [newPost, ...prev]);
      return newPost;
    } catch (e: any) {
      throw new Error(e.message);
    } finally {
      setCreating(false);
    }
  };

  return { posts, loading, error, creating, submitPost, refresh: loadFeed };
}