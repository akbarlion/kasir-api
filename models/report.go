package models

type Report struct {
	TotalRevenue     int            `json:"total_revenue"`
	TotalTransaction int            `json:"total_transaksi"`
	ProdukTerlaris   ProdukTerlaris `json:"produk_terlaris"`
}

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
