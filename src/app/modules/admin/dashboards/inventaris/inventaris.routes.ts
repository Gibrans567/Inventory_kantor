import { Routes } from '@angular/router';
import { InventarisComponent } from './inventaris.component';

export default [
    {
        path: '',
        component: InventarisComponent,
    },
    {
        path: 'view-detail/:id',
        loadChildren: () => import('app/modules/admin/dashboards/inventaris/detail-inventaris/detail-inventaris.routes')
    },

    {
        path: 'add-user',
        loadChildren: () => import('app/modules/admin/dashboards/inventaris/add-inventaris/add-inventaris.routes')
    },
    // {
    //     path: 'edit-user/:name',
    //     loadChildren: () => import('app/modules/admin/dashboards/customer/form-edit/form-edit.routes'),
    // }
] as Routes;
