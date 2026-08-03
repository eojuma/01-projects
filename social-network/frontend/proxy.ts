import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function proxy(request: NextRequest) {
  const token = request.cookies.get("session_id");
  const { pathname } = request.nextUrl;

  // Public paths
  const publicPaths = ["/login", "/register"];
  const isPublicPath = publicPaths.some(path => pathname.startsWith(path));

  // If no token and trying to access protected route, redirect to login
  if (!token && !isPublicPath) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // If has token and trying to access public path, redirect to home
  if (token && isPublicPath) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};