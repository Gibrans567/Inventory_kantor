import { CommonModule, Location } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { IonIcon, IonLoading, Platform } from '@ionic/angular/standalone';
import { ApiService } from 'app/services/api.service';
import { ToastrService } from 'ngx-toastr';
import { InventoryService } from 'app/modules/component/table/table-inventaris/table-inventaris.service';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { MatDatepicker, MatDatepickerInput, MatDatepickerToggle } from '@angular/material/datepicker';
import { MatOption } from '@angular/material/core';
import { DateTime } from 'luxon';
import { LoadingController } from '@ionic/angular';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';


@Component({
  selector: 'app-add-inventaris',
  standalone: true,
  imports: [
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    ReactiveFormsModule,
    CommonModule,
    MatIconModule,
    MatDatepickerInput,
    MatDatepickerToggle,
    MatDatepicker,
    MatProgressSpinnerModule

  ],
  templateUrl: './add-inventaris.component.html',
  styleUrl: './add-inventaris.component.scss'
})

export class AddInventarisComponent {
    inventarisForm!: FormGroup
    inventarisDatas: any[] = []
    isLoading = false;
    gudangList: any[] = []
    kategoriList: any[] = []
    divisiList: any[] = []
    userList: any[] = []


    constructor(
        private platform: Platform,
        private fb: FormBuilder,
        private _apiService: ApiService,
        private router: Router,
        private loadingController: LoadingController,
        private fuseConfirmationService: FuseConfirmationService,
        private _tableUsersService: InventoryService,
        public dialogRef : MatDialogRef<AddInventarisComponent>,
        @Inject(MAT_DIALOG_DATA) public data : any,
    ){}

    ngOnInit(): void {
        this.inventarisForm = this.fb.group({
            nama_barang: ['', Validators.required],
            qty_barang: [null, [Validators.required, Validators.min(1)]],
            harga_pembelian: [null, [Validators.required, Validators.min(1)]],
            spesifikasi: ['', Validators.required],
            tanggal_pembelian: ['', Validators.required],
            gudang_id: [null, Validators.required],
            kategori_id: [null, Validators.required],
            divisi_id: [null, Validators.required],
            user_id: [null, Validators.required]
        });

        this.platform.backButton.subscribeWithPriority(10, () => {
            this.router.navigate(['/dashboard/inventaris']);
        });

        this.getGudangList()
        this.getKategoriList()
        this.getDivisiList()
        this.getUserList()
    }


    onNoClick(): void {
        this.dialogRef.close();
    }

    async getRole() {
        const roleData = await this._apiService.get('/api/mikrotik/get-profiles')
        this.inventarisDatas = roleData.profiles
    }


    async onSubmit() {
        if (!this.inventarisForm.valid) {
            console.warn('Form tidak valid');
            return;
        }

        this.isLoading = true;

        try {
            const formData = {
                nama_barang: this.inventarisForm.value.nama_barang,
                qty_barang: this.inventarisForm.value.qty_barang,
                harga_pembelian: this.inventarisForm.value.harga_pembelian,
                spesifikasi: this.inventarisForm.value.spesifikasi,
                tanggal_pembelian: this.convertToISOString(this.inventarisForm.value.tanggal_pembelian),
                gudang_id: Number(this.inventarisForm.value.gudang_id),
                kategori_id: Number(this.inventarisForm.value.kategori_id),
                divisi_id: Number(this.inventarisForm.value.divisi_id),
                user_id: Number(this.inventarisForm.value.user_id)
            };

            const dataPost = await this._apiService.post('/inventaris', formData);
            console.log(dataPost);

            this.dialogRef.close('refresh');
            this._tableUsersService.fetchData();
        } catch (error) {
            console.error('Gagal submit form', error);
        } finally {
            this.isLoading = false;
        }
    }

    convertToISOString(date: any): string {
        if (date) {
          const formattedDate = new Date(date);
          return DateTime.fromJSDate(new Date(date)).toISODate(); // Mengonversi ke format ISO string
        }
        return '';
      }

    preventInvalidChars(event: any) {
        const inputValue = event.target.value;
        // Hanya izinkan angka (0-9)
        event.target.value = inputValue.replace(/[^0-9]/g, '');
    }

    async getGudangList() {
        try {
            const res = await this._apiService.get('/gudang');
            this.gudangList = res; // sesuaikan struktur response kalau beda
        } catch (error) {
            console.error('Gagal ambil data gudang', error);
        }
    }

    async getKategoriList() {
        try {
            const res = await this._apiService.get('/kategori');
            this.kategoriList = res.data;
        } catch (error) {
            console.error('Gagal ambil data kategori', error);
        }
    }

    async getDivisiList() {
        try {
            const res = await this._apiService.get('/divisi');
            this.divisiList = res.data;
        } catch (error) {
            console.error('Gagal ambil data divisi', error);
        }
    }

    async getUserList() {
        try {
            const res = await this._apiService.get('/user');
            this.userList = res;
        } catch (error) {
            console.error('Gagal ambil data user', error);
        }
    }

    formatCurrency(event: any): void {
        let value = event.target.value;

        // Remove all non-digit characters
        value = value.replace(/\D/g, '');

        // Convert to number and format with thousand separators
        if (value) {
            // Convert to number
            const numValue = Number(value);

            // Format with thousand separators (using locale ID for Indonesian formatting)
            const formattedValue = numValue.toLocaleString('id-ID');

            // Update the form control with the raw numeric value (without formatting)
            this.inventarisForm.get('harga_pembelian').setValue(numValue, { emitEvent: false });

            // Update the displayed value with formatting
            event.target.value = formattedValue;
        } else {
            this.inventarisForm.get('harga_pembelian').setValue(null, { emitEvent: false });
        }
    }

}
