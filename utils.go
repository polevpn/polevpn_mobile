package polevpnmobile

import (
	"net"
	"os"
	"strings"
	"time"
)

func getAppName() string {
	appName := os.Args[0]
	slash := strings.LastIndex(appName, string(os.PathSeparator))
	if slash >= 0 {
		appName = appName[slash+1:]
	}
	return appName
}

func getTimeNowDate() string {
	return time.Now().Format("2006-01-02")
}

func GetRouteIpsFromDomain(domains []string) []string {

	ips := make([]string, 0)
	for _, domain := range domains {

		netips, err := net.LookupIP(domain)
		if err != nil {
			continue
		}
		for _, netip := range netips {
			if !strings.Contains(netip.String(), ":") {
				ips = append(ips, netip.String()+"/32")
			}
		}
	}
	return ips
}
