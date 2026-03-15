package gluetun

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	z "github.com/Oudwins/zog"
)

type singleStringCase struct {
	name      string
	input     string
	wantValid bool
}

func runSingleStringSchemaCases(t *testing.T, schema *z.StringSchema[string], cases []singleStringCase) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dest string

			errs := schema.Parse(tc.input, &dest)

			gotValid := len(errs) == 0
			if gotValid != tc.wantValid {
				t.Fatalf("schema.Parse(%q) valid = %v, want %v", tc.input, gotValid, tc.wantValid)
			}
		})
	}
}

func TestListeningAddressSchema(t *testing.T) {
	schema := z.String().TestFunc(isValidListeningAddress)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_port", input: ":8080", wantValid: true},
		{name: "valid_zero_port", input: ":0", wantValid: true},
		{name: "missing_colon", input: "8080", wantValid: false},
		{name: "empty_string", input: "", wantValid: false},
		{name: "missing_port_digits", input: ":", wantValid: false},
		{name: "non_numeric_port", input: ":abc", wantValid: false},
		{name: "port_out_of_range", input: ":70000", wantValid: false},
	})
}

func TestFilePathSchema(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.txt")
	if err := os.WriteFile(filePath, []byte("vpnx"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	schema := z.String().TestFunc(isValidFilePath)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "existing_file", input: filePath, wantValid: true},
		{name: "existing_file_with_whitespace", input: " " + filePath + " ", wantValid: true},
		{name: "missing_file", input: filepath.Join(tempDir, "missing.txt"), wantValid: false},
		{name: "directory_path", input: tempDir, wantValid: false},
		{name: "null_byte", input: "bad\x00path", wantValid: false},
		{name: "empty_string", input: "", wantValid: false},
	})
}

func TestTimePeriodSchema(t *testing.T) {
	schema := z.String().TestFunc(isValidTimePeriod)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_seconds", input: "15s", wantValid: true},
		{name: "valid_minutes", input: "45m", wantValid: true},
		{name: "valid_hours", input: "24h", wantValid: true},
		{name: "invalid_suffix", input: "15d", wantValid: false},
		{name: "non_numeric_period", input: "xxm", wantValid: false},
		{name: "empty_string", input: "", wantValid: false},
		{name: "missing_suffix", input: "15", wantValid: false},
		{name: "seconds_out_of_range", input: "61s", wantValid: false},
		{name: "hours_out_of_range", input: "25h", wantValid: false},
		{name: "negative_period", input: "-1m", wantValid: false},
	})
}

func TestHostnameSchema(t *testing.T) {
	schema := z.String().TestFunc(isValidHostname)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_hostname", input: "example.com", wantValid: true},
		{name: "valid_trailing_dot", input: "sub.example.com.", wantValid: true},
		{name: "missing_dot", input: "localhost", wantValid: false},
		{name: "leading_hyphen", input: "-bad.example.com", wantValid: false},
		{name: "invalid_character", input: "bad_name.example.com", wantValid: false},
		{name: "label_too_long", input: strings.Repeat("a", 64) + ".com", wantValid: false},
		{name: "empty_string", input: "", wantValid: false},
	})
}

func TestPlainAddressSchema(t *testing.T) {
	schema := z.String().TestFunc(isValidPlainAddress)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_plain_address", input: "1.2.3.4:53", wantValid: true},
		{name: "empty_string", input: "", wantValid: false},
		{name: "missing_port", input: "1.2.3.4", wantValid: false},
		{name: "empty_host", input: ":53", wantValid: false},
		{name: "hostname_instead_of_ipv4", input: "example.com:53", wantValid: false},
		{name: "ipv6_rejected", input: "2001:db8::1:53", wantValid: false},
		{name: "non_numeric_port", input: "1.2.3.4:http", wantValid: false},
	})
}

func TestSubnetSliceSchema(t *testing.T) {
	schema := z.Slice(z.String().TestFunc(isValidSubnet))

	testCases := []struct {
		name      string
		input     []string
		wantValid bool
	}{
		{name: "valid_ipv4_subnet", input: []string{"192.168.0.1/24"}, wantValid: true},
		{name: "valid_multiple_ipv4_subnets", input: []string{"192.168.0.1/24", "10.0.0.0/8"}, wantValid: true},
		{name: "invalid_multiple_subnets", input: []string{"192.168.0.1/24", ""}, wantValid: false},
		{name: "ipv6_subnet_rejected", input: []string{"2001:db8::/32"}, wantValid: false},
		{name: "invalid_subnet_string", input: []string{"invalid"}, wantValid: false},
		{name: "empty_subnet", input: []string{""}, wantValid: false},
		{name: "empty_slice_is_valid", input: []string{}, wantValid: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dest []string

			errs := schema.Parse(tc.input, &dest)

			gotValid := len(errs) == 0
			if gotValid != tc.wantValid {
				t.Fatalf("schema.Parse(%q) valid = %v, want %v", tc.input, gotValid, tc.wantValid)
			}

			if tc.wantValid && len(dest) != len(tc.input) {
				t.Fatalf("expected %d parsed subnets, got %d", len(tc.input), len(dest))
			}
		})
	}
}

