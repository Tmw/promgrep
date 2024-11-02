package exposition

type TokenType string

const (
	TokenTypeHelp       = "help"
	TokenTypeType       = "type"
	TokenTypeMetric     = "metric"
	TokenTypeLabelName  = "labelname"
	TokenTypeLabelValue = "labelvalue"
	TokenTypeNumber     = "number"
)

const TYPE = "TYPE"
const HELP = "HELP"

type Token struct {
	Typ TokenType
	Str string
}
