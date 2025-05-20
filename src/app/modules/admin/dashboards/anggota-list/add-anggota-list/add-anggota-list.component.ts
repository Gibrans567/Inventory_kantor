import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { Router } from '@angular/router';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { Platform } from '@ionic/angular';
import { LoadingController } from '@ionic/angular/standalone';
import { InventoryService } from 'app/modules/component/table/table-inventaris/table-inventaris.service';
import { ApiService } from 'app/services/api.service';

@Component({
  selector: 'app-add-anggota-list',
  standalone: true,
  imports: [
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    CommonModule,
    ReactiveFormsModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSelectModule
  ],
  templateUrl: './add-anggota-list.component.html',
  styleUrl: './add-anggota-list.component.scss'
})
export class AddAnggotaListComponent {


    userForm!: FormGroup;
    isLoading = false;
    divisiList: any[] = []; // Tambahkan ini


    constructor(
        private platform: Platform,
        private fb: FormBuilder,
        private _apiService: ApiService,
        private router: Router,
        private loadingController: LoadingController,
        private fuseConfirmationService: FuseConfirmationService,
        private _tableUsersService: InventoryService,
        public dialogRef: MatDialogRef<AddAnggotaListComponent>,
    ){}

    ngOnInit(): void {
        this.userForm = this.fb.group({
        id_divisi: ['', Validators.required],
        email: ['', [Validators.required, Validators.email]],
        password: ['', Validators.required],
        nama_user: ['', Validators.required],
        role: ['', Validators.required] // Ubah dari role_id ke role
    });

        this.getDivisiList(); // Panggil fungsi untuk mendapatkan daftar divisi
    }

     async getDivisiList() {
        try {
            const res = await this._apiService.get('/divisi', true);
            this.divisiList = res.data || [];
        } catch (error) {
            console.error('Failed to fetch divisi data', error);
            this.showErrorDialog('Gagal mengambil data divisi');
        }
    }


    onNoClick(): void {
        this.dialogRef.close();
    }

    async onSubmit() {
        if (!this.userForm.valid) {
            console.warn('Form tidak valid');
            return;
        }

        this.isLoading = true;

        try {
            const formData = {
            id_divisi: this.userForm.value.id_divisi,
            email: this.userForm.value.email,
            password: this.userForm.value.password,
            nama_user: this.userForm.value.nama_user, // Perbaiki dari nama_user ke nama_user
            role: this.userForm.value.role
        };


            const dataPost = await this._apiService.post('/user', formData);
            console.log(dataPost);

            this.dialogRef.close('refresh');
            this._tableUsersService.fetchData();

        } catch (error) {
            console.error('Gagal submit form', error);
        } finally {
            this.isLoading = false;
        }
    }

    showErrorDialog(message: string): void {
        this.fuseConfirmationService.open({
            title: 'Error',
            message,
            actions: {
                confirm: {
                    label: 'OK',
                    color: 'warn'
                },
                cancel: {
                    show: false
                }
            }
        });
    }

}
