package natpmp

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// Implement the NAT-PMP protocol, typically supported by Apple routers and open source
// routers such as DD-WRT and Tomato.
//
// See https://tools.ietf.org/rfc/rfc6886.txt
//
// Usage:
//
//    client := natpmp.NewClient(gatewayIP)
//    response, err := client.GetExternalAddress()

// The recommended mapping lifetime for AddPortMapping.
const RECOMMENDED_MAPPING_LIFETIME_SECONDS = 3600

type Option func(*Client)

func OptionDialer(d dialer) Option {
	return func(c *Client) {
		c.dialer = d
	}
}

func OptionTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

// Client is a NAT-PMP protocol client.
type Client struct {
	gateway netip.Addr
	dialer  dialer
	timeout time.Duration
}

// Create a NAT-PMP client for the NAT-PMP server at the gateway.
// Uses default timeout which is around 128 seconds.
func NewClient(gateway netip.Addr, options ...Option) (nat *Client) {
	return langx.Autoptr(langx.Clone(Client{
		timeout: 0,
		gateway: gateway,
		dialer:  &net.Dialer{},
	}, options...))
}

// Results of the NAT-PMP GetExternalAddress operation.
type GetExternalAddressResult struct {
	SecondsSinceStartOfEpoc uint32
	ExternalIPAddress       [4]byte
}

// Get the external address of the router.
//
// Note that this call can take up to 128 seconds to return.
func (n *Client) GetExternalAddress() (result *GetExternalAddressResult, err error) {
	msg := make([]byte, 2)
	msg[0] = 0 // Version 0
	msg[1] = 0 // OP Code 0
	response, err := n.rpc(msg, 12)
	if err != nil {
		return
	}
	result = &GetExternalAddressResult{}
	result.SecondsSinceStartOfEpoc = readNetworkOrderUint32(response[4:8])
	copy(result.ExternalIPAddress[:], response[8:12])
	return
}

// Results of the NAT-PMP AddPortMapping operation
type AddPortMappingResult struct {
	SecondsSinceStartOfEpoc      uint32
	InternalPort                 uint16
	MappedExternalPort           uint16
	PortMappingLifetimeInSeconds uint32
}

// Add (or delete) a port mapping. To delete a mapping, set the requestedExternalPort and lifetime to 0.
// Note that this call can take up to 128 seconds to return.
func (n *Client) AddPortMapping(protocol string, internalPort, requestedExternalPort int, lifetime int) (result *AddPortMappingResult, err error) {
	var opcode byte
	switch protocol {
	case "udp":
		opcode = 1
	case "tcp":
		opcode = 2
	default:
		err = fmt.Errorf("unknown protocol %v", protocol)
		return
	}
	msg := make([]byte, 12)
	msg[0] = 0 // Version 0
	msg[1] = opcode

	// [2:3] is reserved.
	writeNetworkOrderUint16(msg[4:6], uint16(internalPort))
	writeNetworkOrderUint16(msg[6:8], uint16(requestedExternalPort))
	writeNetworkOrderUint32(msg[8:12], uint32(lifetime))
	response, err := n.rpc(msg, 16)
	if err != nil {
		return nil, err
	}
	result = &AddPortMappingResult{}
	result.SecondsSinceStartOfEpoc = readNetworkOrderUint32(response[4:8])
	result.InternalPort = readNetworkOrderUint16(response[8:10])
	result.MappedExternalPort = readNetworkOrderUint16(response[10:12])
	result.PortMappingLifetimeInSeconds = readNetworkOrderUint32(response[12:16])
	return result, nil
}

func (n *Client) rpc(msg []byte, resultSize int) (result []byte, err error) {
	c := network{
		dialer:  n.dialer,
		gateway: n.gateway,
	}
	result, err = c.call(context.Background(), msg, n.timeout)
	if err != nil {
		return nil, err
	}
	return result, protocolChecks(msg, resultSize, result)
}

func protocolChecks(msg []byte, resultSize int, result []byte) (err error) {
	if len(result) != resultSize {
		return fmt.Errorf("unexpected result size %d, expected %d", len(result), resultSize)
	}
	if result[0] != 0 {
		return fmt.Errorf("unknown protocol version %d", result[0])
	}
	expectedOp := msg[1] | 0x80
	if result[1] != expectedOp {
		return fmt.Errorf("unexpected opcode %d. Expected %d", result[1], expectedOp)
	}
	resultCode := readNetworkOrderUint16(result[2:4])
	if resultCode != 0 {
		return fmt.Errorf("non-zero result code %d", resultCode)
	}

	// If we got here the RPC is good.
	return nil
}

func writeNetworkOrderUint16(buf []byte, d uint16) {
	buf[0] = byte(d >> 8)
	buf[1] = byte(d)
}

func writeNetworkOrderUint32(buf []byte, d uint32) {
	buf[0] = byte(d >> 24)
	buf[1] = byte(d >> 16)
	buf[2] = byte(d >> 8)
	buf[3] = byte(d)
}

func readNetworkOrderUint16(buf []byte) uint16 {
	return (uint16(buf[0]) << 8) | uint16(buf[1])
}

func readNetworkOrderUint32(buf []byte) uint32 {
	return (uint32(buf[0]) << 24) | (uint32(buf[1]) << 16) | (uint32(buf[2]) << 8) | uint32(buf[3])
}
