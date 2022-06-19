package digitaltwin

import "testing"

func TestQueryResultGeneric_HasContinuationToken(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"EmptyToken", "", false},
		{"NonEmptyToken", "AToken", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := QueryResultGeneric{ContinuationToken: test.content}

			if input.HasContinuationToken() != test.expected {
				t.Errorf("Expected %t but got %t", test.expected, input.HasContinuationToken())
			}
		})
	}
}
