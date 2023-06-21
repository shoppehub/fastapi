package session

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// 获取ip
func GetIP(r *http.Request) string {
	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	logrus.Error("no ip!!!")
	return ip
}

// 获取域名
func GetDomain(host string) string {
	host = strings.TrimSpace(host)

	hostPorts := strings.Split(host, ":")

	hostParts := strings.Split(hostPorts[0], ".")
	lengthOfHostParts := len(hostParts)

	if lengthOfHostParts == 1 {
		return hostParts[0] // scenario C
	} else {
		_, err := strconv.ParseFloat(hostParts[0], 64)
		if err == nil {
			return host
		} else {
			return strings.Join([]string{hostParts[lengthOfHostParts-2], hostParts[lengthOfHostParts-1]}, ".")
		}
	}
}

// 获取运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Error(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Encrypt encrypts the input string using AES-128 CBC mode with PKCS7 padding.
// The key and iv must be 16 bytes long.
func Encrypt(inputStr string, key string, iv string) (string, error) {
	if iv == "" {
		iv = "chemball"
	}
	// Create AES cipher instance
	block, err := aes.NewCipher(pkcs7Pad([]byte(key), aes.BlockSize))
	if err != nil {
		return "", err
	}

	// Pad the input string
	paddedInput := pkcs7Pad([]byte(inputStr), aes.BlockSize)

	// Create CBC mode cipher
	cbc := cipher.NewCBCEncrypter(block, pkcs7Pad([]byte(iv), aes.BlockSize))

	// Encrypt the padded input
	encrypted := make([]byte, len(paddedInput))
	cbc.CryptBlocks(encrypted, paddedInput)

	// Encode the encrypted data using base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// pkcs7Pad pads the input data using PKCS7 padding scheme.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := strings.Repeat(" ", padding)
	return append(data, []byte(padText)...)
}
