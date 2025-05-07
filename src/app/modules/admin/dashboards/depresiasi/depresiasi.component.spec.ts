import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DepresiasiComponent } from './depresiasi.component';

describe('DepresiasiComponent', () => {
  let component: DepresiasiComponent;
  let fixture: ComponentFixture<DepresiasiComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DepresiasiComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DepresiasiComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
