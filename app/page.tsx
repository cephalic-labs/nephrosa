import { Suspense } from "react";
import Link from "next/link";
import { requireUser, apiFetch } from "./api-client";
import { ActionForm } from "./action-form";
import { SiteHeader } from "./site-header";

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

const workflow = [
  "Order placed",
  "Test completed",
  "Report received",
  "Next action recorded",
];

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

async function Dashboard({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  const user = await requireUser();
  const history = (await searchParams).history === "true";
  const data = await getOverview(history);
  const connected = data !== null;
  const episodes = data?.episodes ?? [];
  const today = new Date().toISOString().slice(0, 10);

  return (
    <main className="dashboard-shell">
      <SiteHeader connected={connected} />
      <div className="workspace-toolbar"><span>{user.name} · {user.role}</span>{user.role === "clinician" && <Link className="retry-button" href="/episodes/new">Create episode</Link>}<Link href={history ? "/" : "/?history=true"}>{history ? "Open episodes" : "Include closed episodes"}</Link><ActionForm operation="logout" label="Sign out" /></div>

      <section className="hero" id="dashboard" aria-labelledby="dashboard-title">
        <div className="hero-copy">
          <h1 id="dashboard-title">Care continues after the order.</h1>
          <p>Keep uACR and eGFR monitoring moving—from the clinician’s order through report review and a documented next action.</p>
        </div>
        <section className="workflow-panel" id="workflow" aria-labelledby="workflow-title">
          <h2 id="workflow-title">The follow-through</h2>
          <p>One clear path for every monitoring episode.</p>
          <ol className="workflow-list">
            {workflow.map((step, index) => (
              <li key={step}>
                <span className="workflow-index">0{index + 1}</span>
                <span className="workflow-label">{step}</span>
              </li>
            ))}
          </ol>
        </section>
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
            <p>Current follow-through across the care team.</p>
          </div>
          <span className="episode-total">{connected ? episodes.length : "—"}<span> displayed</span></span>
        </div>
        {!connected ? (
          <div className="episode-state" role="status">
            <strong>Live episodes couldn’t be loaded.</strong>
            <span>Check the dashboard service, then try again.</span>
            <form action="/" method="get">
              <button className="retry-button" type="submit">Try again</button>
            </form>
          </div>
        ) : episodes.length === 0 ? (
          <div className="episode-state">
            <strong>No open monitoring episodes.</strong>
            <span>New episodes will appear here when they are recorded.</span>
          </div>
        ) : (
          <div className="episode-list" role="list" aria-label="Open monitoring episodes">
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

export default function Home({ searchParams }: { searchParams: Promise<{ history?: string }> }) {
  return (
    <Suspense fallback={<main className="loading-screen" role="status">Loading dashboard…</main>}>
      <Dashboard searchParams={searchParams} />
    </Suspense>
  );
}
