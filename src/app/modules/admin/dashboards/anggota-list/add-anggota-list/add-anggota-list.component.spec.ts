import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddAnggotaListComponent } from './add-anggota-list.component';

describe('AddAnggotaListComponent', () => {
  let component: AddAnggotaListComponent;
  let fixture: ComponentFixture<AddAnggotaListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddAnggotaListComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddAnggotaListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
