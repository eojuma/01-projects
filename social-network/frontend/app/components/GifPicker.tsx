"use client";

import { useState, useEffect } from "react";

interface GifPickerProps {
  onSelect: (gifUrl: string) => void;
}

const GIF_LIBRARY: Record<string, string[]> = {
  trending: [
    "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
    "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
    "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
    "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
    "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
    "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
    "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
    "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
    "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
    "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
    "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
    "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
  ],
  happy: [
    "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
    "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
    "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
    "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
    "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
    "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
    "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
    "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
    "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
    "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
    "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
    "https://media1.giphy.com/media/26tknCqiJrBQG6bxC/200.gif",
  ],
  sad: [
    "https://media1.giphy.com/media/d2lcHJTG5Tscg/200.gif",
    "https://media1.giphy.com/media/OPU6wzx8JrHna/200.gif",
    "https://media1.giphy.com/media/9Y5BbDSkSTiY8/200.gif",
    "https://media1.giphy.com/media/3o6wrvdHFbwBrUFenu/200.gif",
    "https://media1.giphy.com/media/l2JhtKtDWYNKdRpoA/200.gif",
    "https://media1.giphy.com/media/3ohzdIuqJoo8QdKlnW/200.gif",
    "https://media1.giphy.com/media/26ufnwz3wDUli7GU0/200.gif",
    "https://media1.giphy.com/media/ISOckXUybVfQ4/200.gif",
    "https://media1.giphy.com/media/BEob5qwFkSJ7G/200.gif",
    "https://media1.giphy.com/media/L95W4wv8nnb9K/200.gif",
    "https://media1.giphy.com/media/ROF8OQvDmxytW/200.gif",
    "https://media1.giphy.com/media/26FPy3QZQqGtDcrja/200.gif",
  ],
  laugh: [
    "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
    "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
    "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
    "https://media1.giphy.com/media/3o7aCTfyhYawMw0wVi/200.gif",
    "https://media1.giphy.com/media/10JhviFuU2gWD6/200.gif",
    "https://media1.giphy.com/media/ZqlvCTNHpqrio/200.gif",
    "https://media1.giphy.com/media/26BRBKqUiq586bRVm/200.gif",
    "https://media1.giphy.com/media/bC9czlgCMtw4cj8RgH/200.gif",
    "https://media1.giphy.com/media/Q7ozWVYCR0nyW2rvPW/200.gif",
    "https://media1.giphy.com/media/l3q2K5jinAlChoCLS/200.gif",
    "https://media1.giphy.com/media/xUA7aM09ByyR1w5YWc/200.gif",
    "https://media1.giphy.com/media/26tOZ42Mg6r9CUEBq/200.gif",
  ],
  love: [
    "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
    "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
    "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
    "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
    "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
    "https://media1.giphy.com/media/26FLdmIp6wJr91JAI/200.gif",
    "https://media1.giphy.com/media/3oEdv6sy3ulljPMGdy/200.gif",
    "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
    "https://media1.giphy.com/media/M90mJvfWfd5mbUuULX/200.gif",
    "https://media1.giphy.com/media/l0ExncehJzexFpRHq/200.gif",
    "https://media1.giphy.com/media/26tknCqiJrBQG6bxC/200.gif",
    "https://media1.giphy.com/media/3oEjHV0z8S7WM4MwnK/200.gif",
  ],
  cool: [
    "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
    "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
    "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
    "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
    "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
    "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
    "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
    "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
    "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
    "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
    "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
    "https://media1.giphy.com/media/3o6Zt6KHxJTbXCnSvu/200.gif",
  ],
  wow: [
    "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
    "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
    "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
    "https://media1.giphy.com/media/5VKbvrjxpVJCM/200.gif",
    "https://media1.giphy.com/media/l0MYEqEzwMWFCg8rm/200.gif",
    "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
    "https://media1.giphy.com/media/26AHvEoTlR5PGWsJq/200.gif",
    "https://media1.giphy.com/media/xT0GqtpF1NWd9BcuI0/200.gif",
    "https://media1.giphy.com/media/3o7abKhOpu0NwenH3O/200.gif",
    "https://media1.giphy.com/media/l3q2K5jinAlChoCLS/200.gif",
    "https://media1.giphy.com/media/26n6Gx9moCgs1pUuk/200.gif",
    "https://media1.giphy.com/media/26BRBKqUiq586bRVm/200.gif",
  ],
  thumbsup: [
    "https://media1.giphy.com/media/111ebonMs90YLu/200.gif",
    "https://media1.giphy.com/media/l0HlBO7eyXzSZkJri/200.gif",
    "https://media1.giphy.com/media/3o7TKSjRrfIPjeiVyM/200.gif",
    "https://media1.giphy.com/media/l0MYC0LajbaPoEADu/200.gif",
    "https://media1.giphy.com/media/xT0xeJpnrWC3XWblEk/200.gif",
    "https://media1.giphy.com/media/26BRv0ThflsHCqDrG/200.gif",
    "https://media1.giphy.com/media/5GoVLqeAOo6PK/200.gif",
    "https://media1.giphy.com/media/26ufdipQqU2lhNA4g/200.gif",
    "https://media1.giphy.com/media/l4FGuhL4U2WyjdkaY/200.gif",
    "https://media1.giphy.com/media/3oEjI6SIIHBdRxXI40/200.gif",
    "https://media1.giphy.com/media/xT9IgG50Fb7Mi0prBC/200.gif",
    "https://media1.giphy.com/media/l0MYt5jPR6QX5pnqM/200.gif",
  ],
};

