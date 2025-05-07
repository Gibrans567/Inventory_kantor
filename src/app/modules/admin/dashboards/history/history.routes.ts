import { inject } from '@angular/core';
import { Routes } from '@angular/router';
import { FinanceComponent } from 'app/modules/admin/dashboards/finance/finance.component';
import { FinanceService } from 'app/modules/admin/dashboards/finance/finance.service';
import { HistoryComponent } from './history.component';

export default [
    {
        path: '',
        component: HistoryComponent,
    },
] as Routes;
