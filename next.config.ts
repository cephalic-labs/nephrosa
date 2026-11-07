import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [{
      source: "/api/v1/sarvam/webhooks/:secret",
      destination: `${process.env.RELAY_API_URL ?? "http://127.0.0.1:8080"}/api/v1/sarvam/webhooks/:secret`,
    }];
  },
  /* config options here */
  experimental: {
    agentFeedback: true,
    useTypeScriptCli: false,
  },
  cacheComponents: true,
  partialPrefetching: true,
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
};

export default nextConfig;
