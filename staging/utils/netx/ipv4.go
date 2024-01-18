package netx

import (
	"net"
	"reflect"
)

type IPv4 net.IPNet

func MustIPv4FromCIDR(cidr string) IPv4 {
	ip, err := IPv4FromCIDR(cidr)
	if err != nil {
		panic(err)
	}

	return ip
}

func IPv4FromCIDR(cidr string) (IPv4, error) {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return IPv4{}, err
	}

	// NB(thxCode): returns result that defined in RFC 4632 and RFC 4291.
	return IPv4FromIPNet(*n), nil
}

func IPv4FromIP(ip net.IP) IPv4 {
	return IPv4(net.IPNet{
		IP:   ip.To4(),
		Mask: ip.DefaultMask(),
	})
}

func IPv4FromIPNet(n net.IPNet) IPv4 {
	return IPv4(net.IPNet{
		IP:   n.IP.To4().Mask(n.Mask),
		Mask: n.Mask,
	})
}

// Contains returns true if the given net.IP is contained in the IPv4.
func (l IPv4) Contains(ip net.IP) bool {
	return (*net.IPNet)(&l).Contains(ip)
}

// String returns the string representation of the IPv4.
func (l IPv4) String() string {
	return (*net.IPNet)(&l).String()
}

// IPNet returns the net.IPNet of the IPv4.
func (l IPv4) IPNet() net.IPNet {
	return net.IPNet(l)
}

// Equal returns true if two IPv4 are equal.
func (l IPv4) Equal(r IPv4) bool {
	return reflect.DeepEqual(l.Mask, r.Mask) && (*net.IPNet)(&l).Contains(r.IP) && (*net.IPNet)(&r).Contains(l.IP)
}

// Overlap returns true if two IPv4 overlap.
func (l IPv4) Overlap(r IPv4) bool {
	return (*net.IPNet)(&l).Contains(r.IP) || (*net.IPNet)(&r).Contains(l.IP)
}

// Next returns a IPv4 with the same mask.
//
// For example,
// returns 172.16.192.0/18 if given 172.16.128.0/18,
// returns 172.17.0.0/18 if given 172.16.192.0/18.
func (l IPv4) Next() IPv4 {
	var (
		o, s = l.Mask.Size()
		i, j int
	)
	{
		for i, j = range []int{24, 16, 8, 0} {
			if o > j {
				break
			}
		}
		i = 3 - i
		j = 8 + j - o
	}

	r := IPv4{
		IP:   []byte{l.IP[0], l.IP[1], l.IP[2], l.IP[3]},
		Mask: net.CIDRMask(o, s),
	}
	r.IP[i] += 1 << j

	if r.IP[i] == 0 && i > 0 {
		r.IP[i-1] += 1
	}

	return r
}
