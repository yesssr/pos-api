CREATE TYPE pembayaran_metode AS ENUM ('cash', 'qris', 'debit', 'kredit');
CREATE TYPE pembayaran_status AS ENUM ('pending', 'lunas', 'gagal');

CREATE TABLE penjualan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_user UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    id_pelanggan UUID REFERENCES pelanggan(id) ON DELETE SET NULL,
    tanggal TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    total NUMERIC(12,2) NOT NULL,
    metode_pembayaran pembayaran_metode NOT NULL DEFAULT 'cash',
    status_pembayaran pembayaran_status NOT NULL DEFAULT 'pending',
    id_transaksi_gateway VARCHAR(100)
);
