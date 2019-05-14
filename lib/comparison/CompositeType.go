package comparison

import(
	"fmt"
	"reflect"
	"strings"
	"github.com/google/go-cmp/cmp"
)

func DeepDiff(x interface{}, y interface{}) (bool, string) {
	diff := cmp.Diff(x, y)
	return diff != "", diff
}

func IsPartOf(part interface{}, whole interface{}) (bool, string) {
	var r DiffReporter
	diff := cmp.Diff(part, whole, cmp.Reporter(&r))
	return !r.HasDiffs(), diff
}

type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		if vx.IsValid() {
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
		}
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

func (r *DiffReporter) HasDiffs() bool {
	return len(r.diffs) > 0
}

func (r *DiffReporter) GetDiffs() []string {
	return r.diffs
}

func IsZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func VariableInfo(label string, val interface{}) {
	fmt.Printf(" - %s: [%v], type: %s\n", label, val, reflect.ValueOf(val).Type().String())
}
