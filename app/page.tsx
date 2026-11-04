import { Suspense } from "react";
import Link from "next/link";

type Episode = {
  id: string;
  patientReference: string;
  patientDisplayName: string;
  status: string;
  dueAt: string;
  ownerDisplayName: string;
  tests: string[];
  openBarrierCategory?: string;
  nextTask?: string;
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

function SplitMark() {
  return (
    <svg className="split-mark" viewBox="0 0 120 120" aria-hidden="true">
      <defs>
        <pattern id="relay-dots" width="7" height="7" patternUnits="userSpaceOnUse">
          <circle cx="1.5" cy="1.5" r="1.3" fill="currentColor" />
        </pattern>
      </defs>
      <path d="M0 0h120v120H0z" fill="url(#relay-dots)" />
      <path d="M0 0h57v57H0zM63 0h57v57H63zM0 63h57v57H0z" fill="var(--ink)" />
    </svg>
  );
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
          <span className="wordmark-symbol" aria-hidden="true"><span /><span /><span /><span /></span>
          <span>renal relay</span>
        </Link>
        <nav className="primary-nav" aria-label="Main navigation">
          <a href="#dashboard" aria-current="page">Dashboard</a>
          <a href="#workflow">Workflow</a>
          <a href="#episodes">Episodes</a>
        </nav>
        <div className={`connection-state ${connected ? "is-connected" : "is-offline"}`} role="status">
          <span className="connection-dot" />
          {connected ? "Live connection" : "Connection unavailable"}
        </div>
      </header>

      <section className="split-grid" id="dashboard" aria-label="Care team dashboard">
        <div className="grid-cell introduction">
          <h1>Care continues after the order.</h1>
          <p>Keep kidney monitoring moving from uACR and eGFR orders to clinician acknowledgement and a documented next action.</p>
          <span className="intro-rule" aria-hidden="true" />
        </div>

        <section className="grid-cell workflow-panel" id="workflow" aria-labelledby="workflow-title">
          <div className="workflow-topline">
            <h2 id="workflow-title">The follow-through</h2>
            <SplitMark />
          </div>
          <p className="workflow-copy">One visible path from ordered test to recorded next step.</p>
          <ol className="workflow-list">
            {workflow.map((step, index) => (
              <li key={step}>
                <span className="workflow-index">0{index + 1}</span>
                <span>{step}</span>
                {index < workflow.length - 1 && <span className="workflow-arrow" aria-hidden="true" />}
              </li>
            ))}
          </ol>
        </section>

        <section className="grid-cell metrics-panel" aria-label="Episode overview">
          <h2>Team overview</h2>
          <div className="metrics-grid">
            {metrics.map(([label, key]) => (
              <div className="metric" key={key}>
                <span className="metric-label">{label}</span>
                <strong className="metric-value">{connected ? data.summary[key] : "—"}</strong>
              </div>
            ))}
          </div>
          <p className="metric-note">{connected ? `${data.summary.dueWithinSevenDays} due within seven days` : "Counts appear when the dashboard service is available."}</p>
        </section>

        <section className="grid-cell episodes-panel" id="episodes" aria-labelledby="episodes-title">
          <div className="episodes-heading">
            <div>
              <h2 id="episodes-title">Monitoring episodes</h2>
              <p>Current follow-through across the care team.</p>
            </div>
            <span className="episode-total">{connected ? episodes.length : "—"} <span>open</span></span>
          </div>
          {!connected ? (
            <div className="episode-empty" role="status">
              <strong>Live episodes couldn’t be loaded.</strong>
              <span>Check the service connection and try again.</span>
              <form action="/" method="get">
                <button className="retry-button" type="submit">Try again</button>
              </form>
            </div>
          ) : episodes.length === 0 ? (
            <div className="episode-empty">
              <strong>No open monitoring episodes.</strong>
              <span>New episodes will appear here when they are recorded.</span>
            </div>
          ) : (
            <div className="episode-list">
              {episodes.map((episode) => {
                const overdue = episode.dueAt.slice(0, 10) < today;
                return (
                  <article className="episode-row" key={episode.id}>
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
