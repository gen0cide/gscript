package engine

type VMExecResponse struct {
	Stdout   []string `json:"stdout"`
	Stderr   []string `json:"stderr"`
	Success  bool     `json:"success"`
	PID      int      `json:"pid"`
	ErrorMsg string   `json:"error_message"`
}
