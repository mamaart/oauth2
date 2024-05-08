package p256

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func Get(name string) tls.Certificate {
	cert, key := fmt.Sprintf("%s.crt", name), fmt.Sprintf("%s.key", name)
	return must(tls.LoadX509KeyPair(cert, key))
}

func Generate(name string) {
	key := must(ecdsa.GenerateKey(elliptic.P256(), rand.Reader))
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Martin Maartensson"},
			CommonName:   "localhost",
		},
		EmailAddresses:        []string{"martinmaartensson@gmail.com"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	if err := pem.Encode(
		must(os.Create(fmt.Sprintf("%s.crt", name))),
		&pem.Block{
			Type: "CERTIFICATE",
			Bytes: must(
				x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key),
			),
		},
	); err != nil {
		panic(err)
	}
	if err := pem.Encode(
		must(os.Create(fmt.Sprintf("%s.key", name))),
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: must(x509.MarshalECPrivateKey(key)),
		},
	); err != nil {
		panic(err)
	}
}

func Read(name string) {}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
