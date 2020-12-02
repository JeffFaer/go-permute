package permute

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFact(t *testing.T) {
	s := []string{"1", "2", "3"}
	want := []fact{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{0, 2},
		{1, 2},
	}
	wantPerm := []perm{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for i := range want {
		f, ok := newFact(int64(i), len(s))
		if !ok {
			t.Errorf("newFact(%d, %d) was not ok", i, len(s))
			continue
		}
		if d := cmp.Diff(want[i], f); d != "" {
			t.Errorf("newFact(%d, %d) -want +got:\n%s", i, len(s), d)
		}

		p := f.perm()
		if d := cmp.Diff(wantPerm[i], p); d != "" {
			t.Errorf("newFact(%d,%d).perm() -want +got:\n%s", i, len(s), d)
		}
	}
}

func TestFactTooBig(t *testing.T) {
	// 7 > 3!
	if f, ok := newFact(7, 3); ok {
		t.Errorf("newFact(7, 3) was ok: %v", f)
	}
}

func TestPerm(t *testing.T) {
	s := []string{"0", "1", "2"}
	p := perm([]int{0, 1, 2})

	s2 := append([]string(nil), s...)
	p.apply(slice{s2})
	if g := strings.Join(s2, ""); g != "012" {
		t.Errorf("got %s, want 012", g)
	}

	s2 = append([]string(nil), s...)
	p = perm([]int{1, 2, 0})
	p.apply(slice{s2})
	if g := strings.Join(s2, ""); g != "120" {
		t.Errorf("got %s, want 120", g)
	}
}

func TestPermuter(t *testing.T) {
	for _, tc := range []struct {
		in   []string
		want []string
	}{
		{
			in: []string{"1", "2"},
			want: []string{
				"12",
				"21",
			},
		},
		{
			in: []string{"1", "2", "3"},
			want: []string{
				"123",
				"132",
				"213",
				"231",
				"312",
				"321",
			},
		},
		{
			in: []string{"1", "2", "3", "4"},
			want: []string{
				"1234",
				"1243",
				"1324",
				"1342",
				"1423",
				"1432",
				"2134",
				"2143",
				"2314",
				"2341",
				"2413",
				"2431",
				"3124",
				"3142",
				"3214",
				"3241",
				"3412",
				"3421",
				"4123",
				"4132",
				"4213",
				"4231",
				"4312",
				"4321",
			},
		},
	} {
		t.Run(fmt.Sprintf("len(in)=%d", len(tc.in)), func(t *testing.T) {
			origIn := append([]string(nil), tc.in...)

			var got []string
			for p := NewSlicePermuter(tc.in); p.Permute(); {
				got = append(got, strings.Join(tc.in, ""))
			}

			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("Permuter did not generate correct permutations -want +got:\n%s", d)
			}
			if d := cmp.Diff(origIn, tc.in); d != "" {
				t.Errorf("Permuter did not return input to its original order -want +got:\n%s", d)
			}
		})
	}
}
