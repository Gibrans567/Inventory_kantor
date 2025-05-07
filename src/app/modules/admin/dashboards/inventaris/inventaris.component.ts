import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { MatButtonModule } from '@angular/material/button';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TableInventarisComponent } from 'app/modules/component/table/table-inventaris/table-inventaris.component';
import { RouterLink } from '@angular/router';



@Component({
  selector: 'app-inventaris',
  standalone: true,
  imports: [IonicModule,MatButtonModule,MatTableModule,MatSortModule,MatFormFieldModule,
            MatInputModule,MatPaginatorModule,ReactiveFormsModule,CommonModule,TableInventarisComponent,RouterLink],
  templateUrl: './inventaris.component.html',
  styleUrl: './inventaris.component.scss'
})
export class InventarisComponent {

}
