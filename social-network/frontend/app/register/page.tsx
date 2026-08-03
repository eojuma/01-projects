"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function Register() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
    first_name: "",
    last_name: "",
    date_of_birth: "",
    nickname: "",
    bio: "",
    avatar_url: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const response = await fetch("/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });
      if (response.ok) {
        router.push("/login");
      } else {
        const text = await response.text();
        setError(text || "Registration failed");
      }
    } catch {
      setError("Failed to connect to server");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-logo">
          Social<span>Net</span>
        </div>
        <h1 className="auth-title">Create account</h1>
        <p className="auth-subtitle">Join the community today</p>

        {error && (
          <div className="alert-error" style={{ marginBottom: "16px" }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Email *</label>
            <input
              id="reg-email"
              type="email"
              name="email"
              className="input"
              value={formData.email}
              onChange={handleChange}
              placeholder="you@example.com"
              required
              autoFocus
            />
          </div>

          <div className="form-group">
            <label className="form-label">Username *</label>
            <input
              id="reg-username"
              type="text"
              name="username"
              className="input"
              value={formData.username}
              onChange={handleChange}
              placeholder="Choose a username"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Password *</label>
            <input
              id="reg-password"
              type="password"
              name="password"
              className="input"
              value={formData.password}
              onChange={handleChange}
              placeholder="Minimum 8 characters"
              required
              minLength={8}
            />
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px" }}>
            <div className="form-group">
              <label className="form-label">First Name</label>
              <input
                id="reg-first-name"
                type="text"
                name="first_name"
                className="input"
                value={formData.first_name}
                onChange={handleChange}
                placeholder="First"
              />
            </div>
            <div className="form-group">
              <label className="form-label">Last Name</label>
              <input
                id="reg-last-name"
                type="text"
                name="last_name"
                className="input"
                value={formData.last_name}
                onChange={handleChange}
                placeholder="Last"
              />
            </div>
          </div>

          <div className="form-group">
            <label className="form-label">Date of Birth</label>
            <input
              id="reg-dob"
              type="date"
              name="date_of_birth"
              className="input"
              value={formData.date_of_birth}
              onChange={handleChange}
            />
          </div>

          <div className="form-group">
            <label className="form-label">Nickname (optional)</label>
            <input
              id="reg-nickname"
              type="text"
              name="nickname"
              className="input"
              value={formData.nickname}
              onChange={handleChange}
              placeholder="What do people call you?"
            />
          </div>

          <div className="form-group">
            <label className="form-label">Choose an Avatar (optional)</label>
            <div style={{ display: "grid", gridTemplateColumns: "repeat(6, 1fr)", gap: "8px", marginBottom: "8px" }}>
              {[
                "https://api.dicebear.com/7.x/bottts/svg?seed=Felix",
                "https://api.dicebear.com/7.x/bottts/svg?seed=Aneka",
                "https://api.dicebear.com/7.x/bottts/svg?seed=Jack",
                "https://api.dicebear.com/7.x/bottts/svg?seed=Boots",
                "https://api.dicebear.com/7.x/bottts/svg?seed=Coco",
                "https://api.dicebear.com/7.x/bottts/svg?seed=Zoe",
              ].map((url) => (
                <img
                  key={url}
                  src={url}
                  alt="Avatar option"
                  onClick={() => setFormData({ ...formData, avatar_url: url })}
                  style={{
                    width: "100%",
                    aspectRatio: "1",
                    borderRadius: "50%",
                    cursor: "pointer",
                    border: formData.avatar_url === url ? "3px solid #4f8ef7" : "2px solid rgba(255,255,255,0.1)",
                    objectFit: "cover",
                    transition: "border-color 0.2s, transform 0.2s",
                    transform: formData.avatar_url === url ? "scale(1.1)" : "scale(1)",
                    background: "rgba(255,255,255,0.05)",
                  }}
                />
              ))}
            </div>
            {formData.avatar_url && (
              <div style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "12px" }}>
                <span style={{ color: "#4f8ef7" }}>✓ Avatar selected</span>
                <button
                  type="button"
                  onClick={() => setFormData({ ...formData, avatar_url: "" })}
                  style={{ background: "none", border: "none", color: "rgba(255,255,255,0.5)", cursor: "pointer", fontSize: "11px", textDecoration: "underline" }}
                >
                  Clear
                </button>
              </div>
            )}
          </div>

          <div className="form-group">
            <label className="form-label">About Me (optional)</label>
            <input
              id="reg-bio"
              type="text"
              name="bio"
              className="input"
              value={formData.bio}
              onChange={handleChange}
              placeholder="Tell us about yourself..."
            />
          </div>

          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
            style={{ width: "100%", justifyContent: "center", marginTop: "8px" }}
          >
            {loading ? "Creating account…" : "Create Account"}
          </button>
        </form>

        <p className="auth-footer">
          Already have an account?{" "}
          <Link href="/login" className="auth-link">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}