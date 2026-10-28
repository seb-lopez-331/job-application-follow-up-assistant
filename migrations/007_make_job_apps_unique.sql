ALTER TABLE job_applications
ADD CONSTRAINT job_applications_user_company_role_unique UNIQUE (user_id, company_name, role_title);