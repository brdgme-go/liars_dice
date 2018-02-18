package liars_dice

import (
	"errors"

	"github.com/brdgme-go/brdgme"
)

var (
	MinBidQuantity = 1
	MinBidValue    = 1
	MaxBidValue    = 6
)

type bidCommand struct {
	Quantity, Value int
}

type callCommand struct{}

func (g *Game) Command(player int, input string, players []string) (brdgme.CommandResponse, error) {
	parseOutput, err := g.Parser(player).Parse(input, players)
	if err != nil {
		return brdgme.CommandResponse{}, err
	}
	switch value := parseOutput.Value.(type) {
	case bidCommand:
		return g.BidCommand(player, value.Quantity, value.Value, parseOutput.Remaining)
	case callCommand:
		return g.CallCommand(player, parseOutput.Remaining)
	}
	return brdgme.CommandResponse{}, errors.New("inexhaustive command handler")
}

func (g *Game) Parser(player int) brdgme.Parser {
	oneOf := brdgme.OneOf{}
	if g.CanBid(player) {
		oneOf = append(oneOf, bidParser)
	}
	if g.CanCall(player) {
		oneOf = append(oneOf, callParser)
	}
	return oneOf
}

var callParser = brdgme.Map{
	Parser: brdgme.Token("call"),
	Func: func(value interface{}) interface{} {
		return callCommand{}
	},
}

var bidParser = brdgme.Map{
	Parser: brdgme.Chain([]brdgme.Parser{
		brdgme.Token("bid"),
		brdgme.AfterSpace(brdgme.Int{
			Min: &MinBidQuantity,
		}),
		brdgme.AfterSpace(brdgme.Int{
			Min: &MinBidValue,
			Max: &MaxBidValue,
		}),
	}),
	Func: func(value interface{}) interface{} {
		values := value.([]interface{})
		return bidCommand{
			Quantity: values[1].(int),
			Value:    values[2].(int),
		}
	},
}

func (g *Game) CommandSpec(player int) *brdgme.Spec {
	spec := g.Parser(player).ToSpec()
	return &spec
}
