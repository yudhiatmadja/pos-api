-- Enable UUID extension (already in 001 but good to ensure availability)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ROLES & RBAC
CREATE TABLE roles (
    code VARCHAR(50) PRIMARY KEY, -- ADMIN, KASIR, KITCHEN, SUPERVISOR
    name VARCHAR(100) NOT NULL,
    description TEXT
);

INSERT INTO roles (code, name) VALUES 
('ADMIN', 'Administrator'),
('KASIR', 'Cashier'),
('KITCHEN', 'Kitchen Staff'),
('SUPERVISOR', 'Supervisor');

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_code VARCHAR(50) NOT NULL REFERENCES roles(code) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, role_code)
);

-- OUTLETS
CREATE TABLE outlets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- RESTAURANT (TABLES)
CREATE TABLE tables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL, -- e.g., "Table 1", "VIP 2"
    capacity INT DEFAULT 4,
    qr_code VARCHAR(255) UNIQUE, -- Unique string for QR generation
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- TABLE SESSIONS (QR ORDERING)
CREATE TABLE table_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE, -- Session token for client
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- INVENTORY (CATEGORIES & UPDATES)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE, -- Optional if categories are global
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE products 
ADD COLUMN outlet_id UUID REFERENCES outlets(id) ON DELETE CASCADE,
ADD COLUMN category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
ADD COLUMN image_url TEXT,
ADD COLUMN is_available BOOLEAN DEFAULT TRUE;

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL, -- positive for add, negative for deduct
    type VARCHAR(50) NOT NULL, -- IN, OUT, ADJUSTMENT, SALE
    reference_id UUID, -- order_id or other ref
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ORDERS
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    table_session_id UUID REFERENCES table_sessions(id) ON DELETE SET NULL, -- Nullable for POS orders
    cashier_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Nullable for QR orders initially
    
    order_number VARCHAR(50) NOT NULL, -- Human readable ID
    status VARCHAR(50) NOT NULL DEFAULT 'NEW', -- NEW, ACCEPTED, COOKING, READY, DONE, VOIDED
    payment_status VARCHAR(50) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PAID, REFUNDED
    
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    final_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_name VARCHAR(100) NOT NULL, -- Snapshot name
    product_price DECIMAL(10, 2) NOT NULL, -- Snapshot price
    quantity INT NOT NULL DEFAULT 1,
    total_price DECIMAL(10, 2) NOT NULL,
    note TEXT
);

-- PAYMENTS
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    payment_method VARCHAR(50) NOT NULL, -- CASH, QRIS, PAY_LATER
    amount DECIMAL(10, 2) NOT NULL,
    reference_number VARCHAR(100), -- Transaction ID from Gateway
    status VARCHAR(50) NOT NULL, -- PENDING, SUCCESS, FAILED
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- SHIFTS (KASIR)
CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    outlet_id UUID NOT NULL REFERENCES outlets(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    start_cash DECIMAL(10, 2) NOT NULL DEFAULT 0,
    end_cash DECIMAL(10, 2),
    expected_cash DECIMAL(10, 2)
);

-- AUDIT LOGS
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_name VARCHAR(100),
    entity_id UUID,
    details JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_orders_outlet_status ON orders(outlet_id, status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_table_sessions_token ON table_sessions(token);
