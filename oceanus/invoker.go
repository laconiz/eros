package oceanus

//
// type Thread struct {
// 	invoker interface{}
// }
//
// func (t *Thread) Start() {
//
// }
//
// func (t *Thread) OnMessage(message *Message) {
//
// }
//
// func (t *Thread) Stop() {
//
// }
//
// func NewThread(invoker interface{}) (*Thread, error) {
//
// 	typo := reflect.TypeOf(invoker)
// 	if typo == nil {
// 		return nil, errors.New("register a nil invoker")
// 	}
//
// 	value := reflect.ValueOf(invoker)
//
// 	inject.New().Invoke(value.MethodByName("OnStart").Interface())
//
// }
