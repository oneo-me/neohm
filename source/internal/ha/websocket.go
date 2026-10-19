package ha

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// maxMessage bounds a message from Home Assistant; the registries of a large
// home take a few megabytes.
const maxMessage = 64 << 20

// WebSocket opcodes.
const (
	opContinue = 0x0
	opText     = 0x1
	opBinary   = 0x2
	opClose    = 0x8
	opPing     = 0x9
	opPong     = 0xA
)

// webSocket is one end of a WebSocket connection (RFC 6455), enough of it to
// exchange JSON messages with Home Assistant. A client masks what it sends.
type webSocket struct {
	r      *bufio.Reader
	conn   io.ReadWriteCloser
	client bool
}

// acceptKey is the answer a server gives to the key of a client's handshake.
func acceptKey(key string) string {
	h := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h[:])
}

// dialWebSocket opens Home Assistant's WebSocket API. net/http's transport
// makes the connection, through any proxy and TLS, and hands it over once
// upgraded; http.Client would hide it behind its timeout. Once open, the
// connection outlives ctx: its user closes it.
func (c *Client) dialWebSocket(ctx context.Context) (*webSocket, error) {
	var nonce [16]byte
	rand.Read(nonce[:])
	key := base64.StdEncoding.EncodeToString(nonce[:])
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/websocket", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", key)
	transport := c.HTTP.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	conn, ok := res.Body.(io.ReadWriteCloser)
	if res.StatusCode != http.StatusSwitchingProtocols || !ok || res.Header.Get("Sec-WebSocket-Accept") != acceptKey(key) {
		res.Body.Close()
		return nil, fmt.Errorf("websocket: %s", res.Status)
	}
	return &webSocket{r: bufio.NewReader(conn), conn: conn, client: true}, nil
}

func (ws *webSocket) Close() error { return ws.conn.Close() }

// write sends one frame of opcode op.
func (ws *webSocket) write(op byte, payload []byte) error {
	frame := make([]byte, 0, 14+len(payload))
	frame = append(frame, 0x80|op)
	var mask byte
	if ws.client {
		mask = 0x80
	}
	switch n := len(payload); {
	case n < 126:
		frame = append(frame, mask|byte(n))
	case n <= 0xFFFF:
		frame = binary.BigEndian.AppendUint16(append(frame, mask|126), uint16(n))
	default:
		frame = binary.BigEndian.AppendUint64(append(frame, mask|127), uint64(n))
	}
	if !ws.client {
		_, err := ws.conn.Write(append(frame, payload...))
		return err
	}
	var key [4]byte
	rand.Read(key[:])
	frame = append(frame, key[:]...)
	for i, b := range payload {
		frame = append(frame, b^key[i%4])
	}
	_, err := ws.conn.Write(frame)
	return err
}

// read returns the next message, answering pings on the way.
func (ws *webSocket) read() ([]byte, error) {
	var msg []byte
	for {
		var head [2]byte
		if _, err := io.ReadFull(ws.r, head[:]); err != nil {
			return nil, err
		}
		fin, op, masked := head[0]&0x80 != 0, head[0]&0x0F, head[1]&0x80 != 0
		n := uint64(head[1] & 0x7F)
		switch n {
		case 126:
			var b [2]byte
			if _, err := io.ReadFull(ws.r, b[:]); err != nil {
				return nil, err
			}
			n = uint64(binary.BigEndian.Uint16(b[:]))
		case 127:
			var b [8]byte
			if _, err := io.ReadFull(ws.r, b[:]); err != nil {
				return nil, err
			}
			n = binary.BigEndian.Uint64(b[:])
		}
		if n > maxMessage-uint64(len(msg)) {
			return nil, errors.New("websocket: message too large")
		}
		var key [4]byte
		if masked {
			if _, err := io.ReadFull(ws.r, key[:]); err != nil {
				return nil, err
			}
		}
		payload := make([]byte, n)
		if _, err := io.ReadFull(ws.r, payload); err != nil {
			return nil, err
		}
		if masked {
			for i := range payload {
				payload[i] ^= key[i%4]
			}
		}
		switch op {
		case opClose:
			return nil, io.EOF
		case opPing:
			if err := ws.write(opPong, payload); err != nil {
				return nil, err
			}
		case opPong:
		case opText, opBinary, opContinue:
			msg = append(msg, payload...)
			if fin {
				return msg, nil
			}
		default:
			return nil, fmt.Errorf("websocket: unknown opcode %d", op)
		}
	}
}

func (ws *webSocket) writeJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return ws.write(opText, data)
}

func (ws *webSocket) readJSON(v any) error {
	data, err := ws.read()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
