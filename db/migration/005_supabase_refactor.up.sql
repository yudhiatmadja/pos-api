-- 1. Simulate Supabase Auth
CREATE SCHEMA IF NOT EXISTS auth;
CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    encrypted_password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Profiles (Replaces/Links to Users)
-- We'll rename the existing 'users' table to 'profiles' to preserve data, 
-- but we need to adhere to the structure linked to auth.users.
-- Strategy: Use existing users table as public.profiles, but we need to link it to auth.users.
-- For this refactor, let's create a clean slate or migrate. 
-- Simplest: Rename users -> profiles. Add email from username if meaningful.
ALTER TABLE users RENAME TO profiles;

-- Adjust profiles columns
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS full_name VARCHAR(100);
-- 'role' column already exists in profiles (was users)
-- 'store_id' needs to be added (was not in users)

-- 3. Outlets -> Stores
ALTER TABLE outlets RENAME TO stores;

-- Link Profiles to Stores
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS store_id UUID REFERENCES stores(id) ON DELETE SET NULL;

-- 4. Products refinements
ALTER TABLE products RENAME COLUMN outlet_id TO store_id;

-- 5. Orders refinements
ALTER TABLE orders RENAME COLUMN outlet_id TO store_id;

-- 6. Shifts Refinements
ALTER TABLE shifts RENAME COLUMN start_time TO opened_at;
ALTER TABLE shifts RENAME COLUMN end_time TO closed_at;
ALTER TABLE shifts RENAME COLUMN start_cash TO opening_cash;
ALTER TABLE shifts RENAME COLUMN end_cash TO closing_cash;

-- 7. Audit Logs Refinements
ALTER TABLE audit_logs RENAME COLUMN entity_name TO entity;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS "before" JSONB;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS "after" JSONB;

-- 8. Payments Refinements (QRIS)
ALTER TABLE payments ADD COLUMN IF NOT EXISTS qris_url TEXT;

-- 9. Cleanup / Indexes
-- Drop old foreign keys if names confuse, but Postgres handles renames usually.
-- Ensure RBAC roles are correct
INSERT INTO roles (code, name) VALUES ('super_admin', 'Super Admin'), ('store_owner', 'Store Owner'), ('staff', 'Staff'), ('supplier', 'Supplier') ON CONFLICT (code) DO NOTHING;
-- We might want to migrate existing users to have auth.users counterparts if we were strictly simulating.
-- For now, we assume direct usage of profiles or mapped auth.
