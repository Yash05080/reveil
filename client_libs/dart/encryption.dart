import 'dart:convert';
import 'dart:typed_data';
import 'dart:math';

// Requires: encrypt package (pubspec.yaml: encrypt: ^5.0.0)
// This is a plain file, user must add dependency.
// For standard library only implementation, it's very hard in Dart without packages.
// Assuming 'encrypt' package is available or using pointycastle.
// I will provide code using 'encrypt' package as it is the standard.

// To verify: AES GCM Mode in 'encrypt' package.
// If not available easily, I will implement using PointyCastle logic if needed, 
// but sticking to high level 'encrypt' package is safer for users.
// Note: 'encrypt' package AES mode.gcm might be tricky.
// Let's use 'pointycastle' directly if we want pure compatibility without bulky wrappers?
// Actually 'encrypt' wraps pointycastle.

// Implementation using 'encrypt' package:
/*
import 'package:encrypt/encrypt.dart' as enc;

class EncryptionUtil {
  static String encrypt(String keyBase64, String plaintext) {
    final key = enc.Key.fromBase64(keyBase64);
    final iv = enc.IV.fromLength(12); // Random IV
    final encrypter = enc.Encrypter(enc.AES(key, mode: enc.AESMode.gcm));

    final encrypted = encrypter.encrypt(plaintext, iv: iv);
    
    // Output: Nonce (IV) + Ciphertext + Tag
    // Encrypted.bytes contains ciphertext + tag (usually).
    // Wait, encrypt package GCM keeps tag separate or inside?
    // It's safer to check documentation. 
  }
}
*/

// Since I cannot verify external packages easily, I will write this assuming the user adds `encrypt` package
// and handles the concatenation correctly. The `encrypt` package `Encrypted` object holds bytes. 
// GCM mode in `encrypt` puts tag at the end of bytes automatically in recent versions?
// Let's write a standard class that users can drop in, but warn about dependencies.

import 'package:encrypt/encrypt.dart';
import 'package:crypto/crypto.dart';

class EncryptionUtil {
  /// Encrypts plaintext using AES-GCM (256).
  /// Required pubspec: encrypt: ^5.0.3
  static String encrypt(String keyBase64, String plaintext) {
    final key = Key.fromBase64(keyBase64);
    
    // Generate 12 byte IV (Nonce)
    final iv = IV.fromSecureRandom(12);
    
    final encrypter = Encrypter(AES(key, mode: AESMode.gcm));
    
    final encrypted = encrypter.encrypt(plaintext, iv: iv);
    
    // encrypted.bytes usually includes the tag for GCM in this package.
    // Concatenate IV + EncryptedBytes
    final combined = Uint8List(iv.bytes.length + encrypted.bytes.length);
    combined.setAll(0, iv.bytes);
    combined.setAll(iv.bytes.length, encrypted.bytes);
    
    return base64.encode(combined);
  }

  /// Decrypts ciphertext.
  static String decrypt(String keyBase64, String ciphertextBase64) {
    final key = Key.fromBase64(keyBase64);
    final combined = base64.decode(ciphertextBase64);
    
    if (combined.length < 12 + 16) throw Exception("Ciphertext too short");
    
    // Extract IV (first 12 bytes)
    final iv = IV(combined.sublist(0, 12));
    
    // Extract Ciphertext+Tag (rest)
    final encryptedBytes = combined.sublist(12);
    final encrypted = Encrypted(encryptedBytes);
    
    final encrypter = Encrypter(AES(key, mode: AESMode.gcm));
    
    return encrypter.decrypt(encrypted, iv: iv);
  }
}
