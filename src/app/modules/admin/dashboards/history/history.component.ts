import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FuseCardComponent } from '@fuse/components/card';
import { ApiService } from 'app/services/api.service';

@Component({
  selector: 'app-history',
  standalone: true,
  imports: [
    RouterLink,
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
  ],
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.scss']
})
export class HistoryComponent implements OnInit {
  displayedColumns: string[] = ['kategori', 'keterangan']; // kolom yang akan ditampilkan

  dataSource = new MatTableDataSource<any>([]);
  isLoading: boolean = false;
  isNotDataFound: boolean = false;

  logAktivitas: any[] = []; // Hold fetched data


  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;


  constructor(
    private _apiService: ApiService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(() => {
      this.getHistory();
    });
  }


  async getHistory(): Promise<void> {
    this.isLoading = true;
    this.isNotDataFound = false;

    try {
      const response = await this._apiService.get(`/histories`,true);
      this.logAktivitas = Array.isArray(response) ? response : [];

      this.dataSource.data = this.logAktivitas;
      this.isNotDataFound = this.logAktivitas.length === 0;
    } catch (error) {
      console.error(error);
      this.isNotDataFound = true;
    } finally {
      this.isLoading = false;
    }
  }


  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;

    // Custom filter: hanya berdasarkan kategori dan keterangan
    this.dataSource.filterPredicate = (data: any, filter: string) => {
      const transformedFilter = filter.trim().toLowerCase();
      return (
        data.kategori?.toLowerCase().includes(transformedFilter) ||
        data.keterangan?.toLowerCase().includes(transformedFilter)
      );
    };
  }


  // This method will trigger when the page is changed
  onPageChange(event: any): void {
    this.dataSource.paginator = this.paginator;  // Ensure paginator is applied when page changes
  }

  applySearchFilter(search: string): void {
    this.dataSource.filter = search.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }


  applyFilterUserProfile(): void {
    // Add logic for filtering user profiles (if needed)
  }
}
