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

const demoOverview: Overview = {
  summary: {
    openEpisodes: 2,
    overdueEpisodes: 1,
    awaitingReport: 0,
    awaitingClinicianReview: 1,
    dueWithinSevenDays: 1,
  },
  episodes: [
    {
      id: "demo-episode-1",
      patientReference: "RR-1042",
      patientDisplayName: "Maya R.",
      status: "awaiting_booking",
      dueAt: relativeDate(-5),
      ownerDisplayName: "A. Patel",
      tests: ["uACR", "eGFR"],
      openBarrierCategory: "time_off_work",
      nextTask: "Offer an appointment outside work hours",
    },
    {
      id: "demo-episode-2",
      patientReference: "RR-1078",
      patientDisplayName: "Arun K.",
      status: "awaiting_clinician_review",
      dueAt: relativeDate(2),
      ownerDisplayName: "J. Thomas",
      tests: ["uACR", "eGFR"],
      nextTask: "Review received report and record next action",
    },
  ],
};

function relativeDate(offsetDays: number) {
  const date = new Date();
  date.setUTCDate(date.getUTCDate() + offsetDays);
  return date.toISOString().slice(0, 10);
}

async function getOverview(): Promise<Overview> {
  const apiUrl = process.env.RELAY_API_URL ?? "http://127.0.0.1:8080";

  try {
    const response = await fetch(`${apiUrl}/api/v1/dashboard/overview`, {
      cache: "no-store",
      signal: AbortSignal.timeout(1500),
    });
    if (!response.ok) throw new Error("Dashboard service unavailable");
    return (await response.json()) as Overview;
  } catch {
    return demoOverview;
  }
}

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    ordered: "Ordered",
    outreach_pending: "Outreach pending",
    awaiting_booking: "Booking needed",
    booked: "Booked",
    awaiting_test_completion: "Awaiting test",
    awaiting_report: "Awaiting report",
    awaiting_clinician_review: "Clinician review",
  };
  return labels[status] ?? status.replaceAll("_", " ");
}

function testLabel(test: string) {
  if (test.toLowerCase() === "uacr") return "uACR";
  if (test.toLowerCase() === "egfr") return "eGFR";
  return test;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  }).format(new Date(value));
}

function Metric({ label, value, detail, tone = "default" }: { label: string; value: number; detail: string; tone?: "default" | "urgent" }) {
  return (
    <article className={`metric metric-${tone}`}>
      <p className="metric-label">{label}</p>
      <p className="metric-value">{value}</p>
      <p className="metric-detail">{detail}</p>
    </article>
  );
}

