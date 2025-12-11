const crypto = require('crypto');

/**
 * AES-GCM Encryption Utility (Node.js)
 * Compatible with Go implementation: AES-256, GCM, 12-byte Nonce prepended.
 */
class EncryptionUtil {
    static ALGORITHM = 'aes-256-gcm';
    static NONCE_LENGTH = 12;

    /**
     * Encrypts plaintext.
     * @param {string} keyBase64 - 32-byte key in Base64
     * @param {string} plaintext - Text to encrypt
     * @returns {string} Base64 encoded (Nonce + Ciphertext + Tag)
     */
    static encrypt(keyBase64, plaintext) {
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
     * @param {string} keyBase64 - 32-byte key in Base64
     * @param {string} ciphertextBase64 - Base64 encoded (Nonce + Ciphertext + Tag)
     * @returns {string} Decrypted plaintext
     */
    static decrypt(keyBase64, ciphertextBase64) {
        const key = Buffer.from(keyBase64, 'base64');
        const input = Buffer.from(ciphertextBase64, 'base64');

        if (input.length < this.NONCE_LENGTH + 16) throw new Error('Ciphertext too short');

        const nonce = input.subarray(0, this.NONCE_LENGTH);
        const tag = input.subarray(input.length - 16); // Tag is last 16 bytes
        const ciphertext = input.subarray(this.NONCE_LENGTH, input.length - 16);

        const decipher = crypto.createDecipheriv(this.ALGORITHM, key, nonce);
        decipher.setAuthTag(tag);

        let plaintext = decipher.update(ciphertext, null, 'utf8');
        plaintext += decipher.final('utf8');

        return plaintext;
    }
}

module.exports = EncryptionUtil;
