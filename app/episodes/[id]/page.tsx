import Link from "next/link";
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { Suspense } from "react";
import type { ReactNode } from "react";
import { requireUser, apiFetch } from "../../api-client";
import { EpisodeControls } from "../../episode-controls";
import { SiteHeader } from "../../site-header";

export const metadata: Metadata = {
  title: "Episode details | Renal Relay",
  description: "Follow-through details for a Renal Relay monitoring episode.",
};

type EpisodeDetail = {
  id: string;
  patientReference: string;
  patientDisplayName: string;
  status: string;
  consent: string;
  language: string;
  phone: string;
  dueAt: string;
  ownerDisplayName: string;
  tests: {
    id: string;
    testType: string;
    orderedAt: string;
    dueAt: string;
    completion?: { completedAt: string; verificationSource: string; recordedBy: string };
  }[];
  contacts: { channel: string; outcome: string; language: string; attemptedAt: string; recordedBy: string }[];
  barriers: { category: string; description: string; status: string; createdAt: string; resolvedAt?: string }[];
  tasks: {
    id: string;
    type: string;
    summary: string;
    priority: string;
    status: string;
    dueAt?: string;
    assignedToDisplayName?: string;
  }[];
  reports: { receivedAt: string; reference: string; recordedBy: string }[];
  clinicalReviews: { acknowledgedAt: string; clinician: string; nextAction: string }[];
  auditEvents: { eventType: string; occurredAt: string; actor: string; note?: string }[];
};

type LoadResult = { detail: EpisodeDetail; connected: true } | { detail: null; connected: false };

