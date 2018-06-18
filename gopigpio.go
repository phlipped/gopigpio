package gopigpio

import (
	"net"
)

type Pigpio struct {
  conn net.Conn
}

func New(network, address string) (Pigpio, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return Pigpio{}, err
	}

	return Pigpio{conn: conn}, nil
}

func (p *Pigpio) Close() error {
	return p.conn.Close()
}
