package main

import (
	"errors"
	"fmt"
	"go/format"
	"reflect"
	"strings"
	"syscall/js"
	"time"

	"github.com/quasilyte/gophers-and-dragons/game"
	"github.com/quasilyte/gophers-and-dragons/wasm/gamedata"
	"github.com/quasilyte/gophers-and-dragons/wasm/sim"
	"github.com/quasilyte/gophers-and-dragons/wasm/simstep"
	"github.com/traefik/yaegi/interp"
	// "github.com/traefik/yaegi/stdlib"
)

func main() {
	js.Global().Set("gominify", js.FuncOf(gominify))
	js.Global().Set("gofmt", js.FuncOf(gofmt))
	js.Global().Set("evalGo", js.FuncOf(evalGo))
	js.Global().Set("runSimulation", js.FuncOf(runSimulationJS))
	js.Global().Set("getCreepStats", js.FuncOf(getCreepStats))
	js.Global().Set("getCardStats", js.FuncOf(getCardStats))

	ch := make(chan struct{})
	<-ch
}

func gofmt(this js.Value, inputs []js.Value) interface{} {
	code := inputs[0].String()
	pretty, err := format.Source([]byte(code))
	if err != nil {
		return "error: " + err.Error()
	}
	return string(pretty)
}

func gominify(this js.Value, inputs []js.Value) interface{} {
	code := inputs[0].String()
	return code

	// fset := token.NewFileSet()
	// f, err := parser.ParseFile(fset, "gophers-and-dragons.go", []byte(code), 0)
	// if err != nil {
	// 	panic(err)
	// }
	// return string(minformat.Node(f))
}

func evalGo(this js.Value, inputs []js.Value) interface{} {
	code := inputs[0].String()
	i := interp.New(interp.Options{})
	res, err := i.Eval(code)
	if err != nil {
		panic(err)
	}
	return res.Interface()
}

func getCardStats(this js.Value, inputs []js.Value) interface{} {
	name := inputs[0].String()

	var typ game.CardType

	switch name {
	case "Attack":
		typ = game.CardAttack
	case "PowerAttack":
		typ = game.CardPowerAttack
	case "Stun":
		typ = game.CardStun
	case "MagicArrow":
		typ = game.CardMagicArrow
	case "Firebolt":
		typ = game.CardFirebolt

	case "Retreat":
		typ = game.CardRetreat
	case "Rest":
		typ = game.CardRest
	case "Heal":
		typ = game.CardHeal
	case "Parry":
		typ = game.CardParry

	default:
		return nil
	}

	return cardStatsToJS(gamedata.GetCardStats(typ))
}

func getCreepStats(this js.Value, inputs []js.Value) interface{} {
	name := inputs[0].String()

	var typ game.CreepType

	switch name {
	case "Cheepy":
		typ = game.CreepCheepy
	case "Imp":
		typ = game.CreepImp
	case "Lion":
		typ = game.CreepLion
	case "Fairy":
		typ = game.CreepFairy
	case "Mummy":
		typ = game.CreepMummy
	case "Dragon":
		typ = game.CreepDragon
	default:
		return nil
	}

	return creepStatsToJS(gamedata.GetCreepStats(typ))
}

