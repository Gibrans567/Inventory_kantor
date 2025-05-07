import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddKategoriComponent } from './add-kategori.component';

describe('AddKategoriComponent', () => {
  let component: AddKategoriComponent;
  let fixture: ComponentFixture<AddKategoriComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddKategoriComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddKategoriComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
