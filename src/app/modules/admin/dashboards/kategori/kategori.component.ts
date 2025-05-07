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
import { AddKategoriComponent } from 'app/modules/admin/dashboards/kategori/add-kategori/add-kategori.component';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-kategori',
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
  templateUrl: './kategori.component.html',
  styleUrls: ['./kategori.component.scss']
})
export class KategoriComponent {
    displayedColumns: string[] = ['nama_kategori', 'action']; // Tampilkan kolom nama kategori dan aksi

    dataSource = new MatTableDataSource<any>([]);
    isLoading: boolean = false;
    isNotDataFound: boolean = false;

    logAktivitas: any[] = []; // Menyimpan data kategori
    kategoriItems: { name: string, checked: boolean }[] = [];


    @ViewChild(MatPaginator) paginator!: MatPaginator;
    @ViewChild(MatSort) sort!: MatSort;

    constructor(private apiService: ApiService, private route: ActivatedRoute, private matDialog: MatDialog,
      private fuseConfirmationService: FuseConfirmationService) {}

    ngOnInit(): void {
      this.route.paramMap.subscribe(() => {
        this.getKategori(); // Mengambil data kategori saat komponen diinisialisasi
      });
    }

    async getKategori(): Promise<void> {
      this.isLoading = true;
      this.isNotDataFound = false;

      try {
        const response = await this.apiService.get(`/kategori`);

        // Memastikan response berhasil dan mengakses data kategori
        const kategoriData = response?.data ?? [];  // Mengakses 'data' dari response

        // Memperbarui logAktivitas dan dataSource
        this.logAktivitas = kategoriData.map((item) => ({
          id: item.id,
          nama_kategori: item.nama_kategori,
          createdAt: item.CreatedAt,
          updatedAt: item.UpdatedAt,
        }));

        // Update dataSource dengan data kategori terbaru
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

      // Custom filter berdasarkan 'nama_kategori' dan 'createdAt'
      this.dataSource.filterPredicate = (data: any, filter: string) => {
        const transformedFilter = filter.trim().toLowerCase();
        return (
          data.nama_kategori?.toLowerCase().includes(transformedFilter) ||
          data.createdAt?.toLowerCase().includes(transformedFilter) ||
          data.updatedAt?.toLowerCase().includes(transformedFilter)
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

    addNewCategory() {
      const dialogRef = this.matDialog.open(AddKategoriComponent, {
        width: window.innerWidth < 600 ? '90%' : '50%',
        maxWidth: '100vw',
      });

      dialogRef.afterClosed().subscribe((result) => {
        // Mengecek apakah hasil dialog adalah 'refresh' yang menandakan kategori berhasil ditambahkan
        if (result === 'refresh') {
          this.getKategori(); // Memanggil fungsi getKategori untuk mendapatkan data terbaru
        }
      });
    }

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
            this.getKategori(); // Memanggil fungsi untuk memperbarui data kategori dengan data terbaru

            // Memberi feedback ke user bahwa kategori berhasil dihapus
            console.log(`Kategori "${namaKategori}" berhasil dihapus`);
          } catch (error) {
            console.error('Gagal menghapus kategori', error);
          }
        }
      });
    }
}
