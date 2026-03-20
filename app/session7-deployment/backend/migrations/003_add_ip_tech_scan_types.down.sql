ALTER TABLE scan_jobs DROP CONSTRAINT IF EXISTS scan_jobs_scan_type_check;
ALTER TABLE scan_jobs ADD CONSTRAINT scan_jobs_scan_type_check
    CHECK (scan_type IN ('subdomain','dns','whois','port','asn','ssl','all'));