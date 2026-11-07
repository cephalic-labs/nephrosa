"use client";

import { useEffect, useRef, useTransition } from "react";
import { useRouter } from "next/navigation";

export function CallUpdates({ unresolved }: { unresolved: boolean | null }) {
  const router = useRouter();
  const [refreshing, startTransition] = useTransition();
  const lastKnownUnresolved = useRef(false);

  useEffect(() => {
    if (unresolved !== null) lastKnownUnresolved.current = unresolved;
    if (!lastKnownUnresolved.current || refreshing) return;

    function refreshWhenVisible() {
      if (document.visibilityState === "visible") startTransition(() => router.refresh());
    }

    const timer = window.setTimeout(refreshWhenVisible, 5000);
    document.addEventListener("visibilitychange", refreshWhenVisible);
    return () => {
      window.clearTimeout(timer);
      document.removeEventListener("visibilitychange", refreshWhenVisible);
    };
  }, [unresolved, refreshing, router]);

  return (
    <button
      className="retry-button"
      type="button"
      aria-disabled={refreshing}
      onClick={() => {
        if (!refreshing) startTransition(() => router.refresh());
      }}
    >
      {refreshing ? "Checking for update…" : "Check for update"}
    </button>
  );
}
