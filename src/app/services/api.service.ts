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

     private getHeaders() {
        return {
            headers: {
                'Authorization': `Bearer ${this._cryptoService.getItem("accessToken")}`,
                'Divisi_id': `${this._cryptoService.getItem("divisionId")}`,
                "Accept": "application/json",
                "Content-Type": "application/json"
            }
        };
    }

    async get(endPoint: string, headers: boolean) {
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

    async delete(endPoint: string, headers: boolean) {
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
        try {
            const response = await axios.post(this.apiUrl + endPoint, data);
            return response.data;
        } catch (error) {
            this._loadingService.hide()
            return { data: error, status: error.response?.status };
        } finally {
            this._loadingService.hide()
        }
    }

    async put(endPoint: string, data: any) {
        this._loadingService.show();
        return await axios.put(this.apiUrl + endPoint, data).then(
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
