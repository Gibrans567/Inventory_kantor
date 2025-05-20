import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FuseCardComponent } from '@fuse/components/card';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { IonicModule } from '@ionic/angular';
import { ApiService } from 'app/services/api.service';
import { AddGudangComponent } from './add-gudang/add-gudang.component';

@Component({
  selector: 'app-gudang',
  standalone: true,
  imports: [
    RouterLink,
    MatTableModule,
    MatPaginatorModule,
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
  templateUrl: './gudang.component.html',
  styleUrl: './gudang.component.scss'
})
export class GudangComponent {
  displayedColumns: string[] = ['nama_gudang', 'lokasi_gudang','action'];

  dataSource = new MatTableDataSource<any>([]);
  isLoading: boolean = false;
  isNotDataFound: boolean = false;

  gudangList: any[] = [];

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
      this.getGudang();
    });
  }

  async getGudang(): Promise<void> {
    this.isLoading = true;
    this.isNotDataFound = false;

    try {
      const response = await this.apiService.get(`/gudang`,true);

      const gudangData = response ?? [];

      this.gudangList = gudangData.map((item: any) => ({
        id: item.id,
        nama_gudang: item.nama_gudang,
        lokasi_gudang: item.lokasi_gudang,
        createdAt: item.CreatedAt,
        updatedAt: item.UpdatedAt,
      }));

      this.dataSource.data = this.gudangList;
      this.isNotDataFound = this.gudangList.length === 0;
    } catch (error) {
      console.error('Gagal mengambil data gudang:', error);
      this.isNotDataFound = true;
    } finally {
      this.isLoading = false;
    }
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;

    this.dataSource.filterPredicate = (data: any, filter: string) => {
      const transformedFilter = filter.trim().toLowerCase();
      return (
        data.nama_gudang?.toLowerCase().includes(transformedFilter) ||
        data.lokasi_gudang?.toLowerCase().includes(transformedFilter) ||
        data.createdAt?.toLowerCase().includes(transformedFilter)
      );
    };
  }

  onPageChange(event: any): void {
    this.dataSource.paginator = this.paginator;
  }

  applySearchFilter(search: string): void {
    this.dataSource.filter = search.trim().toLowerCase();
    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  addNewGudang() {
      const dialogRef = this.matDialog.open(AddGudangComponent, {
        width: window.innerWidth < 600 ? '90%' : '50%',
        maxWidth: '100vw',
      });

      dialogRef.afterClosed().subscribe((result) => {
        // Mengecek apakah hasil dialog adalah 'refresh' yang menandakan divisi berhasil ditambahkan
        if (result === 'refresh') {
          this.getGudang(); // Memanggil fungsi getDivisi untuk mendapatkan data terbaru
        }
      });
    }

  deleteGudangById(id: number): void {
    const confirm = this.fuseConfirmationService.open({
      title: 'Konfirmasi Hapus',
      message: `Apakah Anda yakin ingin menghapus gudang dengan ID ${id}?`,
      actions: {
        confirm: { label: 'Hapus' },
        cancel: { label: 'Batal' }
      }
    });

    confirm.afterClosed().subscribe(async (result) => {
      if (result === 'confirmed') {
        try {
          await this.apiService.delete(`/gudang/${id}`,true);
          this.getGudang();
          console.log(`Gudang dengan ID ${id} berhasil dihapus`);
        } catch (error) {
          console.error('Gagal menghapus gudang:', error);
        }
      }
    });
  }
}
