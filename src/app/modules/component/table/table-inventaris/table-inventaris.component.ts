import { CommonModule } from '@angular/common';
import { Component, ViewChild, OnInit, Inject, inject, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuContent, MatMenuModule } from '@angular/material/menu';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { FuseCardComponent } from '@fuse/components/card';
import { ToastrService } from 'ngx-toastr';
import { ApiService } from 'app/services/api.service';
import { IonLoading } from '@ionic/angular/standalone';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { MatDialog } from '@angular/material/dialog';
import { InventoryService } from './table-inventaris.service';
import { AddInventarisComponent } from 'app/modules/admin/dashboards/inventaris/add-inventaris/add-inventaris.component';
import { AddDivisiComponent } from 'app/modules/admin/dashboards/divisi/add-divisi/add-divisi.component';
import { RouterLink } from '@angular/router';



@Component({
  selector: 'app-table-inventaris',
  standalone: true,
  imports: [CommonModule,FuseCardComponent,MatButtonModule,MatIconModule,MatInputModule,
    MatPaginatorModule,MatSortModule,MatTableModule,MatMenuModule,MatSelectModule,MatCheckboxModule,
    FormsModule,RouterLink,
  ],
  templateUrl: './table-inventaris.component.html',
  styleUrl: './table-inventaris.component.scss'
})
export class TableInventarisComponent implements OnInit{
    displayedColumns = [
        'tanggal_pembelian',
        'nama_barang',
        'nama_divisi',
        'harga_pembelian',
        'qty_barang',
        'total_nilai',
        'action'
      ];

    dataSource = new MatTableDataSource<any>([]);
    isLoadingDelete: boolean = false;
    isNotDataFound: boolean

    filterValues = {
        divisi: [],
        search: ''
      }

      divisiItems: { name: string, checked: boolean }[] = []

        @ViewChild(MatPaginator) paginator: MatPaginator;
        @ViewChild(MatSort) sort: MatSort;

        constructor(
            private fuseConfirmationService: FuseConfirmationService,
            private _matDialog: MatDialog,
            private apiService: ApiService,
            private _tableUserService: InventoryService,
          )
          {effect(() => {
                // Akses data dari service
                this.dataSource.data = this._tableUserService.inventoryItems();
                this.isNotDataFound = this._tableUserService.isNotFound();
              });
          }

          ngOnInit(): void {
            this._tableUserService.fetchData()
            this.getdivisi()

            this.dataSource.sortingDataAccessor = (item, property) => {
                switch (property) {
                  case 'tanggal_pembelian':
                    return new Date(item.tanggal_pembelian);
                  case 'harga_pembelian':
                    return this.parseCurrency(item.harga_pembelian);
                  case 'total_nilai':
                  case 'total_harga':
                    return this.parseCurrency(item.total_nilai || item.total_harga);
                  case 'qty_barang':
                  case 'qty':
                    return Number(item.qty_barang || item.qty);
                  default:
                    return item[property];
                }
              };

          }

          parseCurrency(value: string): number {
            // Menghapus IDR, titik, dan koma lalu konversi ke angka
            if (!value) return 0;
            return Number(value.replace(/[^\d]/g, ''));
          }


          async getdivisi() {
            try {
              const divisiItem = await this.apiService.get('/divisi');

              // Create a Set to track unique division names
              const uniqueDivisiNames = new Set<string>();

              // Clear the existing divisiItems to avoid duplicates
              this.divisiItems = [];

              // Only add division names that haven't been seen before
              divisiItem.data.forEach((data) => {
                if (data.nama_divisi && !uniqueDivisiNames.has(data.nama_divisi)) {
                  uniqueDivisiNames.add(data.nama_divisi);
                  this.divisiItems.push({ name: data.nama_divisi, checked: false });
                }
              });
            } catch (error) {
              console.error("Failed to fetch data:", error);
            }
          }


          get user() {
            return this._tableUserService.inventoryItems();
          }

          get isLoading() {
            return this._tableUserService.isLoading()
          }

          get isNotFound() {
            return this._tableUserService.isNotFound()
          }

          ngAfterViewInit() {
            this.dataSource.paginator = this.paginator;
            this.dataSource.sort = this.sort;

            this.dataSource.filterPredicate = (data: any, filter: string): boolean => {
              const filters = JSON.parse(filter);

              const search = filters.search.toLowerCase();
              const divisiFilter = filters.divisi || [];

              const searchMatch =
                !search ||
                data.nama_barang.toLowerCase().includes(search) ||
                data.nama_divisi.toLowerCase().includes(search);

              const divisiMatch =
                divisiFilter.length === 0 ||
                divisiFilter.includes(data.nama_divisi.toLowerCase());

              return searchMatch && divisiMatch;
            };

            this.dataSource.sortingDataAccessor = (item, property) => {
                switch (property) {
                  case 'harga_pembelian':
                    return item.harga_pembelian;
                  case 'total_nilai':
                    return item.total_nilai;
                  case 'qty_barang':
                    return item.qty_barang;
                  case 'tanggal_pembelian':
                    return new Date(item.tanggal_pembelian);
                  default:
                    return item[property];
                }
            };
          }


          TambahBarang() {
            this._matDialog.open(AddInventarisComponent, {
              width: window.innerWidth < 600 ? '90%' : '50%',
              maxWidth: '100vw'
            });
          }

