package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestBufferTime(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTime(Step(2)),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTime(Step(4)),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTime(Step(6)),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.BufferTime(Step(1)),
		),
		ErrTest,
	).TestAll()

	t.Log("----------")

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan: Step(8),
			}.Make(),
			ToString(),
		),
		"[A B C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:      Step(8),
				MaxBufferSize: 3,
			}.Make(),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:      Step(8),
				MaxBufferSize: 2,
			}.Make(),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:      Step(8),
				MaxBufferSize: 1,
			}.Make(),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", "[]", Completed,
	).TestAll()

	t.Log("----------")

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:         Step(2),
				CreationInterval: Step(2),
			}.Make(),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:         Step(2),
				CreationInterval: Step(4),
			}.Make(),
			ToString(),
		),
		"[A]", "[C]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.BufferTimeConfig{
				TimeSpan:         Step(4),
				CreationInterval: Step(2),
			}.Make(),
			ToString(),
		),
		"[A B]", "[B C]", "[C D]", "[D E]", "[E]", Completed,
	).TestAll()
}
