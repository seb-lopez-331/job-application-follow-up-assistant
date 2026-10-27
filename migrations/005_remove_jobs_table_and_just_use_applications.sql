DROP TABLE jobs;

ALTER TABLE job_applications
    ADD role_title TEXT NOT NULL,
    ADD company_name TEXT NOT NULL,
    ADD date_posted TIMESTAMPTZ NOT NULL,
    ADD loc TEXT,
    ADD descr TEXT NOT NULL,
    ADD job_url TEXT,
    ADD work_arrangement TEXT CHECK (
        work_arrangement IN ('remote', 'hybrid', 'on-site')
    ),
    DROP job_id,
    ADD CONSTRAINT job_applications_user_role_company_unique UNIQUE (
        user_id, role_title, company_name
    );
