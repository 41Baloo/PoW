package server

type IP_INFORMATION struct {
	Served     int64
	Attempts   int
	PublicSalt string // will be given to the client. Used so attackers cant pre-cache all possible hashes
	Salt       string
	Solution   string
	Challenge  string
}
