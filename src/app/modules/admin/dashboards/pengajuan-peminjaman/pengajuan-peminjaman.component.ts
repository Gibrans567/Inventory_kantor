import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { RouterLink } from '@angular/router';
import { FuseCardComponent } from '@fuse/components/card';
import { ApiService } from 'app/services/api.service';
import { MatButtonModule } from '@angular/material/button';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatDialog } from '@angular/material/dialog';
import { AddSebaranBarangComponent } from 'app/modules/admin/dashboards/sebaran-barang/add-sebaran-barang/add-sebaran-barang.component';
import { AddBarangRusakComponent } from 'app/modules/admin/dashboards/barang-rusak/add-barang-rusak/add-barang-rusak.component';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { AddPengajuanPeminjamanComponent } from './add-pengajuan-peminjaman/add-pengajuan-peminjaman.component';

@Component({
  selector: 'app-pengajuan-peminjaman',
  standalone: true,
  imports: [
    RouterLink,
    MatButtonModule,
    MatTableModule,
    MatSortModule,
    MatFormFieldModule,
    MatInputModule,
    MatPaginatorModule,
    ReactiveFormsModule,
    CommonModule,
    MatMenuModule,
    FuseCardComponent,
    MatIconModule,
    MatSortModule,
    FormsModule,
    MatTableModule,
  ],
  templateUrl: './pengajuan-peminjaman.component.html',
  styleUrl: './pengajuan-peminjaman.component.scss'
})
export class PengajuanPeminjamanComponent {

    displayedColumns: string[] = ['nama_barang', 'name', 'qty_barang', 'nama_divisi', 'status_permohonan', 'note', 'Yang Menyetujui','action'];
    dataSource = new MatTableDataSource<any>([]);

    @ViewChild(MatPaginator) paginator: MatPaginator;
    @ViewChild(MatSort) sort: MatSort;

    pengajuanData: any[] = [];
    isLoading: boolean = false;

    constructor(
        private _apiService: ApiService,
        private _matDialog: MatDialog,
        private _fuseConfirmationService: FuseConfirmationService,
    ) { }

    ngOnInit(): void {
        this.getPengajuanData();
    }

    async getPengajuanData() {
        try {
            this.isLoading = true;
            const response = await this._apiService.get('/pengajuan');

            if (response.status === 'success' && response.data) {
                this.pengajuanData = response.data.map(item => ({
                    id: item.id || 'N/A',
                    nama_barang: item.nama_barang || 'N/A',
                    name: item.name || 'N/A',
                    nama_divisi: item.nama_divisi || 'N/A',
                    id_barang: item.id_barang || 'N/A',
                    id_user: item.id_user || 'N/A',
                    id_divisi: item.id_divisi || 'N/A',
                    id_approver: item.id_approver || null,
                    nama_approver: item.nama_approver || 'Belum ada yang Menyetujui',
                    status_kepemilikan: item.status_kepemilikan || '',
                    tanggal_pengajuan: item.tanggal_pengajuan || 'N/A',
                    status_permohonan: item.status_permohonan || 'Menunggu Approve',
                    status_pengembalian: item.status_pengembalian || '',
                    qty_barang: item.qty_barang || 0,
                    note: item.note || '',
                    posisi_akhir: item.posisi_akhir || 'N/A',
                    createdAt: item.CreatedAt || 'N/A',
                    updatedAt: item.UpdatedAt || 'N/A'
                }));

                // Set data ke dataSource
                this.dataSource.data = this.pengajuanData;
                this.dataSource.paginator = this.paginator;
                this.dataSource.sort = this.sort;
            } else {
                console.error('Format response tidak sesuai:', response);
                this.pengajuanData = [];
            }
        } catch (error) {
            console.error('Gagal memuat data pengajuan:', error);
            this.pengajuanData = [];
        } finally {
            this.isLoading = false;
        }
    }

