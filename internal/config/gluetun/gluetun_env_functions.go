package gluetun

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	z "github.com/Oudwins/zog"
)

func stringEnumSchema[T ~string](values []T) *z.StringSchema[T] {
	return z.StringLike[T]().OneOf(values)
}

func intEnumSchema[T ~int](values []T) *z.NumberSchema[T] {
	return z.IntLike[T]().OneOf(values)
}

func floatEnumSchema[T ~float64](values []T) *z.NumberSchema[T] {
	return z.FloatLike[T]().OneOf(values)
}

func isValidListeningAddress[T ~string](data *T, ctx z.Ctx) bool {
	if len(*data) == 0 || (*data)[0] != ':' {
		ctx.AddIssue(ctx.Issue().SetMessage("Listening Addess needs to start with ':'"))
		return false
	}

	port, err := strconv.Atoi(string((*data)[1:]))
	if err != nil || port < 0 || port > 65535 {
		ctx.AddIssue(ctx.Issue().SetMessage("Listening Address needs to be a valid port number"))
		return false
	}

	return true
}

func isValidFilePath[T ~string](data *T, ctx z.Ctx) bool {
	p := strings.TrimSpace(string(*data))
	if p == "" {
		ctx.AddIssue(ctx.Issue().SetMessage("file path cannot be empty"))
		return false
	}

	if strings.ContainsRune(p, '\x00') {
		ctx.AddIssue(ctx.Issue().SetMessage("file path contains invalid null byte"))
		return false
	}

	abs, err := filepath.Abs(filepath.Clean(p))
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage("file path is invalid for this OS"))
		return false
	}

	info, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ctx.AddIssue(ctx.Issue().SetMessage("file does not exist"))
		} else {
			ctx.AddIssue(ctx.Issue().SetMessage("file path is not accessible"))
		}
		return false
	}

	if info.IsDir() {
		ctx.AddIssue(ctx.Issue().SetMessage("expected a file path, got a directory"))
		return false
	}

	return true
}

type timePeriodSuffix rune

const (
	timePeriodSuffixSecond timePeriodSuffix = 's'
	timePeriodSuffixMinute timePeriodSuffix = 'm'
	timePeriodSuffixHour   timePeriodSuffix = 'h'
)

func isValidTimePeriod[T ~string](data *T, ctx z.Ctx) bool {
	s := string(*data)
	if len(s) < 2 {
		ctx.AddIssue(ctx.Issue().SetMessage("Time period must contain a valid integer before the time period suffix. Valid suffixes are s, m, h"))
		return false
	}

	period, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage("Time period must contain a valid integer before the time period suffix. Valid suffixes are s, m, h"))
		return false
	}

	suffix := timePeriodSuffix(s[len(s)-1])

	switch suffix {
	case timePeriodSuffixSecond, timePeriodSuffixMinute:
		if period < 0 || period > 60 {
			ctx.AddIssue(ctx.Issue().SetMessage("Time period must be between 0 and 60 seconds"))
			return false
		}
	case timePeriodSuffixHour:
		if period < 0 || period > 24 {
			ctx.AddIssue(ctx.Issue().SetMessage("Time period must be between 0 and 24 hours"))
			return false
		}
	default:
		ctx.AddIssue(ctx.Issue().SetMessage("Time period must contain a valid integer before the time period suffix. Valid suffixes are s, m, h"))
		return false
	}

	return true
}

func isValidHostname[T ~string](data *T, ctx z.Ctx) bool {
	s := strings.TrimSpace(string(*data))
	s = strings.TrimSuffix(s, ".")
	if len(s) == 0 || len(s) > 253 {
		ctx.AddIssue(ctx.Issue().SetMessage("invalid hostname length"))
		return false
	}

	labels := strings.Split(s, ".")
	if len(labels) < 2 {
		ctx.AddIssue(ctx.Issue().SetMessage("hostname must contain at least one dot"))
		return false
	}

	for _, l := range labels {
		if len(l) == 0 || len(l) > 63 {
			ctx.AddIssue(ctx.Issue().SetMessage("invalid hostname label length"))
			return false
		}
		if l[0] == '-' || l[len(l)-1] == '-' {
			ctx.AddIssue(ctx.Issue().SetMessage("hostname labels cannot start or end with '-'"))
			return false
		}
		for _, r := range l {
			if !(r == '-' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
				ctx.AddIssue(ctx.Issue().SetMessage("hostname contains invalid characters"))
				return false
			}
		}
	}

	return true
}

func isValidPlainAddress[T ~string](data *T, ctx z.Ctx) bool {
	s := strings.TrimSpace(string(*data))
	if len(s) == 0 {
		ctx.AddIssue(ctx.Issue().SetMessage("plain address cannot be empty"))
		return false
	}

	labels := strings.Split(s, ":")
	if len(labels) != 2 {
		ctx.AddIssue(ctx.Issue().SetMessage("plain address must contain a port"))
		return false
	}

	host := labels[0]
	if len(host) == 0 {
		ctx.AddIssue(ctx.Issue().SetMessage("plain address host cannot be empty"))
		return false
	}
	if len(host) > 253 {
		ctx.AddIssue(ctx.Issue().SetMessage("plain address host cannot exceed 253 characters"))
		return false
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("plain address must start with a valid IPv4 address: %s", host)))
		return false
	}
	if !addr.Is4() {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("plain address must start with a valid IPv4 address: %s", host)))
		return false
	}

	portStr := labels[1]
	if _, err := strconv.Atoi(portStr); err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("invalid plain address port: %s", portStr)))
		return false
	}

	return true
}

