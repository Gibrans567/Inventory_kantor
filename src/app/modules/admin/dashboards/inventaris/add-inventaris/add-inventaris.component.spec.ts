import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddInventarisComponent } from './add-inventaris.component';

describe('AddInventarisComponent', () => {
  let component: AddInventarisComponent;
  let fixture: ComponentFixture<AddInventarisComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddInventarisComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddInventarisComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
