package utils

func LuhnValid(order string) bool {
	sum := 0
	parity := len(order) % 2
	for i := len(order) - 1; i >= 0; i-- {
		dig := int(order[i]) - 48
		if i%2 == parity {
			dig *= 2
			if dig > 9 {
				dig -= 9
			}
		}
		sum += dig
	}
	return sum%10 == 0
}
