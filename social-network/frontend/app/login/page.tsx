"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "../context/AuthContext";

export default function Login() {
  const router = useRouter();
  const { login } = useAuth();
  const [emailOrUsername, setEmailOrUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(emailOrUsername, password);
      router.push("/");
    } catch (err: any) {
      setError(err.message || "Invalid credentials");
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
        <h1 className="auth-title">Welcome back</h1>
        <p className="auth-subtitle">Sign in to your account</p>

        {error && (
          <div className="alert-error" style={{ marginBottom: "16px" }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Email or Username</label>
            <input
              id="login-email"
              type="text"
              className="input"
              value={emailOrUsername}
              onChange={(e) => setEmailOrUsername(e.target.value)}
              placeholder="Enter email or username"
              required
              autoFocus
            />
          </div>
          <div className="form-group">
            <label className="form-label">Password</label>
            <input
              id="login-password"
              type="password"
              className="input"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
            />
          </div>
          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
            style={{ width: "100%", justifyContent: "center", marginTop: "8px" }}
          >
            {loading ? "Signing in…" : "Sign In"}
          </button>
        </form>

        <p className="auth-footer">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="auth-link">
            Sign up
          </Link>
        </p>
      </div>
    </div>
  );
}