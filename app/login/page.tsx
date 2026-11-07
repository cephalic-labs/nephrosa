import Link from "next/link";
import { ActionForm } from "../action-form";
import { ThemeToggle } from "../theme-toggle";
export default function Login() {
  return <main className="login-page"><header className="login-header"><Link className="wordmark" href="/"><span className="wordmark-symbol" aria-hidden="true"><span/><span/><span/></span><span>renal relay</span></Link><ThemeToggle/></header><section className="login-layout"><div className="login-story"><div className="login-art" role="img" aria-label="Illustration of a clinician and patient reviewing a care plan"><span className="login-sun"/><span className="login-person login-person-one"/><span className="login-person login-person-two"/><span className="login-table"/></div><p>Thoughtful follow-through.<br/>For every person, every step.</p></div><section className="workspace-panel login-panel"><p className="hero-kicker">Care team access</p><h1>Welcome back.</h1><p>Sign in to continue coordinating kidney monitoring.</p>
    <ActionForm operation="login" label="Sign in">
      <label>Sign in as<select name="login" required><option value="clinician">Clinician</option><option value="coordinator">Care coordinator</option></select></label>
      <label>Password<input name="password" type="password" autoComplete="current-password" minLength={12} maxLength={72} required /></label>
    </ActionForm>
  </section></section><footer className="login-footer">Renal Relay · Care teams remain responsible for clinical decisions.</footer></main>;
}