        // deleteUser(no_hp: string) {
        //     const confirm = this.fuseConfirmationService.open({
        //       title: 'Konfirmasi Hapus',
        //       message: `Apakah Anda yakin ingin menghapus akun ${no_hp}?`,
        //       actions: {
        //         confirm: {
        //           label: 'Delete',
        //         },
        //       },
        //     });

        //     confirm.afterClosed().subscribe((result) => {
        //       if (result === 'confirmed') {
        //         this.deleteUserByPhone(no_hp);
        //       }
        //     });
        //   }

        // async deleteUserByPhone(name: string) {
        //     try {
        //       this.isLoadingDelete = true
        //       const deleteData = await this.apiService.delete(`/api/mikrotik/deleteExpiredHotspotUsersByPhone/${name}`, true)
        //       if (deleteData.message === "Hotspot user deleted successfully") {
        //         this.toast.error('Delete User Succesfully', 'Delete User')
        //         this._tableUserService.fetchData()
        //         this.isLoadingDelete = false
        //       }
        //     } catch (error) {
        //       this.toast.error('Error Delete Users', 'Delete User')
        //       this.isLoadingDelete = false
        //       throw error
        //     }
        //   }

        applyFilter() {
            this.dataSource.filterPredicate = (data: any, filter: string) => {
                let filters = JSON.parse(filter);

                // Filter berdasarkan pencarian umum
                let searchMatch =
                    filters.search === '' ||
                    data.nama_barang.toLowerCase().includes(filters.search) ||
                    data.nama_divisi.toLowerCase().includes(filters.search);

                // Filter berdasarkan divisi yang dipilih
                let divisiMatch = true;
                if (filters.divisi && filters.divisi.length > 0) {
                    divisiMatch = filters.divisi.includes(data.nama_divisi.toLowerCase());
                }

                // Kembalikan true jika data memenuhi kedua kondisi filter
                return searchMatch && divisiMatch;
            };

            // Terapkan filter berdasarkan filterValues yang telah disesuaikan
            this.dataSource.filter = JSON.stringify(this.filterValues);

            // Cek apakah data yang difilter ada atau tidak
            this.isNotDataFound = this.dataSource.filteredData.length === 0;
        }


          applySearchFilter(filterValue: string): void {
            this.filterValues.search = filterValue.trim().toLowerCase();
            this.dataSource.filter = JSON.stringify(this.filterValues);

            // Reset ke halaman pertama jika ada paginator
            if (this.dataSource.paginator) {
              this.dataSource.paginator.firstPage();
            }

            this.isNotDataFound = this.dataSource.filteredData.length === 0;
          }


          applyFilterUserProfile(): void {
            // Get selected division names and convert to lowercase for case-insensitive matching
            this.filterValues.divisi = this.divisiItems
                .filter(item => item.checked)
                .map(item => item.name.toLowerCase());

            // Apply the updated filter
            this.dataSource.filter = JSON.stringify(this.filterValues);

            // Reset to first page if paginator exists
            if (this.dataSource.paginator) {
                this.dataSource.paginator.firstPage();
            }

            // Update the no data found flag
            this.isNotDataFound = this.dataSource.filteredData.length === 0;
        }

        addNewDivision() {
            const dialogRef = this._matDialog.open(AddDivisiComponent, {
              width: window.innerWidth < 600 ? '90%' : '50%',
              maxWidth: '100vw'
            });

            dialogRef.afterClosed().subscribe(result => {
              // Mengecek apakah hasil dialog adalah 'refresh' yang menandakan divisi berhasil ditambahkan
              if (result === 'refresh') {
                // Memanggil fungsi getDivisi untuk mendapatkan data terbaru
                this.getdivisi(); // Ganti dengan fungsi getDivisi milikmu
              }
            });
          }

          // Fungsi untuk menghapus divisi
          deleteSelectedDivisi() {
            // Menyaring divisi yang dipilih untuk dihapus
            const divisiToDelete = this.divisiItems.filter(item => item.checked);

            // Jika tidak ada divisi yang dipilih
            if (divisiToDelete.length === 0) {
              console.log('No divisi selected for deletion');
              return;
            }

            // Menanyakan konfirmasi penghapusan
            const confirm = this.fuseConfirmationService.open({
              title: 'Konfirmasi Hapus',
              message: `Apakah Anda yakin ingin menghapus ${divisiToDelete.length} divisi yang dipilih?`,
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
                  // Panggil API untuk menghapus divisi yang dipilih
                  for (let item of divisiToDelete) {
                    await this.apiService.delete(`/divisi/${item.name}`);
                  }

                  // Menghapus divisi yang dipilih dari array divisiItems
                  this.divisiItems = this.divisiItems.filter(item => !item.checked);

                  // Mengambil data divisi terbaru setelah penghapusan
                  // Update divisiItems only with fresh data (avoid appending or duplication)
                  this.getdivisi(); // Ensure this method replaces divisiItems with fresh data, not appending
                  this._tableUserService.fetchData(); // Refresh the inventory data
                  console.log(`${divisiToDelete.length} divisi berhasil dihapus`);
                } catch (error) {
                  console.error('Gagal menghapus divisi', error);
                }
              }
            });
          }





}
