package miio

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/fffonion/mi-vacuum-exporter/miio/packet"
)

func Discovery(timeout time.Duration) ([]*MiioClientConfig, error) {
	heloBytes := make([]byte, 32)
	hex.Decode(
		heloBytes,
		[]byte("21310020ffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	)

	ip := net.ParseIP("255.255.255.255")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 54321}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		return nil, err
	}
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	_, err = conn.WriteToUDP(heloBytes, dstAddr)
	if err != nil {
		return nil, err
	}

	ta := time.After(timeout)
	tmp := make([]byte, 1024)
	ret := []*MiioClientConfig{}
	for {
		select {
		case <-ta:
			break
		default:
			conn.SetReadDeadline(time.Now().Add(timeout))
			n, addr, err := conn.ReadFrom(tmp)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return ret, nil
				}
				return nil, err
			}
			if n <= 0 {
				return nil, fmt.Errorf("Read returned %d bytes", n)
			}
			_, err = packet.Decode(tmp, nil)
			if err != nil {
				return nil, err
			}
			addrIP, addrPort, err := net.SplitHostPort(addr.String())
			if err != nil {
				return nil, err
			} else if addrPort != "54321" {
				continue
			}
			ret = append(ret, &MiioClientConfig{
				Host: addrIP,
			})
		}
	}
	return ret, nil
}
