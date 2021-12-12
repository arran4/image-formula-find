package dna1

import "testing"

func TestSplitString2(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		LHS  string
		RHS  string
	}{
		{
			name: "Basic even",
			arg:  string([]byte{chars[len(chars)/2]}) + "ABCDEFGHIJ",
			LHS:  "ABCDE",
			RHS:  "FGHIJ",
		},
		{
			name: "Basic even smol",
			arg:  string([]byte{chars[len(chars)/2]}) + "AB",
			LHS:  "A",
			RHS:  "B",
		},
		{
			name: "Basic uneven",
			arg:  string([]byte{chars[len(chars)/2]}) + "ABCDEFGHIJK",
			LHS:  "ABCDE",
			RHS:  "FGHIJK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lhs, rhs := SplitString2(tt.arg)
			if lhs != tt.LHS {
				t.Errorf("SplitString2() lhs = %v, LHS %v", lhs, tt.LHS)
			}
			if rhs != tt.RHS {
				t.Errorf("SplitString2() rhs = %v, LHS %v", rhs, tt.RHS)
			}
		})
	}
}

func TestSplitString3(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		P1   string
		P2   string
		P3   string
	}{
		{
			name: "Basic equal",
			arg:  string([]byte{chars[len(chars)/3], chars[2*len(chars)/3]}) + "AAABBBCCC",
			P1:   "AAA",
			P2:   "BBB",
			P3:   "CCC",
		},
		{
			name: "Basic equal smol",
			arg:  string([]byte{chars[len(chars)/3], chars[2*len(chars)/3]}) + "ABC",
			P1:   "A",
			P2:   "B",
			P3:   "C",
		},
		{
			name: "Basic unequal",
			arg:  string([]byte{chars[len(chars)/3], chars[2*len(chars)/3]}) + "AAABBBBCCC",
			P1:   "AAA",
			P2:   "BBBB",
			P3:   "CCC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1, p2, p3 := SplitString3(tt.arg)
			if p1 != tt.P1 {
				t.Errorf("SplitString2() p1 = %v, P1 %v", p1, tt.P1)
			}
			if p2 != tt.P2 {
				t.Errorf("SplitString2() p2 = %v, P2 %v", p2, tt.P2)
			}
			if p3 != tt.P3 {
				t.Errorf("SplitString2() p3 = %v, P3 %v", p3, tt.P3)
			}
		})
	}
}
