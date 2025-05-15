import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddPengajuanPeminjamanComponent } from './add-pengajuan-peminjaman.component';

describe('AddPengajuanPeminjamanComponent', () => {
  let component: AddPengajuanPeminjamanComponent;
  let fixture: ComponentFixture<AddPengajuanPeminjamanComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddPengajuanPeminjamanComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddPengajuanPeminjamanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
