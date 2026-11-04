import { Suspense } from "react";
import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";

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

async function getOverview(): Promise<Overview | null> {
  const apiUrl = process.env.RELAY_API_URL ?? "http://127.0.0.1:8080";

  try {
    const response = await fetch(`${apiUrl}/api/v1/dashboard/overview`, {
      cache: "no-store",
      signal: AbortSignal.timeout(1500),
    });
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

async function Dashboard() {
  const data = await getOverview();
  const connected = data !== null;
  const episodes = data?.episodes ?? [];
  const today = new Date().toISOString().slice(0, 10);

  return (
    <main className="dashboard-shell">
      <header className="site-header">
        <Link className="wordmark" href="/" aria-label="Renal Relay dashboard">
          <span className="wordmark-symbol" aria-hidden="true"><span /><span /><span /></span>
          <span>renal relay</span>
        </Link>
        <nav className="primary-nav" aria-label="Main navigation">
          <a href="#dashboard" aria-current="page">Dashboard</a>
          <a href="#workflow">Workflow</a>
          <a href="#episodes">Episodes</a>
        </nav>
        <ThemeToggle />
        <div className={`connection-state ${connected ? "is-connected" : "is-offline"}`} role="status">
          <span className="connection-dot" />
          {connected ? "Live connection" : "Connection unavailable"}
        </div>
      </header>

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
          <span className="episode-total">{connected ? episodes.length : "—"}<span> open</span></span>
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
                    <strong>{episode.patientDisplayName}</strong>
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

export default function Home() {
  return (
    <Suspense fallback={<main className="loading-screen" role="status">Loading dashboard…</main>}>
      <Dashboard />
    </Suspense>
  );
}
