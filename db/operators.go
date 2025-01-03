//lint:file-ignore ST1001 Allow using dot imports following Jet's convention
package db

import (
	. "github.com/go-jet/jet/v2/postgres"
)

func ILIKE(lhs, rhs StringExpression) BoolExpression {
	return BoolExp(BinaryOperator(lhs, rhs, "ILIKE"))
}
