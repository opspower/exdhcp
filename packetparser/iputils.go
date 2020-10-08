package main

import (
	"fmt"
	"net"
	"time"
)

func GetLinkLocalAddr(ifname string) (*net.IP, *net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	var iface *net.Interface
	var linkLocalAddr *net.IP
	for _, ifi := range ifaces {
		if ifi.Name == ifname {
			iface = &ifi
			break
		}
	}
	// build the addr from the interface
	hwa := iface.HardwareAddr
	linkLocalAddr = &net.IP{
		0xfe, 0x80, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		hwa[0] ^ 2, hwa[1], hwa[2], 0xff,
		0xfe, hwa[3], hwa[4], hwa[5],
	}
	m := net.CIDRMask(64, 128)
	linkLocalNet := net.IPNet{IP: linkLocalAddr.Mask(m), Mask: m} // a /64
	return linkLocalAddr, &linkLocalNet, nil
}

// Wait for an interface to be up. Will return an error if it is not up before
// the timeout expires. The status is synchronously polled every 100msec
func WaitForInterfaceStatusUp(ifname string, timeout time.Duration) error {
	// FIXME should use netlink events rather than polling like this
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Timed out while waiting for interface to be up")
		}
		for _, ifi := range ifaces {
			if ifi.Flags&net.FlagUp != 0 {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
