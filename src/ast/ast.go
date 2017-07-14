package ast

type Node interface {
	TokenLiteral() string
	String() string
}
