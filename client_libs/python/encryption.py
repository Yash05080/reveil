import base64
import os
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
# Try-catch import for back-compatibility if needed, but standardizing on pyca/cryptography

class EncryptionUtil:
    """
    AES-GCM Encryption Utility (Python)
    Compatible with Go implementation: AES-256, GCM, 12-byte Nonce prepended.
    Requires: pip install cryptography
    """
    
    @staticmethod
    def encrypt(key_base64: str, plaintext: str) -> str:
        key = base64.b64decode(key_base64)
        if len(key) != 32:
            raise ValueError("Key must be 32 bytes (AES-256)")
            
        nonce = os.urandom(12)
        
        # In GCM, authentication tag is generated automatically by encryptor
        cipher = Cipher(algorithms.AES(key), modes.GCM(nonce))
        encryptor = cipher.encryptor()
        
        ciphertext = encryptor.update(plaintext.encode('utf-8')) + encryptor.finalize()
        tag = encryptor.tag
        
        # Go Format: Nonce + Ciphertext + Tag
        # Note: Python cryptography GCM output is just ciphertext. Tag is separate.
        
        combined = nonce + ciphertext + tag
        return base64.b64encode(combined).decode('utf-8')

    @staticmethod
    def decrypt(key_base64: str, ciphertext_base64: str) -> str:
        key = base64.b64decode(key_base64)
        data = base64.b64decode(ciphertext_base64)
        
        if len(data) < 12 + 16:
            raise ValueError("Ciphertext too short")
            
        nonce = data[:12]
        
        # Tag is the last 16 bytes
        tag = data[-16:]
        
        # Ciphertext is between Nonce and Tag
        ciphertext = data[12:-16]
        
        cipher = Cipher(algorithms.AES(key), modes.GCM(nonce, tag))
        decryptor = cipher.decryptor()
        
        return (decryptor.update(ciphertext) + decryptor.finalize()).decode('utf-8')

# Example usage
if __name__ == "__main__":
    # Generate random key for testing
    # test_key = base64.b64encode(os.urandom(32)).decode('utf-8')
    # enc = EncryptionUtil.encrypt(test_key, "Hello Python")
    # print(enc)
    # print(EncryptionUtil.decrypt(test_key, enc))
    pass
