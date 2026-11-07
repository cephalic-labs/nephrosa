import { apiFetch, type SessionUser } from "./api-client";
import { ActionForm } from "./action-form";

type Call = { id: string; attemptId: string; status: string; outcome: string; createdAt: string; note: string; followUpNote: string };
export async function EpisodeControls({ id, closed, role, consent, phone, tests, tasks }: { id: string; closed: boolean; role: string; consent: string; phone: string; tests: { id: string; testType: string; completion?: object }[]; tasks: { id: string; summary: string; status: string }[] }) {
  const [callsResponse, usersResponse] = await Promise.all([apiFetch(`/episodes/${id}/calls`), apiFetch("/users")]);
  const calls: Call[] | null = callsResponse.ok ? await callsResponse.json() : null;
  const users: SessionUser[] = usersResponse.ok ? await usersResponse.json() : [];
  const activeCall = calls?.some(c => ["requesting", "pending", "uncertain"].includes(c.status));
  return <section className="workspace-panel"><h2>Care-team actions</h2><p>Voice consent: <strong>{consent}</strong> · Contact: {phone || "Not recorded"}</p>
    {!closed && <div className="controls-grid">
      <section><h3>Voice follow-up</h3><p>Call results are patient-reported. Refresh the page after the call to load its outcome.</p>
        {consent === "granted" && phone && calls !== null && !activeCall ? <ActionForm operation="call" episodeId={id} label="Call patient" /> : <p>Calling requires granted consent, a phone number, and no unresolved call.</p>}
        <ActionForm operation="workflow" episodeId={id} action="consent" label="Record consent"><label>Voice consent<select name="consent"><option value="withdrawn">Withdrawn</option><option value="granted">Granted</option></select></label><label>Source<input name="source" maxLength={120} required /></label></ActionForm>
      </section>
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
    </div>}
    <h3>Call history</h3>{calls === null ? <p role="alert">Call history is unavailable. Refresh before requesting another call.</p> : calls.length === 0 ? <p>No calls recorded.</p> : <ul className="detail-list">{calls.map(call => <li key={call.id}><strong>{call.status.replaceAll("_", " ")}</strong><p>{call.outcome || "Awaiting outcome"} · {call.createdAt}</p>{call.followUpNote && <p>Patient-reported: {call.followUpNote}</p>}<p>Provider attempt: {call.attemptId || "Not confirmed"}</p>{call.note && <p>Reconciliation: {call.note}</p>}
      {!closed && ["requesting", "pending", "uncertain"].includes(call.status) && <ActionForm operation="reconcile" episodeId={id} label="Confirm provider call is terminal"><input name="callId" type="hidden" value={call.id} /><p>Check Samvaad logs first. Confirm the call has ended or was never placed. Do not release an active call.</p><label>Provider attempt ID, if available<input name="attemptId" maxLength={200} /></label><label>What you verified<input name="note" maxLength={240} required /></label></ActionForm>}
    </li>)}</ul>}
  </section>;
}
