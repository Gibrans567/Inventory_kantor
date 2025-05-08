import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BarangRusakComponent } from './barang-rusak.component';

describe('BarangRusakComponent', () => {
  let component: BarangRusakComponent;
  let fixture: ComponentFixture<BarangRusakComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BarangRusakComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(BarangRusakComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
