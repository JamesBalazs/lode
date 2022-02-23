package report

import (
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"time"
)

const timingResolution = 1 * time.Millisecond

type Timing struct {
	DnsStart     time.Time
	DnsDone      time.Time
	ConnectStart time.Time
	ConnectDone  time.Time
	TlsStart     time.Time
	TlsDone      time.Time
	GotConn      time.Time
	FirstByte    time.Time
	Done         time.Time
}

func (t *Timing) DnsLookupDuration() time.Duration {
	return t.DnsDone.Sub(t.DnsStart)
}

func (t *Timing) TcpConnectDuration() time.Duration {
	if t.DnsDone.IsZero() { // did not do DNS lookup (connecting to IP)
		return t.ConnectDone.Sub(t.ConnectStart)
	} else {
		return t.ConnectDone.Sub(t.DnsDone)
	}
}

func (t *Timing) TlsHandshakeDuration() time.Duration {
	return t.TlsDone.Sub(t.TlsStart)
}

func (t *Timing) ServerDuration() time.Duration {
	return t.FirstByte.Sub(t.GotConn)
}

func (t *Timing) ResponseTransferDuration() time.Duration {
	return t.Done.Sub(t.FirstByte)
}

func (t *Timing) TotalDuration() time.Duration {
	var start time.Time

	if !t.DnsStart.IsZero() {
		start = t.DnsStart
	} else if !t.ConnectStart.IsZero() { // did not do DNS lookup (connecting to IP) )
		start = t.ConnectStart
	} else if !t.TlsStart.IsZero() { // reused existing connection, new TLS handshake
		start = t.TlsStart
	} else if !t.GotConn.IsZero() { // reused existing connection
		start = t.GotConn
	} else {
		return time.Duration(0)
	}

	return t.Done.Sub(start)
}

func (t *Timing) String() string {
	return fmt.Sprintf(`<=>             DNS Lookup:        %s
   <=>          TCP Connection:    %s
      <=>       TLS Handshake:     %s
         <=>    Server:            %s
            <=> Response Transfer: %s
<=============> Total:             %s`,
		t.DnsLookupDuration().Truncate(timingResolution),
		t.TcpConnectDuration().Truncate(timingResolution),
		t.TlsHandshakeDuration().Truncate(timingResolution),
		t.ServerDuration().Truncate(timingResolution),
		t.ResponseTransferDuration().Truncate(timingResolution),
		t.TotalDuration().Truncate(timingResolution))
}

func NewTrace(timing *Timing) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { timing.DnsStart = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { timing.DnsDone = time.Now() },
		ConnectStart:         func(_, _ string) { timing.ConnectStart = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { timing.ConnectDone = time.Now() },
		GotConn:              func(_ httptrace.GotConnInfo) { timing.GotConn = time.Now() },
		GotFirstResponseByte: func() { timing.FirstByte = time.Now() },
		TLSHandshakeStart:    func() { timing.TlsStart = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { timing.TlsDone = time.Now() },
	}
}
