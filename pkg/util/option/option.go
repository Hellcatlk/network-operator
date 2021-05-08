package option

type Option struct {
	Name  string
	Value interface{}
}

func IsExist(options []Option, name string) bool {
	for _, option := range options {
		if option.Name == name {
			return true
		}
	}
	return false
}

func Get(options []Option, name string) interface{} {
	for _, option := range options {
		if option.Name == name {
			return option.Value
		}
	}
	return nil
}
