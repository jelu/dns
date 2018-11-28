package dns

// MsgAcceptFunc is used early in the server code to accept or reject a message with RcodeFormatError.
// There are to booleans to be returned, once signaling the rejection and another to signal if
// a reply is to be send back (you want to prevent DNS ping-pong and not reply to a response for instance).
type MsgAcceptFunc func(dh Header) MsgAcceptAction

// DefaultMsgAcceptFunc checks the request and will reject if:
//
// * isn't a request (don't respond in that case).
// * opcode isn't OpcodeQuery or OpcodeNotify
// * Zero bit isn't zero
// * has more than 1 question in the question section
// * has more than 0 RRs in the Answer section
// * has more than 0 RRs in the Authority section
// * has more than 2 RRs in the Additional section
var DefaultMsgAcceptFunc MsgAcceptFunc = defaultMsgAcceptFunc

// MsgAcceptAction represents the action to be taken.
type MsgAcceptAction int

const (
	MsgAccept MsgAcceptAction = iota
	MsgReject
	MsgIgnore
)

var defaultMsgAcceptFunc = func(dh Header) MsgAcceptAction {
	if isResponse := dh.Bits&_QR != 0; isResponse {
		return MsgIgnore
	}

	// Don't allow dynamic updates, because then the sections can contain a whole bunch of RRs.
	opcode := int(dh.Bits>>11) & 0xF
	if opcode != OpcodeQuery && opcode != OpcodeNotify {
		return MsgReject
	}

	if isZero := dh.Bits&_Z != 0; isZero {
		return MsgReject
	}
	if dh.Qdcount != 1 {
		return MsgReject
	}
	if dh.Ancount != 0 {
		return MsgReject
	}
	if dh.Nscount != 0 {
		return MsgReject
	}
	if dh.Arcount > 2 {
		return MsgReject
	}
	return MsgAccept
}
