-- Add missing columns to Profiles
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;

-- Add missing columns to Products
ALTER TABLE products ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS sku VARCHAR(100) UNIQUE;

-- Ensure Shifts outlet_id rename validation (in case 006 failed or partially ran)
-- (No-op if already renamed)

-- Fix Table Sessions
-- table_sessions links to tables. tables has outlet_id -> store_id.
-- table_sessions doesn't have outlet_id directly, it has table_id.
-- But sessions.sql might be joining or selecting something.
-- Let's check session query content next.
