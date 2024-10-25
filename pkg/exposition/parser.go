package exposition

import (
	"fmt"
	"iter"
	"strconv"
)

type Parser struct {
	tokenSeq iter.Seq[Token]
	idx      int
}

func NewParser(tokensSeq iter.Seq[Token]) *Parser {
	return &Parser{
		tokenSeq: tokensSeq,
		idx:      0,
	}
}

func (p *Parser) Parse() ([]MetricFamily, error) {
	var (
		metrics          = []MetricFamily{}
		currentToken     *MetricFamily
		currentLabelName *string
	)

	for tok := range p.tokenSeq {
		if tok.typ == TokenTypeHelp || tok.typ == TokenTypeType {
			// we dont care about help and type for now
			continue
		}

		if tok.typ == TokenTypeMetric {
			// if we were already working on an entry,
			// this is the sign we are done and starting a new one..
			if currentToken != nil {
				metrics = append(metrics, *currentToken)
			}

			currentToken = &MetricFamily{
				Name:   tok.str,
				Labels: make(map[string]string),
				Val:    0,
			}
		}

		if tok.typ == TokenTypeLabelName {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label name, expected metric name first")
			}

			currentLabelName = &tok.str
		}

		if tok.typ == TokenTypeLabelValue {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label value, expected metric name first")
			}

			if currentLabelName == nil {
				return nil, fmt.Errorf("illegal token label value, expected label name first")
			}

			currentToken.Labels[*currentLabelName] = tok.str
		}

		if tok.typ == TokenTypeNumber {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label name, expected metric name first")
			}

			val, err := strconv.ParseFloat(tok.str, 64)
			if err != nil {
				return nil, fmt.Errorf("parse error: unable to parse %s as number value: %w", tok.str, err)
			}

			currentToken.Val = val
		}
	}

	return metrics, nil
}
