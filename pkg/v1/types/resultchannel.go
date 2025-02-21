package types

type ResultChan chan bool

func NewResultChan() ResultChan {
	return make(ResultChan)
}

func (c *ResultChan) WriteFailure() {
	*c <- false
}

func (c *ResultChan) WriteSuccess() {
	*c <- true
}
