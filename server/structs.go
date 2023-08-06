package server

type IP_INFORMATION struct {
	Served    int64
	Attempts  int
	Salt      string
	Solution  string
	Challenge string
}
