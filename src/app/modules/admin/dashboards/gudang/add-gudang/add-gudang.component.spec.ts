import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddGudangComponent } from './add-gudang.component';

describe('AddGudangComponent', () => {
  let component: AddGudangComponent;
  let fixture: ComponentFixture<AddGudangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddGudangComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddGudangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
