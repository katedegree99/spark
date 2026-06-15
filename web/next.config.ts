import type { NextConfig } from "next";

const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001";

const nextConfig: NextConfig = {
	async rewrites() {
		return [
			{
				source: "/auth/:path*",
				destination: `${apiBaseUrl}/auth/:path*`,
			},
		];
	},
};

export default nextConfig;
