package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

var (
	errDivisionByZero = errors.New("division by zero")

	diceItemRegexpText = `([\+\-\*\/])?(\d+)(d([\d%]+))?`
	diceTextRegexp     = regexp.MustCompile(`^(` + diceItemRegexpText + `)+$`)
	diceItemRegexp     = regexp.MustCompile(diceItemRegexpText)
)

func (p *Plugin) handleDice(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tools.dice.too-few")
		return
	}

	diceText := event.Fields()[1]
	if !diceTextRegexp.MatchString(diceText) {
		event.Respond("tools.dice.invalid")
		return
	}

	result, err := parseDiceText(diceText)
	if err != nil {
		if errors.Is(err, errDivisionByZero) {
			event.Respond("tools.dice.division-by-zero")
			return
		}
		event.Except(err)
		return
	}

	// remove 0s at the end after the decimal point
	resultText := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", result), "0"), ".")

	_, err = event.Respond("tools.dice.result", "result", resultText)
	event.Except(err)
}

func parseDiceText(text string) (float64, error) {
	parts := diceItemRegexp.FindAllStringSubmatch(text, -1)

	dices := make([]dice, len(parts))
	for i, part := range parts {
		dices[i] = parseDiceItem(part)
	}

	var totalResult, diceResult float64
	for _, dice := range dices {
		diceResult = 0

		if dice.fixed > 0 {
			diceResult = float64(dice.fixed)
		} else if dice.amount > 0 {
			for i := 0; i < dice.count; i++ {
				diceResult += float64(rand.Intn(dice.amount) + 1)
			}
		}

		switch dice.prevOp {
		case operatorPlus:
			totalResult = totalResult + diceResult
		case operatorMinus:
			totalResult = totalResult - diceResult
		case operatorMultiplication:
			totalResult = totalResult * diceResult
		case operatorDivision:
			if diceResult != 0 {
				totalResult = totalResult / diceResult
			} else {
				return 0, errDivisionByZero
			}
		}
	}

	return totalResult, nil
}

func parseDiceItem(item []string) dice {

	count, _ := strconv.Atoi(item[2])

	amount, _ := strconv.Atoi(item[4])
	if item[4] == "%" {
		amount = 100
	}

	var fixedAmount int
	if item[3] == "" {
		count = 0
		amount = 0
		fixedAmount, _ = strconv.Atoi(item[2])
	}

	var op operator
	switch item[1] {
	case "+":
		op = operatorPlus
	case "-":
		op = operatorMinus
	case "*":
		op = operatorMultiplication
	case "/":
		op = operatorDivision
	}

	return dice{
		count:  count,
		amount: amount,
		fixed:  fixedAmount,
		prevOp: op,
	}
}

type operator int

const (
	operatorPlus operator = iota
	operatorMinus
	operatorMultiplication
	operatorDivision
)

type dice struct {
	count  int
	amount int
	fixed  int
	prevOp operator
}
