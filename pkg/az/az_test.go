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
			name: "happy path with symbol",
			args: args{
				length:     32,
				numDigits:  8,
				numSymbols: 2,
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
				// expectedLettersNumber is the total length of generated Password, where first character is a letter due to Azure password requirements
				// number of letters is a difference between total length and the rest of possible symbols
				expectedLettersNumber := tt.args.length - (tt.args.numDigits + tt.args.numSymbols)
				if letters != expectedLettersNumber {
					t.Errorf(`GeneratePassword() generated = %v. It has %d letters, but expected was %d letters`, got, letters, expectedLettersNumber)
				}
				if digits != tt.args.numDigits {
					t.Errorf(`GeneratePassword() generated = %v. It has %d digits, but expected was %d digits`, got, digits, tt.args.numDigits)
				}
				if others != tt.args.numSymbols {
					t.Errorf(`GeneratePassword() generated = %v. It has %d others, but expected was %d others`, got, others, tt.args.numSymbols)
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
