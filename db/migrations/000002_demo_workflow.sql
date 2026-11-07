ALTER TABLE users ADD COLUMN login text UNIQUE, ADD COLUMN password_hash text;
ALTER TABLE patients ADD COLUMN phone_number text CHECK (phone_number ~ '^\+[1-9][0-9]{7,14}$');
CREATE TABLE sessions (
 token_hash text PRIMARY KEY, user_id uuid NOT NULL REFERENCES users(id), expires_at timestamptz NOT NULL
);
CREATE TABLE call_attempts (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(), episode_id uuid NOT NULL REFERENCES episodes(id),
 requested_by uuid NOT NULL REFERENCES users(id), attempt_id text UNIQUE,
 callback_hash text NOT NULL UNIQUE, agent_phone_number text NOT NULL,
 status text NOT NULL CHECK (status IN ('requesting','pending','uncertain','connected','no_answer','busy','failed','reconciled')),
 follow_up_note text NOT NULL DEFAULT '' CHECK(length(follow_up_note)<=240),
 outcome text NOT NULL DEFAULT '', interaction_id text NOT NULL DEFAULT '',
 created_at timestamptz NOT NULL DEFAULT now(), finished_at timestamptz,
 reconciliation_note text NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX one_active_call_per_episode ON call_attempts(episode_id) WHERE status IN ('requesting','pending','uncertain');
CREATE INDEX consents_latest_idx ON consents(patient_id, channel, created_at DESC, id);
