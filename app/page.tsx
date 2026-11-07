import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";

const steps = [
  ["Order", "A clinician places a uACR and eGFR monitoring order with a due date."],
  ["Follow through", "The care team helps the patient arrange testing and records verified completion."],
  ["Review", "The report reaches the clinician, who records the next action."],
];

export default function LandingPage() {
  return <main className="landing-page">
    <header className="landing-header"><Link className="wordmark" href="/" aria-label="Renal Relay home"><span className="wordmark-symbol" aria-hidden="true"><span/><span/><span/></span><span>renal relay</span></Link><nav aria-label="Main"><a href="#how-it-works">How it works</a><Link className="landing-sign-in" href="/login">Care team sign in <span aria-hidden="true">↗</span></Link><ThemeToggle/></nav></header>
    <section className="landing-hero">
      <div className="landing-photo" role="img" aria-label="A care coordinator listens to an older adult patient in a bright clinic"><div className="photo-caption">Care that follows through.</div></div>
      <div className="landing-copy"><p className="hero-kicker">Kidney monitoring, with people at the center</p><h1>Better follow-through starts with a conversation.</h1><p>Renal Relay helps care teams keep clinician-ordered kidney monitoring moving—from the first order to a reviewed report and a clear next step.</p><Link className="landing-cta" href="/login">Open care team workspace <span aria-hidden="true">↗</span></Link><p className="landing-note">Built for care teams. Clinical decisions stay with clinicians.</p></div>
    </section>
    <section className="landing-process" id="how-it-works"><div className="process-intro"><p className="hero-kicker">A clear path through care</p><h2>Every order has a next step.</h2><p>Keep the handoffs visible so patients and care teams know what still needs attention.</p></div><ol className="process-list">{steps.map(([title, copy], index)=><li key={title}><span className="process-step">{String(index+1).padStart(2,"0")}</span><div><h3>{title}</h3><p>{copy}</p></div></li>)}</ol></section>
    <footer className="landing-footer"><span>Renal Relay</span><span>A care-team follow-through workspace</span><span>Demonstration only · use fictional records</span></footer>
  </main>;
}