async function Dashboard() {
  const data = await getOverview();
  const todayISO = new Date().toISOString().slice(0, 10);
  const today = new Intl.DateTimeFormat("en", {
    weekday: "long",
    month: "long",
    day: "numeric",
  }).format(new Date());

  return (
    <main className="workspace">
      <aside className="sidebar" aria-label="Main navigation">
        <Link className="brand" href="/" aria-label="Renal Relay home">
          <span className="brand-mark" aria-hidden="true">R</span>
          <span className="brand-name">renal<span>relay</span></span>
        </Link>
        <p className="nav-heading">WORKSPACE</p>
        <nav className="main-nav">
          <a className="nav-item active" href="#overview" aria-current="page">
            <svg viewBox="0 0 20 20" aria-hidden="true"><rect x="3" y="3" width="5" height="5" rx="1"/><rect x="12" y="3" width="5" height="5" rx="1"/><rect x="3" y="12" width="5" height="5" rx="1"/><rect x="12" y="12" width="5" height="5" rx="1"/></svg>
            Overview
          </a>
          <a className="nav-item" href="#episodes">
            <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M5 4.5h10M5 9.5h10M5 14.5h6"/><path d="M3.5 4.5h.01M3.5 9.5h.01M3.5 14.5h.01"/></svg>
            Monitoring episodes <span className="nav-count">{data.summary.openEpisodes}</span>
          </a>
          <a className="nav-item" href="#follow-through">
            <svg viewBox="0 0 20 20" aria-hidden="true"><circle cx="10" cy="10" r="7"/><path d="M10 6v4l2.5 1.5"/></svg>
            Follow-through
          </a>
        </nav>
        <div className="sidebar-bottom">
          <div className="clinic-switcher"><span className="clinic-avatar">C</span><span><strong>Care team</strong><small>Prototype workspace</small></span><span className="chevron">⌄</span></div>
          <p className="sidebar-footnote">Fictional records for demonstration</p>
        </div>
      </aside>

      <section className="main-panel" id="overview">
        <header className="topbar">
          <div className="breadcrumb">Workspace <span>/</span> Overview</div>
          <div className="topbar-meta"><span className="today-date">{today}</span><span className="avatar" aria-label="Signed in as Alex Patel">AP</span></div>
        </header>

        <div className="content">
          <div className="page-heading">
            <div>
              <h1>Good morning, Alex</h1>
              <p>Keep ordered kidney monitoring moving toward a documented next step.</p>
            </div>
            <button className="button-primary" type="button" disabled title="Episode creation is not available in this prototype">＋ <span>New episode</span></button>
          </div>

          <div className="demo-notice" role="status"><span className="notice-dot" /> Prototype workspace · Use fictional records only.</div>

          <section className="metrics" aria-label="Episode overview">
            <Metric label="Open episodes" value={data.summary.openEpisodes} detail="Across your care team" />
            <Metric label="Overdue" value={data.summary.overdueEpisodes} detail="Needs follow-through" tone="urgent" />
            <Metric label="Awaiting report" value={data.summary.awaitingReport} detail="Completed test, report not logged" />
            <Metric label="Clinician review" value={data.summary.awaitingClinicianReview} detail="Next action to document" />
          </section>

          <section className="episode-section" id="episodes" aria-labelledby="episode-title">
            <div className="section-heading">
              <div><h2 id="episode-title">Monitoring episodes</h2><p>Ordered tests, open barriers, and the next accountable step.</p></div>
            </div>
            <div className="table-wrap">
              <table>
                <thead><tr><th scope="col">PATIENT</th><th scope="col">MONITORING</th><th scope="col">STATUS</th><th scope="col">DUE DATE</th><th scope="col">OWNER</th></tr></thead>
                <tbody>
                  {data.episodes.map((episode) => (
                    <tr key={episode.id}>
                      <td><div className="patient-cell"><span className="patient-avatar">{episode.patientDisplayName.slice(0, 1)}</span><span><strong>{episode.patientDisplayName}</strong><small>{episode.patientReference}</small></span></div></td>
                      <td><div className="test-list">{episode.tests.map((test) => <span className="test-tag" key={test}>{testLabel(test)}</span>)}</div></td>
                      <td><span className={`status-pill ${episode.dueAt.slice(0, 10) < todayISO ? "status-overdue" : episode.status === "awaiting_clinician_review" ? "status-review" : "status-pending"}`}><span />{episode.dueAt.slice(0, 10) < todayISO ? `Overdue · ${statusLabel(episode.status)}` : statusLabel(episode.status)}</span></td>
                      <td className="date-cell">{formatDate(episode.dueAt)}</td>
                      <td className="owner-cell">{episode.ownerDisplayName}</td>
                    </tr>
                  ))}
                  {data.episodes.length === 0 && <tr><td className="empty-state" colSpan={5}>No open monitoring episodes are assigned to this team.</td></tr>}
                </tbody>
              </table>
            </div>
          </section>

          <section className="follow-through" id="follow-through" aria-labelledby="follow-title">
            <div className="follow-heading"><div><h2 id="follow-title">Next steps to keep moving</h2><p>Open work linked to a monitoring episode.</p></div><span className="follow-count">{data.episodes.length} active</span></div>
            <div className="task-list">
              {data.episodes.slice(0, 3).map((episode) => (
                <article className="task-row" key={`${episode.id}-task`}>
                  <span className="task-check" aria-hidden="true" />
                  <div className="task-copy"><strong>{episode.nextTask || `Continue follow-through for ${episode.patientDisplayName}`}</strong><span>{episode.patientDisplayName} <span aria-hidden="true">·</span> {episode.openBarrierCategory ? `Barrier: ${episode.openBarrierCategory.replaceAll("_", " ")}` : statusLabel(episode.status)}</span></div>
                  <span className="task-owner">{episode.ownerDisplayName}</span>
                </article>
              ))}
              {data.episodes.length === 0 && <p className="empty-tasks">When an episode needs follow-through, its next step will appear here.</p>}
            </div>
          </section>
          <footer className="page-footer">Renal Relay prototype <span>·</span> Care teams retain clinical review and next-action authority.</footer>
        </div>
      </section>
    </main>
  );
}

export default function Home() {
  return (
    <Suspense fallback={<main className="workspace loading-state" role="status">Loading the care-team overview…</main>}>
      <Dashboard />
    </Suspense>
  );
}
