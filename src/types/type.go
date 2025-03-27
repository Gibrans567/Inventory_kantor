package types

import "time"

type Penyusutan struct {
	NilaiPenyusutan int `json:"nilai_penyusutan"`
	HargaPenyusutan int `json:"harga_penyusutan"`
}

type Gudang struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    NamaGudang     string    `json:"nama_gudang"`
    LokasiGudang   string    `json:"lokasi_gudang"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
    Inventaris     []Inventaris `gorm:"foreignKey:GudangID"`  // One to Many
    Depresiasi     []Depresiasi `gorm:"foreignKey:IdGudang"`   // One to Many
}

type Inventaris struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    TanggalPembelian  time.Time `json:"tanggal_pembelian"`
    GudangID          uint      `json:"gudang_id"`
    KategoriID        uint      `json:"kategori_id"`
    NamaBarang        string    `json:"nama_barang"`
    QtyBarang         int       `json:"qty_barang"`
    HargaPembelian    float64   `json:"harga_pembelian"`
    Spesifikasi       string    `json:"spesifikasi"`
    QtyTersedia       int       `json:"qty_tersedia"`
    QtyTerpakai       int       `json:"qty_terpakai"`
    TotalNilai        float64   `json:"total_nilai"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    Gudang            Gudang    `gorm:"foreignKey:GudangID"`  // One to Many
    Kategori          Kategori  `gorm:"foreignKey:KategoriID"`  // One to Many
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdBarang"`   // Many to Many
    Depresiasi        []Depresiasi `gorm:"foreignKey:IdBarang"`   // One to Many
}

type Divisi struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    NamaDivisi        string    `json:"nama_divisi"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    User              []User      `gorm:"foreignKey:IdDivisi"`  // One to Many
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdDivisi"`  // One to Many
}

type User struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    IdDivisi          uint      `json:"id_divisi"`
    Username          string    `json:"username"`
    Password          string    `json:"password"`
    NamaUser          string    `json:"nama_user"`
    Role              string    `json:"role"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    Divisi            Divisi     `gorm:"foreignKey:IdDivisi"`  // Many to One
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdUser"`  // One to Many
}

type SebaranBarang struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    IdDivisi          uint      `json:"id_divisi"`
    IdBarang          uint      `json:"id_barang"`
    IdUser            uint      `json:"id_user"`
    QtyBarang         int       `json:"qty_barang"`
    PosisiAwal        string    `json:"posisi_awal"`
    PosisiAkhir       string    `json:"posisi_akhir"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    Divisi            Divisi     `gorm:"foreignKey:IdDivisi"`  // Many to One
    User              User      `gorm:"foreignKey:IdUser"`    // Many to One
    Inventaris        Inventaris `gorm:"foreignKey:IdBarang"` // Many to Many
}

type Kategori struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    NamaKategori     string    `json:"nama_kategori"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    Inventaris       []Inventaris `gorm:"foreignKey:KategoriID"`  // One to Many
}

type Depresiasi struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    IdGudang         uint      `json:"id_gudang"`
    IdBarang         uint      `json:"id_barang"`
    HargaDepresiasi  int       `json:"harga_depresiasi"`
    Perbulan         int       `json:"perbulan"`
    Tahun            int       `json:"tahun"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    Gudang           Gudang    `gorm:"foreignKey:IdGudang"`  // Many to One
    Inventaris       Inventaris `gorm:"foreignKey:IdBarang"` // Many to One
}
