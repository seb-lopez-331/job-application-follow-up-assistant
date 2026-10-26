CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    preferences JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    CONSTRAINT preferences_follow_up_schedule_check CHECK (
        preferences @> '{"follow_up_schedule": {}}'::jsonb AND
        jsonb_typeof(preferences->'follow_up_schedule') = 'object'
    );
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id)
        ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '1 hour',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_title TEXT NOT NULL,
    company_name TEXT NOT NULL,
    date_posted TIMESTAMPZ NOT NULL,
    loc TEXT,
    descr TEXT NOT NULL,
    job_url TEXT,
    work_arrangement TEXT CHECK (
        work_arrangement IN ('remote', 'hybrid', 'on-site')
    ),
    created_at TIMESTAMPZ NOT NULL DEFAULT now()
    updated_at TIMESTAMPZ NOT NULL DEFAULT now()
);

CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action_type TEXT NOT NULL CHECK (
        action_type IN (
            'job_saved',
            'job_applied',
            'interview_scheduled',
            'offer_received',
            'application_withdrawn',
            'application_rejected'
        )
    ),
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE TABLE follow_ups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_application_id UUID NOT NULL
        REFERENCES job_applications(id) ON DELETE CASCADE,
    due_on TIMESTAMPTZ NOT NULL,
    message_draft TEXT,
    follow_up_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN (
            'pending', 'completed', 'cancelled'
        )),
    sent_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE TABLE job_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id)
        ON DELETE CASCADE,
    job_id UUID NOT NULL REFERENCES jobs(id)
        ON DELETE CASCADE,
    contact_id UUID REFERENCES contacts(id)
        ON DELETE SET NULL,
    stage TEXT NOT NULL DEFAULT 'applied'
        CHECK(status IN (
            'saved',
            'applied',
            'recruiter_screen',
            'interviewing',
            'offer',
            'rejected',
            'withdrawn'
        )),
    notes TEXT,
    applied_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    next_follow_up_on TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

create TABLE inteviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_application_id UUID NOT NULL REFERENCES job_applications(id)
        ON DELETE CASCADE,
    interview_date TIMESTAMPTZ NOT NULL,
    interview_type TEXT CHECK (
        interview_type IN ('phone', 'video', 'in-person')
    ),
    interviewer_name TEXT,
    interviewer_contact TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE TABLE job_application_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_application_id UUID NOT NULL REFERENCES job_applications(id)
        ON DELETE CASCADE,
    activity_id UUID NOT NULL REFERENCES activities(id)
        ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE job_application_follow_ups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_application_id UUID NOT NULL REFERENCES job_applications(id)
        ON DELETE CASCADE,
    follow_up_id UUID NOT NULL REFERENCES follow_ups(id)
        ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_applications_user_status
    ON applications (user_id, status);

CREATE INDEX idx_applications_user_follow_up
    ON applications (user_id, next_follow_up_on);

CREATE INDEX idx_follow_ups_due
    ON follow_ups (due_on)
    WHERE status = 'pending';

CREATE INDEX idx_activities_application_time
    ON activities (application_id, occurred_at DESC);
