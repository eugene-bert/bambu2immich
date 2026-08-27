package bambu

import "crypto/tls"

// printerTLS is used for MQTT and FTPS. Bambu printers present a self-signed
// X.509 v1 certificate that standard verification rejects. Traffic stays on
// the local network; see the README security notes.
func printerTLS() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec // Bambu self-signed X.509 v1
		MinVersion:         tls.VersionTLS12,
	}
}
