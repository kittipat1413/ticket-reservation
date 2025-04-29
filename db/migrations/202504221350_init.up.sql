-- 202504221337_init.up.sql

-- Concerts Table
CREATE TABLE concerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    venue TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER concerts_updated_at_modtime BEFORE UPDATE ON concerts FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Zones Table
CREATE TABLE zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    concert_id UUID NOT NULL REFERENCES concerts(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER zones_updated_at_modtime BEFORE UPDATE ON zones FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Seats Table
CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    seat_number TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('available', 'pending', 'booked')),
    locked_until TIMESTAMPTZ,
    locked_by_session_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(zone_id, seat_number)
);
CREATE TRIGGER seats_updated_at_modtime BEFORE UPDATE ON seats FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Reservations Table
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seat_id UUID NOT NULL REFERENCES seats(id) ON DELETE CASCADE,
    session_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'expired')),
    reserved_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER reservations_updated_at_modtime BEFORE UPDATE ON reservations FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Payments Table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('initiated', 'paid', 'failed')),
    amount DECIMAL(10, 2),
    paid_at TIMESTAMPTZ,
    payment_method TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER payments_updated_at_modtime BEFORE UPDATE ON payments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();