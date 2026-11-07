import "server-only";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function apiFetch(path: string, init: RequestInit = {}) {
  const session = (await cookies()).get("relay-session")?.value;
  const response = await fetch(`${process.env.RELAY_API_URL ?? "http://127.0.0.1:8080"}/api/v1${path}`, {
    ...init, cache: "no-store", signal: AbortSignal.timeout(25000),
    headers: { "Content-Type": "application/json", ...(session ? { Authorization: `Bearer ${session}` } : {}), ...init.headers },
  });
  return response;
}
export type SessionUser = { id: string; name: string; role: string };
export async function requireUser(): Promise<SessionUser> {
  const response = await apiFetch("/session");
  if (response.status === 401) redirect("/login");
  if (!response.ok) throw new Error("The care-team service is unavailable. Please retry.");
  return response.json();
}
