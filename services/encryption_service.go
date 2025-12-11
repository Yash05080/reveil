package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/google/uuid"
)

// EncryptionService handles encrypt/decrypt of post content
type EncryptionService struct {
	db      *sql.DB
	masterK []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(db *sql.DB, masterKey []byte) *EncryptionService {
	return &EncryptionService{
		db:      db,
		masterK: masterKey,
	}
}

// EncryptContent encrypts plaintext using a community-scoped AES key
func (s *EncryptionService) EncryptContent(communityID uuid.UUID, plaintext string) (string, error) {
	key, err := s.getOrCreateCommunityKey(communityID)
	if err != nil {
		return "", err
	}
	ciphertext, err := encryptAESGCM(key, []byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptContent decrypts content using the community-scoped AES key
func (s *EncryptionService) DecryptContent(communityID uuid.UUID, encoded string) (string, error) {
	key, err := s.getOrCreateCommunityKey(communityID)
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("invalid base64 content: %w", err)
	}
	plaintext, err := decryptAESGCM(key, raw)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// getOrCreateCommunityKey fetches or creates an AES-256 key for a community in encryption_keys
func (s *EncryptionService) getOrCreateCommunityKey(communityID uuid.UUID) ([]byte, error) {
	var encKey string
	err := s.db.QueryRow(`
		SELECT encrypted_key 
		FROM public.encryption_keys 
		WHERE community_id = $1
	`, communityID).Scan(&encKey)
	if err == sql.ErrNoRows {
		// create new key
		rawKey := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, rawKey); err != nil {
			return nil, fmt.Errorf("failed generating key: %w", err)
		}
		// For now, store raw key base64-encoded; later you can wrap with master key.
		encKey = base64.StdEncoding.EncodeToString(rawKey)

		_, err = s.db.Exec(`
			INSERT INTO public.encryption_keys (community_id, encrypted_key) 
			VALUES ($1, $2)
		`, communityID, encKey)
		if err != nil {
			return nil, fmt.Errorf("failed inserting encryption key: %w", err)
		}
		return rawKey, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed loading encryption key: %w", err)
	}

	rawKey, err := base64.StdEncoding.DecodeString(encKey)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption_key encoding: %w", err)
	}
	if len(rawKey) != 32 {
		return nil, fmt.Errorf("invalid key length: %d", len(rawKey))
	}
	return rawKey, nil
}

func encryptAESGCM(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return aesgcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAESGCM(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aesgcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ciphertext[:aesgcm.NonceSize()]
	data := ciphertext[aesgcm.NonceSize():]
	return aesgcm.Open(nil, nonce, data, nil)
}
