// types/InventoryItem.ts
export interface InventoryItem {
    id: number,
    tanggal_pembelian: string,
    gudang_id: number,
    kategori_id: number,
    divisi_id: number,
    user_id: number,
    nama_barang: string,
    qty_barang: number,
    harga_pembelian: number,
    spesifikasi: string,
    qty_tersedia: number,
    qty_terpakai: number,
    total_nilai: number,
    upload_nota: string
}
