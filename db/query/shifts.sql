-- name: CreateShift :one
INSERT INTO shifts (
    user_id, outlet_id, opening_cash
) VALUES (
    $1, $2, $3
) RETURNING *;
-- Note: outlet_id in shifts table? 
-- Migration 005: ALTER TABLE shifts RENAME COLUMN outlet_id TO store_id?
-- Wait, I missed renaming outlet_id in shifts in migration 005.
-- I should check the migration file I wrote.
-- I wrote: ALTER TABLE products RENAME ..., ALTER TABLE orders RENAME ...
-- I did NOT write ALTER TABLE shifts RENAME COLUMN outlet_id TO store_id.
-- I only renamed start_time etc.
-- I need to fix this or use outlet_id in query. Ideally fix.
-- I will create a fix migration or just use outlet_id if time constrained.
-- "Clean Architecture" "Supabase Refactor" implies I should correct it.
-- But I cannot edit 005.up.sql after running it easily without down/up.
-- I will blindly use `store_id` in query here and assume I will fix DB in next step if error,
-- OR I will look at `003` which created `shifts`. 
-- 003: outlet_id UUID.
-- 005: didn't rename it.
-- So it is still `outlet_id` in DB.
-- I should strictly rename it to `store_id` for consistency.
-- I will create 006 migration immediately to rename remaining columns.

-- name: CloseShift :one
UPDATE shifts
SET closed_at = NOW(),
    closing_cash = $2,
    expected_cash = $3
WHERE id = $1
RETURNING *;

-- name: GetCurrentShift :one
SELECT * FROM shifts
WHERE user_id = $1 AND closed_at IS NULL
LIMIT 1;
