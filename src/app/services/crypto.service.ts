import { Injectable } from '@angular/core';
import * as CryptoJS from 'crypto-js';
import { environment } from 'environment/environment-dev';

@Injectable({
    providedIn: 'root',
})
export class CryptoService {
    private iv = CryptoJS.enc.Utf8.parse('abcdefghijklmnop');

    /**
     * Enkripsi kunci sebelum digunakan di localStorage
     */
    private encryptKey(key: string): string {
        return CryptoJS.HmacSHA256(key, environment.keyCrypto).toString(CryptoJS.enc.Base64);
    }

    /**
     * Simpan data dengan kunci terenkripsi
     */
    setItem(key: string, value: string): void {
        const encryptedKey = this.encryptKey(key);
        localStorage.setItem(encryptedKey, value);
    }

    /**
     * Ambil data dengan kunci terenkripsi
     */
    getItem(key: string): string | null {
        const encryptedKey = this.encryptKey(key);
        return localStorage.getItem(encryptedKey);
    }

    /**
     * Hapus data dari localStorage
     */
    removeItem(key: string): void {
        const encryptedKey = this.encryptKey(key);
        localStorage.removeItem(encryptedKey);
    }
}
