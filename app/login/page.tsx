import { ActionForm } from "../action-form";
export default function Login() {
  return <main className="dashboard-shell"><section className="workspace-panel"><h1>Sign in to Renal Relay</h1><p>Use your care-team demonstration account.</p>
    <ActionForm operation="login" label="Sign in">
      <label>Login<select name="login" required><option value="clinician">Clinician</option><option value="coordinator">Coordinator</option></select></label>
      <label>Password<input name="password" type="password" autoComplete="current-password" minLength={12} maxLength={72} required /></label>
    </ActionForm>
  </section></main>;
}
