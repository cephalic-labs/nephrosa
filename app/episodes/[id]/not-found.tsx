import Link from "next/link";
import { SiteHeader } from "../../site-header";

export default function EpisodeNotFound() {
  return (
    <main className="dashboard-shell">
      <SiteHeader />
      <section className="episode-detail" aria-labelledby="not-found-title">
        <Link className="back-link" href="/#episodes">← All episodes</Link>
        <div className="detail-unavailable">
          <h1 id="not-found-title">Episode not found</h1>
          <p>This episode does not exist or is unavailable.</p>
          <Link className="retry-button detail-return" href="/#episodes">Return to episodes</Link>
        </div>
      </section>
      <footer className="site-footer">
        <span>Renal Relay</span>
        <span>Care teams retain responsibility for clinical review and next-action decisions.</span>
      </footer>
    </main>
  );
}
