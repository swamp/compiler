package generate_ir

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateIntLiteral(integer *decorated.IntegerLiteral, genContext *generateContext) (value.Value, error) {
	return constant.NewInt(types.I32, int64(integer.Value())), nil
}
