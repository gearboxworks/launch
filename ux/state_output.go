package ux

import (
	"fmt"
	"regexp"
	"strings"
)


func (p *State) SetRunState(rs string) {
	p.RunState = rs
}
func (p *State) GetRunState() string {
	return p.RunState
}

func (p *State) RunStateEquals(format string, args ...interface{}) bool {
	var ret bool

	for range OnlyOnce {
		s := fmt.Sprintf(format, args...)
		if strings.Compare(p.RunState, s) == 0 {
			ret = true
		}
	}

	return ret
}


func (p *State) SetOutput(data ...interface{}) {
	for range OnlyOnce {
		p.Output = ""
		p.OutputArray = []string{}

		p.OutputAppend(data...)
	}
}

func (p *State) OutputAppend(data ...interface{}) {
	for range OnlyOnce {
		if p._Separator == "" {
			p._Separator = DefaultSeparator
		}

		for _, d := range data {
			//value := reflect.ValueOf(d)
			//switch value.Kind() {
			//	case reflect._Output:
			//		p._Array = append(p._Array, value._Output())
			//	case reflect.Array:
			//		p._Array = append(p._Array, d.([]string)...)
			//	case reflect.Slice:
			//		p._Array = append(p._Array, d.([]string)...)
			//}

			var sa []string
			switch d.(type) {
				case []string:
					for _, s := range d.([]string) {
						sa = append(sa, strings.Split(s, p._Separator)...)
					}
				case string:
					sa = append(sa, strings.Split(d.(string), p._Separator)...)
				case []byte:
					ts := d.([]byte)
					sa = append(sa, strings.Split(string(ts), p._Separator)...)
			}

			p.OutputArray = append(p.OutputArray, sa...)
		}
		p.Output = strings.Join(p.OutputArray, p._Separator)
	}
}

func (p *State) GetOutput() string {
	if p._Separator == "" {
		p._Separator = DefaultSeparator
	}

	return strings.Join(p.OutputArray, p._Separator)
}

func (p *State) GetOutputArray() []string {
	return p.OutputArray
}

func (p *State) SetSeparator(separator string) {
	for range OnlyOnce {
		p._Separator = separator
		p.OutputArray = strings.Split(p.Output, p._Separator)
	}
}

func (p *State) GetSeparator() string {
	return p._Separator
}

func (p *State) OutputTrim() {
	for range OnlyOnce {
		p.Output = strings.TrimSpace(p.Output)
		p.OutputArray = strings.Split(p.Output, p._Separator)
	}
}

func (p *State) OutputArrayTrim() {
	for range OnlyOnce {
		for _, s := range p.OutputArray {
			p.OutputArray = append(p.OutputArray, strings.Split(s, p._Separator)...)
		}
		p.Output = strings.Join(p.OutputArray, p._Separator)
	}
}

func (p *State) OutputEquals(format string, args ...interface{}) bool {
	var ret bool

	for range OnlyOnce {
		s := fmt.Sprintf(format, args...)
		if strings.Compare(p.Output, s) == 0 {
			ret = true
		}
	}

	return ret
}

func (p *State) OutputParse(format string, args ...interface{}) bool {
	var ret bool

	for range OnlyOnce {
		s := fmt.Sprintf(format, args...)

		ret = strings.Contains(p.Output, s)
	}

	return ret
}

func (p *State) OutputArrayGrep(format string, a ...interface{}) []string {
	var ret []string

	for range OnlyOnce {
		if len(p.OutputArray) == 0 {
			break
		}

		res := fmt.Sprintf(format, a...)
		re := regexp.MustCompile(res)
		for _, line := range p.OutputArray {
			if re.MatchString(line) {
				ret = append(ret, line)
			}
		}
	}

	return ret
}

func (p *State) OutputGrep(format string, a ...interface{}) string {
	var ret string

	for range OnlyOnce {
		if p.Output == "" {
			break
		}

		res := fmt.Sprintf(format, a...)
		re := regexp.MustCompile(res)
		if re.MatchString(p.Output) {
			ret = p.Output
		}
	}

	return ret
}