func isValidSubnet[T ~string](data *T, ctx z.Ctx) bool {
	s := strings.TrimSpace(string(*data))
	if len(s) == 0 {
		ctx.AddIssue(ctx.Issue().SetMessage("subnet cannot be empty"))
		return false
	}

	labels := strings.Split(s, "/")
	if len(labels) != 2 {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("invalid subnet: %s", s)))
		return false
	}

	addr, err := netip.ParseAddr(labels[0])
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("invalid subnet address: %s", labels[0])))
		return false
	}
	if !addr.Is4() {
		ctx.AddIssue(ctx.Issue().SetMessage("subnet must be an IPv4 address"))
		return false
	}

	prefix, err := strconv.Atoi(labels[1])
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage("invalid subnet prefix"))
		return false
	}
	if prefix < 0 || prefix > 32 {
		ctx.AddIssue(ctx.Issue().SetMessage("subnet prefix must be between 0 and 32"))
		return false
	}

	return true
}

func isValidCIDRAddress[T ~string](data *T, ctx z.Ctx) bool {
	s := strings.TrimSpace(string(*data))
	if len(s) == 0 {
		ctx.AddIssue(ctx.Issue().SetMessage("address cannot be empty"))
		return false
	}

	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("invalid CIDR address: %s", s)))
		return false
	}

	addr := prefix.Addr()
	if !addr.Is4() {
		ctx.AddIssue(ctx.Issue().SetMessage("address must be an IPv4 CIDR"))
		return false
	}

	return true
}

func isValidBase64PEM(data *string, ctx z.Ctx) bool {
	s := strings.TrimSpace(*data)
	if s == "" {
		ctx.AddIssue(ctx.Issue().SetMessage("value cannot be empty"))
		return false
	}

	blocks, ok := decodePEMBlocks([]byte(s))
	if !ok {
		compact := strings.Join(strings.Fields(s), "")

		decoded, err := base64.StdEncoding.DecodeString(compact)
		if err != nil {
			decoded, err = base64.RawStdEncoding.DecodeString(compact)
			if err != nil {
				ctx.AddIssue(ctx.Issue().SetMessage("value must be PEM or base64-encoded PEM"))
				return false
			}
		}

		blocks, ok = decodePEMBlocks(decoded)
		if !ok {
			ctx.AddIssue(ctx.Issue().SetMessage("invalid PEM payload"))
			return false
		}
	}

	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			if _, err := x509.ParseCertificate(block.Bytes); err != nil {
				ctx.AddIssue(ctx.Issue().SetMessage("invalid certificate PEM block"))
				return false
			}
		case "PRIVATE KEY":
			if _, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
				ctx.AddIssue(ctx.Issue().SetMessage("invalid PKCS#8 private key PEM block"))
				return false
			}
		case "RSA PRIVATE KEY":
			if isLegacyEncryptedPEMBlock(block) {
				ctx.AddIssue(ctx.Issue().SetMessage("legacy encrypted RSA private keys are not supported; use PKCS#8 ENCRYPTED PRIVATE KEY"))
				return false
			}
			if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
				ctx.AddIssue(ctx.Issue().SetMessage("invalid RSA private key PEM block"))
				return false
			}
		case "EC PRIVATE KEY":
			if isLegacyEncryptedPEMBlock(block) {
				ctx.AddIssue(ctx.Issue().SetMessage("legacy encrypted EC private keys are not supported; use PKCS#8 ENCRYPTED PRIVATE KEY"))
				return false
			}
			if _, err := x509.ParseECPrivateKey(block.Bytes); err != nil {
				ctx.AddIssue(ctx.Issue().SetMessage("invalid EC private key PEM block"))
				return false
			}
		case "ENCRYPTED PRIVATE KEY":
			if len(block.Bytes) == 0 {
				ctx.AddIssue(ctx.Issue().SetMessage("invalid encrypted private key PEM block"))
				return false
			}
		default:
			ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("unsupported PEM block type: %s", block.Type)))
			return false
		}
	}

	return true
}

func decodePEMBlocks(data []byte) ([]*pem.Block, bool) {
	rest := bytes.TrimSpace(data)
	blocks := make([]*pem.Block, 0, 1)

	for len(rest) > 0 {
		block, next := pem.Decode(rest)
		if block == nil {
			return nil, false
		}

		blocks = append(blocks, block)
		rest = bytes.TrimSpace(next)
	}

	return blocks, len(blocks) > 0
}

func isLegacyEncryptedPEMBlock(block *pem.Block) bool {
	if block == nil {
		return false
	}

	if _, ok := block.Headers["DEK-Info"]; ok {
		return true
	}

	if procType, ok := block.Headers["Proc-Type"]; ok {
		return strings.Contains(procType, "ENCRYPTED")
	}

	return false
}
