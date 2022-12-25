package snafu

type Num int

func nDigitMax(n int) int {
	if n <= 0 {
		return 0
	}
	acc := 0
	for i := 0; i < n; i++ {
		acc += 2 * pow(5, i)
	}
	return acc
}

func (n Num) Int() int { return int(n) }

func (n Num) Snafu() string {
	switch n {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	}

	digits := ""

	var (
		digit1n = 0
	)

	for ; ; digit1n++ {
		if nDigitMax(digit1n) >= n.Int() {
			break
		}
	}

	num := int(n)

	// derive 1st digit, which must be 1 or 2
	if num > pow(5, digit1n-1)+nDigitMax(digit1n-1) {
		digits += "2"
		num -= 2 * pow(5, digit1n-1)
	} else {
		digits += "1"
		num -= pow(5, digit1n-1)
	}

	// derive remaining to 2nd last
	for i := digit1n - 2; i > 0; i-- {
		for _, d := range Digits {
			testVal := d.Value()*pow(5, i) + nDigitMax(i)
			if num <= testVal {
				digits += string(d)
				num -= d.Value() * pow(5, i)
				goto nextDigit
			}
		}
		digits += string(Digit2)
		num -= 2 * pow(5, i)
	nextDigit:
	}

	// derive last digit
	for _, d := range Digits {
		if num == d.Value() {
			digits += string(d)
			break
		}
	}

	return digits
}

type Digit string

const (
	Digit2           Digit = "2"
	Digit1                 = "1"
	Digit0                 = "0"
	DigitMinus             = "-"
	DigitDoubleMinus       = "="
)

var (
	Digits    = []Digit{DigitDoubleMinus, DigitMinus, Digit0, Digit1, Digit2}
	PosDigits = []Digit{Digit1, Digit2}
)

func (d Digit) Value() int {
	switch d {
	case Digit2:
		return 2
	case Digit1:
		return 1
	case Digit0:
		return 0
	case DigitMinus:
		return -1
	case DigitDoubleMinus:
		return -2
	}
	return 0
}

func FromString(s string) Num {
	num := 0
	for i := 0; i < len(s); i++ {
		digit := Digit(s[len(s)-i-1])
		num += digit.Value() * pow(5, i)
	}
	return Num(num)
}

func FromInt(n int) Num {
	return Num(n)
}

func pow(x, y int) int {
	if y == 0 {
		return 1
	}
	return x * pow(x, y-1)
}
