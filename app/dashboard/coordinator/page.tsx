import { Suspense } from "react";
import { RoleDashboard } from "../role-dashboard";

export default function CoordinatorDashboard({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  return <Suspense fallback={<main className="loading-screen" role="status">Loading coordination dashboard…</main>}><RoleDashboard expectedRole="coordinator" searchParams={searchParams}/></Suspense>;
}