func runSimulation(config js.Value, code string) (actions []simstep.Action, err error) {
	i := interp.New(interp.Options{})

	i.Use(map[string]map[string]reflect.Value{
		"github.com/quasilyte/gophers-and-dragons/game": {
			"State":          reflect.ValueOf((*game.State)(nil)),
			"Avatar":         reflect.ValueOf((*game.Avatar)(nil)),
			"AvatarStats":    reflect.ValueOf((*game.AvatarStats)(nil)),
			"Card":           reflect.ValueOf((*game.Card)(nil)),
			"CardStats":      reflect.ValueOf((*game.CardStats)(nil)),
			"CardType":       reflect.ValueOf((*game.CardType)(nil)),
			"Creep":          reflect.ValueOf((*game.Creep)(nil)),
			"CreepStats":     reflect.ValueOf((*game.CreepStats)(nil)),
			"CreepType":      reflect.ValueOf((*game.CreepType)(nil)),
			"CreepTrait":     reflect.ValueOf((*game.CreepTrait)(nil)),
			"CreepTraitList": reflect.ValueOf((*game.CreepTraitList)(nil)),
			"IntRange":       reflect.ValueOf((*game.IntRange)(nil)),

			"CreepCheepy": reflect.ValueOf(game.CreepCheepy),
			"CreepImp":    reflect.ValueOf(game.CreepImp),
			"CreepLion":   reflect.ValueOf(game.CreepLion),
			"CreepFairy":  reflect.ValueOf(game.CreepFairy),
			"CreepMummy":  reflect.ValueOf(game.CreepMummy),
			"CreepDragon": reflect.ValueOf(game.CreepDragon),

			"TraitCoward":        reflect.ValueOf(game.TraitCoward),
			"TraitMagicImmunity": reflect.ValueOf(game.TraitMagicImmunity),
			"TraitWeakToFire":    reflect.ValueOf(game.TraitWeakToFire),
			"TraitSlow":          reflect.ValueOf(game.TraitSlow),
			"TraitRanged":        reflect.ValueOf(game.TraitRanged),

			"CardMagicArrow":  reflect.ValueOf(game.CardMagicArrow),
			"CardAttack":      reflect.ValueOf(game.CardAttack),
			"CardPowerAttack": reflect.ValueOf(game.CardPowerAttack),
			"CardStun":        reflect.ValueOf(game.CardStun),
			"CardFirebolt":    reflect.ValueOf(game.CardFirebolt),
			"CardRetreat":     reflect.ValueOf(game.CardRetreat),
			"CardRest":        reflect.ValueOf(game.CardRest),
			"CardHeal":        reflect.ValueOf(game.CardHeal),
			"CardParry":       reflect.ValueOf(game.CardParry),
		},
	})
	// i.Use(stdlib.Symbols)

	if _, err := i.Eval(code); err != nil {
		return nil, err
	}

	pkg := inferPackage(code)
	userFuncSym := "ChooseCard"
	if pkg != "" {
		userFuncSym = pkg + ".ChooseCard"
	}

	res, err := i.Eval(userFuncSym)
	if err != nil {
		return nil, errors.New("can't find proper ChooseCard definition")
	}

	userFunc, ok := res.Interface().(func(game.State) game.CardType)
	if !ok {
		return nil, errors.New("can't find proper ChooseCard definition")
	}

	seed := config.Get("seed")
	simConfig := &sim.Config{
		Rounds:   config.Get("rounds").Int(),
		AvatarHP: config.Get("avatarHP").Int(),
		AvatarMP: config.Get("avatarMP").Int(),
	}
	if seed.Type() == js.TypeNumber {
		simConfig.Seed = int64(seed.Int())
	} else {
		simConfig.Seed = time.Now().UnixNano()
	}

	return sim.Run(simConfig, userFunc), nil
}

func runSimulationJS(this js.Value, inputs []js.Value) interface{} {
	config := inputs[0]
	code := inputs[1].String()

	actions, err := runSimulation(config, code)
	if err != nil {
		return []interface{}{
			(simstep.RedLog{Message: fmt.Sprintf("Error: %s", err.Error())}).Fields(),
		}
	}

	jsResult := make([]interface{}, len(actions))
	for i, a := range actions {
		jsResult[i] = a.Fields()
	}
	return jsResult
}

func inferPackage(s string) string {
	newline := strings.IndexByte(s, '\n')
	if newline == -1 {
		return ""
	}
	line := s[:newline]
	if !strings.HasPrefix(line, "package ") {
		return ""
	}
	packageName := line[len("package "):]
	return packageName
}

func creepStatsToJS(stats game.CreepStats) map[string]interface{} {
	var traits []interface{}
	for _, x := range stats.Traits {
		traits = append(traits, x.String())
	}
	return map[string]interface{}{
		"maxHP":       stats.MaxHP,
		"damage":      []interface{}{stats.Damage[0], stats.Damage[1]},
		"scoreReward": stats.ScoreReward,
		"cardsReward": stats.CardsReward,
		"traits":      traits,
	}
}

func cardStatsToJS(stats game.CardStats) map[string]interface{} {
	return map[string]interface{}{
		"mp":          stats.MP,
		"isMagic":     stats.IsMagic,
		"isOffensive": stats.IsOffensive,
		"power":       []interface{}{stats.Power[0], stats.Power[1]},
		"effect":      stats.Effect,
	}
}
