package NetworkScanner

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	mysqlProtocolVersion = 10

	MySQLHandshakePacketLength = 30
)

type ScannerInterface interface {
	FindMySQLInstance() bool
}

type NetworkScanner struct {
	host string
	port int
}

// SQLInfo is an information struct, automatically returns wether or not MySQL is active and fills with as much information
// as is available.
type SQLInfo struct {
	OriginalHost string `default:"" json:"original_host,omitempty"`
	OriginalPort int    `json:"original_port,omitempty"`

	MySQLBanner          string `default:"null" json:"my_sql_banner,omitempty"`
	MySQLServerVersion   string `default:"null" json:"my_sql_version,omitempty"`
	MySQLProtocolVersion int    `json:"my_sql_protocol_version,omitempty"`

	MySQLActive bool `json:"my_sql_active,omitempty"`
}

// MySQLHandshakePacket is a handshake packet based off of the handshake packet expected in the docs
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_v10.html
type MySQLHandshakePacket struct {
	MySQLProtocol   int
	MySQLVersion    string
	ConnectionID    int
	AuthPluginData  []byte
	Flags           int
	CharacterSet    int
	StatusFlags     int
	CapabilityFlags int
	AuthPluginName  string
	MySQLBanner     string
}

// NewNetworkScanner generates a new NetworkScanner on a chosen port and host
func NewNetworkScanner(host string, port int) *NetworkScanner {
	return &NetworkScanner{
		host: host,
		port: port,
	}
}

// extractHandshake recives the correct size of information, grabbing the mysqlprotocolversion from the 4th index
func extractHandshake(data []byte) ([]byte, error) {
	// Find the index of the uint 10 in the data
	uint10Index := bytes.IndexByte(data, mysqlProtocolVersion)
	// Check if the uint 10 is present in the data.
	if uint10Index == -1 {
		return nil, fmt.Errorf("The uint 10 is not present in the data")
	}

	// Extract the total byte length starting from the first character.
	extractedBytes := data[uint10Index:]

	return extractedBytes, nil
}

func (n *NetworkScanner) FindMySqlInstance() *SQLInfo {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.host, n.port))
	if err != nil {
		fmt.Println(err)
		return &SQLInfo{
			MySQLActive: false,
		}
	}
	reader := bufio.NewReader(conn)

	// Read the entire handshake packet, including the server banner field.
	handshakePacket := make([]byte, MySQLHandshakePacketLength)
	_, err = reader.Read(handshakePacket)
	if err != nil {
		fmt.Println(err)

	}

	newData, err := extractHandshake(handshakePacket)
	if err != nil {
		fmt.Println(err)
		return &SQLInfo{
			MySQLActive: false,
		}
	}
	Packet := GenerateSQLHandshakePacket(newData)

	if Packet.MySQLVersion != "" || Packet.MySQLBanner != "" {
		return &SQLInfo{
			OriginalHost:         n.host,
			OriginalPort:         n.port,
			MySQLBanner:          Packet.MySQLBanner,
			MySQLProtocolVersion: Packet.MySQLProtocol,
			MySQLServerVersion:   Packet.MySQLVersion,
			MySQLActive:          true,
		}
	}
	return &SQLInfo{
		MySQLActive: false,
	}
}

// GenerateSQLHandshakePacket uses the protocol for handshake from mysql v10 to build a handshake packet that is sent to our TCP connection on connect
func GenerateSQLHandshakePacket(data []byte) *MySQLHandshakePacket {
	handshakePacket := &MySQLHandshakePacket{}

	// Set the protocol version.
	handshakePacket.MySQLProtocol = int(data[0])

	// Set the server version.
	serverVersion := string(data[1:5])
	handshakePacket.MySQLVersion = serverVersion

	// Set the connection ID.
	connectionID := binary.LittleEndian.Uint32(data[5:9])
	handshakePacket.ConnectionID = int(connectionID)

	// Set the auth plugin data.
	authPluginData := data[9:15]
	handshakePacket.AuthPluginData = authPluginData

	// Set the flags.
	flags := binary.LittleEndian.Uint16(data[15:17])
	handshakePacket.Flags = int(flags)

	// Set the character set.
	characterSet := int(data[17])
	handshakePacket.CharacterSet = characterSet

	// Set the status flags.
	statusFlags := binary.LittleEndian.Uint16(data[18:20])
	handshakePacket.StatusFlags = int(statusFlags)

	// Set the capability flags.
	capabilityFlags := binary.LittleEndian.Uint32(data[20:24])
	handshakePacket.CapabilityFlags = int(capabilityFlags)

	// Set the auth plugin name.
	authPluginName := string(data[24:26])
	handshakePacket.AuthPluginName = authPluginName

	// Set the server banner.
	serverBanner := string(data[26:])
	handshakePacket.MySQLBanner = serverBanner

	return handshakePacket
}
