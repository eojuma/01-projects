"use client";

import { useState, useMemo } from "react";

interface EmojiPickerProps {
  onSelect: (emoji: string) => void;
}

const EMOJI_CATEGORIES: Record<string, string[]> = {
  "😀 Smileys": ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🙂", "😊", "😇", "🥰", "😍", "🤩", "😘", "😋", "😛", "😜", "🤪", "😝", "🤗", "🤔", "🤭", "🤫", "😶", "😐", "😑", "😬", "😏", "😌", "😴", "🤤", "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "🥳", "😎", "🤓", "🧐", "😕", "😟", "🙁", "😮", "😯", "😲", "😳", "🥺", "😭", "😱", "😖", "😤", "😡", "🤬", "😈", "👿", "💀", "☠️", "💩", "🤡", "👹", "👻", "👽", "🤖"],
  "👋 Gestures": ["👋", "🤚", "✋", "🖐️", "👌", "🤌", "🤏", "✌️", "🤞", "🤟", "🤘", "🤙", "👈", "👉", "👆", "👇", "☝️", "👍", "👎", "✊", "👊", "🤛", "🤜", "👏", "🙌", "👐", "🤲", "🤝", "🙏", "💪", "🦾", "🖕"],
  "❤️ Hearts": ["❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍", "🤎", "💔", "❣️", "💕", "💞", "💓", "💗", "💖", "💘", "💝", "💟", "♥️"],
  "🐶 Animals": ["🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵", "🐔", "🐧", "🐦", "🦅", "🦆", "🦉", "🐴", "🦄", "🐝", "🐛", "🦋", "🐌", "🐙", "🐬", "🐳", "🐠", "🦈"],
  "🍕 Food": ["🍎", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓", "🫐", "🍒", "🍑", "🥭", "🍍", "🥝", "🍅", "🥑", "🌽", "🌶️", "🥕", "🧀", "🍕", "🍔", "🍟", "🌭", "🍿", "🧂", "🥚", "🍳", "🥞", "🧇", "🍰", "🎂", "☕"],
  "⚽ Activities": ["⚽", "🏀", "🏈", "⚾", "🥎", "🎾", "🏐", "🏉", "🥏", "🎱", "🏓", "🎯", "🎮", "🕹️", "🎲", "🧩", "🎭", "🎨", "🎤", "🎧", "🎼", "🎹", "🎸", "🎺", "🥁", "🎬", "🏆", "🥇", "🥈", "🥉", "🎪", "🎟️"],
  "🚀 Objects": ["🔥", "⭐", "🌟", "✨", "💫", "🎉", "🎊", "🎈", "💡", "🔔", "🔊", "📢", "💬", "💭", "🗯️", "📱", "💻", "🖥️", "📷", "📹", "📺", "📻", "⏰", "🔑", "🔒", "📌", "📎", "✂️", "📝", "📚", "🚀", "✈️"],
  "🏁 Symbols": ["✅", "❌", "❓", "❗", "‼️", "⁉️", "💯", "🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "⚫", "⚪", "🟤", "♻️", "✳️", "❇️", "🔆", "🔅", "⚜️", "🔱", "🏴", "🏳️", "🏁", "🚩", "🎌", "🏴‍☠️", "🇺🇸", "🇬🇧", "🇪🇺"],
};

export function EmojiPicker({ onSelect }: EmojiPickerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [activeCategory, setActiveCategory] = useState(Object.keys(EMOJI_CATEGORIES)[0]);

  const categoryKeys = Object.keys(EMOJI_CATEGORIES);
  const categoryIcon = (key: string) => key.split(" ")[0];

  const displayedEmojis = useMemo(() => {
    return EMOJI_CATEGORIES[activeCategory] || [];
  }, [activeCategory]);

  return (
    <div style={{ position: "relative", display: "inline-block" }}>
      <button
        type="button"
        className="btn btn-ghost"
        onClick={() => setIsOpen(!isOpen)}
        style={{ padding: "8px", fontSize: "16px", minWidth: "auto", height: "auto" }}
        title="Insert Emoji"
      >
        😀
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
              bottom: "40px",
              right: 0,
              zIndex: 1000,
              padding: "0",
              width: "320px",
              maxHeight: "360px",
              boxShadow: "var(--shadow-lg)",
              background: "var(--bg-elevated)",
              border: "1px solid var(--border)",
              display: "flex",
              flexDirection: "column",
              overflow: "hidden",
              borderRadius: "12px",
            }}
          >
            {/* Category tabs */}
            <div style={{
              display: "flex",
              borderBottom: "1px solid var(--border)",
              padding: "6px 6px 0",
              gap: "2px",
              overflowX: "auto",
              flexShrink: 0,
            }}>
              {categoryKeys.map((key) => (
                <button
                  key={key}
                  type="button"
                  onClick={() => setActiveCategory(key)}
                  title={key.split(" ").slice(1).join(" ")}
                  style={{
                    background: activeCategory === key ? "var(--bg-hover)" : "transparent",
                    border: "none",
                    borderBottom: activeCategory === key ? "2px solid var(--accent)" : "2px solid transparent",
                    padding: "6px 8px",
                    fontSize: "16px",
                    cursor: "pointer",
                    borderRadius: "6px 6px 0 0",
                    flexShrink: 0,
                  }}
                >
                  {categoryIcon(key)}
                </button>
              ))}
            </div>

            {/* Emoji grid */}
            <div style={{
              overflowY: "auto",
              padding: "8px",
              display: "grid",
              gridTemplateColumns: "repeat(8, 1fr)",
              gap: "2px",
              flex: 1,
            }}>
              {displayedEmojis.map((emoji, idx) => (
                <button
                  key={`${emoji}-${idx}`}
                  type="button"
                  onClick={() => {
                    onSelect(emoji);
                    setIsOpen(false);
                  }}
                  style={{
                    background: "transparent",
                    border: "none",
                    fontSize: "20px",
                    cursor: "pointer",
                    padding: "4px",
                    borderRadius: "6px",
                    transition: "background 0.15s",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    height: "34px",
                  }}
                  onMouseOver={(e) => (e.currentTarget.style.background = "var(--bg-hover)")}
                  onMouseOut={(e) => (e.currentTarget.style.background = "transparent")}
                >
                  {emoji}
                </button>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
