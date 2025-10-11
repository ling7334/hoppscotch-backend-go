package graph

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dto"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func getLimit(take *int) int {
	if take != nil {
		return *take
	}
	return 10
}

func contains(s []string, elem string) bool {
	for _, a := range s {
		if a == elem {
			return true
		}
	}
	return false
}

func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

// Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func isServiceConfigured(provider dto.AuthProvider, cfgMap map[string]string) bool {
	switch provider {
	case "GOOGLE":
		return cfgMap["GOOGLE_CLIENT_ID"] != "" &&
			cfgMap["GOOGLE_CLIENT_SECRET"] != "" &&
			cfgMap["GOOGLE_CALLBACK_URL"] != "" &&
			cfgMap["GOOGLE_SCOPE"] != ""
	case "GITHUB":
		return cfgMap["GITHUB_CLIENT_ID"] != "" &&
			cfgMap["GITHUB_CLIENT_SECRET"] != "" &&
			cfgMap["GITHUB_CALLBACK_URL"] != "" &&
			cfgMap["GITHUB_SCOPE"] != ""
	case "MICROSOFT":
		return cfgMap["MICROSOFT_CLIENT_ID"] != "" &&
			cfgMap["MICROSOFT_CLIENT_SECRET"] != "" &&
			cfgMap["MICROSOFT_CALLBACK_URL"] != "" &&
			cfgMap["MICROSOFT_SCOPE"] != "" &&
			cfgMap["MICROSOFT_TENANT"] != ""
	case "EMAIL":
		if cfgMap["MAILER_SMTP_ENABLE"] != "true" {
			return false
		}
		if cfgMap["MAILER_USE_CUSTOM_CONFIGS"] == "true" {
			return cfgMap["MAILER_SMTP_HOST"] != "" &&
				cfgMap["MAILER_SMTP_PORT"] != "" &&
				cfgMap["MAILER_SMTP_SECURE"] != "" &&
				cfgMap["MAILER_SMTP_USER"] != "" &&
				cfgMap["MAILER_SMTP_PASSWORD"] != "" &&
				cfgMap["MAILER_TLS_REJECT_UNAUTHORIZED"] != "" &&
				cfgMap["MAILER_ADDRESS_FROM"] != ""
		}
		return cfgMap["MAILER_SMTP_URL"] != "" && cfgMap["MAILER_ADDRESS_FROM"] != ""
	default:
		return false
	}
}

const EncryptionAlgorithm = "aes-256-cbc"

// Encrypt 加密文本，返回 iv:密文（十六进制字符串）
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) == 0 {
		envKey := os.Getenv("DATA_ENCRYPTION_KEY")
		if envKey == "" {
			return "", errors.New("ENV_NOT_FOUND_KEY_DATA_ENCRYPTION_KEY")
		}
		key = []byte(envKey)
	}

	if plaintext == "" {
		return plaintext, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// CBC 模式需要填充
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return hex.EncodeToString(iv) + ":" + hex.EncodeToString(ciphertext), nil
}

// Decrypt 解密文本，输入格式为 iv:密文
func Decrypt(encrypted string, key []byte) (string, error) {
	if len(key) == 0 {
		envKey := os.Getenv("DATA_ENCRYPTION_KEY")
		if envKey == "" {
			return "", errors.New("ENV_NOT_FOUND_KEY_DATA_ENCRYPTION_KEY")
		}
		key = []byte(envKey)
	}

	if encrypted == "" {
		return encrypted, nil
	}

	parts := bytes.SplitN([]byte(encrypted), []byte(":"), 2)
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted data format")
	}

	iv, err := hex.DecodeString(string(parts[0]))
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(string(parts[1]))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	unpadded, err := pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}

// pkcs7Pad 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad 去填充
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid PKCS#7 padding (length)")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid PKCS#7 padding (value)")
	}
	return data[:len(data)-padding], nil
}

func getDefaultInfraConfigs(excluded []string) (res map[dto.InfraConfigEnum]string) {
	for _, k := range dto.AllInfraConfigEnum {
		if contains(excluded, k.String()) {
			continue
		}
		val := os.Getenv(k.String())
		if val != "" {
			res[k] = val
		}
	}
	return
}
