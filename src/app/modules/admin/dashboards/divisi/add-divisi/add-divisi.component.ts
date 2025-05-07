import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { Router } from '@angular/router';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { Platform } from '@ionic/angular';
import { LoadingController } from '@ionic/angular/standalone';
import { InventoryService } from 'app/modules/component/table/table-inventaris/table-inventaris.service';
import { ApiService } from 'app/services/api.service';
import { MatDialogRef } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-divisi',
  standalone: true,
  imports: [
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    CommonModule,
    ReactiveFormsModule,
    MatIconModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: './add-divisi.component.html',
  styleUrl: './add-divisi.component.scss'
})
export class AddDivisiComponent {
    categoryForm!: FormGroup;
    isLoading = false;

    constructor(
        private platform: Platform,
        private fb: FormBuilder,
        private _apiService: ApiService,
        private router: Router,
        private loadingController: LoadingController,
        private fuseConfirmationService: FuseConfirmationService,
        private _tableUsersService: InventoryService,
        public dialogRef: MatDialogRef<AddDivisiComponent>,
    ){}

    ngOnInit(): void {
        this.categoryForm = this.fb.group({
            nama_divisi: ['', Validators.required]
        });

        this.platform.backButton.subscribeWithPriority(10, () => {
            this.router.navigate(['/dashboard/inventaris']);
        });
    }

    onNoClick(): void {
        this.dialogRef.close();
    }

    async onSubmit() {
        if (!this.categoryForm.valid) {
            console.warn('Form tidak valid');
            return;
        }

        this.isLoading = true;

        try {
            const formData = {
                nama_divisi: this.categoryForm.value.nama_divisi
            };

            const dataPost = await this._apiService.post('/divisi', formData);
            console.log(dataPost);

            this.dialogRef.close('refresh');
            this._tableUsersService.fetchData();

        } catch (error) {
            console.error('Gagal submit form', error);
        } finally {
            this.isLoading = false;
        }
    }
}
