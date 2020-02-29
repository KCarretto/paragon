package cdn

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kcarretto/paragon/pkg/auth"
)

type Client struct {
	HTTP *http.Client
	URL  string

	Service    string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// Upload a file to the CDN.
func (cdn *Client) Upload(name string, file io.Reader) error {
	url := fmt.Sprintf("%s/%s", cdn.URL, "cdn/upload")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("fileContent", name)
	if err != nil {
		return fmt.Errorf("failed to create form: %w", err)
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}

	if err = writer.WriteField("fileName", name); err != nil {
		return fmt.Errorf("failed to write file name: %w", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if cdn.HTTP == nil {
		cdn.HTTP = http.DefaultClient
	}

	// Sign http request
	epoch := fmt.Sprintf("%d", time.Now().Unix())
	sig, err := cdn.sign([]byte(epoch))
	if err != nil {
		panic(fmt.Errorf("failed to sign request: %w", err))
	}
	req.Header.Set(auth.HeaderService, cdn.Service)
	req.Header.Set(auth.HeaderIdentity, base64.StdEncoding.EncodeToString(cdn.PublicKey))
	req.Header.Set(auth.HeaderEpoch, epoch)
	req.Header.Set(auth.HeaderSignature, base64.StdEncoding.EncodeToString(sig))

	if _, err = cdn.HTTP.Do(req); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// Download a file from the CDN.
func (cdn *Client) Download(name string) (io.Reader, error) {
	url := fmt.Sprintf("%s/%s/%s", cdn.URL, "cdn/download", name)
	if cdn.HTTP == nil {
		cdn.HTTP = http.DefaultClient
	}

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	// Sign http request
	epoch := fmt.Sprintf("%d", time.Now().Unix())
	sig, err := cdn.sign([]byte(epoch))
	if err != nil {
		panic(fmt.Errorf("failed to sign request: %w", err))
	}
	httpReq.Header.Set(auth.HeaderService, cdn.Service)
	httpReq.Header.Set(auth.HeaderIdentity, base64.StdEncoding.EncodeToString(cdn.PublicKey))
	httpReq.Header.Set(auth.HeaderEpoch, epoch)
	httpReq.Header.Set(auth.HeaderSignature, base64.StdEncoding.EncodeToString(sig))

	resp, err := cdn.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return resp.Body, nil
}

func (cdn *Client) sign(msg []byte) ([]byte, error) {
	// If nil, try loading from the environment
	if cdn.PublicKey == nil || cdn.PrivateKey == nil {
		cdn.PublicKey, cdn.PrivateKey = cdn.keyFromEnv()
	}

	// Still nil, generate a keypair
	if cdn.PublicKey == nil || cdn.PrivateKey == nil {
		pubKey, privKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate keypair: %w", err)
		}
		cdn.PublicKey = pubKey
		cdn.PrivateKey = privKey
	}

	// Sign the message using the client's private key
	return cdn.PrivateKey.Sign(nil, msg, crypto.Hash(0))
}

func (cdn *Client) keyFromEnv() (ed25519.PublicKey, ed25519.PrivateKey) {
	if key := os.Getenv("PG_SVC_KEY"); key != "" {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid format for PG_SVC_KEY, expected b64PubKey:b64PrivKey"))
		}
		pubKey, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			panic(fmt.Errorf("invalid base64 provided for PubKey in PG_SVC_KEY: %w", err))
		}
		privKey, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			panic(fmt.Errorf("invalid base64 provided for PrivKey in PG_SVC_KEY: %w", err))
		}
		return pubKey, privKey
	}
	return nil, nil
}
