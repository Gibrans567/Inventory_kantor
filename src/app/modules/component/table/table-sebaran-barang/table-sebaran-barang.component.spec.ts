import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableSebaranBarangComponent } from './table-sebaran-barang.component';

describe('TableSebaranBarangComponent', () => {
  let component: TableSebaranBarangComponent;
  let fixture: ComponentFixture<TableSebaranBarangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TableSebaranBarangComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TableSebaranBarangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
