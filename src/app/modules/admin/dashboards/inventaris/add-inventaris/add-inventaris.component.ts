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
    selectedFile: File | null = null;
    imageError: string | null = null;


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
            user_id: [null, Validators.required],
            upload_nota: [null ,Validators.required] // Tambahkan field untuk upload_nota
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

            // Kirim data ke API untuk membuat inventaris
            const dataPost = await this._apiService.post('/inventaris', formData);
            console.log(dataPost);

            // Ambil ID inventaris yang baru dibuat dari respons
            const inventarisId = dataPost.inventaris.id;  // Pastikan `id` ada di sini

            if (inventarisId) {
                // Jika ada gambar yang dipilih, lakukan upload
                if (this.selectedFile) {
                    await this.uploadImage(inventarisId);  // Kirim ID yang benar untuk upload
                }
            }

            this.dialogRef.close('refresh');
            this._tableUsersService.fetchData();
        } catch (error) {
            console.error('Gagal submit form', error);
        } finally {
            this.isLoading = false;
        }
    }


    // Fungsi untuk menangani pemilihan file
    onFileSelected(event: any) {
        const file = event.target.files[0];
        if (file) {
            // Cek apakah file yang dipilih adalah gambar
            if (file.type.startsWith('image/')) {
                this.selectedFile = file;
                this.imageError = null;  // Reset error jika file valid

                // Menyimpan file yang dipilih ke dalam form control 'upload_nota'
                this.inventarisForm.get('upload_nota')?.setValue(file);
            } else {
                this.selectedFile = null;
                this.imageError = 'Hanya file gambar yang diperbolehkan';
            }
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

        // Menghapus semua karakter non-digit, termasuk koma dan titik
        value = value.replace(/\D/g, '');  // Menghapus semua karakter non-digit

        // Jika ada value yang dimasukkan, konversi menjadi angka
        if (value) {
            const numValue = Number(value);

            // Menampilkan angka tanpa format ribuan di form
            this.inventarisForm.get('harga_pembelian').setValue(numValue, { emitEvent: false });

            // Update tampilan dengan format ribuan (untuk menampilkan di input)
            const formattedValue = numValue.toLocaleString('id-ID');
            event.target.value = formattedValue;
        } else {
            this.inventarisForm.get('harga_pembelian').setValue(null, { emitEvent: false });
        }
    }



    async uploadImage(inventarisId: number) {
        if (!this.selectedFile) {
            console.error('No file selected');
            return;
        }

        const formData = new FormData();
        formData.append('upload_nota', this.selectedFile, this.selectedFile.name);

        try {
            // Upload gambar ke API
            const response = await this._apiService.post(`/upload?id=${inventarisId}`, formData);

            // Log response to check the data
            console.log('Upload response:', response);

            if (response && response.file) {
                // Update the form with the URL returned from the server
                this.inventarisForm.get('upload_nota')?.setValue(response.file);

                // Show success message
                this.showSuccessDialog('Gambar berhasil diupload');
            }
        } catch (error) {
            console.error('Gagal mengupload gambar:', error);
            this.imageError = 'Gagal mengupload gambar';
        }
    }


    showSuccessDialog(message: string) {
        this.fuseConfirmationService.open({
            title: 'Berhasil',
            message: message,
            icon: {
                show: true,
                name: 'heroicons_outline:check-circle',
                color: 'success'
            },
            actions: {
                confirm: {
                    label: 'OK',
                    color: 'primary'
                },
                cancel: {
                    show: false
                }
            },
            dismissible: true
        });
    }



}
