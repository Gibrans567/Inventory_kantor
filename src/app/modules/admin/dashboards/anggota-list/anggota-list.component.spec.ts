import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AnggotaListComponent } from './anggota-list.component';

describe('AnggotaListComponent', () => {
  let component: AnggotaListComponent;
  let fixture: ComponentFixture<AnggotaListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AnggotaListComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AnggotaListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
