package IPSpaceManager

import (
	"net"
)

type IPManager struct {
	network net.IPNet
	gateway net.IP
	dnsA net.IP
	dnsB net.IP
	allocatedIPs []net.IP
}

func New(cidr string) (*IPManager) {
	x := new(IPManager)
	if(x.claimIP(x.GetGateway()) == false) {
		return nil
	}
	return x
}

func (x* IPManager) SetGateway(gw string) bool {

	return true
}

func (x* IPManager) SetDNS(a, b string) bool {

}

func (x* IPManager) GetGateway() net.IP {
	return ""
}

func (x* IPManager) getFreeIP() (net.IP, error) {

}

func (x* IPManager) claimIP(ip net.IP) bool {

}

func (x* IPManager) AllocateIP() (net.IP, error) {
	var ip net.IP
	ipx, err := x.getFreeIP()
	if err != nil {
		return ip, errors.New("No free IPs")
	}
	if(x.claimIP(ipx) == false) {
		return ip, errors.New("Failed to allocate IP")
	}
	x.allocatedIPs = append(x.allocatedIPs, ip)
	return ip, nil
}

func (x* IPManager) FreeIP(ip net.IP) bool {
	return true
}
