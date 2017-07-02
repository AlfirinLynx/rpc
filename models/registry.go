package models

var registry map[string]interface{}

func Registry() map[string]interface{} {
	if registry == nil {
		registry = make(map[string]interface{})
	}
	return registry
}

func register(name string, model interface{}) {
	if _, ok := Registry()[name]; !ok {
		Registry()[name] = model
	}
}