    TambahPengajuan() {
        this._matDialog.open(AddPengajuanPeminjamanComponent, {
            width: window.innerWidth < 600 ? '90%' : '50%',
            maxWidth: '100vw',
            data: { mode: 'pengajuan' } // Pass mode pengajuan to the dialog
        }).afterClosed().subscribe(result => {
            // Refresh data if needed
            if (result === 'refresh') {
                this.getPengajuanData();
            }
        });
    }

    applySearchFilter(search: string): void {
        this.dataSource.filter = search.trim().toLowerCase();

        if (this.dataSource.paginator) {
            this.dataSource.paginator.firstPage();
        }
    }

    deletePengajuan(pengajuan: any) {
        const confirm = this._fuseConfirmationService.open({
          title: 'Konfirmasi Hapus',
          message: `Apakah Anda ingin menghapus pengajuan ${pengajuan.nama_barang} oleh ${pengajuan.name}?`,
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
              // Hapus data pengajuan
              await this._apiService.delete(`/pengajuan/${pengajuan.id}`);

              // Refresh data
              this.getPengajuanData();

              console.log(`Pengajuan "${pengajuan.nama_barang}" berhasil dihapus`);
            } catch (error) {
              console.error('Gagal menghapus pengajuan', error);
            }
          }
        });
    }

    editPengajuan(id: number) {
        // Ambil data pengajuan berdasarkan ID
        const pengajuan = this.pengajuanData.find(item => item.id === id);

        if (pengajuan) {
            this._matDialog.open(AddBarangRusakComponent, {
                width: window.innerWidth < 600 ? '90%' : '50%',
                maxWidth: '100vw',
                data: pengajuan // Kirim data pengajuan yang sesuai
            }).afterClosed().subscribe(result => {
                // Refresh data jika perlu
                if (result === 'refresh') {
                    this.getPengajuanData();
                }
            });
        } else {
            console.error('Pengajuan dengan ID tersebut tidak ditemukan');
        }
    }

    approvePengajuan(pengajuan: any) {
        const confirm = this._fuseConfirmationService.open({
          title: 'Konfirmasi Approve',
          message: `Apakah Anda ingin menyetujui pengajuan ${pengajuan.nama_barang} oleh ${pengajuan.name}?`,
          actions: {
            confirm: {
              label: 'Approve',
            },
            cancel: {
              label: 'Batal',
            },
          },
        });

        confirm.afterClosed().subscribe(async (result) => {
          if (result === 'confirmed') {
            try {
              // Update status pengajuan menjadi disetujui
              await this._apiService.put(`/pengajuan/${pengajuan.id}`, {
                ...pengajuan,
                status_permohonan: 'Disetujui',
                id_approver: 1 // Ganti dengan ID user yang sedang login
              });

              // Refresh data
              this.getPengajuanData();

              console.log(`Pengajuan "${pengajuan.nama_barang}" berhasil disetujui`);
            } catch (error) {
              console.error('Gagal menyetujui pengajuan', error);
            }
          }
        });
    }

    rejectPengajuan(pengajuan: any) {
        const confirm = this._fuseConfirmationService.open({
          title: 'Konfirmasi Tolak',
          message: `Apakah Anda ingin menolak pengajuan ${pengajuan.nama_barang} oleh ${pengajuan.name}?`,
          actions: {
            confirm: {
              label: 'Tolak',
            },
            cancel: {
              label: 'Batal',
            },
          },
        });

        confirm.afterClosed().subscribe(async (result) => {
          if (result === 'confirmed') {
            try {
              // Update status pengajuan menjadi ditolak
              await this._apiService.put(`/pengajuan/${pengajuan.id}`, {
                ...pengajuan,
                status_permohonan: 'Ditolak',
                id_approver: 1 // Ganti dengan ID user yang sedang login
              });

              // Refresh data
              this.getPengajuanData();

              console.log(`Pengajuan "${pengajuan.nama_barang}" berhasil ditolak`);
            } catch (error) {
              console.error('Gagal menolak pengajuan', error);
            }
          }
        });
    }
}
