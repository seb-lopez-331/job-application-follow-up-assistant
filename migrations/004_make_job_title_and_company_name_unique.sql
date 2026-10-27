ALTER TABLE jobs
ADD CONSTRAINT jobs_role_company_unique UNIQUE (role_title, company_name);