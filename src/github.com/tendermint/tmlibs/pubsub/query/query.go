package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/tmlibs/pubsub"
)

type Query struct {
	str	string
	parser	*QueryParser
}

type Condition struct {
	Tag	string
	Op	Operator
	Operand	interface{}
}

func New(s string) (*Query, error) {
	p := &QueryParser{Buffer: fmt.Sprintf(`"%s"`, s)}
	p.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}
	return &Query{str: s, parser: p}, nil
}

func MustParse(s string) *Query {
	q, err := New(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %s: %v", s, err))
	}
	return q
}

func (q *Query) String() string {
	return q.str
}

type Operator uint8

const (
	OpLessEqual	Operator	= iota

	OpGreaterEqual

	OpLess

	OpGreater

	OpEqual

	OpContains
)

func (q *Query) Conditions() []Condition {
	conditions := make([]Condition, 0)

	buffer, begin, end := q.parser.Buffer, 0, 0

	var tag string
	var op Operator

	for _, token := range q.parser.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruletag:
			tag = buffer[begin:end]
		case rulele:
			op = OpLessEqual
		case rulege:
			op = OpGreaterEqual
		case rulel:
			op = OpLess
		case ruleg:
			op = OpGreater
		case ruleequal:
			op = OpEqual
		case rulecontains:
			op = OpContains
		case rulevalue:

			valueWithoutSingleQuotes := buffer[begin+1 : end-1]
			conditions = append(conditions, Condition{tag, op, valueWithoutSingleQuotes})
		case rulenumber:
			number := buffer[begin:end]
			if strings.Contains(number, ".") {
				value, err := strconv.ParseFloat(number, 64)
				if err != nil {
					panic(fmt.Sprintf("got %v while trying to parse %s as float64 (should never happen if the grammar is correct)", err, number))
				}
				conditions = append(conditions, Condition{tag, op, value})
			} else {
				value, err := strconv.ParseInt(number, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("got %v while trying to parse %s as int64 (should never happen if the grammar is correct)", err, number))
				}
				conditions = append(conditions, Condition{tag, op, value})
			}
		case ruletime:
			value, err := time.Parse(time.RFC3339, buffer[begin:end])
			if err != nil {
				panic(fmt.Sprintf("got %v while trying to parse %s as time.Time / RFC3339 (should never happen if the grammar is correct)", err, buffer[begin:end]))
			}
			conditions = append(conditions, Condition{tag, op, value})
		case ruledate:
			value, err := time.Parse("2006-01-02", buffer[begin:end])
			if err != nil {
				panic(fmt.Sprintf("got %v while trying to parse %s as time.Time / '2006-01-02' (should never happen if the grammar is correct)", err, buffer[begin:end]))
			}
			conditions = append(conditions, Condition{tag, op, value})
		}
	}

	return conditions
}

func (q *Query) Matches(tags pubsub.TagMap) bool {
	if tags.Len() == 0 {
		return false
	}

	buffer, begin, end := q.parser.Buffer, 0, 0

	var tag string
	var op Operator

	for _, token := range q.parser.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruletag:
			tag = buffer[begin:end]
		case rulele:
			op = OpLessEqual
		case rulege:
			op = OpGreaterEqual
		case rulel:
			op = OpLess
		case ruleg:
			op = OpGreater
		case ruleequal:
			op = OpEqual
		case rulecontains:
			op = OpContains
		case rulevalue:

			valueWithoutSingleQuotes := buffer[begin+1 : end-1]

			if !match(tag, op, reflect.ValueOf(valueWithoutSingleQuotes), tags) {
				return false
			}
		case rulenumber:
			number := buffer[begin:end]
			if strings.Contains(number, ".") {
				value, err := strconv.ParseFloat(number, 64)
				if err != nil {
					panic(fmt.Sprintf("got %v while trying to parse %s as float64 (should never happen if the grammar is correct)", err, number))
				}
				if !match(tag, op, reflect.ValueOf(value), tags) {
					return false
				}
			} else {
				value, err := strconv.ParseInt(number, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("got %v while trying to parse %s as int64 (should never happen if the grammar is correct)", err, number))
				}
				if !match(tag, op, reflect.ValueOf(value), tags) {
					return false
				}
			}
		case ruletime:
			value, err := time.Parse(time.RFC3339, buffer[begin:end])
			if err != nil {
				panic(fmt.Sprintf("got %v while trying to parse %s as time.Time / RFC3339 (should never happen if the grammar is correct)", err, buffer[begin:end]))
			}
			if !match(tag, op, reflect.ValueOf(value), tags) {
				return false
			}
		case ruledate:
			value, err := time.Parse("2006-01-02", buffer[begin:end])
			if err != nil {
				panic(fmt.Sprintf("got %v while trying to parse %s as time.Time / '2006-01-02' (should never happen if the grammar is correct)", err, buffer[begin:end]))
			}
			if !match(tag, op, reflect.ValueOf(value), tags) {
				return false
			}
		}
	}

	return true
}

func match(tag string, op Operator, operand reflect.Value, tags pubsub.TagMap) bool {

	value, ok := tags.Get(tag)
	if !ok {
		return false
	}
	switch operand.Kind() {
	case reflect.Struct:
		operandAsTime := operand.Interface().(time.Time)
		v, ok := value.(time.Time)
		if !ok {
			return false
		}
		switch op {
		case OpLessEqual:
			return v.Before(operandAsTime) || v.Equal(operandAsTime)
		case OpGreaterEqual:
			return v.Equal(operandAsTime) || v.After(operandAsTime)
		case OpLess:
			return v.Before(operandAsTime)
		case OpGreater:
			return v.After(operandAsTime)
		case OpEqual:
			return v.Equal(operandAsTime)
		}
	case reflect.Float64:
		operandFloat64 := operand.Interface().(float64)
		var v float64

		switch vt := value.(type) {
		case float64:
			v = vt
		case float32:
			v = float64(vt)
		case int:
			v = float64(vt)
		case int8:
			v = float64(vt)
		case int16:
			v = float64(vt)
		case int32:
			v = float64(vt)
		case int64:
			v = float64(vt)
		default:
			panic(fmt.Sprintf("Incomparable types: %T (%v) vs float64 (%v)", value, value, operandFloat64))
		}
		switch op {
		case OpLessEqual:
			return v <= operandFloat64
		case OpGreaterEqual:
			return v >= operandFloat64
		case OpLess:
			return v < operandFloat64
		case OpGreater:
			return v > operandFloat64
		case OpEqual:
			return v == operandFloat64
		}
	case reflect.Int64:
		operandInt := operand.Interface().(int64)
		var v int64

		switch vt := value.(type) {
		case int64:
			v = vt
		case int8:
			v = int64(vt)
		case int16:
			v = int64(vt)
		case int32:
			v = int64(vt)
		case int:
			v = int64(vt)
		case float64:
			v = int64(vt)
		case float32:
			v = int64(vt)
		default:
			panic(fmt.Sprintf("Incomparable types: %T (%v) vs int64 (%v)", value, value, operandInt))
		}
		switch op {
		case OpLessEqual:
			return v <= operandInt
		case OpGreaterEqual:
			return v >= operandInt
		case OpLess:
			return v < operandInt
		case OpGreater:
			return v > operandInt
		case OpEqual:
			return v == operandInt
		}
	case reflect.String:
		v, ok := value.(string)
		if !ok {
			return false
		}
		switch op {
		case OpEqual:
			return v == operand.String()
		case OpContains:
			return strings.Contains(v, operand.String())
		}
	default:
		panic(fmt.Sprintf("Unknown kind of operand %v", operand.Kind()))
	}

	return false
}
