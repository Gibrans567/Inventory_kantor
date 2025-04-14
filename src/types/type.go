package types

import ("time"
    "github.com/golang-jwt/jwt/v5"
)
type Gudang struct {
    ID           uint        `json:"id" gorm:"primaryKey"`
    NamaGudang   string      `json:"nama_gudang"`
    LokasiGudang string      `json:"lokasi_gudang"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Inventaris   []Inventaris `gorm:"foreignKey:GudangID" json:"-"` // Tidak akan muncul di JSON
    Depresiasi   []Depresiasi `gorm:"foreignKey:IdGudang" json:"-"`
}


type Inventaris struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    TanggalPembelian  time.Time `json:"tanggal_pembelian"`
    GudangID          uint      `json:"gudang_id"`
    KategoriID        uint      `json:"kategori_id"`
    DivisiID          uint      `json:"divisi_id"`
    UserID            int      `json:"user_id"`
    Role              string    `json:"role"`
    NamaBarang        string    `json:"nama_barang"`
    QtyBarang         int       `json:"qty_barang"`
    HargaPembelian    int`json:"harga_pembelian"`
    Spesifikasi       string    `json:"spesifikasi"`
    QtyTersedia       int       `json:"qty_tersedia"`
    QtyTerpakai       int       `json:"qty_terpakai"`
    TotalNilai        int `json:"total_nilai"`
    UploadNota        string    `json:"upload_nota"`  // Menyimpan path file nota
    CreatedAt         time.Time
    UpdatedAt         time.Time
    User              *User      `gorm:"foreignKey:UserID" json:"-"`    // One to Many
    Gudang            Gudang    `gorm:"foreignKey:GudangID" json:"-"`  // One to Many
    Divisi            Divisi    `gorm:"foreignKey:DivisiID" json:"-"`  // One to Many
    Kategori          Kategori  `gorm:"foreignKey:KategoriID" json:"-"`  // One to Many
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdBarang" json:"-"`   // Many to Many
}

type Divisi struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    NamaDivisi        string    `json:"nama_divisi"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    User              []User      `gorm:"foreignKey:IdDivisi" json:"-"`  // One to Many
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdDivisi" json:"-"`  // One to Many
    Inventaris        []Inventaris `gorm:"foreignKey:DivisiID" json:"-"`  // One to Many
}

type User struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    IdDivisi          uint      `json:"id_divisi"`
    Email             string    `json:"email" binding:"required,email"`
    Password          string    `json:"password"`
    Name              string    `json:"nama_user"`
    Role              string    `json:"role"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    Divisi            Divisi     `gorm:"foreignKey:IdDivisi" json:"-"`  // Many to One
    SebaranBarang     []SebaranBarang `gorm:"foreignKey:IdUser" json:"-"`  // One to Many
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
    Divisi            Divisi     `gorm:"foreignKey:IdDivisi" json:"-"`  // Many to One
    User              *User      `gorm:"foreignKey:IdUser" json:"-"`    // Many to One
    Inventaris        Inventaris `gorm:"foreignKey:IdBarang" json:"-"` // Many to Many
}

type Kategori struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    NamaKategori     string    `json:"nama_kategori"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    Inventaris       []Inventaris `gorm:"foreignKey:KategoriID" json:"-"`  // One to Many
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
    Gudang           Gudang    `gorm:"foreignKey:IdGudang" json:"-"`  // Many to One
    Inventaris       Inventaris `gorm:"foreignKey:IdBarang" json:"-"` // Many to One
}

type History struct {
    ID               uint      `json:"id" gorm:"primaryKey"`
    Kategori         string      `json:"kategori"`
    Keterangan       string      `json:"keterangan"`
    CreatedAt        time.Time   `json:"created_at"`
}

type JadwalDepresiasi struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	IdBarang         uint      `json:"id_barang"`
	NextRun          time.Time `json:"next_run"`
}

type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type DeleteAllByTimeframeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

//audit

//ekstrak 

//dizip lalu disimpen di foldering untuk file upload atau data dari database

//Halaman data barang masing masing

//docker hub