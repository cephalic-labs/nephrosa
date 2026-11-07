"use client";
import { useActionState, type ReactNode } from "react";
import { submitAction } from "./actions";

export function ActionForm({ operation, episodeId, action, label, children }: { operation: string; episodeId?: string; action?: string; label: string; children?: ReactNode }) {
  const [state, formAction, pending] = useActionState(submitAction, { message: "" });
  return <form action={formAction} className="workflow-form">
    <input type="hidden" name="operation" value={operation} />
    {episodeId && <input type="hidden" name="episodeId" value={episodeId} />}
    {action && <input type="hidden" name="action" value={action} />}
    {children}
    <button className="retry-button" type="submit" disabled={pending}>{pending ? "Working…" : label}</button>
    <p role="status" aria-live="polite">{state.message}</p>
  </form>;
}
