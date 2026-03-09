package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Pattern matches ${VAR}, ${VAR:-default}, ${VAR:?error}
var interpolationPattern = regexp.MustCompile(`\$\{([^}:]+)(?:(:[-?])([^}]*))?\}`)

// Interpolator handles variable interpolation in values.
type Interpolator struct{}

// New creates a new Interpolator instance.
func New() *Interpolator {
	return &Interpolator{}
}

// Interpolate processes a string and replaces variable references with their values.
// Supported syntax:
//   - ${VAR} - replaced with env var value, empty string if not set
//   - ${VAR:-default} - replaced with env var value, or default if not set
//   - ${VAR:?error} - replaced with env var value, or returns error if not set
func (ip *Interpolator) Interpolate(value string) (string, error) {
	if !strings.Contains(value, "${") {
		return value, nil
	}

	var interpolationErr error

	result := interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
		if interpolationErr != nil {
			return match
		}

		submatch := interpolationPattern.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		varName := strings.TrimSpace(submatch[1])
		operator := submatch[2]
		operand := submatch[3]

		envValue := os.Getenv(varName)

		switch operator {
		case ":-":
			// Default value: use operand if env var is empty
			if envValue == "" {
				return operand
			}
			return envValue
		case ":?":
			// Error if empty: return error with operand as message
			if envValue == "" {
				msg := operand
				if msg == "" {
					msg = fmt.Sprintf("variable %s is not set", varName)
				}
				interpolationErr = fmt.Errorf("interpolation error: %s", msg)
				return match
			}
			return envValue
		default:
			// Simple substitution: ${VAR}
			return envValue
		}
	})

	if interpolationErr != nil {
		return "", interpolationErr
	}

	return result, nil
}

// InterpolateList processes a list of strings and interpolates variables in each.
func (ip *Interpolator) InterpolateList(values []string) ([]string, error) {
	result := make([]string, len(values))
	for i, v := range values {
		interpolated, err := ip.Interpolate(v)
		if err != nil {
			return nil, fmt.Errorf("interpolation failed for value at index %d: %w", i, err)
		}
		result[i] = interpolated
	}
	return result, nil
}
