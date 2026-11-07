import { Suspense } from "react";
import Link from "next/link";
import { requireUser, apiFetch } from "../api-client";
import { ActionForm } from "../action-form";
import { SiteHeader } from "../site-header";
import { redirect } from "next/navigation";
import { io } from "next/cache";

type Episode = {
  id: string;
  patientReference: string;
  patientDisplayName: string;
  status: string;
  dueAt: string;
  ownerDisplayName: string;
  tests: string[];
};

type Overview = {
  summary: {
    openEpisodes: number;
    overdueEpisodes: number;
    awaitingReport: number;
    awaitingClinicianReview: number;
    dueWithinSevenDays: number;
  };
  episodes: Episode[];
};

async function getOverview(history: boolean): Promise<Overview | null> {

  try {
    const response = await apiFetch(`/dashboard/overview?history=${history}`);
    if (!response.ok) return null;
    return (await response.json()) as Overview;
  } catch {
    return null;
  }
}

const metrics = [
  ["Open episodes", "openEpisodes"],
  ["Overdue", "overdueEpisodes"],
  ["Awaiting report", "awaitingReport"],
  ["Clinician review", "awaitingClinicianReview"],
  ["Due within seven days", "dueWithinSevenDays"],
] as const;

function readable(value: string) {
  return value.replaceAll("_", " ");
}

function dateLabel(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "Date not set";
  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  }).format(date);
}

export async function RoleDashboard({ searchParams, expectedRole }: { searchParams?: Promise<{ history?: string }>; expectedRole: "clinician" | "coordinator" }) {
  const user = await requireUser();
  if (user.role !== expectedRole) redirect(`/dashboard/${user.role}`);
  const history = (await searchParams)?.history === "true";
  const data = await getOverview(history);
  const connected = data !== null;
  const episodes = data?.episodes ?? [];
  await io();
  const today = new Date().toISOString().slice(0, 10);

  const clinician = user.role === "clinician";
  const rolePath = `/dashboard/${user.role}`;
  const urgent = episodes.filter(episode => episode.dueAt.slice(0,10) < today).length;
  return (
    <main className="dashboard-shell">
      <SiteHeader connected={connected} role={user.role} />
      <div className="workspace-toolbar"><span>{user.name} · {clinician ? "Clinician" : "Care coordinator"}</span>{clinician && <Link className="retry-button" href="/episodes/new">Create episode</Link>}<Link href={history ? rolePath : `${rolePath}?history=true`}>{history ? "Open episodes" : "Include closed episodes"}</Link><ActionForm operation="logout" label="Sign out" /></div>

      <section className="role-welcome" id="dashboard" aria-labelledby="dashboard-title">
        <div><p className="hero-kicker">{clinician ? "Clinical review workspace" : "Care coordination workspace"}</p><h1 id="dashboard-title">{clinician ? "Your clinic, at a glance." : "Your follow-through queue."}</h1><p>{clinician ? "Review monitoring progress, record next actions, and keep orders moving." : "Start with the people who need a call, a test, or a handoff."}</p></div>
        <aside className="today-focus"><span>{history ? "History view" : "Today’s focus"}</span><strong>{urgent}</strong><p>{urgent === 1 ? "episode is past its due date" : "episodes are past their due date"}</p></aside>
      </section>

      <section className="overview-panel" aria-labelledby="overview-title">
        <div className="overview-heading">
          <h2 id="overview-title">Team overview</h2>
          <span className="overview-context">{connected ? "Current workload" : "Awaiting live data"}</span>
        </div>
        <div className="metrics-grid">
          {metrics.map(([label, key]) => (
            <div className="metric" key={key}>
              <strong className="metric-value">{connected ? data.summary[key] : "—"}</strong>
              <span className="metric-label">{label}</span>
            </div>
          ))}
        </div>
      </section>

      <section className="episodes-panel" id="episodes" aria-labelledby="episodes-title">
        <div className="episodes-heading">
          <div>
            <h2 id="episodes-title">Monitoring episodes</h2>
            <p>{clinician ? "Orders needing clinical attention and decisions." : "Patients needing practical help to complete ordered testing."}</p>
          </div>
          <span className="episode-total">{connected ? episodes.length : "—"}<span> displayed</span></span>
        </div>
        {!connected ? (
          <div className="episode-state" role="status">
            <strong>Live episodes couldn’t be loaded.</strong>
            <span>Check the dashboard service, then try again.</span>
            <form action={rolePath} method="get">
              <button className="retry-button" type="submit">Try again</button>
            </form>
          </div>
        ) : episodes.length === 0 ? (
          <div className="episode-state">
            <strong>No open monitoring episodes.</strong>
            <span>New episodes will appear here when they are recorded.</span>
          </div>
        ) : (
          <div className="episode-list" role="list" aria-label="Open monitoring episodes"><div className="episode-columns" aria-hidden="true"><span>Patient and reference</span><span>Ordered tests</span><span>Current step</span><span>Due date</span><span>Care owner</span><span/></div>
            {episodes.map((episode) => {
              const overdue = episode.dueAt.slice(0, 10) < today;
              return (
                <article className="episode-row" key={episode.id} role="listitem">
                  <div className="episode-person">
                    <Link href={`/episodes/${episode.id}`} aria-label={`Open episode details for ${episode.patientDisplayName}`}>
                    <strong>{episode.patientDisplayName}</strong>
                    </Link>
                    <span>{episode.patientReference}</span>
                  </div>
                  <div className="episode-tests">{episode.tests.join(" · ")}</div>
                  <div className="episode-status">
                    <span className={overdue ? "status-dot is-overdue" : "status-dot"} />
                    {overdue ? "Overdue · " : ""}{readable(episode.status)}
                  </div>
                  <div className="episode-due">Due {dateLabel(episode.dueAt)}</div>
                  <div className="episode-owner">{episode.ownerDisplayName}</div>
                  <Link className="row-open" href={`/episodes/${episode.id}`} aria-label={`Open episode for ${episode.patientDisplayName}`}>Open record <span aria-hidden="true">↗</span></Link>
                </article>
              );
            })}
          </div>
        )}
      </section>

      <footer className="site-footer">
        <span>Renal Relay</span>
        <span>Care teams retain responsibility for clinical review and next-action decisions.</span>
      </footer>
    </main>
  );
}

async function DashboardRedirect({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  const user = await requireUser();
  const params = await searchParams;
  redirect(`/dashboard/${user.role}${params.history === "true" ? "?history=true" : ""}`);
  return null;
}

export default function Home({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  return <Suspense fallback={<main className="loading-screen" role="status">Opening your dashboard…</main>}><DashboardRedirect searchParams={searchParams}/></Suspense>;
}
