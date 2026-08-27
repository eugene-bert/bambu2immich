package bambu

import "testing"

func TestHandleMessageFinishTransition(t *testing.T) {
	t.Parallel()
	state := &PrintState{}
	var calls int
	onFinish := func() { calls++ }

	handleMessage([]byte(`{"print":{"gcode_state":"RUNNING"}}`), state, onFinish)
	if calls != 0 {
		t.Fatalf("RUNNING should not fire, got %d", calls)
	}

	handleMessage([]byte(`{"print":{"gcode_state":"FINISH"}}`), state, onFinish)
	if calls != 1 {
		t.Fatalf("RUNNING→FINISH should fire once, got %d", calls)
	}

	handleMessage([]byte(`{"print":{"gcode_state":"FINISH"}}`), state, onFinish)
	if calls != 1 {
		t.Fatalf("duplicate FINISH should not fire, got %d", calls)
	}
}

func TestHandleMessageFirstFinishIgnored(t *testing.T) {
	t.Parallel()
	state := &PrintState{}
	var calls int
	handleMessage([]byte(`{"print":{"gcode_state":"FINISH"}}`), state, func() { calls++ })
	if calls != 0 {
		t.Fatal("first observed state FINISH should not fire")
	}
}

func TestHandleMessageIgnoresNoise(t *testing.T) {
	t.Parallel()
	state := &PrintState{}
	var calls int
	onFinish := func() { calls++ }
	handleMessage([]byte(`not json`), state, onFinish)
	handleMessage([]byte(`{}`), state, onFinish)
	handleMessage([]byte(`{"print":{}}`), state, onFinish)
	if calls != 0 {
		t.Fatalf("noise fired onFinish: %d", calls)
	}
}
