CREATE TABLE checkout_order (
    id BIGSERIAL PRIMARY KEY,
    order_id text NOT NULL UNIQUE,
    session_id text,
    status text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE checkout_order_item (
    id BIGSERIAL PRIMARY KEY,
    order_id bigint NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    amount bigint NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);
