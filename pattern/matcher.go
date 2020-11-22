package pattern

import (
	"errors"
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

var (
	ErrInvalidReturnValType = errors.New("invalid return value type")
)

// Matcher 用于检查 data 是否匹配表达式
type Matcher struct {
	expression string
	vm         *vm.Program
}

// NewMatcher 创建一个新的 Matcher
func NewMatcher(expression string, dataType interface{}) (*Matcher, error) {
	vmm, err := expr.Compile(expression, expr.Env(dataType), expr.AsBool())
	if err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	return &Matcher{
		expression: expression,
		vm:         vmm,
	}, nil
}

// Match 执行匹配，检查 data 是否能够通过 expression 的评估
func (matcher *Matcher) Match(data interface{}) (bool, error) {
	rs, err := expr.Run(matcher.vm, data)
	if err != nil {
		return false, fmt.Errorf("evaluate expression failed: %w", err)
	}

	if matched, ok := rs.(bool); ok {
		return matched, nil
	}

	return false, ErrInvalidReturnValType
}

// Match 检查 data 是否匹配表达式
func Match(expression string, data interface{}) (bool, error) {
	matcher, err := NewMatcher(expression, data)
	if err != nil {
		return false, err
	}

	return matcher.Match(data)
}

// StringMatch 检查 data 是否匹配表达式
func StringMatch(expression string, data string) (bool, error) {
	matcher, err := NewMatcher(expression, &Data{Data: data})
	if err != nil {
		return false, err
	}

	return matcher.Match(&Data{Data: data})
}

// Evaluator 用于计算 data 在应用 expression 后的结果
type Evaluator struct {
	expression string
	vm         *vm.Program
}

// NewEvaluator 创建一个新的 Evaluator
func NewEvaluator(expression string, dataType interface{}) (*Evaluator, error) {
	vmm, err := expr.Compile(expression, expr.Env(dataType))
	if err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	return &Evaluator{
		expression: expression,
		vm:         vmm,
	}, nil
}

// Eval 对 data 应用 expression 表达式，返回评估后的结果
func (eval *Evaluator) Eval(data interface{}) (string, error) {
	rs, err := expr.Run(eval.vm, data)
	if err != nil {
		return "", fmt.Errorf("evaluate expression failed: %w", err)
	}

	return fmt.Sprintf("%v", rs), nil
}

// Eval 对 data 应用表达式
func Eval(expression string, data interface{}) (string, error) {
	evaluator, err := NewEvaluator(expression, data)
	if err != nil {
		return "", err
	}

	return evaluator.Eval(data)
}

// StringEval 对 data 应用表达式
func StringEval(expression string, data string) (string, error) {
	evaluator, err := NewEvaluator(expression, &Data{Data: data})
	if err != nil {
		return "", err
	}

	return evaluator.Eval(&Data{Data: data})
}
