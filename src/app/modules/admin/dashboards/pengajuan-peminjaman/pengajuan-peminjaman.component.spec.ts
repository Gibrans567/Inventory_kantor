import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PengajuanPeminjamanComponent } from './pengajuan-peminjaman.component';

describe('PengajuanPeminjamanComponent', () => {
  let component: PengajuanPeminjamanComponent;
  let fixture: ComponentFixture<PengajuanPeminjamanComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PengajuanPeminjamanComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PengajuanPeminjamanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
