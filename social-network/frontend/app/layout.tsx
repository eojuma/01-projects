import type { Metadata } from "next";
import "./globals.css";
import { WebSocketProvider } from "./context/WebSocketContext";
import { AuthProvider } from "./context/AuthContext";
import { NotificationProvider } from "./context/NotificationContext";

export const metadata: Metadata = {
  title: "SocialNet",
  description: "Connect with friends and communities",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <WebSocketProvider>
            <NotificationProvider>
              {children}
            </NotificationProvider>
          </WebSocketProvider>
        </AuthProvider>
      </body>
    </html>
  );
}