package simstep

type Action interface {
	Fields() []interface{}
}

type Wait struct{}

func (a Wait) Fields() []interface{} {
	return []interface{}{"wait"}
}

type SetScore struct {
	Value int
}

func (a SetScore) Fields() []interface{} {
	return []interface{}{"setScore", a.Value}
}
