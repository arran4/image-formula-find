package dna1

import (
	"testing"
)

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

func TestUnshiftMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnshiftMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("UnshiftMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestAppendMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("AppendMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestPopMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PopMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("PopMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestShiftMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShiftMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("ShiftMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestDeleteMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("DeleteMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestInsertMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("InsertMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}
func TestPositionMutate(t *testing.T) {
	tests := []struct {
		name string
		a    string
		rlen int
	}{
		{name: "basic", a: "basic", rlen: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PositionMutate(tt.a); len(got) != tt.rlen {
				t.Errorf("PositionMutate() = %v %#v, want len %v", len(got), got, tt.rlen)
			}
		})
	}
}

//func TestBugToMany0s(t *testing.T) {
//	var dna = "hzLo3uKJVAvnLGzQMF8IrP+uMG9YJ3VL7Kx/k6Go5/iLSUqKCkVE"
//	rd, _, _ := SplitString3(dna)
//	rf := ParseFunction(rd)
//	switch rf.Equals.LHS.(type) {
//	}
//}
