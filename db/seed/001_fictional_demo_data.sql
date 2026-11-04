INSERT INTO users (id, role, display_name) VALUES
    ('00000000-0000-4000-8000-000000000001', 'clinician', 'Demo Clinician'),
    ('00000000-0000-4000-8000-000000000002', 'coordinator', 'Demo Coordinator')
ON CONFLICT (id) DO NOTHING;

INSERT INTO patients (id, external_reference, display_name, preferred_language) VALUES
    ('00000000-0000-4000-8000-000000000101', 'DEMO-001', 'Fictional Patient 01', 'en'),
    ('00000000-0000-4000-8000-000000000102', 'DEMO-002', 'Fictional Patient 02', 'en')
ON CONFLICT (id) DO NOTHING;

INSERT INTO consents (id, patient_id, channel, status) VALUES
    ('00000000-0000-4000-8000-000000000201', '00000000-0000-4000-8000-000000000101', 'voice', 'pending'),
    ('00000000-0000-4000-8000-000000000202', '00000000-0000-4000-8000-000000000102', 'voice', 'pending')
ON CONFLICT (id) DO NOTHING;

INSERT INTO episodes (id, patient_id, owner_id, status, due_at, created_by) VALUES
    (
        '00000000-0000-4000-8000-000000000301',
        '00000000-0000-4000-8000-000000000101',
        '00000000-0000-4000-8000-000000000002',
        'awaiting_booking',
        CURRENT_DATE - 5,
        '00000000-0000-4000-8000-000000000001'
    ),
    (
        '00000000-0000-4000-8000-000000000302',
        '00000000-0000-4000-8000-000000000102',
        '00000000-0000-4000-8000-000000000001',
        'awaiting_clinician_review',
        CURRENT_DATE + 2,
        '00000000-0000-4000-8000-000000000001'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO test_orders (id, episode_id, test_type, ordered_at, due_at) VALUES
    ('00000000-0000-4000-8000-000000000401', '00000000-0000-4000-8000-000000000301', 'uacr', now() - interval '12 days', CURRENT_DATE - 5),
    ('00000000-0000-4000-8000-000000000402', '00000000-0000-4000-8000-000000000301', 'egfr', now() - interval '12 days', CURRENT_DATE - 5),
    ('00000000-0000-4000-8000-000000000403', '00000000-0000-4000-8000-000000000302', 'uacr', now() - interval '8 days', CURRENT_DATE + 2),
    ('00000000-0000-4000-8000-000000000404', '00000000-0000-4000-8000-000000000302', 'egfr', now() - interval '8 days', CURRENT_DATE + 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO barriers (id, episode_id, category, description, recorded_by) VALUES
    (
        '00000000-0000-4000-8000-000000000501',
        '00000000-0000-4000-8000-000000000301',
        'time_off_work',
        'Fictional demo barrier: patient needs a lab time outside work hours.',
        '00000000-0000-4000-8000-000000000002'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, episode_id, assigned_to, task_type, summary, priority, due_at) VALUES
    (
        '00000000-0000-4000-8000-000000000601',
        '00000000-0000-4000-8000-000000000301',
        '00000000-0000-4000-8000-000000000002',
        'resolve_barrier',
        'Discuss a lab time that works around the patient schedule.',
        'attention',
        CURRENT_DATE + 2
    ),
    (
        '00000000-0000-4000-8000-000000000602',
        '00000000-0000-4000-8000-000000000302',
        '00000000-0000-4000-8000-000000000001',
        'review_report',
        'Review the received report and record the next action.',
        'routine',
        CURRENT_DATE + 1
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO test_completions (id, test_order_id, completed_at, verification_source, recorded_by) VALUES
    ('00000000-0000-4000-8000-000000000701', '00000000-0000-4000-8000-000000000403', now() - interval '2 days', 'staff_confirmed', '00000000-0000-4000-8000-000000000002'),
    ('00000000-0000-4000-8000-000000000702', '00000000-0000-4000-8000-000000000404', now() - interval '2 days', 'staff_confirmed', '00000000-0000-4000-8000-000000000002')
ON CONFLICT (id) DO NOTHING;

INSERT INTO reports (id, episode_id, received_at, reference, recorded_by) VALUES
    (
        '00000000-0000-4000-8000-000000000801',
        '00000000-0000-4000-8000-000000000302',
        now() - interval '1 day',
        'DEMO-REPORT-002',
        '00000000-0000-4000-8000-000000000002'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO audit_events (actor_id, episode_id, event_type, metadata)
SELECT
        '00000000-0000-4000-8000-000000000001',
        '00000000-0000-4000-8000-000000000301',
        'demo_episode_created',
        '{"source":"fictional_seed","patient_reference":"DEMO-001"}'::jsonb
WHERE NOT EXISTS (
    SELECT 1 FROM audit_events
    WHERE episode_id = '00000000-0000-4000-8000-000000000301'
      AND event_type = 'demo_episode_created'
);

INSERT INTO audit_events (actor_id, episode_id, event_type, metadata)
SELECT
        '00000000-0000-4000-8000-000000000001',
        '00000000-0000-4000-8000-000000000302',
        'demo_episode_created',
        '{"source":"fictional_seed","patient_reference":"DEMO-002"}'::jsonb
WHERE NOT EXISTS (
    SELECT 1 FROM audit_events
    WHERE episode_id = '00000000-0000-4000-8000-000000000302'
      AND event_type = 'demo_episode_created'
);
