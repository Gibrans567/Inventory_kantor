import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FuseCardComponent } from '@fuse/components/card';
import { ApiService } from 'app/services/api.service';
import { MatButtonModule } from '@angular/material/button';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatDialog } from '@angular/material/dialog';
import { AddSebaranBarangComponent } from 'app/modules/admin/dashboards/sebaran-barang/add-sebaran-barang/add-sebaran-barang.component';




@Component({
  selector: 'app-detail-inventaris',
  standalone: true,
  imports: [ RouterLink,
    MatButtonModule,
    MatTableModule,
    MatSortModule,
    MatFormFieldModule,
    MatInputModule,
    MatPaginatorModule, // Make sure MatPaginatorModule is imported
    ReactiveFormsModule,
    CommonModule,
    MatMenuModule,
    FuseCardComponent,
    MatIconModule,
    MatSortModule,
    FormsModule,
    MatTableModule,
  ],
  templateUrl: './detail-inventaris.component.html',
  styleUrl: './detail-inventaris.component.scss'
})
export class DetailInventarisComponent {
    // name: string | null = null;
    // usersData: any
    name: string | null; // Simulated route parameter for example
    barangData: any;
    sebaranBarang: any;
    id: string| null = null;
    isLoading: boolean = false;


    constructor(
        private _apiService: ApiService,
        private route: ActivatedRoute,
        private _matDialog: MatDialog,
    ) { }

    ngOnInit(): void {
        this.route.paramMap.subscribe(params => {
            this.id = params.get("id")

            if (this.id) {
                this.getBarangById(this.id)
                this.getRiwayatBarangById(this.id)
            }
        })
    }



    async getBarangById(id: string) {
        try {
            // Fetch data from the API
            const barangData = await this._apiService.get(`/inventaris/barang/${id}`);
            console.log(barangData);

            // Mapping API data to output format
            this.barangData = {
                id: barangData.id || 'Memuat Data....',
                nama_barang: barangData.nama_barang || 'Memuat Data....',
                spesifikasi: barangData.spesifikasi || 'Memuat Data....',
                harga_pembelian: barangData.harga_pembelian || 'Memuat Data....',
                qty_barang: barangData.qty_barang || 'Memuat Data....',
                qty_tersedia: barangData.qty_tersedia || 'Barang Sudah Terpakai Semua',
                qty_terpakai: barangData.qty_terpakai || 'Tidak Ada Barang Yang Terpakai',
                kategori_nama: barangData.kategori_nama || 'Memuat Data....',
                nama_divisi: barangData.nama_divisi || 'Memuat Data....',
                nama_gudang: barangData.nama_gudang || 'Memuat Data....',
                tanggal_pembelian: barangData.tanggal_pembelian || 'Memuat Data....',
                total_nilai: barangData.total_nilai || 'Memuat Data....',
                sebaran_barang: barangData.sebaran_barang || [],
                upload_nota: "http://localhost:8080/" + barangData.upload_nota || 'Memuat Data....',
                created_at: barangData.created_at || 'Memuat Data....',
                updated_at: barangData.updated_at || 'Memuat Data....',
                divisi_id: barangData.divisi_id || 'Memuat Data....',
                gudang_id: barangData.gudang_id || 'Memuat Data....',
                kategori_id: barangData.kategori_id || 'Memuat Data....'
            };
        } catch (error) {
            // Default values when data is loading or there's an error
            this.barangData = {
                id: 'Memuat Data....',
                nama_barang: 'Memuat Data....',
                spesifikasi: 'Memuat Data....',
                harga_pembelian: 'Memuat Data....',
                qty_barang: 'Memuat Data....',
                qty_tersedia: 'Memuat Data....',
                qty_terpakai: 'Memuat Data....',
                kategori_nama: 'Memuat Data....',
                nama_divisi: 'Memuat Data....',
                nama_gudang: 'Memuat Data....',
                tanggal_pembelian: 'Memuat Data....',
                total_nilai: 'Memuat Data....',
                sebaran_barang: [],
                upload_nota: 'Memuat Data....',
                created_at: 'Memuat Data....',
                updated_at: 'Memuat Data....',
                divisi_id: 'Memuat Data....',
                gudang_id: 'Memuat Data....',
                kategori_id: 'Memuat Data....'
            };
        }
    }

    async getRiwayatBarangById(id: string) {
        try {
            // Fetch data from the API
            const riwayatData = await this._apiService.get(`/sebaranBarang/${id}`);
            console.log(riwayatData);

            // Mapping API data to output format
            this.sebaranBarang = riwayatData.map(item => ({
                barang: item.nama_barang || 'Belum Ada Barang',
                createdAt: item.created_at || 'Belum Ada Barang',
                divisi: item.nama_divisi || 'Belum Ada Barang',
                posisi_akhir: item.posisi_akhir || 'Belum Ada Barang',
                posisi_awal: item.posisi_awal || 'Tidak Memiliki Posisi Awal',
                qty_barang: item.qty_barang || 'Belum Ada Barang',
                updatedAt: item.updated_at || 'Belum Ada Barang', // Jika ada, kalau tidak bisa dihapus
                user: item.nama || 'Belum Ada Barang'
            }));

        } catch (error) {
            console.error('Gagal memuat data riwayat barang:', error);

            // Default values when data is loading or there's an error
            this.sebaranBarang = [{
                barang: 'Belum Ada Barang',
                createdAt: 'Belum Ada Barang',
                divisi: 'Belum Ada Barang',
                posisi_akhir: 'Belum Ada Barang',
                posisi_awal: 'Belum Ada Barang',
                qty_barang: 'Belum Ada Barang',
                updatedAt: 'Belum Ada Barang',
                user: 'Belum Ada Barang'
            }];
        }
    }

    TambahSebaranBarang() {
        // Pastikan ID barang dan data tersedia
        if (this.id && this.barangData) {
            this._matDialog.open(AddSebaranBarangComponent, {
                width: window.innerWidth < 600 ? '90%' : '50%',
                maxWidth: '100vw',
                data: this.barangData // Pass the entire barangData object to the dialog
            }).afterClosed().subscribe(result => {
                // Refresh data if needed
                if (result === 'refresh') {
                    this.getBarangById(this.id!);
                    this.getRiwayatBarangById(this.id!);
                }
            });
        } else {
            console.error('ID barang atau data tidak tersedia');
        }
    }

}
