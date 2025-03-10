package common

import (
	"regexp"
	"strconv"
	"time"

	"github.com/mvgrimes/when/rules"
)

/*

- DD/MM/YYYY
- 11/3/2015
- 11/3/2015
- 11/3

also with "\", gift for windows' users

https://play.golang.org/p/29LkTfe1Xr
*/

func SlashDMY(s rules.Strategy) rules.Rule {

	return &rules.F{
		RegExp: regexp.MustCompile("(?i)(?:\\W|^)" +
			"([0-3]{0,1}[0-9]{1})" +
			"[\\/\\\\]" +
			"([0-3]{0,1}[0-9]{1})" +
			"(?:[\\/\\\\]" +
			"((?:1|2)[0-9]{3})\\s*)?" +
			"(?:\\W|$)"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			if (c.Day != nil || c.Month != nil || c.Year != nil) && s != rules.Override {
				return false, nil
			}

			day, _ := strconv.Atoi(m.Captures[0])
			month, _ := strconv.Atoi(m.Captures[1])
			year := -1
			if m.Captures[2] != "" {
				year, _ = strconv.Atoi(m.Captures[2])
			}

			if day == 0 {
				return false, nil
			}

			if month > 12 {
				return false, nil
			}

			if day > GetDays(ref.Year(), month) {
				// invalid date: day is after last day of the month
				return false, nil
			}

		WithYear:
			if year != -1 {
				c.Year = &year
				c.Month = &month
				c.Day = &day
				return true, nil
			}

			if o.WantPast {
				if month > int(ref.Month()) {
					year = ref.Year() - 1
				} else if month == int(ref.Month()) {
					if day <= ref.Day() {
						year = ref.Year()
					} else {
						year = ref.Year() - 1
					}
				} else {
					year = ref.Year()
				}

			} else {
				if month < int(ref.Month()) {
					year = ref.Year() + 1
				} else if month == int(ref.Month()) {
					if day >= ref.Day() {
						year = ref.Year()
					} else {
						year = ref.Year() + 1
					}
				} else {
					year = ref.Year()
				}
			}

			goto WithYear
		},
	}
}
