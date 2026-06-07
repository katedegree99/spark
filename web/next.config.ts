import type { NextConfig } from "next";

const apiBaseUrl = process.env.API_BASE_URL ?? "http://localhost:3001";

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
