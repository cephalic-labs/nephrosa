CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    role text NOT NULL CHECK (role IN ('clinician', 'coordinator', 'administrator')),
    display_name text NOT NULL CHECK (length(trim(display_name)) BETWEEN 1 AND 120),
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE patients (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    external_reference text NOT NULL UNIQUE CHECK (length(trim(external_reference)) BETWEEN 1 AND 80),
    display_name text NOT NULL CHECK (length(trim(display_name)) BETWEEN 1 AND 120),
    preferred_language text NOT NULL CHECK (preferred_language IN ('en', 'hi', 'ta')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE consents (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    channel text NOT NULL CHECK (channel IN ('voice', 'sms', 'email')),
    status text NOT NULL CHECK (status IN ('pending', 'granted', 'withdrawn')),
    recorded_at timestamptz,
    source text CHECK (source IS NULL OR length(source) <= 120),
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((status = 'pending' AND recorded_at IS NULL) OR (status <> 'pending' AND recorded_at IS NOT NULL))
);

CREATE TABLE episodes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id uuid NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status text NOT NULL CHECK (status IN (
        'ordered', 'outreach_pending', 'awaiting_booking', 'booked',
        'awaiting_test_completion', 'awaiting_report', 'awaiting_clinician_review', 'closed'
    )),
    due_at date NOT NULL,
    created_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE test_orders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    test_type text NOT NULL CHECK (test_type IN ('uacr', 'egfr')),
    ordered_at timestamptz NOT NULL DEFAULT now(),
    due_at date NOT NULL,
    UNIQUE (episode_id, test_type)
);

CREATE TABLE contact_attempts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    channel text NOT NULL CHECK (channel IN ('voice', 'sms', 'email', 'manual')),
    outcome text NOT NULL CHECK (outcome IN (
        'reached', 'no_answer', 'callback_requested', 'opted_out', 'wrong_contact', 'other'
    )),
    language text NOT NULL CHECK (language IN ('en', 'hi', 'ta')),
    attempted_at timestamptz NOT NULL,
    recorded_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE barriers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    category text NOT NULL CHECK (category IN (
        'cost', 'travel', 'time_off_work', 'unclear_instructions', 'language_access', 'other'
    )),
    description text NOT NULL CHECK (length(trim(description)) BETWEEN 1 AND 240),
    status text NOT NULL CHECK (status IN ('open', 'resolved')) DEFAULT 'open',
    recorded_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    resolved_at timestamptz,
    CHECK ((status = 'open' AND resolved_at IS NULL) OR (status = 'resolved' AND resolved_at IS NOT NULL))
);

CREATE TABLE tasks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    assigned_to uuid REFERENCES users(id) ON DELETE RESTRICT,
    task_type text NOT NULL CHECK (task_type IN (
        'contact_patient', 'resolve_barrier', 'verify_test', 'confirm_report', 'review_report'
    )),
    summary text NOT NULL CHECK (length(trim(summary)) BETWEEN 1 AND 240),
    priority text NOT NULL CHECK (priority IN ('routine', 'attention', 'urgent')) DEFAULT 'routine',
    status text NOT NULL CHECK (status IN ('open', 'in_progress', 'completed', 'cancelled')) DEFAULT 'open',
    due_at date,
    created_at timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz
);

CREATE TABLE test_completions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    test_order_id uuid NOT NULL UNIQUE REFERENCES test_orders(id) ON DELETE RESTRICT,
    completed_at timestamptz NOT NULL,
    verification_source text NOT NULL CHECK (verification_source IN (
        'lab_document', 'lab_portal', 'staff_confirmed'
    )),
    recorded_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    received_at timestamptz NOT NULL,
    reference text NOT NULL CHECK (length(trim(reference)) BETWEEN 1 AND 120),
    recorded_by uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE clinical_reviews (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id uuid NOT NULL REFERENCES episodes(id) ON DELETE RESTRICT,
    clinician_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    acknowledged_at timestamptz NOT NULL,
    next_action text NOT NULL CHECK (length(trim(next_action)) BETWEEN 1 AND 1000),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE audit_events (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_id uuid REFERENCES users(id) ON DELETE RESTRICT,
    episode_id uuid REFERENCES episodes(id) ON DELETE RESTRICT,
    event_type text NOT NULL CHECK (length(trim(event_type)) BETWEEN 1 AND 100),
    occurred_at timestamptz NOT NULL DEFAULT now(),
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX episodes_status_due_at_idx ON episodes (status, due_at);
CREATE INDEX episodes_owner_status_idx ON episodes (owner_id, status);
CREATE INDEX tasks_assignee_status_due_at_idx ON tasks (assigned_to, status, due_at);
CREATE INDEX barriers_episode_status_idx ON barriers (episode_id, status);
CREATE INDEX audit_events_episode_occurred_at_idx ON audit_events (episode_id, occurred_at DESC);

CREATE FUNCTION reject_audit_event_mutation() RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    RAISE EXCEPTION 'audit events are append-only';
END;
$function$;

CREATE TRIGGER audit_events_append_only
BEFORE UPDATE OR DELETE ON audit_events
FOR EACH ROW EXECUTE FUNCTION reject_audit_event_mutation();
