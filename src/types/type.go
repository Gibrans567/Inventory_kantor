package types

import "time"


type Penyusutan struct {
	NilaiPenyusutan int `json:"nilai_penyusutan"`
	HargaPenyusutan int `json:"harga_penyusutan"`
}

type Inventory struct {
    ID            uint           `json:"id" gorm:"primaryKey"`
    NamaPenggunaAwal     string  `json:"nama_admin"`
    TanggalBeli   time.Time      `json:"tanggal_beli"`
    NamaBarang    string         `json:"nama_barang"`
    Qty           int            `json:"qty"`
    Nilai         float64        `json:"nilai"`
    PosisiAwal    string         `json:"posisi_awal"`
    NamaPenggunaAkhir    *string `json:"nama_pengguna_akhir"`
    TanggalMutasi time.Time      `json:"tanggal_mutasi"`
    PosisiAkhir   *string        `json:"posisi_akhir"`
    Status       string          `json:"status"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}