package util

import (
	"fmt"
	"net"
	"strings"
)

// privateBlocks contains non-forwardable address blocks which are used
// for private networks. RFC 6890 provides an overview of special
// address blocks.
var privateBlocks = []*net.IPNet{
	parseCIDR("10.0.0.0/8"),     // RFC 1918 IPv4 private network address
	parseCIDR("100.64.0.0/10"),  // RFC 6598 IPv4 shared address space
	parseCIDR("127.0.0.0/8"),    // RFC 1122 IPv4 loopback address
	parseCIDR("169.254.0.0/16"), // RFC 3927 IPv4 link local address
	parseCIDR("172.16.0.0/12"),  // RFC 1918 IPv4 private network address
	parseCIDR("192.0.0.0/24"),   // RFC 6890 IPv4 IANA address
	parseCIDR("192.0.2.0/24"),   // RFC 5737 IPv4 documentation address
	parseCIDR("192.168.0.0/16"), // RFC 1918 IPv4 private network address
	parseCIDR("::1/128"),        // RFC 1884 IPv6 loopback address
	parseCIDR("fe80::/10"),      // RFC 4291 IPv6 link local addresses
	parseCIDR("fc00::/7"),       // RFC 4193 IPv6 unique local addresses
	parseCIDR("fec0::/10"),      // RFC 1884 IPv6 site-local addresses
	parseCIDR("2001:db8::/32"),  // RFC 3849 IPv6 documentation address
}

func parseCIDR(s string) *net.IPNet {
	_, block, err := net.ParseCIDR(s)
	if err != nil {
		panic(fmt.Sprintf("Bad CIDR %s: %s", s, err))
	}
	return block
}

// 检测是否是私有ip
func isPrivate(ip net.IP) bool {
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

func getIPNets() (ips []*net.IPNet) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			ips = append(ips, ipNet)
		}
	}
	return ips
}

// 获取ipv4地址
func GetPrivateIPv4s() (ips []string) {
	nets := getIPNets()
	if len(nets) == 0 {
		return
	}
	for _, oneNet := range nets {
		if oneNet.IP.To4() != nil && isPrivate(oneNet.IP) {
			ips = append(ips, oneNet.IP.String())
		}
	}
	return
}

// 获取ipv4地址
func GetPrivateIPv4() string {
	list := GetPrivateIPv4s()
	if len(list) == 0 {
		return ""
	}
	for _, l := range list {
		if strings.HasPrefix(l, "10.") {
			return l
		}
	}
	for _, l := range list {
		if strings.HasPrefix(l, "192.168.") {
			return l
		}
	}
	return list[0]
}

// 获取一个mac地址
func GetMacAddr() string {
	addrs := GetMacAddrs()
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0]
}

// 获取mac地址
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}
