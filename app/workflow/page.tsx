import Link from "next/link";
import { Suspense } from "react";
import { requireUser } from "../api-client";
import { SiteHeader } from "../site-header";
import { ActionForm } from "../action-form";

const pathway = [
  { name: "Order placed", owner: "Clinician", detail: "Record the ordered uACR and eGFR tests, due date, and care owner." },
  { name: "Testing arranged", owner: "Care team", detail: "Contact the patient with consent, understand practical barriers, and record a useful next step." },
  { name: "Completion verified", owner: "Care coordinator", detail: "Record test completion only when it has been confirmed and note its source." },
  { name: "Report received", owner: "Care coordinator", detail: "Record the report reference after ordered tests have been verified complete." },
  { name: "Review and next action", owner: "Clinician", detail: "Acknowledge the report, document the clinical next action, then close when ready." },
];

async function WorkflowContent() {
  const user = await requireUser();
  return <main className="dashboard-shell"><SiteHeader active="workflow" connected role={user.role}/><div className="workspace-toolbar"><span>{user.name} · {user.role === "clinician" ? "Clinician" : "Care coordinator"}</span><Link href={`/dashboard/${user.role}`}>Back to dashboard</Link><ActionForm operation="logout" label="Sign out"/></div><section className="pathway-page"><p className="hero-kicker">How an episode moves forward</p><h1>One clear path.<br/>A team at every step.</h1><p className="pathway-intro">Renal Relay keeps responsibility and unfinished work visible from order through clinician review.</p><ol className="pathway-list">{pathway.map((step,index)=><li key={step.name}><span className="pathway-number">{String(index+1).padStart(2,"0")}</span><div className="pathway-content"><span className="pathway-owner">{step.owner}</span><h2>{step.name}</h2><p>{step.detail}</p></div></li>)}</ol><aside className="pathway-note"><strong>Clinical decisions stay with the clinician.</strong><span>Contact attempts never count as completed tests, report receipt, or clinical review.</span></aside><Link className="landing-cta" href={`/dashboard/${user.role}#episodes`}>Go to episodes <span aria-hidden="true">↗</span></Link></section></main>;
}

export default function WorkflowPage() {
  return <Suspense fallback={<main className="loading-screen" role="status">Loading care pathway…</main>}><WorkflowContent/></Suspense>;
}
