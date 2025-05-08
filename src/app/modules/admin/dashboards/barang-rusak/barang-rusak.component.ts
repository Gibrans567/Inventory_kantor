import { FuseConfirmationService } from './../../../../../@fuse/services/confirmation/confirmation.service';
import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FuseCardComponent } from '@fuse/components/card';
import { IonicModule } from '@ionic/angular';
import { ApiService } from 'app/services/api.service';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-barang-rusak',
  standalone: true,
  imports: [
    RouterLink,
    MatTableModule,
    MatPaginator,
    MatSort,
    MatIconModule,
    CommonModule,
    FuseCardComponent,
    ReactiveFormsModule,
    MatTableModule,
    MatMenuModule,
    IonicModule,
    MatDialogModule,
    FormsModule,
    MatButtonModule,
  ],
  templateUrl: './barang-rusak.component.html',
  styleUrl: './barang-rusak.component.scss'
})
export class BarangRusakComponent {
  // Mengubah kolom yang ditampilkan di tabel
  displayedColumns: string[] = ['nama_barang', 'note','qty', 'status', 'action']; // Tampilkan kolom nama barang, note, status, dan aksi

  dataSource = new MatTableDataSource<any>([]);
  isLoading: boolean = false;
  isNotDataFound: boolean = false;

  logAktivitas: any[] = []; // Menyimpan data barang
  kategoriItems: { name: string, checked: boolean }[] = [];

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor(
    private apiService: ApiService, 
    private route: ActivatedRoute, 
    private matDialog: MatDialog,
    private fuseConfirmationService: FuseConfirmationService
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(() => {
      this.getBarangStatus(); // Mengambil data barang status saat komponen diinisialisasi
    });
  }

  // Fungsi untuk mengambil data barang status dari API
  async getBarangStatus(): Promise<void> {
    this.isLoading = false;
    this.isNotDataFound = false;

    try {
      // Mengambil data barang status dari API
      const response = await this.apiService.get(`/barangStatus`);
      
      // Memastikan response berhasil dan mengakses data barang
      const barangData = response?.dataBarang ?? [];  // Mengakses 'dataBarang' dari response

      // Memperbarui logAktivitas dan dataSource
      this.logAktivitas = barangData.map((item) => ({
        id: item.id,
        nama_barang: item.nama_barang,
        note: item.note,  // Menambahkan kolom note
        status: item.status,
      }));

      // Update dataSource dengan data barang terbaru
      this.dataSource.data = this.logAktivitas;

      // Set flag jika tidak ada data ditemukan
      this.isNotDataFound = this.logAktivitas.length === 0;
    } catch (error) {
      console.error('Failed to fetch data:', error);
      this.isNotDataFound = true;
    } finally {
      this.isLoading = false;
    }
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;

    // Custom filter berdasarkan 'nama_barang', 'note', dan 'status'
    this.dataSource.filterPredicate = (data: any, filter: string) => {
      const transformedFilter = filter.trim().toLowerCase();
      return (
        data.nama_barang?.toLowerCase().includes(transformedFilter) ||
        data.note?.toLowerCase().includes(transformedFilter) || // Filter berdasarkan note
        data.status?.toLowerCase().includes(transformedFilter)
      );
    };
  }

  // Fungsi untuk menangani perubahan halaman
  onPageChange(event: any): void {
    this.dataSource.paginator = this.paginator; // Pastikan paginator diterapkan saat halaman berubah
  }

  // Fungsi untuk mencari data berdasarkan filter
  applySearchFilter(search: string): void {
    this.dataSource.filter = search.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  // Fungsi untuk menambahkan kategori baru (Anda bisa sesuaikan dengan fungsionalitas yang diinginkan)

  // Fungsi untuk menghapus kategori (sesuaikan sesuai dengan API dan fungsionalitas Anda)
  deleteCategoryByName(namaKategori: string) {
    // Menanyakan konfirmasi penghapusan
    const confirm = this.fuseConfirmationService.open({
      title: 'Konfirmasi Hapus',
      message: `Apakah Anda yakin ingin menghapus kategori "${namaKategori}"?`,
      actions: {
        confirm: {
          label: 'Hapus',
        },
        cancel: {
          label: 'Batal',
        },
      },
    });

    confirm.afterClosed().subscribe(async (result) => {
      if (result === 'confirmed') {
        try {
          // Panggil API untuk menghapus kategori berdasarkan nama
          await this.apiService.delete(`/kategori/${namaKategori}`);

          // Mengambil data kategori terbaru setelah penghapusan
          this.getBarangStatus(); // Memanggil fungsi untuk memperbarui data kategori dengan data terbaru

          // Memberi feedback ke user bahwa kategori berhasil dihapus
          console.log(`Kategori "${namaKategori}" berhasil dihapus`);
        } catch (error) {
          console.error('Gagal menghapus kategori', error);
        }
      }
    });
  }
}
