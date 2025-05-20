import { inject, Injectable, signal } from '@angular/core';
import { ApiService } from 'app/services/api.service';
import { ToastrService } from 'ngx-toastr';
import { InventoryItem } from 'types/inventaris';

@Injectable({
    providedIn: 'root'
})
export class InventoryService {
    inventoryItems = signal<InventoryItem[]>([]);
    isLoading = signal<boolean>(false);
    isNotFound = signal<boolean>(false)
    private _apiService = inject(ApiService)

    fetchData(): void {
        this.isLoading.set(true)
        this._apiService.get("/inventaris",true)
            .then(response => {
                this.isLoading.set(false)
                this.isNotFound.set(false)
                this.inventoryItems.set(response.data ?? []);
            })
            .catch(error => {
                this.isLoading.set(false)
                this.isNotFound.set(true)
                console.error("Error fetching inventory:", error);
            });
    }
}
