package com.reveil.utils

import java.security.SecureRandom
import java.util.Base64
import javax.crypto.Cipher
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.SecretKeySpec

object EncryptionUtil {
    private const val ALGORITHM = "AES/GCM/NoPadding"
    private const val TAG_LENGTH_BIT = 128 // 16 bytes * 8
    private const val NONCE_LENGTH = 12    // 12 bytes

    /**
     * Encrypts plaintext using AES-GCM.
     * Matches Go's: aesgcm.Seal(nonce, nonce, plaintext, nil)
     *
     * @param keyBase64: Base64 string of the 32-byte AES key
     * @param plaintext: The string to encrypt
     * @return Base64 string of (Nonce + Ciphertext + Tag)
     */
    fun encrypt(keyBase64: String, plaintext: String): String {
        // 1. Decode Key
        val keyBytes = Base64.getDecoder().decode(keyBase64)
        if (keyBytes.size != 32) throw IllegalArgumentException("Key must be 32 bytes (AES-256)")
        val secretKey = SecretKeySpec(keyBytes, "AES")

        // 2. Generate Nonce
        val nonce = ByteArray(NONCE_LENGTH)
        SecureRandom().nextBytes(nonce)

        // 3. Initialize Cipher
        val cipher = Cipher.getInstance(ALGORITHM)
        val spec = GCMParameterSpec(TAG_LENGTH_BIT, nonce)
        cipher.init(Cipher.ENCRYPT_MODE, secretKey, spec)

        // 4. Encrypt
        val ciphertext = cipher.doFinal(plaintext.toByteArray(Charsets.UTF_8))

        // 5. Combine Nonce + Ciphertext (Tag is included in ciphertext by doFinal in Java)
        // Wait, verifying if Java GCM includes tag in doFinal output. Yes it does.
        // Go's Seal: Appends ciphertext+tag to dst. We prepended nonce manually.
        // Final layout: [Nonce][Ciphertext][Tag]
        
        val combined = ByteArray(nonce.size + ciphertext.size)
        System.arraycopy(nonce, 0, combined, 0, nonce.size)
        System.arraycopy(ciphertext, 0, combined, nonce.size, ciphertext.size)

        return Base64.getEncoder().encodeToString(combined)
    }

    /**
     * Decrypts ciphertext using AES-GCM.
     *
     * @param keyBase64: Base64 string of the 32-byte AES key
     * @param ciphertextBase64: Base64 string of (Nonce + Ciphertext + Tag)
     * @return Decrypted String
     */
    fun decrypt(keyBase64: String, ciphertextBase64: String): String {
        // 1. Decode Key & Input
        val keyBytes = Base64.getDecoder().decode(keyBase64)
        val combined = Base64.getDecoder().decode(ciphertextBase64)

        if (keyBytes.size != 32) throw IllegalArgumentException("Key must be 32 bytes (AES-256)")
        if (combined.size < NONCE_LENGTH + (TAG_LENGTH_BIT / 8)) throw IllegalArgumentException("Ciphertext too short")

        val secretKey = SecretKeySpec(keyBytes, "AES")

        // 2. Extract Nonce
        val nonce = ByteArray(NONCE_LENGTH)
        System.arraycopy(combined, 0, nonce, 0, NONCE_LENGTH)

        // 3. Extract Ciphertext + Tag
        val ciphertextLength = combined.size - NONCE_LENGTH
        val ciphertext = ByteArray(ciphertextLength)
        System.arraycopy(combined, NONCE_LENGTH, ciphertext, 0, ciphertextLength)

        // 4. Decrypt
        val cipher = Cipher.getInstance(ALGORITHM)
        val spec = GCMParameterSpec(TAG_LENGTH_BIT, nonce)
        cipher.init(Cipher.DECRYPT_MODE, secretKey, spec)

        val plaintextBytes = cipher.doFinal(ciphertext)
        return String(plaintextBytes, Charsets.UTF_8)
    }
}

// Example Usage Main (Commented out)
/*
fun main() {
    val key = "uT5v/1x...32bytes...=="
    val text = "Hello World"
    val enc = EncryptionUtil.encrypt(key, text)
    println("Encrypted: $enc")
    val dec = EncryptionUtil.decrypt(key, enc)
    println("Decrypted: $dec")
}
*/
