/*
 * Copyright (c) 2019 MengYX.
 */

package wrap

import (
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	"go.uber.org/zap"
)

type Subject pkix.Name

func CreateSubject(name string) *Subject {
	return &Subject{CommonName: name}
}

func (s *Subject) SetLocation(country, province, city string) {
	s.Country = []string{country}
	if province != "" {
		s.Province = []string{province}
	}
	if city != "" {
		s.Locality = []string{city}
	}
}

func (s *Subject) SetOrg(org, orgUnit string) {
	s.Organization = []string{org}
	if orgUnit != "" {
		s.OrganizationalUnit = []string{orgUnit}
	}
}

func (s *Subject) SetSerial(serial string) {
	s.SerialNumber = serial
}

func (s *Subject) GetRaw() *pkix.Name {
	return (*pkix.Name)(s)
}

func (s *Subject) GetKeyId() []byte {
	idHash := sha256.New()
	data, err := asn1.Marshal(*s)
	if err != nil {
		utils.Log().Error("compile subject failed", zap.Error(err))
	}
	idHash.Write(data)
	toEnc := "p2pNG-User-Id:" + s.CommonName
	idHash.Write([]byte(toEnc))
	keyId := idHash.Sum(nil)
	return keyId
}
