"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useNotificationContext } from "../context/NotificationContext";

const icons: Record<string, React.ReactNode> = {
  feed: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  ),
  profile: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  ),
  groups: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  ),
  notifications: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
      <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0" />
    </svg>
  ),
  messages: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
    </svg>
  ),
  users: (
    <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "18px", height: "18px" }}>
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  )
};

const navItems = [
  {
    section: "MAIN",
    links: [
      { href: "/", label: "Feed", icon: "feed" },
      { href: "/profile", label: "Profile", icon: "profile" },
      { href: "/groups", label: "Groups", icon: "groups" },
    ],
  },
  {
    section: "SOCIAL",
    links: [
      { href: "/users", label: "Find Friends", icon: "users" },
      { href: "/notifications", label: "Notifications", icon: "notifications", badge: true },
      { href: "/chat", label: "Messages", icon: "messages", badge: true },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const { unreadCount, unreadChatCount } = useNotificationContext();

  const handleLogout = async () => {
    try {
      await fetch("/api/auth/logout", {
        method: "POST",
        credentials: "include",
      });
      router.push("/login");
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        Social<span>Net</span>
      </div>

      {navItems.map(({ section, links }) => (
        <div key={section}>
          <p className="sidebar-section-label">{section}</p>
          {links.map(({ href, label, icon, badge }) => {
            const isActive = pathname === href;
            const badgeVal = href === "/notifications" ? unreadCount : (href === "/chat" ? unreadChatCount : 0);
            return (
              <Link
                key={href}
                href={href}
                className={`sidebar-link ${isActive ? "active" : ""}`}
              >
                {icons[icon]}
                <span>{label}</span>
                {badge && badgeVal > 0 && (
                  <span className="badge">{badgeVal > 99 ? "99+" : badgeVal}</span>
                )}
              </Link>
            );
          })}
        </div>
      ))}

      <div className="sidebar-footer">
        <button
          onClick={handleLogout}
          className="btn btn-danger"
          style={{ width: "100%", justifyContent: "flex-start", padding: "10px 12px" }}
        >
          <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ width: "16px", height: "16px", marginRight: "8px" }}>
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>
          Sign Out
        </button>
      </div>
    </aside>
  );
}