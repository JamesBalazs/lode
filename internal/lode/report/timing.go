package report

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

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
	if t.DnsDone.IsZero() { // did not do DNS lookup (connecting to IP)
		return t.Done.Sub(t.ConnectStart)
	} else {
		return t.Done.Sub(t.DnsStart)
	}
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
