package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func expectedResponseMatchParser(context *context) {
	context.log("05 expected-response-match-parser")

	expectedResponse := context.SpecTriplet.ExpectedResponse

	for _, line := range expectedResponse.allLines() {
		parts := strings.Split(line.Text, regexpIdentifier)

		count := len(parts)

		if count == 1 {
			continue
		}

		if (count-4)%3 != 0 {
			errorHandler(
				context,
				fmt.Errorf(
					"regexp must be formed %soptional-capture-name%sregexp%s",
					regexpIdentifier,
					regexpIdentifier,
					regexpIdentifier,
				),
			)

			return
		}

		var re *regexp.Regexp
		var err error

		if parts[0] != "" {
			re, err = regexp.Compile(regexp.QuoteMeta(parts[0]))

			if errorHandler(context, err) {
				return
			}

			line.RegexpNames = append(line.RegexpNames, ":prefix")
			line.Regexps = append(line.Regexps, re)
		}

		for i := 1; i < count-1; i += 3 {
			captureName := parts[i]

			if captureName == ":prefix" || captureName == ":postfix" {
				errorHandler(
					context,
					errors.New("capture names cannot be :prefix or :postfix"),
				)

				return
			}

			reString := "("

			switch parts[i+1] {
			case ":date":
				reString +=
					"(Mon|Tue|Wed|Thu|Fri|Sat|Sun), " +
						"(0\\d|1\\d|2\\d|3[01]) " +
						"(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) " +
						"2\\d{3} " +
						"(0\\d|1\\d|2[0-3]):" +
						"(0\\d|1\\d|2\\d|3\\d|4\\d|5\\d):" +
						"(0\\d|1\\d|2\\d|3\\d|4\\d|5\\d) " +
						"(A|M|N|Y|Z|UT|GMT|[A-Z]{3}|[+-](0\\d|1[012]))"
			case ":b62:22":
				reString += "[0-9A-Za-z]{22}"
			case ":iso8601:µs:z":
				reString += "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}[.]\\d{6}Z"
			case ":uuid":
				reString +=
					"[[:xdigit:]]{8}-" +
						"[[:xdigit:]]{4}-" +
						"[[:xdigit:]]{4}-" +
						"[[:xdigit:]]{4}-" +
						"[[:xdigit:]]{12}"
			default:
				reString += parts[i+1]
			}

			reString += ")"

			re, err = regexp.Compile(reString)

			if errorHandler(context, err) {
				return
			}

			line.RegexpNames = append(line.RegexpNames, parts[i])
			line.Regexps = append(line.Regexps, re)

			postfix := parts[i+2]

			if postfix != "" {
				re, err = regexp.Compile(regexp.QuoteMeta(postfix))

				if errorHandler(context, err) {
					return
				}

				line.RegexpNames = append(line.RegexpNames, ":postfix")
				line.Regexps = append(line.Regexps, re)
			}
		}
	}

	expectedResponseSubstituter(context)
}
