package az

import (
	"testing"
	"unicode"
)

func TestGeneratePassword(t *testing.T) {
	type args struct {
		length     int
		numDigits  int
		numSymbols int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				length:     32,
				numDigits:  10,
				numSymbols: 0,
			},
			wantErr: false,
		},
		{
			name: "too many digits",
			args: args{
				length:     10,
				numDigits:  11,
				numSymbols: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePassword(tt.args.length, tt.args.numDigits, tt.args.numSymbols)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				letters, digits, others := passwordCharactersCounter(got)
				if letters != (tt.args.length-(tt.args.numDigits-1-tt.args.numSymbols)) || digits != (tt.args.numDigits-1-tt.args.numSymbols) || others != tt.args.numSymbols {
					t.Errorf(`GeneratePassword() generated = %v.
It has %d letters, %d digits and %d other characters, but expected was %d letters, %d digits and 0 others`, got, letters, digits, others, tt.args.length-tt.args.numDigits, tt.args.numDigits)
				}
			}
		})
	}
}

func passwordCharactersCounter(password string) (int, int, int) {
	letters := 0
	digits := 0
	others := 0
	for _, r := range password {
		if unicode.IsLetter(r) {
			letters++
		} else if unicode.IsDigit(r) {
			digits++
		} else {
			others++
		}
	}
	return letters, digits, others
}
