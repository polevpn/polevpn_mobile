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

func GetRouteIpsFromDomain(domainStr string) string {

	domains := strings.Split(domainStr, "\n")

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
	return strings.Join(ips, "\n")
}

func GetLocalIP() string {

	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	wifiIp := ""
	cellerIp := ""

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return ""
		}
		addresses, err := byName.Addrs()
		if err != nil {
			return ""
		}
		for _, v := range addresses {
			if strings.Contains(v.String(), ":") {
				continue
			}

			if strings.Contains(v.String(), "169.") || strings.Contains(v.String(), "127.") {
				continue
			}

			ip, _, err := net.ParseCIDR(v.String())
			if err != nil {
				continue
			}

			if i.Name == "pdp_ip0" {
				cellerIp = ip.String()
			}

			if i.Name == "en0" {
				wifiIp = ip.String()
			}

		}
	}

	if wifiIp != "" {
		return wifiIp
	} else {
		return cellerIp
	}
}

func GetSubNetMask(cidr string) string {

	_, network, err := net.ParseCIDR(cidr)

	if err != nil {
		return ""
	}
	mask := net.IPv4(network.Mask[0], network.Mask[1], network.Mask[2], network.Mask[3])
	return mask.String()
}

func Log(level string, msg string) {
	if level == "debug" {
		plog.Debug(msg)
	} else if level == "info" {
		plog.Info(msg)
	} else if level == "warn" {
		plog.Warn(msg)
	} else if level == "error" {
		plog.Error(msg)
	}
	plog.Flush()
}
