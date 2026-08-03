type Listener = (data: { liked: boolean; likeCount: number; commentCount: number }) => void;

class PostStore {
  private liked = new Map<string, boolean>();
  private likeCounts = new Map<string, number>();
  private commentCounts = new Map<string, number>();
  private listeners = new Map<string, Set<Listener>>();

  getLiked(postId: string, initial: boolean = false): boolean {
    if (!this.liked.has(postId)) {
      this.liked.set(postId, initial);
    }
    return this.liked.get(postId)!;
  }

  getLikeCount(postId: string, initial: number = 0): number {
    if (!this.likeCounts.has(postId)) {
      this.likeCounts.set(postId, initial);
    }
    return this.likeCounts.get(postId)!;
  }

  getCommentCount(postId: string, initial: number = 0): number {
    if (!this.commentCounts.has(postId)) {
      this.commentCounts.set(postId, initial);
    }
    return this.commentCounts.get(postId)!;
  }

  setLiked(postId: string, val: boolean) {
    this.liked.set(postId, val);
    this.notify(postId);
  }

  setLikeCount(postId: string, count: number) {
    this.likeCounts.set(postId, count);
    this.notify(postId);
  }

  setCommentCount(postId: string, count: number) {
    this.commentCounts.set(postId, count);
    this.notify(postId);
  }

  subscribe(postId: string, listener: Listener) {
    if (!this.listeners.has(postId)) {
      this.listeners.set(postId, new Set());
    }
    this.listeners.get(postId)!.add(listener);
    return () => {
      this.listeners.get(postId)?.delete(listener);
    };
  }

  private notify(postId: string) {
    const set = this.listeners.get(postId);
    if (set) {
      const data = {
        liked: this.getLiked(postId),
        likeCount: this.getLikeCount(postId),
        commentCount: this.getCommentCount(postId),
      };
      set.forEach((listener) => {
        try {
          listener(data);
        } catch (e) {
          console.error(e);
        }
      });
    }
  }
}

export const postStore = new PostStore();
