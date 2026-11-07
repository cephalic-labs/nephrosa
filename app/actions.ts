"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { apiFetch } from "./api-client";

export type ActionState = { message: string };
export async function submitAction(_previous: ActionState, form: FormData): Promise<ActionState> {
  const operation = String(form.get("operation") ?? "");
  const id = String(form.get("episodeId") ?? "");
  const payload = Object.fromEntries([...form.entries()].filter(([, value]) => typeof value === "string"));
  let path = "";
  let method = "POST";
  let body: object = payload;
  if (operation === "login") path = "/sessions";
  else if (operation === "create") {
    path = "/episodes";
    body = { ...payload, consent: form.get("consent") === "on", tests: form.getAll("tests") };
  } else if (operation === "workflow" && /^[a-f0-9-]{36}$/i.test(id)) path = `/episodes/${id}/actions`;
  else if (operation === "call" && /^[a-f0-9-]{36}$/i.test(id)) path = `/episodes/${id}/calls`;
  else if (operation === "reconcile" && /^[a-f0-9-]{36}$/i.test(id) && /^[a-f0-9-]{36}$/i.test(String(form.get("callId")))) path = `/episodes/${id}/calls/${form.get("callId")}/reconcile`;
  else if (operation === "logout") { path = "/session"; method = "DELETE"; }
  else return { message: "Invalid action." };
  let destination = "";
  try {
    const response = await apiFetch(path, { method, body: method === "DELETE" ? undefined : JSON.stringify(body) });
    if (!response.ok) {
      const error = await response.json().catch(() => null);
      return { message: error?.error?.message ?? "The action could not be completed. Refresh and retry." };
    }
    if (operation === "login") {
      const result = await response.json();
      (await cookies()).set("relay-session", result.token, { httpOnly: true, secure: process.env.NODE_ENV === "production", sameSite: "lax", path: "/", maxAge: 8 * 60 * 60 });
      destination = result.user.role === "clinician" ? "/dashboard/clinician" : "/dashboard/coordinator";
    } else if (operation === "logout") {
      (await cookies()).delete("relay-session"); destination = "/login";
    } else if (operation === "create") destination = `/episodes/${(await response.json()).id}`;
    revalidatePath("/");
    if (id) revalidatePath(`/episodes/${id}`);
  } catch {
    return { message: "Service unavailable. Refresh to check the current state before retrying, especially after a call request." };
  }
  if (destination) redirect(destination);
  return { message: "Saved. The episode has been updated." };
}
