package sinutil

// TODO(sg): move file to right place
import (
	"errors"
	"net"
)

var (
	privateIPBlocks []*net.IPNet
	errInvalidAddr  = errors.New("invalid address")
)

func init() {
	for _, cidr := range []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	} {
		_, block, _ := net.ParseCIDR(cidr)
		if block != nil {
			privateIPBlocks = append(privateIPBlocks, block)
		}
	}
}

// IsPrivateIP check if the address is under private CIDR blocks.
// reference:
// * https://en.wikipedia.org/wiki/Private_network
// * https://en.wikipedia.org/wiki/Link-local_address
func IsPrivateIP(addr string) (bool, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false, errInvalidAddr
	}

	for _, cidr := range privateIPBlocks {
		if cidr.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}
