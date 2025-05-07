import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddSebaranBarangComponent } from './add-sebaran-barang.component';

describe('AddSebaranBarangComponent', () => {
  let component: AddSebaranBarangComponent;
  let fixture: ComponentFixture<AddSebaranBarangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddSebaranBarangComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddSebaranBarangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
