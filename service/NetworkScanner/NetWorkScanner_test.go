package NetworkScanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testpacket = []byte{10, 56, 46, 49, 46, 48, 0, 11, 0, 0, 0, 103, 61, 104, 12, 55, 98, 104, 66, 0, 255, 255, 255, 2, 0, 255}

func TestPacketGeneration(t *testing.T) {

	value := GenerateSQLHandshakePacket(testpacket)

	assert.Equal(t, value.MySQLProtocol, 10, "expected version 10")
	assert.Equal(t, value.MySQLBanner, "", "expected empty banner")
	assert.Equal(t, len(testpacket), 26)
}

func TestNewNetworkScanner(t *testing.T) {
	myNetworkScanner := NewNetworkScanner("localhost", 8889)

	result := myNetworkScanner.FindMySqlInstance()
	assert.Equal(t, false, result.MySQLActive, "expected false return due to lack of service running")

}
