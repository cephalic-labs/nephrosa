import { Suspense } from "react";
import Link from "next/link";
import { requireUser, apiFetch, type SessionUser } from "../../api-client";
import { ActionForm } from "../../action-form";
import { SiteHeader } from "../../site-header";
async function NewEpisode() {
  const user = await requireUser();
  if (user.role !== "clinician") return <main className="dashboard-shell"><p>A clinician must create the episode.</p><Link href="/">Dashboard</Link></main>;
  const response = await apiFetch("/users");
  if (!response.ok) throw new Error("Care owners could not be loaded.");
  const users: SessionUser[] = await response.json();
  return <main className="dashboard-shell"><SiteHeader connected /><section className="workspace-panel"><Link href="/">← Dashboard</Link><h1>Create monitoring episode</h1><p>Use fictional patient information and a consenting demonstration phone.</p>
    <ActionForm operation="create" label="Create episode">
      <label>Patient display name<input name="name" maxLength={120} required /></label>
      <label>Unique patient reference<input name="reference" maxLength={80} required /></label>
      <label>Phone number (E.164)<input name="phone" type="tel" placeholder="+91…" pattern="\+[1-9][0-9]{7,14}" required /></label>
      <label>Language<select name="language"><option value="en">English</option></select></label>
      <label>Care owner<select name="owner" required>{users.map(u => <option key={u.id} value={u.id}>{u.name} · {u.role}</option>)}</select></label>
      <label>Due date<input name="due" type="date" required /></label>
      <fieldset><legend>Ordered tests (select at least one)</legend><label className="check-label"><input name="tests" value="uacr" type="checkbox" defaultChecked /> uACR</label><label className="check-label"><input name="tests" value="egfr" type="checkbox" defaultChecked /> eGFR</label></fieldset>
      <label className="check-label"><input name="consent" type="checkbox" /> Voice-contact consent has been granted</label>
      <label>Consent source or pending-consent note<input name="consentSource" maxLength={120} required /></label>
    </ActionForm>
  </section></main>;
}
export default function Page() { return <Suspense fallback={<p role="status">Loading episode form…</p>}><NewEpisode /></Suspense>; }
