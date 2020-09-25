/*
 * Copyright (c) 2019 MengYX.
 */

package wrap

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math"
	"math/big"
	"time"
)

type CertTemplate x509.Certificate

func CreateTemplate(subject *pkix.Name, keyId []byte, serial int64) CertTemplate {
	template := CertTemplate{
		SerialNumber: big.NewInt(serial),
		Subject:      *subject,
		SubjectKeyId: keyId,
	}
	return template
}

func (cert *CertTemplate) SetExpire(year, month, day int) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	cert.NotBefore = today
	cert.NotAfter = today.AddDate(year, month, day)
}

func (cert *CertTemplate) SetConstraint(isCA, LimitPathLen bool, pathLen int) {
	cert.BasicConstraintsValid = true
	cert.IsCA = isCA
	if isCA {
		if LimitPathLen {
			if 0 == pathLen {
				cert.MaxPathLenZero = true
			} else {
				cert.MaxPathLenZero = false
				cert.MaxPathLen = pathLen
			}
		} else {
			cert.MaxPathLenZero = false
			cert.MaxPathLen = 0
		}
	}
}

func (cert *CertTemplate) SetUsage(usage []int) {
	all := 0
	for idx := range usage {
		all = all | int(math.Pow(2, float64(usage[idx])))
	}
	cert.KeyUsage = x509.KeyUsage(all)
}
func (cert *CertTemplate) SetExtUsage(usage []int) {
	var all []x509.ExtKeyUsage
	for e := range usage {
		all = append(all, x509.ExtKeyUsage(usage[e]))
	}

	cert.ExtKeyUsage = all
}

func (cert *CertTemplate) GetRaw() *x509.Certificate {
	return (*x509.Certificate)(cert)
}

func (cert *CertTemplate) SetAlgorithm(category string, name string) {
	//仅用于验证证书本身 无关加密算法
	algorithm := 0
	switch category {

	case "ecdsa":
		algorithm = 9
	case "rsa":
		algorithm = 12
	case "ed25519":
		algorithm = 16
	}

	switch name {
	case "sha256":
		algorithm += 1
	case "sha384":
		algorithm += 2
	case "sha512":
		algorithm += 3
	}
	cert.SignatureAlgorithm = x509.SignatureAlgorithm(algorithm)
}
