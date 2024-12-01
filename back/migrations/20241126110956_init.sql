-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
                       name varchar(128) NOT NULL,
                       is_active bool DEFAULT true,
                       metadata jsonb,
                       created_at timestamptz DEFAULT now(),
                       updated_at timestamptz DEFAULT now(),
                       deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS users_created_idx ON users (created_at);
CREATE INDEX IF NOT EXISTS users_updated_idx ON users (updated_at);
CREATE INDEX IF NOT EXISTS users_deleted_idx ON users (deleted_at);
CREATE INDEX IF NOT EXISTS users_name_idx ON users (name);
CREATE INDEX IF NOT EXISTS users_id_idx ON users (id);

CREATE FUNCTION update_updated_at()
    RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION set_deleted_at()
    RETURNS trigger AS $$
BEGIN
    NEW.deleted_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE coupons (
                         id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
                         name varchar(255) NOT NULL,
                         reward varchar(255) NOT NULL,
                         max_redemptions integer NOT NULL,
                         times_redeemed integer NOT NULL,
                         redeem_by timestamptz,
                         metadata jsonb,
                         created_at timestamptz DEFAULT now(),
                         updated_at timestamptz DEFAULT now(),
                         deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS coupons_created_idx ON coupons (created_at);
CREATE INDEX IF NOT EXISTS coupons_updated_idx ON coupons (updated_at);
CREATE INDEX IF NOT EXISTS coupons_deleted_idx ON coupons (deleted_at);
CREATE INDEX IF NOT EXISTS coupons_times_redeemed_idx ON coupons (times_redeemed);
CREATE INDEX IF NOT EXISTS coupons_name_idx ON coupons (name);
CREATE INDEX IF NOT EXISTS coupons_id_idx ON coupons (id);

CREATE TRIGGER update_coupons_updated_at
    BEFORE UPDATE
    ON coupons
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at();

CREATE TABLE coupon_relation (
                                 id uuid DEFAULT uuid_generate_v4() NOT NULL,
                                 user_id uuid NOT NULL,
                                 coupon_id uuid NOT NULL,
                                 metadata jsonb,
                                 created_at timestamptz DEFAULT now(),
                                 updated_at timestamptz DEFAULT now(),
                                 deleted_at timestamptz,
                                 PRIMARY KEY (user_id, coupon_id),

                                 CONSTRAINT users_applied_coupon_fk
                                     FOREIGN KEY (user_id)
                                         REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,

                                 CONSTRAINT coupons_applied_user_fk
                                     FOREIGN KEY (coupon_id)
                                         REFERENCES coupons(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS coupon_relation_created_idx ON coupon_relation (created_at);
CREATE INDEX IF NOT EXISTS coupon_relation_updated_idx ON coupon_relation (updated_at);
CREATE INDEX IF NOT EXISTS coupon_relation_deleted_idx ON coupon_relation (deleted_at);
CREATE INDEX IF NOT EXISTS coupon_relation_user_id_idx ON coupon_relation (user_id);
CREATE INDEX IF NOT EXISTS coupon_relation_coupon_id_idx ON coupon_relation (coupon_id);
CREATE INDEX IF NOT EXISTS coupon_relation_id_idx ON coupon_relation (id);

CREATE FUNCTION update_coupon_relation_deleted_at()
    RETURNS trigger AS $$
BEGIN
    IF NEW.deleted_at IS DISTINCT FROM OLD.deleted_at THEN
        UPDATE coupon_relation
        SET deleted_at = NEW.deleted_at
        WHERE user_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_coupon_relation_deleted_at_trigger
    AFTER UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_coupon_relation_deleted_at();

INSERT INTO users (id, name, is_active, metadata)
VALUES
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d01', 'Alice', true, '{"test_name":"Alice_1"}'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d02', 'Mary', true, '{"test_name":"Mary_2"}'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d03', 'John', true, '{"test_name":"John_3"}'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d04', 'Bob', true, '{"test_name":"Bob_4"}');

INSERT INTO coupons (id, name, reward, max_redemptions, times_redeemed, redeem_by)
VALUES
    ('1fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d01', 'Coupon_test_1', 'free coffee 1', 0, 0, null),
    ('2fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d02', 'Coupon_test_2', 'free coffee 2', 0, 0, null),
    ('3fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d03', 'Coupon_test_3', 'free coffee 3', 0, 0, null),
    ('4fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d04', 'Coupon_test_4', 'free coffee 4', 0, 0, null),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d05', 'Coupon_test_5', 'free coffee 5', 1, 1, null),
    ('6fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d06', 'Coupon_test_6', 'free coffee 6', 10, 9, null);

INSERT INTO coupon_relation (user_id, coupon_id)
VALUES
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d01', '1fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d01'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d01', '2fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d02'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d02', '3fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d03'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d03', '4fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d04'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d03', '5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d05'),
    ('5fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d04', '6fad25aa-0e2c-4ceb-8e08-0d4e4f9e4d06');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table users cascade;
drop table coupons cascade;
drop table coupon_relation cascade;
drop function update_updated_at() cascade;
drop function update_coupon_relation_deleted_at() cascade;

-- +goose StatementEnd
