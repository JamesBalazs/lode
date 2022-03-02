package responseTimings

import (
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"time"
)

const TimingResolution = 1 * time.Millisecond

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

func (t Timing) DnsLookupDuration() time.Duration {
	return t.DnsDone.Sub(t.DnsStart).Truncate(TimingResolution)
}

func (t Timing) TcpConnectDuration() (result time.Duration) {
	if t.DnsDone.IsZero() { // did not do DNS lookup (connecting to IP)
		result = t.ConnectDone.Sub(t.ConnectStart)
	} else {
		result = t.ConnectDone.Sub(t.DnsDone)
	}

	return result.Truncate(TimingResolution)
}

func (t Timing) TlsHandshakeDuration() time.Duration {
	return t.TlsDone.Sub(t.TlsStart).Truncate(TimingResolution)
}

func (t Timing) ServerDuration() time.Duration {
	return t.FirstByte.Sub(t.GotConn).Truncate(TimingResolution)
}

func (t Timing) ResponseTransferDuration() time.Duration {
	return t.Done.Sub(t.FirstByte).Truncate(TimingResolution)
}

func (t Timing) TotalDuration() time.Duration {
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

	return t.Done.Sub(start).Truncate(TimingResolution)
}

func (t Timing) String() string {
	return fmt.Sprintf(`<=>             DNS Lookup:        %s
   <=>          TCP Connection:    %s
      <=>       TLS Handshake:     %s
         <=>    Server:            %s
            <=> Response Transfer: %s
<=============> Total:             %s`,
		t.DnsLookupDuration(),
		t.TcpConnectDuration(),
		t.TlsHandshakeDuration(),
		t.ServerDuration(),
		t.ResponseTransferDuration(),
		t.TotalDuration())
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
