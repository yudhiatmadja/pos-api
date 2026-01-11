-- Fix Shifts
ALTER TABLE shifts RENAME COLUMN outlet_id TO store_id;

-- Fix Tables (restaurant tables)
ALTER TABLE tables RENAME COLUMN outlet_id TO store_id;

-- Fix Categories
ALTER TABLE categories RENAME COLUMN outlet_id TO store_id;

-- Fix Audit Log if needed (no store_id there usually)
