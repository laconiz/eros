package oceanus

func NewThread(thread interface{}) {

	std.NewThread(thread)

}

var std = NewProcess()

func init() {
	std.Run()
}
