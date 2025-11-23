CREATE TABLE detail_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_transaction UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    id_product UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    qty INT NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    subtotal NUMERIC(12,2) NOT NULL
);
