CREATE TABLE detail_penjualan (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_penjualan UUID NOT NULL REFERENCES penjualan(id) ON DELETE CASCADE,
    id_barang UUID NOT NULL REFERENCES barang(id) ON DELETE RESTRICT,
    jumlah INT NOT NULL,
    harga_saat_transaksi NUMERIC(12,2) NOT NULL,
    subtotal NUMERIC(12,2) NOT NULL
);
