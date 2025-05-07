import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableInventarisComponent } from './table-inventaris.component';

describe('TableInventarisComponent', () => {
  let component: TableInventarisComponent;
  let fixture: ComponentFixture<TableInventarisComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TableInventarisComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TableInventarisComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
