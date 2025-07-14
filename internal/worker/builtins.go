package worker

import "github.com/expr-lang/expr"

func (s *sStep) exprBuiltins() []expr.Option {
	// TODO: 内置函数或工具链
	return []expr.Option{
		// 预期返回值类型
		expr.AsBool(),
		expr.AllowUndefinedVariables(),
		//expr.Env(map[string]any{
		//	"storage": s.storage,
		//}),
		//expr.Function("test", func(params ...any) (any, error) {
		//	return "test", nil
		//}),
	}
}
