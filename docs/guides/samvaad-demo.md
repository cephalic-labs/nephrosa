# Renal Relay / Samvaad demonstration

## Local values and accounts

Keep non-secret settings in the repository root `.env`; `.env.example` lists the required names. Next.js uses `.env.local` for `RELAY_API_URL`. Go receives configuration explicitly through Compose. Do not add API credentials to any `NEXT_PUBLIC_*` variable.

For a fresh checkout, `scripts/dev-secrets.sh` generates database and demo-account password files and creates an empty `secrets/sarvam_api_key`. Existing installations should create only missing files. Preserve existing credentials and avoid rerunning credential generation by deleting existing secrets. Put the Samvaad workspace API key in `secrets/sarvam_api_key`, with mode 600. Put 12–72 character demonstration passwords in `secrets/demo_clinician_password` and `secrets/demo_coordinator_password`, also mode 600. Do not paste secrets into chat or shell command arguments. Compose mounts these exact paths. The `*_FILE` settings describe equivalent paths for directly running Go.

The seed command sets the logins `clinician` and `coordinator` to the passwords in those files using bcrypt. It preserves existing database passwords and never prints account passwords. Reseeding rotates demo login credentials and invalidates their old sessions. Seeding is refused outside development.

## Configure the agent in Samvaad

In the existing workspace, create an English voice agent called “Renal Relay follow-up”. Configure the following input variables `patient_name`, `due_date`, and `episode_reference`. Configure output variables `outcome` and `follow_up_note`. Commit a version and put its app ID and numeric version in `.env`. Connect an outbound-capable telephony number and put its connection ID and E.164 number in the corresponding settings. Supported providers and rented-number availability depend on your workspace.

Each output variable needs an extraction prompt in the Samvaad builder. Configure `outcome` as an enum or string that extracts one of the allowed outcomes from what the caller actually said. Configure `follow_up_note` to summarize the caller's practical request in at most 240 characters. A callback can be selected when the caller agrees to care-team follow-up. A structured barrier is created by the application when the returned outcome is `barrier`.

Commit agent changes and update `SARVAM_APP_VERSION` in `.env` to the committed version. Recreate the API after changing this value or its key. The API rejects calls that send input variables absent from the selected agent version.

Use this starting instruction for the fictional demonstration. A clinician must review wording before operational use.

> You are Renal Relay's automated care-team follow-up assistant. Introduce yourself as automated and ask whether this is the intended person and whether now is a suitable time. Until they confirm, disclose no monitoring details. If this is the wrong person, politely stop and set outcome to wrong_contact. If they cannot talk, offer a care-team callback. After confirmation, explain that you are following up on clinician-ordered kidney monitoring, due on due_date. Ask whether a lab visit has been arranged or the tests have been done. Capture practical barriers without giving medical advice. You cannot book a real appointment or verify a test yourself. If the person wants human help or asks a clinical question, explain that a care-team member will follow up. Do not interpret results, diagnose, recommend treatment, promise an appointment, or say the episode is complete. If they ask to stop automated calls, acknowledge that and set outcome to opted_out. End politely and keep the conversation concise. Set outcome to exactly one of booked, test_reported_complete, barrier, callback_requested, opted_out, wrong_contact, unable_to_talk. Leave unknown outcomes for human review. Set follow_up_note to a short factual account of the patient's practical request, at most 240 characters. Never include medical interpretations. A booked outcome means only that the patient reports an arrangement.

The application persists the allowlisted outcome and short follow-up note. Configure provider-side retention separately for the demonstration.

## Callback security and recovery

`SARVAM_WEBHOOK_URL` is the public HTTPS base ending in `/api/v1/sarvam/webhooks`. Expose only this callback prefix publicly. Keep database and other Go routes on loopback. Each outbound request receives a different high-entropy callback URL. The database stores only its hash. The receiver authenticates possession of a per-call bearer URL. Configure your reverse proxy to redact callback paths and bodies from access logs and analytics. Do not publish the complete callback URL, including in screenshots. `SARVAM_WEBHOOK_ENABLED` defaults to false and must remain false until this routing and redaction are verified.

The public docs reviewed do not document a signing header or a guaranteed instant-call delivery retry policy. Do not invent either. The receiver validates the secret, attempt ID, phone line, body limit, and instant-call status. Duplicate terminal events have no side effects. A callback arriving before the outbound response can safely bind the attempt ID. The callback endpoint doesn't trust caller-supplied episode IDs or campaign fields. Invalid/missing agent outcomes create human-review work.

After two minutes, an authenticated coordinator can reconcile an outstanding call. Check the Samvaad call logs, confirm the call is terminal or definitely was not placed, and record its attempt ID if available and a factual note. Then use “Confirm provider call is terminal”. This releases the redial lock without claiming a connected call, test completion, or a fabricated outcome. Use the task and verification workflow to capture subsequent human follow-up. A reconciled call's delayed callback cannot overwrite it.

Do not release an ongoing call. Consent withdrawal prevents new calls. Calls already dispatched to the provider must finish through the provider flow. Finish that call with the opt-out branch.

## Testing checkpoint and startup

Before testing, the implementer asks the owner to populate missing root `.env` settings and secret files and explicitly confirm readiness. Keep values out of logs. Start the configured services with the following commands.

