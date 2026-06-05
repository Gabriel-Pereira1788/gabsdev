package utils

import "strconv"

func FmtDate(date string) string {
	if len(date) < 7 {
		return date
	}
	months := [12]string{"jan", "fev", "mar", "abr", "mai", "jun", "jul", "ago", "set", "out", "nov", "dez"}
	month := date[5:7]
	idx := 0
	if m, err := strconv.Atoi(month); err == nil && m >= 1 && m <= 12 {
		idx = m - 1
	}
	return months[idx] + " " + date[:4]
}
