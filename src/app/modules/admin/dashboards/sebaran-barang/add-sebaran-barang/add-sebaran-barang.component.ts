import { CommonModule } from '@angular/common';
import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { Platform } from '@ionic/angular/standalone';
import { ApiService } from 'app/services/api.service';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-add-sebaran-barang',
  standalone: true,
  imports: [
    MatFormFieldModule,
    FormsModule,
    MatInputModule,
    ReactiveFormsModule,
    CommonModule,
    MatIconModule,
    MatSelectModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './add-sebaran-barang.component.html',
  styleUrl: './add-sebaran-barang.component.scss'
})
export class AddSebaranBarangComponent implements OnInit {

    SebaranBarangForm!: FormGroup;
    isLoading = false;
    divisiList: any[] = [];
    userList: any[] = [];

    // Data barang yang diambil dari dialog parent
    barangDetail: any;

    constructor(
      private fb: FormBuilder,
      private _apiService: ApiService,
      private platform: Platform,
      private _fuseConfirmationService: FuseConfirmationService,
      public dialogRef: MatDialogRef<AddSebaranBarangComponent>,
      @Inject(MAT_DIALOG_DATA) public data: any
    ) {}

    ngOnInit(): void {
      // Pastikan data yang diterima dari parent component
      console.log('Data received:', this.data);

      // Simpan data barang dari parent component
      this.barangDetail = this.data;

      // Initialize form controls - hanya field yang diperlukan
      this.SebaranBarangForm = this.fb.group({
        qty_barang: [null, [Validators.required, Validators.min(1), Validators.max(this.barangDetail?.qty_tersedia || 999999)]],
        posisi_awal: [this.barangDetail?.nama_gudang || ''],
        posisi_akhir: ['', Validators.required],
        divisi_id: [this.barangDetail?.divisi_id || null, Validators.required],
        user_id: [null, Validators.required]
      });

      // Fetch lists for dropdowns
      this.getDivisiList();
      this.getUserList();
    }

    onNoClick(): void {
      this.dialogRef.close();
    }

    // Get list of Divisi
    async getDivisiList() {
      try {
        const res = await this._apiService.get('/divisi');
        this.divisiList = res.data || []; // Assuming response has a 'data' field
      } catch (error) {
        console.error('Failed to fetch divisi data', error);
        this.showErrorDialog('Gagal mengambil data divisi');
      }
    }

    // Get list of Users
    async getUserList() {
      try {
        const res = await this._apiService.get('/user');
        this.userList = res || [];
      } catch (error) {
        console.error('Failed to fetch user data', error);
        this.showErrorDialog('Gagal mengambil data pengguna');
      }
    }

    // Handle form submission
    async onSubmit() {
        if (!this.SebaranBarangForm.valid) {
          this.showWarningDialog('Mohon isi semua field yang wajib diisi');
          return;
        }

        // Validasi jumlah barang tidak melebihi yang tersedia
        const qtyBarang = this.SebaranBarangForm.value.qty_barang;
        const qtyTersedia = this.barangDetail?.qty_tersedia;

        if (qtyBarang > qtyTersedia) {
          this.showWarningDialog(`Jumlah barang tidak boleh melebihi stok tersedia (${qtyTersedia})`);
          return;
        }

        // Get user_id as number
        const userId = Number(this.SebaranBarangForm.value.user_id);  // Konversi ke angka

        // Show confirmation dialog before proceeding
        const confirmation = this._fuseConfirmationService.open({
          title: 'Tambah Sebaran Barang',
          message: `Anda yakin ingin menambahkan ${qtyBarang} ${this.barangDetail.nama_barang} ke distribusi?`,
          icon: {
            show: true,
            name: 'heroicons_outline:question-mark-circle',
            color: 'info'
          },
          actions: {
            confirm: {
              label: 'Ya, Tambahkan',
              color: 'primary'
            },
            cancel: {
              label: 'Batal'
            }
          }
        });

        // Handle the result of the confirmation dialog
        confirmation.afterClosed().subscribe(async (result) => {
          if (result === 'confirmed') {
            this.isLoading = true;

            try {
              // Prepare data for the POST request
              const formData = {
                id_barang: this.barangDetail.id,
                id_divisi: this.SebaranBarangForm.value.divisi_id,
                id_user: userId,  // Kirimkan id_user sebagai number
                qty_barang: qtyBarang,
                posisi_awal: this.SebaranBarangForm.value.posisi_awal,
                posisi_akhir: this.SebaranBarangForm.value.posisi_akhir
              };

              // Make POST request to add new sebaran barang
              const response = await this._apiService.post('/sebaranBarang', formData);
              console.log('Response:', response);

              // Show success dialog
              this.showSuccessDialog('Barang berhasil ditambahkan ke distribusi!');

              // Close the dialog and refresh table data
              this.dialogRef.close('refresh');
            } catch (error) {
              console.error('Failed to add barang', error);
              this.showErrorDialog('Gagal menambahkan barang ke distribusi');
            } finally {
              this.isLoading = false;
            }
          }
        });
      }


    // Helper method to prevent invalid characters
    preventInvalidChars(event: any) {
        const inputValue = event.target.value;
        const maxQty = this.barangDetail?.qty_tersedia;

        // Batasi input agar tidak lebih dari qty_tersedia
        if (inputValue > maxQty) {
            event.target.value = maxQty;
        }

        // Hanya izinkan angka (0-9)
        event.target.value = inputValue.replace(/[^0-9]/g, '');
    }


    // Display nama instead of ID for selection
    getDivisiDisplayFn() {
      return (divisiId: number) => {
        const divisi = this.divisiList.find(d => d.id === divisiId);
        return divisi ? divisi.nama_divisi : '';
      };
    }

    getUserDisplayFn() {
      return (userId: number) => {
        const user = this.userList.find(u => u.id === userId);
        return user ? user.nama : '';
      };
    }

    // Get nama from ID
    getDivisiName(divisiId: number): string {
      const divisi = this.divisiList.find(d => d.id === divisiId);
      return divisi ? divisi.nama_divisi : 'Tidak ditemukan';
    }

    getUserName(userId: number): string {
      const user = this.userList.find(u => u.id === userId);
      return user ? user.nama : 'Tidak ditemukan';
    }

    // Show success dialog using FuseConfirmationService
    showSuccessDialog(message: string): void {
      this._fuseConfirmationService.open({
        title: 'Berhasil',
        message,
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

    // Show warning dialog using FuseConfirmationService
    showWarningDialog(message: string): void {
      this._fuseConfirmationService.open({
        title: 'Peringatan',
        message,
        icon: {
          show: true,
          name: 'heroicons_outline:exclamation',
          color: 'warning'
        },
        actions: {
          confirm: {
            label: 'OK',
            color: 'warn'
          },
          cancel: {
            show: false
          }
        },
        dismissible: true
      });
    }

    // Show error dialog using FuseConfirmationService
    showErrorDialog(message: string): void {
      this._fuseConfirmationService.open({
        title: 'Error',
        message,
        icon: {
          show: true,
          name: 'heroicons_outline:x-circle',
          color: 'error'
        },
        actions: {
          confirm: {
            label: 'OK',
            color: 'warn'
          },
          cancel: {
            show: false
          }
        },
        dismissible: true
      });
    }
}
