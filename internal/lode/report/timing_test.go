package report

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"net/http/httptrace"
	"testing"
	"time"
)

var timing = Timing{
	DnsStart:     time.Unix(0, 1),
	DnsDone:      time.Unix(0, 3),
	ConnectStart: time.Unix(0, 7),
	ConnectDone:  time.Unix(0, 13),
	TlsStart:     time.Unix(0, 21),
	TlsDone:      time.Unix(0, 31),
	GotConn:      time.Unix(0, 33),
	FirstByte:    time.Unix(0, 47),
	Done:         time.Unix(0, 63),
}

func TestTiming_DnsLookupDuration(t *testing.T) {
	assert.Equal(t, time.Duration(2), timing.DnsLookupDuration())
}

func TestTiming_TcpConnectDuration(t *testing.T) {
	assert.Equal(t, time.Duration(10), timing.TcpConnectDuration())

	withoutDnsLookup := timing
	withoutDnsLookup.DnsDone = time.Time{} // did not do DNS lookup
	assert.Equal(t, time.Duration(6), withoutDnsLookup.TcpConnectDuration())
}

func TestTiming_TlsHandshakeDuration(t *testing.T) {
	assert.Equal(t, time.Duration(10), timing.TlsHandshakeDuration())
}

func TestTiming_ServerDuration(t *testing.T) {
	assert.Equal(t, time.Duration(14), timing.ServerDuration())
}

func TestTiming_ResponseTransferDuration(t *testing.T) {
	assert.Equal(t, time.Duration(16), timing.ResponseTransferDuration())
}

func TestTiming_TotalDuration(t *testing.T) {
	assert.Equal(t, time.Duration(62), timing.TotalDuration())

	withoutDnsLookup := timing
	withoutDnsLookup.DnsDone = time.Time{} // did not do DNS lookup
	assert.Equal(t, time.Duration(56), withoutDnsLookup.TotalDuration())
}

func TestNewTrace(t *testing.T) {
	assert := assert.New(t)
	timing = Timing{}

	trace := NewTrace(&timing)

	trace.DNSStart(httptrace.DNSStartInfo{})
	assert.NotZero(timing.DnsStart)
	trace.DNSDone(httptrace.DNSDoneInfo{})
	assert.NotZero(timing.DnsDone)
	trace.ConnectStart("", "")
	assert.NotZero(timing.ConnectStart)
	trace.ConnectDone("", "", nil)
	assert.NotZero(timing.ConnectDone)
	trace.GotConn(httptrace.GotConnInfo{})
	assert.NotZero(timing.GotConn)
	trace.GotFirstResponseByte()
	assert.NotZero(timing.FirstByte)
	trace.TLSHandshakeStart()
	assert.NotZero(timing.TlsStart)
	trace.TLSHandshakeDone(tls.ConnectionState{}, nil)
	assert.NotZero(timing.TlsDone)
}
