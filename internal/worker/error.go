package worker

//连接错误号
const errorConnect = 1

//写入错误号
const errorWrite = 2

//读取
const errorRead = 3

//其他
const errorOther = 0

//连接错误
type workerErrorConnect struct {
	err error
}

func newWorkerErrorConnect(err error) error {
	return workerErrorConnect{
		err: err,
	}
}

func (e workerErrorConnect) Error() string {
	return e.err.Error()
}

//写错误
type workerErrorWrite struct {
	err error
}

func newWorkerErrorWrite(err error) error {
	return workerErrorWrite{
		err: err,
	}
}

func (e workerErrorWrite) Error() string {
	return e.err.Error()
}

//读错误
type workerErrorRead struct {
	Err error
}

func newWorkerErrorRead(err error) error {
	return workerErrorRead{
		Err: err,
	}
}

func (e workerErrorRead) Error() string {
	return e.Err.Error()
}

func getWorkerType(err error) int {
	switch err.(type) {
	case workerErrorConnect:
		return errorConnect
	case workerErrorWrite:
		return errorWrite
	case workerErrorRead:
		return errorRead
	}
	return errorOther
}

func IsErrorConnect(err error) bool {
	return getWorkerType(err) == errorConnect
}

func IsErrorWrite(err error) bool {
	return getWorkerType(err) == errorWrite
}

func IsErrorRead(err error) bool {
	return getWorkerType(err) == errorRead
}