func TestCIDRAddressSchema(t *testing.T) {
	schema := z.String().TestFunc(isValidCIDRAddress)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_ipv4_cidr", input: "192.168.1.0/24", wantValid: true},
		{name: "ipv6_cidr_rejected", input: "2001:db8::/32", wantValid: false},
		{name: "missing_prefix", input: "192.168.1.1", wantValid: false},
		{name: "invalid_cidr", input: "not-a-cidr", wantValid: false},
		{name: "empty_string", input: "", wantValid: false},
	})
}

func TestBase64PEMSchema(t *testing.T) {
	fixtures := newPEMFixtures(t)
	schema := z.String().TestFunc(isValidBase64PEM)

	runSingleStringSchemaCases(t, schema, []singleStringCase{
		{name: "valid_certificate_pem", input: fixtures.certificatePEM, wantValid: true},
		{name: "valid_base64_encoded_pem", input: fixtures.base64CertificatePEM, wantValid: true},
		{name: "valid_pkcs8_private_key_pem", input: fixtures.privateKeyPEM, wantValid: true},
		{name: "empty_string", input: "", wantValid: false},
		{name: "not_pem_or_base64", input: "%%%not-valid%%%", wantValid: false},
		{name: "base64_non_pem_payload", input: fixtures.base64NonPEM, wantValid: false},
		{name: "unsupported_pem_block_type", input: fixtures.unsupportedPEM, wantValid: false},
		{name: "legacy_encrypted_rsa_private_key", input: fixtures.legacyEncryptedRSAPrivateKeyPEM, wantValid: false},
	})
}

func TestNumberGluetunEnvTreatsEmptyStringAsUnset(t *testing.T) {
	env := newNumberGluetunEnv("PUID", z.Int())

	if errs := env.validate(""); len(errs) != 0 {
		t.Fatalf("expected empty string to be treated as unset, got %v", errs)
	}

	if errs := env.validate("1000"); len(errs) != 0 {
		t.Fatalf("expected valid int env, got %v", errs)
	}

	if errs := env.validate("abc"); len(errs) == 0 {
		t.Fatal("expected invalid int env to fail")
	}
}

func TestRequiredGluetunEnvTreatsEmptyStringAsMissing(t *testing.T) {
	env := newGluetunEnv("OPENVPN_USER", z.String().Required())

	if errs := env.validate(""); len(errs) == 0 {
		t.Fatal("expected empty required env to fail")
	}
}

func TestBoolGluetunEnvParsesOnOffValues(t *testing.T) {
	env := newBoolGluetunEnv("PUBLIC_IP_ENABLED", z.Bool())

	if errs := env.validate("on"); len(errs) != 0 {
		t.Fatalf("expected on to parse as bool, got %v", errs)
	}

	if errs := env.validate("off"); len(errs) != 0 {
		t.Fatalf("expected off to parse as bool, got %v", errs)
	}
}

func TestCSVGluetunEnvValidatesCommaSeparatedValues(t *testing.T) {
	env := newCSVGluetunEnv[int]("FIREWALL_INPUT_PORTS", z.Slice(z.Int().GTE(0).LTE(65535)))

	if errs := env.validate("80,443"); len(errs) != 0 {
		t.Fatalf("expected valid csv env, got %v", errs)
	}

	if errs := env.validate("80,nope"); len(errs) == 0 {
		t.Fatal("expected invalid csv env to fail")
	}
}

type pemFixtures struct {
	certificatePEM                  string
	base64CertificatePEM            string
	privateKeyPEM                   string
	base64NonPEM                    string
	unsupportedPEM                  string
	legacyEncryptedRSAPrivateKeyPEM string
}

func newPEMFixtures(t *testing.T) pemFixtures {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ECDSA key: %v", err)
	}

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal PKCS#8 private key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "vpnx.test",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	certificateDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	certificatePEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificateDER,
	}))

	privateKeyPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	}))

	unsupportedPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: []byte("unsupported"),
	}))

	return pemFixtures{
		certificatePEM:       certificatePEM,
		base64CertificatePEM: base64.StdEncoding.EncodeToString([]byte(certificatePEM)),
		privateKeyPEM:        privateKeyPEM,
		base64NonPEM:         base64.StdEncoding.EncodeToString([]byte("not a pem document")),
		unsupportedPEM:       unsupportedPEM,
		legacyEncryptedRSAPrivateKeyPEM: strings.Join([]string{
			"-----BEGIN RSA PRIVATE KEY-----",
			"Proc-Type: 4,ENCRYPTED",
			"DEK-Info: AES-256-CBC,0123456789ABCDEF0123456789ABCDEF",
			"",
			base64.StdEncoding.EncodeToString([]byte("encrypted payload")),
			"-----END RSA PRIVATE KEY-----",
			"",
		}, "\n"),
	}
}
