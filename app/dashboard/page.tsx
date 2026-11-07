import { Suspense } from "react";
import { requireUser } from "../api-client";
import { redirect } from "next/navigation";

async function DashboardRedirect({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  const user = await requireUser();
  const params = await searchParams;
  redirect(`/dashboard/${user.role}${params.history === "true" ? "?history=true" : ""}`);
  return null;
}

export default function Home({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  return <Suspense fallback={<main className="loading-screen" role="status">Opening your dashboard…</main>}><DashboardRedirect searchParams={searchParams}/></Suspense>;
}
