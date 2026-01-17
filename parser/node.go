package parser

type NodeType int

const (
	NodeTypeUnknown NodeType = iota

	NodeTypeBinaryOp
	NodeTypeUnaryOp

	NodeTypeConstant
	NodeTypeFuncApplication
)

func (nt NodeType) String() string {
	switch nt {
	case NodeTypeBinaryOp:
		return "BinaryOp"
	case NodeTypeUnaryOp:
		return "UnaryOp"
	case NodeTypeConstant:
		return "Constant"
	case NodeTypeFuncApplication:
		return "FuncApplication"
	default:
		return "Unknown"
	}
}

type Node interface {
	Type() NodeType
	Children() []Node
}

type BinaryOpType int

const (
	BinaryOpTypeUnknown BinaryOpType = iota

	BinaryOpTypeAdd
	BinaryOpTypeSubtract

	BinaryOpTypeMultiply
	BinaryOpTypeDivide

	BinaryOpTypeExponent

	BinaryOpTypeMod
)

func (bot BinaryOpType) String() string {
	switch bot {
	case BinaryOpTypeAdd:
		return "Add"
	case BinaryOpTypeSubtract:
		return "Subtract"
	case BinaryOpTypeMultiply:
		return "Multiply"
	case BinaryOpTypeDivide:
		return "Divide"
	case BinaryOpTypeExponent:
		return "Exponent"
	case BinaryOpTypeMod:
		return "Mod"
	default:
		return "Unknown"
	}
}

type BinaryOpNode struct {
	OpType BinaryOpType
	Left   Node
	Right  Node
}

func (n *BinaryOpNode) Type() NodeType {
	return NodeTypeBinaryOp
}

func (n *BinaryOpNode) Children() []Node {
	return []Node{n.Left, n.Right}
}

type UnaryOpType int

const (
	UnaryOpTypeUnknown UnaryOpType = iota

	UnaryOpTypeNegate
	UnaryOpTypeAbs
	UnaryOpTypeFact
)

func (uot UnaryOpType) String() string {
	switch uot {
	case UnaryOpTypeNegate:
		return "Negate"
	case UnaryOpTypeAbs:
		return "Abs"
	case UnaryOpTypeFact:
		return "Fact"
	default:
		return "Unknown"
	}
}

type UnaryOpNode struct {
	OpType UnaryOpType
	Child  Node
}

func (n *UnaryOpNode) Type() NodeType {
	return NodeTypeUnaryOp
}

func (n *UnaryOpNode) Children() []Node {
	return []Node{n.Child}
}

type ConstantNode struct {
	Name string
}

func (n *ConstantNode) Type() NodeType {
	return NodeTypeConstant
}

func (n *ConstantNode) Children() []Node {
	return nil
}

type FuncApplicationNode struct {
	FuncName string
	Args     []Node
}

func (n *FuncApplicationNode) Type() NodeType {
	return NodeTypeFuncApplication
}

func (n *FuncApplicationNode) Children() []Node {
	return n.Args
}
