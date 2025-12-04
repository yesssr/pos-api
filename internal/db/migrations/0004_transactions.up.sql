CREATE TYPE payment_method AS ENUM ('cash', 'qris', 'debit', 'kredit');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  id_user UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  id_customer UUID REFERENCES customers(id) ON DELETE SET NULL,
  date DATE DEFAULT NOW(),
  total NUMERIC(12,2) NOT NULL,
  payment_method payment_method NOT NULL DEFAULT 'cash',
  payment_status payment_status NOT NULL DEFAULT 'pending',
  id_transaction_gateway VARCHAR(100),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
