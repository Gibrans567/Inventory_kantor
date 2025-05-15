import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { Router } from '@angular/router';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { Platform } from '@ionic/angular';
import { LoadingController } from '@ionic/angular/standalone';
import { InventoryService } from 'app/modules/component/table/table-inventaris/table-inventaris.service';
import { ApiService } from 'app/services/api.service';

@Component({
  selector: 'app-add-gudang',
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
  templateUrl: './add-gudang.component.html',
  styleUrl: './add-gudang.component.scss'
})
export class AddGudangComponent {
    gudangForm!: FormGroup;
    isLoading = false;

    constructor(
        private platform: Platform,
        private fb: FormBuilder,
        private _apiService: ApiService,
        private router: Router,
        private loadingController: LoadingController,
        private fuseConfirmationService: FuseConfirmationService,
        private _tableUsersService: InventoryService,
        public dialogRef: MatDialogRef<AddGudangComponent>,
    ){}

    ngOnInit(): void {
        this.gudangForm = this.fb.group({
            nama_gudang: ['', Validators.required],
            lokasi_gudang: ['', Validators.required]
        });

        this.platform.backButton.subscribeWithPriority(10, () => {
            this.router.navigate(['/dashboard/inventaris']);
        });
    }

    onNoClick(): void {
        this.dialogRef.close();
    }

    async onSubmit() {
        if (!this.gudangForm.valid) {
            console.warn('Form tidak valid');
            return;
        }

        this.isLoading = true;

        try {
            const formData = {
                nama_gudang: this.gudangForm.value.nama_gudang,
                lokasi_gudang: this.gudangForm.value.lokasi_gudang
            };

            const dataPost = await this._apiService.post('/gudang', formData);
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
