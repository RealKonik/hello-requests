package request

import (
	"crypto/tls"

	"github.com/tam7t/hpkp"
)

func GetCertStorage(hosts []string) (hpkp.Storage, error) {
	var err error

	s := hpkp.NewMemStorage()
	for _, host := range hosts {
		if host == "" {
			continue
		}
		conn, err := tls.Dial("tcp", host+":443", nil)
		if err != nil {
			return nil, err
		}
		var pins []string
		for _, cert := range conn.ConnectionState().PeerCertificates {
			pins = append(pins, hpkp.Fingerprint(cert))
		}
		s.Add(host, &hpkp.Header{
			Permanent:  true,
			Sha256Pins: pins,
		})
	}

	return s, err
}
