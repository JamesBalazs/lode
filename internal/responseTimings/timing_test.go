package responseTimings

import (
	"crypto/tls"
	"github.com/stretchr/testify/assert"
	"net/http/httptrace"
	"testing"
	"time"
)

var timing = Timing{
	DnsStart:     time.Unix(0, 1_000_000),
	DnsDone:      time.Unix(0, 3_000_000),
	ConnectStart: time.Unix(0, 7_000_000),
	ConnectDone:  time.Unix(0, 13_000_000),
	TlsStart:     time.Unix(0, 21_000_000),
	TlsDone:      time.Unix(0, 31_000_000),
	GotConn:      time.Unix(0, 33_000_000),
	FirstByte:    time.Unix(0, 47_000_000),
	Done:         time.Unix(0, 63_000_000),
}

func TestTiming_DnsLookupDuration(t *testing.T) {
	assert.Equal(t, 2*time.Millisecond, timing.DnsLookupDuration())
}

func TestTiming_TcpConnectDuration(t *testing.T) {
	assert.Equal(t, 10*time.Millisecond, timing.TcpConnectDuration())

	withoutDnsLookup := timing
	withoutDnsLookup.DnsDone = time.Time{} // did not do DNS lookup
	assert.Equal(t, 6*time.Millisecond, withoutDnsLookup.TcpConnectDuration())
}

func TestTiming_TlsHandshakeDuration(t *testing.T) {
	assert.Equal(t, 10*time.Millisecond, timing.TlsHandshakeDuration())
}

func TestTiming_ServerDuration(t *testing.T) {
	assert.Equal(t, 14*time.Millisecond, timing.ServerDuration())
}

func TestTiming_ResponseTransferDuration(t *testing.T) {
	assert.Equal(t, 16*time.Millisecond, timing.ResponseTransferDuration())
}

func TestTiming_TotalDuration(t *testing.T) {
	assert.Equal(t, 62*time.Millisecond, timing.TotalDuration())

	withoutDnsLookup := timing
	withoutDnsLookup.DnsStart = time.Time{}
	assert.Equal(t, 56*time.Millisecond, withoutDnsLookup.TotalDuration())

	reusingConnectionTls := withoutDnsLookup
	reusingConnectionTls.ConnectStart = time.Time{}
	assert.Equal(t, 42*time.Millisecond, reusingConnectionTls.TotalDuration())

	reusingConnection := reusingConnectionTls
	reusingConnection.TlsStart = time.Time{}
	assert.Equal(t, 30*time.Millisecond, reusingConnection.TotalDuration())

	noTimingData := reusingConnection
	noTimingData.GotConn = time.Time{}
	assert.Equal(t, time.Duration(0), noTimingData.TotalDuration())
}

func TestTiming_String(t *testing.T) {
	assert.Equal(t, `<=>             DNS Lookup:        2ms
   <=>          TCP Connection:    10ms
      <=>       TLS Handshake:     10ms
         <=>    Server:            14ms
            <=> Response Transfer: 16ms
<=============> Total:             62ms`, timing.String())
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
