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
import { AddAnggotaListComponent } from './add-anggota-list/add-anggota-list.component';

@Component({
  selector: 'app-anggota-list',
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
    MatMenuModule,
    IonicModule,
    MatDialogModule,
    FormsModule,
    MatButtonModule,
  ],
  templateUrl: './anggota-list.component.html',
  styleUrl: './anggota-list.component.scss'
})
export class AnggotaListComponent {

 displayedColumns: string[] = ['nama_user', 'email', 'role', 'createdAt', 'action'];
  dataSource = new MatTableDataSource<any>([]);
  isLoading = false;
  isNotDataFound = false;

  userList: any[] = [];

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
      this.getUser();
    });
  }

  async getUser(): Promise<void> {
    this.isLoading = true;
    this.isNotDataFound = false;

    try {
      const response = await this.apiService.get(`/user`, true);
      const userData = response ?? [];

      this.userList = userData.map((item: any) => ({
        id: item.id,
        email: item.email,
        nama_user: item.nama_user,
        role: item.role,
        createdAt: item.CreatedAt,
        updatedAt: item.UpdatedAt
      }));

      this.dataSource.data = this.userList;
      this.isNotDataFound = this.userList.length === 0;
    } catch (error) {
      console.error('Gagal mengambil data user:', error);
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
        data.nama_user?.toLowerCase().includes(transformedFilter) ||
        data.email?.toLowerCase().includes(transformedFilter) ||
        data.role?.toLowerCase().includes(transformedFilter)
      );
    };
  }

  applySearchFilter(search: string): void {
    this.dataSource.filter = search.trim().toLowerCase();
    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  deleteUserById(id: number): void {
    const confirm = this.fuseConfirmationService.open({
      title: 'Konfirmasi Hapus',
      message: `Apakah Anda yakin ingin menghapus user dengan ID ${id}?`,
      actions: {
        confirm: { label: 'Hapus' },
        cancel: { label: 'Batal' }
      }
    });

    confirm.afterClosed().subscribe(async (result) => {
      if (result === 'confirmed') {
        try {
          await this.apiService.delete(`/user/${id}`, true);
          this.getUser();
          console.log(`User dengan ID ${id} berhasil dihapus`);
        } catch (error) {
          console.error('Gagal menghapus user:', error);
        }
      }
    });
  }

  addNewUser() {
        const dialogRef = this.matDialog.open(AddAnggotaListComponent, {
          width: window.innerWidth < 600 ? '90%' : '50%',
          maxWidth: '100vw',
        });

        dialogRef.afterClosed().subscribe((result) => {
          // Mengecek apakah hasil dialog adalah 'refresh' yang menandakan divisi berhasil ditambahkan
          if (result === 'refresh') {
            this.getUser(); // Memanggil fungsi getDivisi untuk mendapatkan data terbaru
          }
        });
      }


      onPageChange(event: any): void {
    this.dataSource.paginator = this.paginator;
  }
}



