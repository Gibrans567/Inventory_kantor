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
import { AddDivisiComponent } from 'app/modules/admin/dashboards/divisi/add-divisi/add-divisi.component';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-divisi',
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
  templateUrl: './divisi.component.html',
  styleUrls: ['./divisi.component.scss'],
})
export class DivisiComponent {
  displayedColumns: string[] = ['nama_divisi','action']; // Display relevant columns

  dataSource = new MatTableDataSource<any>([]);
  isLoading: boolean = false;
  isNotDataFound: boolean = false;

  logAktivitas: any[] = []; // Hold fetched data
  divisiItems: { name: string, checked: boolean }[] = []


  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;



  constructor(private apiService: ApiService, private route: ActivatedRoute,private matDialog: MatDialog,
    private fuseConfirmationService: FuseConfirmationService,
  ) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(() => {
      this.getDivisi();
    });
  }

  async getDivisi(): Promise<void> {
    this.isLoading = true;
    this.isNotDataFound = false;

    try {
      const response = await this.apiService.get(`/divisi`,true);

      // Memastikan response berhasil dan mengakses data divisi
      const divisiData = response?.data ?? [];  // Safe access to 'data' property

      // Update logAktivitas dan dataSource
      this.logAktivitas = divisiData.map((item) => ({
        id: item.id,
        nama_divisi: item.nama_divisi,
        createdAt: item.CreatedAt,
        updatedAt: item.UpdatedAt,
      }));

      // Update dataSource dengan data terbaru
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

    // Custom filter: based on 'nama_divisi' and 'createdAt' (and others as needed)
    this.dataSource.filterPredicate = (data: any, filter: string) => {
      const transformedFilter = filter.trim().toLowerCase();
      return (
        data.nama_divisi?.toLowerCase().includes(transformedFilter) ||
        data.createdAt?.toLowerCase().includes(transformedFilter) ||
        data.updatedAt?.toLowerCase().includes(transformedFilter)
      );
    };
  }


  // Trigger when the page is changed
  onPageChange(event: any): void {
    this.dataSource.paginator = this.paginator; // Ensure paginator is applied when page changes
  }

  // Apply a search filter on the table
  applySearchFilter(search: string): void {
    this.dataSource.filter = search.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  // Add logic for filtering user profiles if needed
  applyFilterUserProfile(): void {
    // Logic for additional user profile filtering
  }

  addNewDivision() {
    const dialogRef = this.matDialog.open(AddDivisiComponent, {
      width: window.innerWidth < 600 ? '90%' : '50%',
      maxWidth: '100vw',
    });

    dialogRef.afterClosed().subscribe((result) => {
      // Mengecek apakah hasil dialog adalah 'refresh' yang menandakan divisi berhasil ditambahkan
      if (result === 'refresh') {
        this.getDivisi(); // Memanggil fungsi getDivisi untuk mendapatkan data terbaru
      }
    });
  }

  deleteDivisiByName(namaDivisi: string) {
    // Menanyakan konfirmasi penghapusan
    const confirm = this.fuseConfirmationService.open({
      title: 'Konfirmasi Hapus',
      message: `Apakah Anda yakin ingin menghapus divisi "${namaDivisi}"?`,
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
          // Panggil API untuk menghapus divisi berdasarkan nama
          await this.apiService.delete(`/divisi/${namaDivisi}`,true);

          // Mengambil data divisi terbaru setelah penghapusan
          this.getDivisi(); // Memanggil fungsi untuk memperbarui data divisi dengan data terbaru

          // Memberi feedback ke user bahwa divisi berhasil dihapus
          console.log(`Divisi "${namaDivisi}" berhasil dihapus`);
        } catch (error) {
          console.error('Gagal menghapus divisi', error);
        }
      }
    });
  }


}
