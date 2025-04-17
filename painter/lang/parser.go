package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ProMKQ/kpi-lab3/painter"
)

type Parser struct {
	state *painter.State
}

func NewParser(state *painter.State) *Parser {
	return &Parser{state: state}
}

func parseInts(args []string, count int) ([]int, error) {
	if len(args) != count {
		return nil, fmt.Errorf("expected %d arguments, got %d", count, len(args))
	}

	result := make([]int, count)
	for i, s := range args {
		value, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid int argument '%s'", s)
		}
		result[i] = value
	}

	return result, nil
}

// ParseLine читає одну команду та повертає відповідну Operation.
func (p *Parser) ParseLine(line string) (painter.Operation, error) {
	fields := strings.Fields(line)
	cmd := fields[0]
	args := fields[1:]

	switch cmd {
	case "white":
		return painter.WhiteFill(p.state), nil
	case "green":
		return painter.GreenFill(p.state), nil
	case "update":
		return painter.UpdateOp, nil
	case "bgrect":
		vals, err := parseInts(args, 4)
		if err != nil {
			return nil, err
		}
		return painter.BgRect(p.state, vals[0], vals[1], vals[2], vals[3]), nil
	case "figure":
		vals, err := parseInts(args, 2)
		if err != nil {
			return nil, err
		}
		return painter.AddShape(p.state, vals[0], vals[1]), nil
	case "move":
		vals, err := parseInts(args, 2)
		if err != nil {
			return nil, err
		}
		return painter.MoveShapes(p.state, vals[0], vals[1]), nil
	case "reset":
		return painter.Reset(p.state), nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

// Parse читає вхідний io.Reader і повертає список операцій.
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	var ops []painter.Operation

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		op, err := p.ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line '%s': %w", line, err)
		}
		ops = append(ops, op)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ops, nil
}
