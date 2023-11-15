package validation

import "testing"

func TestCheckMapStringNoMustache(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[string]any
		wantErr bool
	}{
		{
			name: "valid input",
			input: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "input with mustache",
			input: map[string]any{
				"key1": "{{var}}",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MapStringNoMustache(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckMapStringNoMustache() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestStringNoMustache(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "valid input",
			input: "no mustache here",
		},
		{
			name:    "input with mustache",
			input:   "some {{var}} mustache",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := StringNoMustache(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("StringNoMustache() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
