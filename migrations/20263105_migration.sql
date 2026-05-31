-- SQL migration script for Sentinel database
-- This creates the necessary tables for Domain models

-- Create domains table
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) UNIQUE NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_domains_active ON domains(active);
CREATE INDEX IF NOT EXISTS idx_domains_domain ON domains(domain);

-- Insert sample data (optional - uncomment to use)
-- Domains
-- INSERT INTO domains (domain, active) VALUES 
--     ('yunusemrealpu.netlify.app:443', true),
--     ('google.com', true),
--     ('https://github.com', true)
-- ON CONFLICT (domain) DO NOTHING;

COMMENT ON TABLE domains IS 'Stores monitored SSL/TLS certificate domains';
COMMENT ON COLUMN domains.domain IS 'Domain in any format: domain.com, https://example.com, domain.com:443';