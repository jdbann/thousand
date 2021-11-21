package templates

type Data interface {
	Add(string, interface{}) Data
}

func NewData() Data {
	return make(data)
}

type data map[string]interface{}

func (d data) Add(k string, v interface{}) Data {
	d[k] = v
	return d
}
