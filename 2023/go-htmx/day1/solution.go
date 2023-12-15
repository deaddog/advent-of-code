package day1

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type CalibrationValue int

type Line struct {
	Raw                        string
	Calibration1, Calibration2 Calibration
}
type Calibration struct {
	Elements []Element
	Value    CalibrationValue
}

type Element struct {
	Content  string
	IsNumber bool
}

func Solve(file string) ([]Line, error) {
	lines := []Line{}
	for _, l := range strings.Split(file, "\n") {
		if l == "" {
			continue
		}

		line, err := SolveLine(l)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func Sum1(lines []Line) int {
	s := 0
	for _, l := range lines {
		s += int(l.Calibration1.Value)
	}
	return s
}
func Sum2(lines []Line) int {
	s := 0
	for _, l := range lines {
		s += int(l.Calibration2.Value)
	}
	return s
}

func SolveLine(l string) (Line, error) {
	cal1, err := GetCalibration(l, "[0-9]")
	if err != nil {
		return Line{}, err
	}

	cal2, err := GetCalibration(l, "[0-9]|one|two|three|four|five|six|seven|eight|nine")
	if err != nil {
		return Line{}, err
	}

	return Line{
		Raw:          l,
		Calibration1: cal1,
		Calibration2: cal2,
	}, nil
}

func GetCalibration(l string, pattern string) (Calibration, error) {
	first := regexp.MustCompile(fmt.Sprintf("(%s)", pattern)).FindStringSubmatchIndex(l)
	last := regexp.MustCompile(fmt.Sprintf(".*(%s)", pattern)).FindStringSubmatchIndex(l)

	if len(first) == 0 {
		return Calibration{}, fmt.Errorf("no calibration digits found in line")
	}

	firstNumber, err := convert(l[first[2]:first[3]])
	if err != nil {
		return Calibration{}, err
	}

	lastNumber, err := convert(l[last[2]:last[3]])
	if err != nil {
		return Calibration{}, err
	}

	cal := Calibration{
		Value: CalibrationValue(firstNumber*10 + lastNumber),
	}

	if len(first) == 0 || len(last) == 0 {
		log.Println(l)
	}

	if first[2] != 0 {
		cal.Elements = append(cal.Elements, Element{
			Content:  l[0:first[2]],
			IsNumber: false,
		})
	}

	cal.Elements = append(cal.Elements, Element{
		Content:  l[first[2]:first[3]],
		IsNumber: true,
	})

	if first[2] != last[2] {
		cal.Elements = append(cal.Elements, Element{
			Content:  l[first[3]:last[2]],
			IsNumber: false,
		})
		cal.Elements = append(cal.Elements, Element{
			Content:  l[last[2]:last[3]],
			IsNumber: true,
		})
	}

	if last[2] != len(l) {
		cal.Elements = append(cal.Elements, Element{
			Content:  l[last[3]:],
			IsNumber: false,
		})
	}

	return cal, nil
}

func convert(value string) (int, error) {
	switch value {
	case "one":
		return 1, nil
	case "two":
		return 2, nil
	case "three":
		return 3, nil
	case "four":
		return 4, nil
	case "five":
		return 5, nil
	case "six":
		return 6, nil
	case "seven":
		return 7, nil
	case "eight":
		return 8, nil
	case "nine":
		return 9, nil
	default:
		return strconv.Atoi(value)
	}
}
