package mailyak

import (
	"bytes"
	"net"
)

// senderLegacy connects to the remote SMTP server without trying to
// upgrades the connection to STARTTLS, and sends the email.
type senderLegacy struct {
	hostAndPort string
	hostname    string
	buf         *bytes.Buffer
}

func (s *senderLegacy) Send(m sendableMail) error {
	conn, err := net.Dial("tcp", s.hostAndPort)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	return smtpExchange(m, conn, s.hostname, false)
}

func newSenderLegacy(hostAndPort string) *senderLegacy {
	hostName, _, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		// Really this should be an error, but we can't return it from the New()
		// constructor without breaking compatibility. Fortunately by the time
		// it gets to the dial() the user will get a pretty clear error as this
		// hostAndPort value is almost certainly invalid.
		//
		// This hostname must be split from the port so the correct value is
		// used when performing the SMTP AUTH as the Go SMTP implementation
		// refuses to send credentials over non-localhost plaintext connections,
		// and including the port messes this check up (and is probably the
		// wrong thing to be sending anyway).
		hostName = hostAndPort
	}

	return &senderLegacy{
		hostAndPort: hostAndPort,
		hostname:    hostName,
		buf:         &bytes.Buffer{},
	}
}
