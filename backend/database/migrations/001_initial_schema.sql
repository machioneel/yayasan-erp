-- Migration: 001_initial_schema
-- Description: Create all initial tables for ERP system
-- Date: 2025-11-08

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create all tables (copy from schema.sql core tables)
-- This ensures incremental migration approach

\i ../schema.sql

-- Mark migration as completed
INSERT INTO schema_migrations (version, name, executed_at)
VALUES ('001', 'initial_schema', CURRENT_TIMESTAMP)
ON CONFLICT (version) DO NOTHING;
