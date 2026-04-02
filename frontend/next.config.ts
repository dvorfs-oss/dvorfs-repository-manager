import type { NextConfig } from "next";

const backendUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${backendUrl}/api/v1/:path*`,
      },
      {
        source: "/repository/:path*",
        destination: `${backendUrl}/repository/:path*`,
      },
      {
        source: "/swagger/:path*",
        destination: `${backendUrl}/swagger/:path*`,
      },
    ];
  },
};

export default nextConfig;