export function GifPicker({ onSelect }: GifPickerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [activeCategory, setActiveCategory] = useState("trending");
  const [gifQuery, setGifQuery] = useState("");
  const [gifs, setGifs] = useState<string[]>([]);

  const categories = Object.keys(GIF_LIBRARY);

  const searchGifs = (query: string) => {
    const term = query.trim().toLowerCase();
    if (!term) {
      setGifs(GIF_LIBRARY[activeCategory] || GIF_LIBRARY.trending);
    } else {
      const matchedCategory = categories.find((cat) => cat.includes(term) || term.includes(cat));
      if (matchedCategory) {
        setGifs(GIF_LIBRARY[matchedCategory]);
      } else {
        const allGifs = Object.values(GIF_LIBRARY).flat();
        const unique = Array.from(new Set(allGifs));
        setGifs(unique.slice(0, 12));
      }
    }
  };

  useEffect(() => {
    if (isOpen) {
      searchGifs(gifQuery);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, activeCategory, gifQuery]);

  return (
    <div style={{ position: "relative", display: "inline-block" }}>
      <button
        type="button"
        className="btn btn-ghost"
        onClick={() => setIsOpen(!isOpen)}
        style={{
          padding: "8px",
          fontSize: "13px",
          fontWeight: "600",
          borderRadius: "var(--radius-sm)",
          background: isOpen ? "var(--accent)" : "var(--bg-elevated)",
          color: isOpen ? "#fff" : "var(--accent)",
          border: "1px solid var(--border-light)",
          cursor: "pointer",
          display: "flex",
          alignItems: "center",
          gap: "4px",
          height: "36px",
        }}
        title="Insert GIF"
      >
        <span>👾 GIF</span>
      </button>
      {isOpen && (
        <>
          <div
            onClick={() => setIsOpen(false)}
            style={{ position: "fixed", top: 0, left: 0, right: 0, bottom: 0, zIndex: 999 }}
          />
          <div
            className="card"
            style={{
              position: "absolute",
              bottom: "44px",
              left: 0,
              zIndex: 1000,
              padding: "10px",
              width: "320px",
              maxHeight: "300px",
              boxShadow: "var(--shadow-lg)",
              background: "var(--bg-elevated)",
              border: "1px solid var(--border)",
              display: "flex",
              flexDirection: "column",
              overflow: "hidden",
              borderRadius: "12px",
            }}
          >
            <div style={{ display: "flex", gap: "4px", marginBottom: "8px", flexWrap: "wrap", maxHeight: "60px", overflowY: "auto" }}>
              {categories.map((cat) => (
                <button
                  key={cat}
                  type="button"
                  onClick={() => {
                    setActiveCategory(cat);
                    setGifQuery("");
                  }}
                  style={{
                    padding: "3px 8px",
                    fontSize: "11px",
                    borderRadius: "12px",
                    border: "1px solid var(--border)",
                    background: activeCategory === cat ? "var(--accent)" : "transparent",
                    color: activeCategory === cat ? "#fff" : "var(--text-muted)",
                    cursor: "pointer",
                    textTransform: "capitalize",
                  }}
                >
                  {cat}
                </button>
              ))}
            </div>
            <input
              type="text"
              placeholder="Search GIFs..."
              value={gifQuery}
              onChange={(e) => setGifQuery(e.target.value)}
              className="input"
              style={{ marginBottom: "8px", width: "100%", fontSize: "12px", padding: "6px 10px" }}
            />
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "repeat(3, 1fr)",
                gap: "6px",
                overflowY: "auto",
                flex: 1,
              }}
            >
              {gifs.map((gifUrl, idx) => (
                <img
                  key={idx}
                  src={gifUrl}
                  alt="gif selection option"
                  style={{ width: "100%", height: "64px", objectFit: "cover", borderRadius: "4px", cursor: "pointer" }}
                  onClick={() => {
                    onSelect(gifUrl);
                    setIsOpen(false);
                    setGifQuery("");
                  }}
                />
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
