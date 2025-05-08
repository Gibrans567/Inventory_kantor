import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddBarangRusakComponent } from './add-barang-rusak.component';

describe('AddBarangRusakComponent', () => {
  let component: AddBarangRusakComponent;
  let fixture: ComponentFixture<AddBarangRusakComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddBarangRusakComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddBarangRusakComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
