import { Injectable, inject } from '@angular/core';
import axios from 'axios';
import { environment } from 'environment/environment-dev';
import { FuseLoadingService } from '@fuse/services/loading';
import { CryptoService } from './crypto.service';

@Injectable({
    providedIn: 'root'
})
export class ApiService {
    private apiUrl = environment.apiUrl;
    private _loadingService = inject(FuseLoadingService);
    private _cryptoService = inject(CryptoService)

    async get(endPoint: string) {
        this._loadingService.show();
        return await axios.get(this.apiUrl + endPoint).then(
            async response => {
                this._loadingService.hide()
                return response.data
            }).catch(
                async error => {
                    this._loadingService.hide()
                    console.log(error)
                }
            )
    }

    async delete(endPoint: string) {
        this._loadingService.show();
        return await axios.delete(this.apiUrl + endPoint).then(
            async response => {
                this._loadingService.hide()
                return response.data
            }).catch(
                async error => {
                    this._loadingService.hide()
                    console.log(error)
                }
            )
    }

    async post(endPoint: string, data: any) {
        this._loadingService.show();
        return await axios.post(this.apiUrl + endPoint, data).then(
            async response => {
                this._loadingService.hide()
                return response.data
            }).catch(
                async error => {
                    this._loadingService.hide()
                    console.log(error)
                }
            )
    }
}
