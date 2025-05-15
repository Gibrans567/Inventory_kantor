import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GudangComponent } from './gudang.component';

describe('GudangComponent', () => {
  let component: GudangComponent;
  let fixture: ComponentFixture<GudangComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GudangComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(GudangComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
