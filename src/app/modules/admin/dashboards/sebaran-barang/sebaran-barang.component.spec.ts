import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SebaranBarangComponent } from './sebaran-barang.component';

describe('SebaranBarangComponent', () => {
  let component: SebaranBarangComponent;
  let fixture: ComponentFixture<SebaranBarangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SebaranBarangComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SebaranBarangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
