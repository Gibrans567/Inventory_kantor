import { CommonModule } from '@angular/common';
import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
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
  selector: 'app-update-status',
  standalone: true,
  imports: [
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    CommonModule,
    ReactiveFormsModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatIconModule,
  ],
  templateUrl: './update-status.component.html',
  styleUrl: './update-status.component.scss'
})
export class UpdateStatusComponent implements OnInit {
    barangStatusForm!: FormGroup;
    isLoading = false;

    constructor(
        private platform: Platform,
        private fb: FormBuilder,
        private _apiService: ApiService,
        private router: Router,
        private loadingController: LoadingController,
        private fuseConfirmationService: FuseConfirmationService,
        private _tableUsersService: InventoryService,
        public dialogRef: MatDialogRef<UpdateStatusComponent>,
        @Inject(MAT_DIALOG_DATA) public data: { barangId: number,qty_barang: number}
    ){}

    ngOnInit(): void {
        this.barangStatusForm = this.fb.group({
            status: ['Tersedia', Validators.required],
            qty_barang: [1, [Validators.required, Validators.min(0)]]
        });

        this.platform.backButton.subscribeWithPriority(10, () => {
            this.router.navigate(['/dashboard/inventaris']);
        });
    }

    onNoClick(): void {
        this.dialogRef.close();
    }

    async onSubmit() {
        if (!this.barangStatusForm.valid) {
            console.warn('Form tidak valid');
            return;
        }

        this.isLoading = true;

        try {
            const formData = {
                status: this.barangStatusForm.value.status,
                qty_barang: this.barangStatusForm.value.qty_barang
            };

            // Use the PUT method to update the barang status with the ID from the data object
            const dataPut = await this._apiService.put(`/barangStatus/${this.data.barangId}`, formData);
            console.log(dataPut);

            this.dialogRef.close('refresh');
            this._tableUsersService.fetchData();

        } catch (error) {
            console.error('Gagal update status barang', error);
        } finally {
            this.isLoading = false;
        }
    }

}
