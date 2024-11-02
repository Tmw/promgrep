package metricfamily

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/tmw/promgrep/pkg/exposition"
)

type Parser struct {
	tokenSeq iter.Seq[exposition.Token]
	idx      int
}

func NewParser(tokensSeq iter.Seq[exposition.Token]) *Parser {
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
		if tok.Typ == exposition.TokenTypeHelp || tok.Typ == exposition.TokenTypeType {
			// we dont care about help and type for now
			continue
		}

		if tok.Typ == exposition.TokenTypeMetric {
			// if we were already working on an entry,
			// this is the sign we are done and starting a new one..
			if currentToken != nil {
				metrics = append(metrics, *currentToken)
			}

			currentToken = &MetricFamily{
				Name:   tok.Str,
				Labels: make(map[string]string),
				Val:    0,
			}
		}

		if tok.Typ == exposition.TokenTypeLabelName {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label name, expected metric name first")
			}

			currentLabelName = &tok.Str
		}

		if tok.Typ == exposition.TokenTypeLabelValue {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label value, expected metric name first")
			}

			if currentLabelName == nil {
				return nil, fmt.Errorf("illegal token label value, expected label name first")
			}

			currentToken.Labels[*currentLabelName] = tok.Str
		}

		if tok.Typ == exposition.TokenTypeNumber {
			if currentToken == nil {
				return nil, fmt.Errorf("illegal token label name, expected metric name first")
			}

			val, err := strconv.ParseFloat(tok.Str, 64)
			if err != nil {
				return nil, fmt.Errorf("parse error: unable to parse %s as number value: %w", tok.Str, err)
			}

			currentToken.Val = val
		}
	}

	return metrics, nil
}
