package engine

func (e *Engine) LoadScript(source []byte) error {
	_, err := e.VM.Run(string(source))
	return err
}
