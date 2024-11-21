package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

type secureStore struct {
	sessions.Store
	block     cipher.Block
	options   *sessions.Options
	keyPrefix string
}

func NewSecureStore(store sessions.Store, options Options) (Store, error) {
	block, err := aes.NewCipher(options.SecurityKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	s := &secureStore{
		Store:     store,
		block:     block,
		keyPrefix: options.KeyPrefix,
		options: &sessions.Options{
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   options.MaxAge,
			Secure:   options.Secure,
			HttpOnly: options.HTTPOnly,
			SameSite: options.SameSite,
		},
	}

	return s, nil
}

func (s *secureStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	session, err := s.Store.Get(r, name)
	if err != nil {
		return nil, err
	}

	// Apply default options if not set
	if session.Options == nil {
		session.Options = s.options
	}

	// Decrypt session values if encrypted
	for key, value := range session.Values {
		if encryptedValue, ok := value.(string); ok {
			decryptedValue, err := s.decrypt(encryptedValue)
			if err != nil {
				continue // Skip invalid values
			}
			session.Values[key] = decryptedValue
		}
	}

	return session, nil
}

func (s *secureStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Encrypt session values before saving
	for key, value := range session.Values {
		if str, ok := value.(string); ok {
			encryptedValue, err := s.encrypt(str)
			if err != nil {
				continue // Skip values that can't be encrypted
			}
			session.Values[key] = encryptedValue
		}
	}

	return s.Store.Save(r, w, session)
}

func (s *secureStore) RegenerateID(_ *http.Request, _ http.ResponseWriter) (string, error) {
	// Generate new random session ID
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *secureStore) SetSecure(secure bool) {
	s.options.Secure = secure
}

func (s *secureStore) SetSameSite(mode http.SameSite) {
	s.options.SameSite = mode
}

func (s *secureStore) SetOptions(options *sessions.Options) {
	s.options = options
}

// Cleanup implements session cleanup for stores that support it
func (s *secureStore) Cleanup(threshold time.Time) error {
	if cleaner, ok := s.Store.(interface {
		Cleanup(time.Time) error
	}); ok {
		return cleaner.Cleanup(threshold)
	}
	return nil
}

// encrypt encrypts the given string using AES
func (s *secureStore) encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Create IV
	iv := make([]byte, s.block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(s.block, iv)
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// Return IV + ciphertext as base64
	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...)), nil
}

// decrypt decrypts the given string using AES
func (s *secureStore) decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Extract IV
	if len(data) < s.block.BlockSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := data[:s.block.BlockSize()]
	data = data[s.block.BlockSize():]

	// Decrypt
	plaintext := make([]byte, len(data))
	stream := cipher.NewCFBDecrypter(s.block, iv)
	stream.XORKeyStream(plaintext, data)

	return string(plaintext), nil
}
