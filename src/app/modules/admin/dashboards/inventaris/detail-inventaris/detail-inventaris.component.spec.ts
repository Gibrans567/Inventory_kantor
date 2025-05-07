import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DetailInventarisComponent } from './detail-inventaris.component';

describe('DetailInventarisComponent', () => {
  let component: DetailInventarisComponent;
  let fixture: ComponentFixture<DetailInventarisComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DetailInventarisComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DetailInventarisComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
