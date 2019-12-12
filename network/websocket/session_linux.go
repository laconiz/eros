package websocket

func (ses *Session) invoke(e *network.Event) {

	defer func() {
		if err := recover(); err != nil {
			str := ses.config.Encoder.String(e)
			ses.log.Errorf("panic on invoke %s: %v", str, err)
		}
	}()

	ses.config.Invoker.Invoke(e)
}
