import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";

type SiteHeaderProps = {
  active?: "dashboard" | "workflow" | "episodes";
  connected?: boolean;
};

export function SiteHeader({ active = "dashboard", connected = false }: SiteHeaderProps) {
  return (
    <header className="site-header">
      <Link className="wordmark" href="/" aria-label="Renal Relay dashboard">
        <span className="wordmark-symbol" aria-hidden="true"><span /><span /><span /></span>
        <span>renal relay</span>
      </Link>
      <nav className="primary-nav" aria-label="Main navigation">
        <Link href="/#dashboard" aria-current={active === "dashboard" ? "page" : undefined}>Dashboard</Link>
        <Link href="/#workflow" aria-current={active === "workflow" ? "page" : undefined}>Workflow</Link>
        <Link href="/#episodes" aria-current={active === "episodes" ? "page" : undefined}>Episodes</Link>
      </nav>
      <ThemeToggle />
      <div className={`connection-state ${connected ? "is-connected" : "is-offline"}`} role="status">
        <span className="connection-dot" />
        {connected ? "Live connection" : "Connection unavailable"}
      </div>
    </header>
  );
}
