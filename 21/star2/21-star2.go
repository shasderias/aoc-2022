package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func main() {
	monkeys, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	for _, monkey := range monkeys {
		monkey.ResolveRefs(monkeys)
	}

	monkeys["root"].Operation = OperationEq
	monkeys["humn"].Operation = OperationHuman

	ee := monkeys["root"].Value().(EqualExpr)

	var polyExpr SymExpr
	var constExpr SymExpr

	if ee.Left.(SymExpr).Coeff.Sign() != 0 {
		polyExpr = ee.Left.(SymExpr)
		constExpr = ee.Right.(SymExpr)
	} else {
		polyExpr = ee.Right.(SymExpr)
		constExpr = ee.Left.(SymExpr)
	}

	fmt.Println(polyExpr, constExpr)

	rhs := constExpr.Const.Sub(constExpr.Const, polyExpr.Const)
	rhs = rhs.Quo(rhs, polyExpr.Coeff)

	fmt.Println(rhs.RatString())
}

type Operation string

const (
	OperationNil   Operation = ""
	OperationAdd             = "+"
	OperationSub             = "-"
	OperationMul             = "*"
	OperationDiv             = "/"
	OperationEq              = "="
	OperationHuman           = "h"
)

type Monkey struct {
	Name             string
	Ref1Str, Ref2Str string
	Ref1, Ref2       *Monkey
	Constant         int
	Operation        Operation
}

type Expr interface {
}

type EqualExpr struct {
	Left, Right Expr
}

type SymExpr struct {
	Coeff *big.Rat
	Const *big.Rat
}

func Add(a, b Expr) Expr {
	ae, aIsSymbolic := a.(SymExpr)
	be, bIsSymbolic := b.(SymExpr)

	if !(aIsSymbolic && bIsSymbolic) {
		panic(fmt.Sprintf("cannot add %T and %T", a, b))
	}

	return SymExpr{
		Coeff: new(big.Rat).Add(ae.Coeff, be.Coeff),
		Const: new(big.Rat).Add(ae.Const, be.Const),
	}
}

func Sub(a, b Expr) Expr {
	ae, aIsSymbolic := a.(SymExpr)
	be, bIsSymbolic := b.(SymExpr)

	if !(aIsSymbolic && bIsSymbolic) {
		panic(fmt.Sprintf("cannot sub %T and %T", a, b))
	}
	return SymExpr{
		Coeff: new(big.Rat).Sub(ae.Coeff, be.Coeff),
		Const: new(big.Rat).Sub(ae.Const, be.Const),
	}
}

func Mul(a, b Expr) Expr {
	ae, aIsSymbolic := a.(SymExpr)
	be, bIsSymbolic := b.(SymExpr)

	if !(aIsSymbolic && bIsSymbolic) {
		panic(fmt.Sprintf("cannot mul %T and %T", a, b))
	}

	coeff := new(big.Rat)
	switch {
	case ae.Coeff.Sign() != 0 && be.Coeff.Sign() != 0:
		panic(fmt.Sprintf("cannot mul %s and %s", ae.Coeff.RatString(), be.Coeff.RatString()))
	case ae.Coeff.Sign() != 0: // ax * b
		coeff = coeff.Mul(ae.Coeff, be.Const)
	case be.Coeff.Sign() != 0: // bx * a
		coeff = coeff.Mul(be.Coeff, ae.Const)
	case ae.Coeff.Sign() == 0 && be.Coeff.Sign() == 0:
		coeff = big.NewRat(0, 1)
	default:
		panic("unreachable code")
	}

	return SymExpr{
		Coeff: coeff,
		Const: new(big.Rat).Mul(ae.Const, be.Const),
	}
}

func Div(a, b Expr) Expr {
	ae, aIsSymbolic := a.(SymExpr)
	be, bIsSymbolic := b.(SymExpr)

	if !(aIsSymbolic && bIsSymbolic) {
		panic(fmt.Sprintf("cannot div %T and %T", a, b))
	}

	coeff := new(big.Rat)
	switch {
	case ae.Coeff.Sign() != 0 && be.Coeff.Sign() != 0:
		panic(fmt.Sprintf("cannot div %s and %s", ae.Coeff.RatString(), be.Coeff.RatString()))
	case ae.Coeff.Sign() != 0: // ax / b
		coeff = coeff.Quo(ae.Coeff, be.Const)
	case be.Coeff.Sign() != 0: // a / bx = a * (1 / bx) = a * (a/b)x
		coeff = coeff.Mul(ae.Const, new(big.Rat).Inv(be.Coeff))
	case ae.Coeff.Sign() == 0 && be.Coeff.Sign() == 0:
		coeff = big.NewRat(0, 1)
	default:
		panic("unreachable code")
	}

	return SymExpr{
		Coeff: coeff,
		Const: new(big.Rat).Quo(ae.Const, be.Const),
	}
}

func (m *Monkey) Value() (ret Expr) {
	defer func() {
		fmt.Printf("%s: %v\n", m.Name, ret)
	}()
	switch m.Operation {
	case OperationNil:
		return SymExpr{Coeff: big.NewRat(0, 1), Const: big.NewRat(int64(m.Constant), 1)}
	case OperationAdd:
		return Add(m.Ref1.Value(), m.Ref2.Value())
	case OperationSub:
		return Sub(m.Ref1.Value(), m.Ref2.Value())
	case OperationMul:
		return Mul(m.Ref1.Value(), m.Ref2.Value())
	case OperationDiv:
		return Div(m.Ref1.Value(), m.Ref2.Value())
	case OperationEq:
		return EqualExpr{m.Ref1.Value(), m.Ref2.Value()}
	case OperationHuman:
		return SymExpr{Coeff: big.NewRat(1, 1), Const: big.NewRat(0, 1)}
	default:
		panic(fmt.Sprintf("unknown operation %s", m.Operation))
	}
}

func (m *Monkey) ResolveRefs(monkeys map[string]*Monkey) {
	var ok bool
	if m.Ref1Str != "" {
		m.Ref1, ok = monkeys[m.Ref1Str]
		if !ok {
			panic(fmt.Sprintf("no monkey %s", m.Ref1Str))
		}
	}
	if m.Ref2Str != "" {
		m.Ref2, ok = monkeys[m.Ref2Str]
		if !ok {
			panic(fmt.Sprintf("no monkey %s", m.Ref2Str))
		}
	}
}

func parseInput(inputPath string) (map[string]*Monkey, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	monkeys := make(map[string]*Monkey)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")

		name := parts[0]

		monkey := Monkey{
			Name: name,
		}

		if strings.ContainsAny(parts[1], "+-*/") {
			_, err := fmt.Sscanf(parts[1], "%s %s %s", &monkey.Ref1Str, &monkey.Operation, &monkey.Ref2Str)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s: '%s'", line, parts[1])
			}
		} else {
			_, err := fmt.Sscanf(parts[1], "%d", &monkey.Constant)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s", line)
			}
		}

		monkeys[name] = &monkey
	}

	return monkeys, nil
}
