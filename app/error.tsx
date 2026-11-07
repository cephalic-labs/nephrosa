"use client";
import Link from "next/link";
export default function ErrorPage({ reset }: { reset: () => void }) {
  return <main className="dashboard-shell"><section className="workspace-panel"><h1>Care-team service unavailable</h1><p>The request could not be completed. Retry when the service is available.</p><button className="retry-button" onClick={reset}>Try again</button><p><Link href="/login">Return to sign in</Link></p></section></main>;
}
