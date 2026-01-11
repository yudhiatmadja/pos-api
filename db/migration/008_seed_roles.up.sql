-- Seed Roles
INSERT INTO roles (code, name, description) VALUES
('SUPER_ADMIN', 'Super Admin', 'Platform Owner, full access to all stores'),
('STORE_OWNER', 'Store Owner', 'Owner of a specific store, can approve voids/refunds'),
('KASIR', 'Cashier', 'Front POS operations, payment, shifts'),
('KITCHEN', 'Kitchen Staff', 'View orders and update cooking status'),
('STAFF', 'Waiter/Staff', 'Helper, can only create orders')
ON CONFLICT (code) DO UPDATE 
SET name = EXCLUDED.name, description = EXCLUDED.description;
