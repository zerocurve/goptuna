// Code generated by "stringer -trimprefix TrialState -output stringer_trial_state.go -type=TrialState"; DO NOT EDIT.

package goptuna

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TrialStateRunning-0]
	_ = x[TrialStateComplete-1]
	_ = x[TrialStatePruned-2]
	_ = x[TrialStateFail-3]
}

const _TrialState_name = "RunningCompletePrunedFail"

var _TrialState_index = [...]uint8{0, 7, 15, 21, 25}

func (i TrialState) String() string {
	if i < 0 || i >= TrialState(len(_TrialState_index)-1) {
		return "TrialState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TrialState_name[_TrialState_index[i]:_TrialState_index[i+1]]
}