```sh
python3 scripts/demo-preflight.py
docker compose --env-file .env -f deploy/compose.yaml up -d --build db api
docker compose --env-file .env -f deploy/compose.yaml --profile dev-seed run --rm seed
pnpm dev
```

The explicit `--env-file .env` loads the root settings for Compose interpolation, including the Samvaad enable flag. Check `/health/live` and `/health/ready` locally on port 8080. Run `go test ./...`, `pnpm lint`, and `pnpm build`. Database integration tests additionally require `TEST_DATABASE_URL` pointing at a disposable PostgreSQL database using a schema-creating test role; they create and drop their own schema. Never point this at production.

The preflight checks formatting and local secret permissions without displaying values or placing calls. It does not establish provider account entitlements or reachability of the callback. Check those through the first authorized live call. Use `DEMO_CALL_PHONE_NUMBER` as the number entered in the fictional patient form; the backend never silently substitutes a different destination.

## Presenter script

1. Sign in as clinician. Create `DEMO-<unique reference>` with a fictional display name, demonstrator phone, English, both tests, due date, owner, granted voice consent, and its source.
2. Sign out and sign in as coordinator. Open the episode, choose “Call patient”, and answer the real phone. Describe a scheduling barrier or request a callback.
3. Refresh the episode. Show the provider attempt, outcome, patient-reported note, contact history, task, and audit label. Resolve the task with a note and resolve any open barrier separately.
4. Verify each ordered test with a source, then record a fictional report reference. Attempting report receipt before verification is blocked.
5. Sign in as clinician. Acknowledge the report and document a fictional next action. Close the episode. Demonstrate that unresolved tasks/calls/barriers block closure.
6. Return to the dashboard. Open episodes decrease. “Include closed episodes” retains the completed episode and its audit history.
7. For the no-answer scenario, create another uniquely referenced episode and do not answer. Show the real provider outcome and follow-up task. No successful contact or test completion is invented.

## Verified demonstration

The following checks were completed on 7 November 2026 using fictional clinical data and consenting demonstration recipients.

| Check | Observed result |
| --- | --- |
| Automated backend checks | Unit tests and the PostgreSQL integration suite passed, including permissions, consent withdrawal, duplicate-call protection, duplicate callbacks, reconciliation, audit protection, and closure prerequisites. |
| Frontend checks | ESLint and the production build passed. A browser rehearsal created an episode, rejected premature closure, recorded evidence, acknowledged a report, and closed the episode with a clean console. |
| Live booking flow | A real call connected and its ngrok callback recorded `booked`. Authenticated workflow APIs completed fictional task resolution, test verification, report receipt, clinician acknowledgement, next action, and closure. Next.js detail and history pages rendered the closed episode, which had ten audit events. The first booking call returned an empty note. |
| Busy-line flow | A real attempt returned `busy`. The callback created contact history, an open coordinator task, and an audit event. |
| Live callback flow | Later connected calls returned `callback_requested` with nonempty notes about arranging a lab visit, care-team contact, and requested timing/location or cost information. The application saved a coordinator task and audit event and retained the open episode. |
| Updated configuration | The running API was checked against the updated committed app version. Newly created episodes were checked against the updated demo number before dispatch. |

The live-call results verify provider connectivity, callbacks, and application persistence. Extracted summaries have not been independently checked against recordings or transcripts. The scheduling-barrier rehearsal returned `callback_requested`, so live creation of a structured barrier remains unverified. Calls start in English and contact history is currently labeled `en`. Speaking another language requires appropriate Samvaad agent settings, and multilingual accuracy remains unverified.

## Troubleshooting observed during the demo

- The initial credential was rejected by the Voice Agents API with HTTP 401 and an invalid-key-format message. Use the workspace key from Samvaad Settings → API Key and store its raw value in the owner-only key file.
- Compose initially started the API with the callback flag disabled. Explicit root `--env-file .env` loading and API recreation applied the configured flag. A malformed validation probe then reached the handler through ngrok and returned HTTP 400 as expected.
- Sarvam rejected the selected agent variable setup with HTTP 422 because `patient_name` and `episode_reference` were absent from the selected agent's variables. Adding the input variables, committing a version, and reloading its numeric version enabled calling.
- Sarvam's observed rejection format uses `error.data.details`. The client handles this envelope and retains bounded, redacted diagnostics. Regression tests cover status visibility and redaction of keys, contacts, and callback URLs.
- The local Docker client's saved registry credential had expired. The public images were pulled with an isolated anonymous client configuration, preserving the saved credentials.

## Safe reset

Prefer a fresh `DEMO-...` episode for each rehearsal; historical records remain available. For a full local reset, stop API traffic, export anything needed, and run `docker compose --env-file .env -f deploy/compose.yaml run --rm migrate reset-demo --confirm-fictional-reset` using the migrator role. It refuses non-development environments, a non-renalrelay database, or any patient reference outside `DEMO-`. It truncates workflow/demo records, preserves migration history and schema. Run the seed command afterward. This permanently deletes the fictional demonstration history. Never use it for retained records.

## Sources

- [Instant outbound call API](https://docs.sarvam.ai/voice-agents/api/instant-outbound/create)
- [Instant-call webhook contract](https://docs.sarvam.ai/voice-agents/api/instant-outbound/webhook-payload)
- [Telephony connections](https://docs.sarvam.ai/voice-agents/deploy/telephony/bring-your-own)
- [Agent input variables and output extraction](https://docs.sarvam.ai/voice-agents/build/variables-personalization)
