import { Suspense } from "react";
import { RoleDashboard } from "../role-dashboard";

export default function ClinicianDashboard({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  return <Suspense fallback={<main className="loading-screen" role="status">Loading clinician dashboard…</main>}><RoleDashboard expectedRole="clinician" searchParams={searchParams}/></Suspense>;
}