async function getEpisode(id: string): Promise<LoadResult | null> {
  try {
    const response = await apiFetch(`/episodes/${encodeURIComponent(id)}`);
    if (response.status === 404) return null;
    if (!response.ok) return { detail: null, connected: false };
    return { detail: (await response.json()) as EpisodeDetail, connected: true };
  } catch {
    return { detail: null, connected: false };
  }
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

function readable(value: string) {
  return value.replaceAll("_", " ");
}

function testLabel(value: string) {
  return value.toLowerCase() === "egfr" ? "eGFR" : "uACR";
}

function PageSection({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="detail-section">
      <h2>{title}</h2>
      {children}
    </section>
  );
}

async function EpisodePage({ params }: { params: Promise<{ id: string }> }) {
  const user = await requireUser();
  const { id } = await params;
  const result = await getEpisode(id);
  if (result === null) notFound();
  const detail = result.detail;

  return (
    <main className="dashboard-shell">
      <SiteHeader active="episodes" connected={result.connected} />
      <section className="episode-detail" aria-labelledby="episode-title">
        <Link className="back-link" href="/#episodes">← All episodes</Link>
        {detail ? (
          <>
            <div className="detail-heading">
              <div>
                <h1 id="episode-title">{detail.patientDisplayName}</h1>
                <p>{detail.patientReference} <span aria-hidden="true">·</span> Due {dateLabel(detail.dueAt)}</p>
              </div>
              <span className="detail-status">{readable(detail.status)}</span>
            </div>
            <p className="detail-owner">Care owner <strong>{detail.ownerDisplayName}</strong></p>

            <EpisodeControls id={id} closed={detail.status === "closed"} role={user.role} consent={detail.consent} phone={detail.phone} tests={detail.tests} tasks={detail.tasks} />
            <div className="detail-grid">
              <div className="detail-column">
                <PageSection title="Ordered tests">
                  {detail.tests.length ? (
                    <ul className="detail-list">
                      {detail.tests.map((test) => (
                        <li key={test.id}>
                          <div className="detail-item-heading">
                            <strong>{testLabel(test.testType)}</strong>
                            <span>{test.completion ? "Completion verified" : "Awaiting completion"}</span>
                          </div>
                          <p>Ordered {dateLabel(test.orderedAt)} <span aria-hidden="true">·</span> Due {dateLabel(test.dueAt)}</p>
                          {test.completion && <p>Verified {dateLabel(test.completion.completedAt)} by {test.completion.recordedBy} <span aria-hidden="true">·</span> {readable(test.completion.verificationSource)}</p>}
                        </li>
                      ))}
                    </ul>
                  ) : <p className="detail-empty">No test orders recorded.</p>}
                </PageSection>

                <PageSection title="Coordinator follow-through">
                  {detail.tasks.length ? (
                    <ul className="detail-list">
                      {detail.tasks.map((task) => (
                        <li key={task.id}>
                          <div className="detail-item-heading">
                            <strong>{task.summary}</strong>
                            <span>{readable(task.status)}</span>
                          </div>
                          <p>{readable(task.type)}{task.assignedToDisplayName ? ` · ${task.assignedToDisplayName}` : " · Unassigned"}{task.dueAt ? ` · Due ${dateLabel(task.dueAt)}` : ""}</p>
                        </li>
                      ))}
                    </ul>
                  ) : <p className="detail-empty">No coordinator tasks recorded.</p>}
                  {detail.barriers.length > 0 && (
                    <div className="detail-subsection">
                      <h3>Access barriers</h3>
                      <ul className="detail-list">
                        {detail.barriers.map((barrier, index) => (
                          <li key={`${barrier.category}-${index}`}>
                            <div className="detail-item-heading"><strong>{readable(barrier.category)}</strong><span>{readable(barrier.status)}</span></div>
                            <p>{barrier.description}</p>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </PageSection>

                <PageSection title="Contact history">
                  {detail.contacts.length ? (
                    <ul className="detail-list">
                      {detail.contacts.map((contact, index) => (
                        <li key={`${contact.attemptedAt}-${index}`}>
                          <div className="detail-item-heading"><strong>{readable(contact.outcome)}</strong><span>{dateLabel(contact.attemptedAt)}</span></div>
                          <p>{readable(contact.channel)} · {contact.language.toUpperCase()} · Recorded by {contact.recordedBy}</p>
                        </li>
                      ))}
                    </ul>
                  ) : <p className="detail-empty">No contact attempts recorded.</p>}
                </PageSection>
              </div>

              <div className="detail-column">
                <PageSection title="Report receipt">
                  {detail.reports.length ? (
                    <ul className="detail-list">
                      {detail.reports.map((report, index) => (
                        <li key={`${report.reference}-${index}`}>
                          <div className="detail-item-heading"><strong>Report received</strong><span>{dateLabel(report.receivedAt)}</span></div>
                          <p>Reference {report.reference} · Recorded by {report.recordedBy}</p>
                        </li>
                      ))}
                    </ul>
                  ) : <p className="detail-empty">No report receipt recorded.</p>}
                </PageSection>

                <PageSection title="Clinician review">
                  {detail.clinicalReviews.length ? (
                    <ul className="detail-list">
                      {detail.clinicalReviews.map((review, index) => (
                        <li key={`${review.acknowledgedAt}-${index}`}>
                          <div className="detail-item-heading"><strong>Acknowledged by {review.clinician}</strong><span>{dateLabel(review.acknowledgedAt)}</span></div>
                          <p>Next action: {review.nextAction}</p>
                        </li>
                      ))}
                    </ul>
                  ) : <p className="detail-empty">No clinician acknowledgement recorded.</p>}
                </PageSection>

                <PageSection title="Activity history">
                  {detail.auditEvents.length ? (
                    <ol className="detail-timeline">
                      {detail.auditEvents.map((event, index) => (
                        <li key={`${event.occurredAt}-${index}`}>
                          <span className="timeline-dot" aria-hidden="true" />
                          <strong>{readable(event.eventType)}</strong>
                          <span>{event.actor} · {dateLabel(event.occurredAt)}</span>{event.note && <p>{event.note}</p>}
                        </li>
                      ))}
                    </ol>
                  ) : <p className="detail-empty">No activity recorded.</p>}
                </PageSection>
              </div>
            </div>
          </>
        ) : (
          <div className="detail-unavailable" role="status">
            <h1 id="episode-title">Episode details unavailable</h1>
            <p>The episode service could not be reached. Try again when the service is available.</p>
            <form action={`/episodes/${encodeURIComponent(id)}`} method="get">
              <button className="retry-button" type="submit">Try again</button>
            </form>
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

export default function EpisodeRoute({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<main className="loading-screen" role="status">Loading episode…</main>}>
      <EpisodePage params={params} />
    </Suspense>
  );
}
