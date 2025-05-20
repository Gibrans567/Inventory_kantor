import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { Platform } from '@ionic/angular/standalone';
import { ApiService } from 'app/services/api.service';
import { MatDialogRef } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-add-pengajuan',
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
  templateUrl: './add-pengajuan-peminjaman.component.html',
  styleUrl: './add-pengajuan-peminjaman.component.scss'
})
export class AddPengajuanPeminjamanComponent implements OnInit {

  pengajuanForm!: FormGroup;
  isLoading = false;
  divisiList: any[] = [];
  userList: any[] = [];

  // Data dari API
  barangDetail: any = {};
  inventarisList: any[] = [];
  selectedBarangId: number | null = null;

  constructor(
    private fb: FormBuilder,
    private _apiService: ApiService,
    private platform: Platform,
    private _fuseConfirmationService: FuseConfirmationService,
    public dialogRef: MatDialogRef<AddPengajuanPeminjamanComponent>
  ) {}

  ngOnInit(): void {
    // Initialize form controls sesuai dengan format API pengajuan yang ditampilkan di Postman
    this.pengajuanForm = this.fb.group({
      id_barang: [null, Validators.required],
      id_user: [null, Validators.required],
      qty_barang: [null, [Validators.required, Validators.min(1)]],
      note: ['', Validators.required]
    });

    this.getUserList();
    this.getInventarisList();

    // Tambahkan listener untuk perubahan id_barang
    this.pengajuanForm.get('id_barang')?.valueChanges.subscribe(value => {
      console.log('Selected barang ID changed:', value);
      if (value) {
        this.updateBarangDetail(value);
      } else {
        this.barangDetail = {};
      }
    });
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

  // Get list of Users
  async getUserList() {
    try {
      const res = await this._apiService.get('/user',true);
      this.userList = res || [];
    } catch (error) {
      console.error('Failed to fetch user data', error);
      this.showErrorDialog('Gagal mengambil data pengguna');
    }
  }

  // Get list of Inventaris
  async getInventarisList() {
    try {
      this.isLoading = true;
      const res = await this._apiService.get('/inventaris',true);
      console.log('Raw inventaris response:', res);
      this.inventarisList = res.data || [];
      console.log('Processed inventarisList:', this.inventarisList);
      this.isLoading = false;
    } catch (error) {
      console.error('Failed to fetch inventaris data', error);
      this.showErrorDialog('Gagal mengambil data inventaris');
      this.isLoading = false;
    }
  }

  // Update barang detail when a barang is selected
  updateBarangDetail(barangId: any) {
    // Convert string to number if it's a string from event
    const id = typeof barangId === 'string' ? parseInt(barangId, 10) : barangId;

    console.log('Updating barang detail for ID:', id);
    console.log('Available inventory items:', this.inventarisList);

    const selectedBarang = this.inventarisList.find(item => item.id === id);
    console.log('Selected barang:', selectedBarang);

    if (selectedBarang) {
      this.barangDetail = selectedBarang;
      console.log('Barang detail updated:', this.barangDetail);

      // Update max qty validator
      const qtyControl = this.pengajuanForm.get('qty_barang');
      if (qtyControl) {
        qtyControl.setValidators([
          Validators.required,
          Validators.min(1),
          Validators.max(this.barangDetail.qty_tersedia || 0)
        ]);
        qtyControl.updateValueAndValidity();
      }
    } else {
      this.barangDetail = {};
      console.log('No barang found with ID:', id);
    }
  }

  // Handle form submission
  async onSubmit() {
    if (!this.pengajuanForm.valid) {
      this.showWarningDialog('Mohon isi semua field yang wajib diisi');
      return;
    }

    // Validasi jumlah barang tidak melebihi yang tersedia
    const qtyBarang = this.pengajuanForm.value.qty_barang;
    const qtyTersedia = this.barangDetail.qty_tersedia;

    if (qtyBarang > qtyTersedia) {
      this.showWarningDialog(`Jumlah barang tidak boleh melebihi stok tersedia (${qtyTersedia})`);
      return;
    }

    // Prepare data for API request - ensuring all values are in the correct format
    const requestData = {
      id_barang: Number(this.pengajuanForm.value.id_barang),
      id_user: Number(this.pengajuanForm.value.id_user),
      qty_barang: Number(this.pengajuanForm.value.qty_barang),
      note: this.pengajuanForm.value.note
    };

    // Show confirmation dialog before proceeding
    const confirmation = this._fuseConfirmationService.open({
      title: 'Tambah Pengajuan',
      message: `Anda yakin ingin mengajukan ${qtyBarang} ${this.barangDetail.nama_barang}?`,
      icon: {
        show: true,
        name: 'heroicons_outline:question-mark-circle',
        color: 'info'
      },
      actions: {
        confirm: {
          label: 'Ya, Ajukan',
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
          // POST request ke endpoint pengajuan dengan format yang sesuai
          const response = await this._apiService.post('/pengajuan', requestData);
          console.log('Response:', response);

          // Show success dialog
          this.showSuccessDialog('Pengajuan berhasil dibuat!');

          // Close the dialog and refresh table data
          this.dialogRef.close('refresh');
        } catch (error) {
          console.error('Failed to create pengajuan', error);
          this.showErrorDialog('Gagal membuat pengajuan');
        } finally {
          this.isLoading = false;
        }
      }
    });
  }

  // Helper method to prevent invalid characters
  preventInvalidChars(event: any) {
    const inputValue = event.target.value;
    const maxQty = this.barangDetail.qty_tersedia;

    // Batasi input agar tidak lebih dari qty_tersedia
    if (inputValue > maxQty) {
      event.target.value = maxQty;
    }

    // Hanya izinkan angka (0-9)
    event.target.value = inputValue.replace(/[^0-9]/g, '');
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
