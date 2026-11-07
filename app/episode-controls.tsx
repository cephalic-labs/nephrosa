import { apiFetch, type SessionUser } from "./api-client";
import { ActionForm } from "./action-form";
import { CallUpdates } from "./call-updates";

type Call = { id: string; attemptId: string; status: string; outcome: string; createdAt: string; note: string; followUpNote: string };
const callStatuses: Record<string, string> = {
  requesting: "Requesting call", pending: "Waiting for call result", uncertain: "Call acceptance uncertain",
  connected: "Call ended · Connected", no_answer: "No answer", busy: "Line busy", failed: "Call failed", reconciled: "Call reconciled",
};
const outcomes: Record<string, string> = {
  booked: "Patient reported booking a test",
  test_reported_complete: "Patient reported completing a test",
  barrier: "Patient reported a barrier",
  callback_requested: "Patient requested a callback",
  opted_out: "Patient opted out of voice calls",
  wrong_contact: "Wrong contact",
  unable_to_talk: "Patient was unable to talk",
  needs_review: "Outcome needs staff review",
};

function isUnresolved(call: Call) {
  return ["requesting", "pending", "uncertain"].includes(call.status);
}

export async function getEpisodeCalls(id: string): Promise<Call[] | null> {
  try {
    const response = await apiFetch(`/episodes/${encodeURIComponent(id)}/calls`);
    return response.ok ? await response.json() : null;
  } catch {
    return null;
  }
}

function CallResult({ call }: { call: Call }) {
  return <>
    <p><strong>{callStatuses[call.status] ?? call.status.replaceAll("_", " ")}</strong></p>
    {call.status === "connected" && call.outcome && <p>{outcomes[call.outcome] ?? call.outcome.replaceAll("_", " ")}</p>}
    {(call.status === "connected" || call.followUpNote) && <p className="patient-response"><strong>Patient response summary</strong><br />{call.followUpNote || "No patient response summary was returned."}</p>}
  </>;
}

function CallRecord({ call, episodeId, closed }: { call: Call; episodeId: string; closed: boolean }) {
  return <>
    <p><time dateTime={call.createdAt}>{call.createdAt}</time> · Provider attempt: {call.attemptId || "Not confirmed"}</p>
    {call.note && <p>Staff reconciliation: {call.note}</p>}
    {!closed && isUnresolved(call) && <ActionForm operation="reconcile" episodeId={episodeId} label="Confirm provider call is terminal"><input name="callId" type="hidden" value={call.id} /><p>Check Samvaad logs first. Confirm the call has ended or was never placed. Do not release an active call.</p><label>Provider attempt ID, if available<input name="attemptId" maxLength={200} /></label><label>What you verified<input name="note" maxLength={240} required /></label></ActionForm>}
  </>;
}

export async function EpisodeControls({ id, calls, closed, role, consent, phone, tests, tasks }: { id: string; calls: Call[] | null; closed: boolean; role: string; consent: string; phone: string; tests: { id: string; testType: string; completion?: object }[]; tasks: { id: string; summary: string; status: string }[] }) {
  const usersResponse = !closed && role === "clinician" ? await apiFetch("/users") : null;
  const users: SessionUser[] = usersResponse?.ok ? await usersResponse.json() : [];
  const [latest, ...earlier] = calls ?? [];
  const activeCall = calls === null ? null : calls.some(isUnresolved);
  return <section className="workspace-panel"><h2>Care-team actions</h2><p>Voice consent: <strong>{consent}</strong> · Contact: {phone || "Not recorded"}</p>
    <div className="controls-grid">
      <section><h3>Voice follow-up</h3>
        <h4>Latest call</h4>
        <div className="latest-call" role="status" aria-atomic="true">
          {latest ? <CallResult call={latest} /> : <p>{calls === null ? "Call history is unavailable." : "No calls recorded."}</p>}
        </div>
        <CallUpdates unresolved={activeCall} />
        {latest && <CallRecord key={latest.id} call={latest} episodeId={id} closed={closed} />}
        {!closed && <>
          {consent === "granted" && phone && activeCall === false ? <ActionForm operation="call" episodeId={id} label="Call patient" /> : <p>Calling requires granted consent, a phone number, and no unresolved call.</p>}
          <ActionForm operation="workflow" episodeId={id} action="consent" label="Record consent"><label>Voice consent<select name="consent"><option value="withdrawn">Withdrawn</option><option value="granted">Granted</option></select></label><label>Source<input name="source" maxLength={120} required /></label></ActionForm>
        </>}
      </section>
      {!closed && <>
      <section><h3>Test verification</h3>{tests.filter(t => !t.completion).map(t => <ActionForm key={t.id} operation="workflow" episodeId={id} action="verify_test" label={`Verify ${t.testType}`}><input name="testId" type="hidden" value={t.id} /><label>Evidence source<select name="source"><option value="staff_confirmed">Staff confirmed</option><option value="lab_document">Lab document</option><option value="lab_portal">Lab portal</option></select></label></ActionForm>)}
        <ActionForm operation="workflow" episodeId={id} action="report" label="Record report receipt"><label>Report reference<input name="reference" maxLength={120} required /></label></ActionForm>
      </section>
      <section><h3>Coordinator work</h3>{tasks.filter(t => ["open", "in_progress"].includes(t.status)).map(t => <ActionForm key={t.id} operation="workflow" episodeId={id} action="resolve_task" label="Resolve task"><p>{t.summary}</p><input name="taskId" type="hidden" value={t.id} /><label>Resolution<input name="note" maxLength={240} required /></label></ActionForm>)}
        <ActionForm operation="workflow" episodeId={id} action="resolve_barriers" label="Resolve open barriers"><label>Resolution<input name="note" maxLength={240} required /></label></ActionForm>
      </section>
      {role === "clinician" && <section><h3>Clinician sign-off</h3>
        <ActionForm operation="workflow" episodeId={id} action="review" label="Acknowledge report"><label>Documented next action<textarea name="nextAction" maxLength={1000} required /></label></ActionForm>
        <ActionForm operation="workflow" episodeId={id} action="assign" label="Assign care owner"><label>Owner<select name="owner" required>{users.map(u => <option key={u.id} value={u.id}>{u.name}</option>)}</select></label></ActionForm>
        <p>Closure requires test verification, report receipt, clinician next action, and resolution of outstanding work.</p><ActionForm operation="workflow" episodeId={id} action="close" label="Close episode" />
      </section>}
      </>}
    </div>
    {earlier.length > 0 && <><h3>Call history</h3><ul className="detail-list">{earlier.map(call => <li key={call.id}>
      <CallResult call={call} /><CallRecord call={call} episodeId={id} closed={closed} />
    </li>)}</ul></>}
  </section>;
}
