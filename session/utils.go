package session

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// 获取ip
func GetIP(r *http.Request) string {
	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	logrus.Error("no ip!!!")
	return ip
}

// 获取域名
func GetDomain(host string) string {
	host = strings.TrimSpace(host)

	hostPorts := strings.Split(host, ":")

	hostParts := strings.Split(hostPorts[0], ".")
	lengthOfHostParts := len(hostParts)

	if lengthOfHostParts == 1 {
		return hostParts[0] // scenario C
	} else {
		_, err := strconv.ParseFloat(hostParts[0], 64)
		if err == nil {
			return host
		} else {
			return strings.Join([]string{hostParts[lengthOfHostParts-2], hostParts[lengthOfHostParts-1]}, ".")
		}
	}
}
