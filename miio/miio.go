package miio

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/fffonion/mi-vacuum-exporter/miio/packet"
	"github.com/mitchellh/mapstructure"
)

type MiioClient struct {
	addr   string
	id     uint32
	token  []byte
	crypto packet.Crypto
}

type MiioClientConfig struct {
	Host  string
	Token string
}

type rpcResponse struct {
	Result []interface{} `json:"result"`
	ID     int           `json:"id"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func New(c *MiioClientConfig) (*MiioClient, error) {
	tokenBytes := make([]byte, 16)
	_, err := hex.Decode(tokenBytes, []byte(c.Token))
	if err != nil {
		return nil, err
	}
	return &MiioClient{
		addr:  c.Host + ":54321",
		token: tokenBytes,
	}, nil
}

func (c *MiioClient) Init() error {
	p, err := c.request(packet.NewHello())
	if err != nil {
		return err
	}

	crypto, err := packet.NewCrypto(p.Header.DeviceID, c.token,
		p.Header.Stamp, time.Now().UTC(), clock.New())
	if err != nil {
		return err
	}
	c.crypto = crypto
	c.id = p.Header.DeviceID
	return nil
}

func (c *MiioClient) RPC(method string, params []interface{}, out interface{}) error {
	if c.crypto == nil {
		return fmt.Errorf("client is not initialized, call Init() first")
	}

	payload := &deviceCommand{
		ID:     time.Now().UTC().Unix(),
		Method: method,
		Params: params,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	p, err := c.crypto.NewPacket(payloadBytes)
	if err != nil {
		return err
	}

	p, err = c.request(p)
	if err != nil {
		return err
	}

	err = p.Verify(c.token)
	if err != nil {
		return err
	}

	dec, err := c.crypto.Decrypt(p.Data)
	if err != nil {
		return err
	}

	response := &rpcResponse{}
	err = json.Unmarshal(dec, response)
	if err != nil {
		return err
	}
	if response.Error.Code != 0 {
		return fmt.Errorf("rpc error: %v", response.Error)
	}

	if response.Result == nil || len(response.Result) != 1 {
		return fmt.Errorf("none or more than one result found in packet")
	}

	mapstructure.Decode(response.Result[0], &out)
	return nil

}

func (c *MiioClient) request(payload *packet.Packet) (*packet.Packet, error) {
	conn, err := net.Dial("udp", c.addr)
	defer conn.Close()

	if err != nil {
		return nil, err
	}
	// the tcp server on smart plug can only handle one connection
	// at a time, make sure we don't wait too long to block others
	conn.SetWriteDeadline(time.Now().Add(time.Second))

	_, err = conn.Write(payload.Serialize())
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, 1024)

	var p *packet.Packet
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	ti := time.NewTicker(time.Second)
	defer ti.Stop()

	n, err := conn.Read(tmp)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, fmt.Errorf("Read returned %d bytes", n)
	}
	p, err = packet.Decode(tmp, nil)
	if err != nil {
		log.Fatal(err)
	}
	return p, nil
}

func (c *MiioClient) ID() string {
	return fmt.Sprintf("%x", c.id)
}
