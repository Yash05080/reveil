import * as crypto from 'crypto';

/**
 * AES-GCM Encryption Utility (TypeScript/Node.js)
 * Compatible with Go implementation: AES-256, GCM, 12-byte Nonce prepended.
 */
export class EncryptionUtil {
    private static readonly ALGORITHM = 'aes-256-gcm';
    private static readonly NONCE_LENGTH = 12;

    /**
     * Encrypts plaintext.
     * @param keyBase64 - 32-byte key in Base64
     * @param plaintext - Text to encrypt
     * @returns Base64 encoded (Nonce + Ciphertext + Tag)
     */
    static encrypt(keyBase64: string, plaintext: string): string {
        const key = Buffer.from(keyBase64, 'base64');
        if (key.length !== 32) throw new Error('Key must be 32 bytes (AES-256)');

        const nonce = crypto.randomBytes(this.NONCE_LENGTH);
        const cipher = crypto.createCipheriv(this.ALGORITHM, key, nonce);

        let ciphertext = cipher.update(plaintext, 'utf8');
        const final = cipher.final();
        const tag = cipher.getAuthTag();

        const combined = Buffer.concat([nonce, ciphertext, final, tag]);
        return combined.toString('base64');
    }

    /**
     * Decrypts ciphertext.
     * @param keyBase64 - 32-byte key in Base64
     * @param ciphertextBase64 - Base64 encoded (Nonce + Ciphertext + Tag)
     * @returns Decrypted plaintext
     */
    static decrypt(keyBase64: string, ciphertextBase64: string): string {
        const key = Buffer.from(keyBase64, 'base64');
        const input = Buffer.from(ciphertextBase64, 'base64');

        if (input.length < this.NONCE_LENGTH + 16) throw new Error('Ciphertext too short');

        const nonce = input.subarray(0, this.NONCE_LENGTH);
        const tag = input.subarray(input.length - 16);
        const ciphertext = input.subarray(this.NONCE_LENGTH, input.length - 16);

        const decipher = crypto.createDecipheriv(this.ALGORITHM, key, nonce);
        decipher.setAuthTag(tag);

        let plaintext = decipher.update(ciphertext, undefined, 'utf8'); // Using undefined for binary input
        plaintext += decipher.final('utf8');

        return plaintext;
    }
}
