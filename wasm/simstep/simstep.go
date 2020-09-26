package simstep

type Action interface {
	Fields() []interface{}
}

type Victory struct{}

func (a Victory) Fields() []interface{} {
	return []interface{}{"victory"}
}

type Defeat struct{}

func (a Defeat) Fields() []interface{} {
	return []interface{}{"defeat"}
}

type Wait struct{}

func (a Wait) Fields() []interface{} {
	return []interface{}{"wait"}
}

type NextRound struct{}

func (a NextRound) Fields() []interface{} {
	return []interface{}{"nextRound"}
}

type UpdateScore struct {
	Delta int
}

func (a UpdateScore) Fields() []interface{} {
	return []interface{}{"updateScore", a.Delta}
}

type Log struct {
	Message string
}

func (a Log) Fields() []interface{} {
	return []interface{}{"log", a.Message}
}

type RedLog struct {
	Message string
}

func (a RedLog) Fields() []interface{} {
	return []interface{}{"redLog", a.Message}
}

type GreenLog struct {
	Message string
}

func (a GreenLog) Fields() []interface{} {
	return []interface{}{"greenLog", a.Message}
}

type ChangeCardCount struct {
	Name  string
	Delta int
}

func (a ChangeCardCount) Fields() []interface{} {
	return []interface{}{"changeCardCount", a.Name, a.Delta}
}

type UpdateHP struct {
	Delta int
}

func (a UpdateHP) Fields() []interface{} {
	return []interface{}{"updateHP", a.Delta}
}

type UpdateMP struct {
	Delta int
}

func (a UpdateMP) Fields() []interface{} {
	return []interface{}{"updateMP", a.Delta}
}

type UpdateCreepHP struct {
	Delta int
}

func (a UpdateCreepHP) Fields() []interface{} {
	return []interface{}{"updateCreepHP", a.Delta}
}

type SetCreep struct {
	Name string
	HP   int
}

func (a SetCreep) Fields() []interface{} {
	return []interface{}{"setCreep", a.Name, a.HP}
}

type SetNextCreep struct {
	Name string
	HP   int
}

func (a SetNextCreep) Fields() []interface{} {
	return []interface{}{"setNextCreep", a.Name, a.HP}
}
