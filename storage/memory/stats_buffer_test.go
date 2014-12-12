package memory

import (
	"testing"

	"github.com/google/cadvisor/info"
	"github.com/google/cadvisor/storage/memory"
)

func createStats(id int32) *info.ContainerStats {
	return &info.ContainerStats{
		Cpu: info.CpuStats{
			Load: id,
		},
	}
}

func expectSize(t *testing.T, sb *memory.StatsBuffer, expectedSize int) {
	if sb.Size() != expectedSize {
		t.Errorf("Expected size %v, got %v", expectedSize, sb.Size())
	}
}

func expectElements(t *testing.T, sb *memory.StatsBuffer, expected []int32) {
	res := sb.FirstN(sb.Size())
	if len(res) != len(expected) {
		t.Errorf("Expected elements %v, got %v", expected, res)
		return
	}
	for i, el := range res {
		if el.Cpu.Load != expected[i] {
			t.Errorf("Expected elements %v, got %v", expected, res)
		}
	}
}


func TestAddAndFirstN(t *testing.T) {
	sb := memory.NewStatsBuffer(5)

	// Add 1.
	sb.Add(createStats(1))
	expectSize(t, sb, 1)
	expectElements(t, sb, []int32{1})

	// Fill the buffer.
	for i := 1; i <= 5; i++ {
		expectSize(t, sb, i)
		sb.Add(createStats(int32(i)))
	}
	expectSize(t, sb, 5)
	expectElements(t, sb, []int32{1, 2, 3, 4, 5})

	// Add more than is available in the buffer
	sb.Add(createStats(6))
	expectSize(t, sb, 5)
	expectElements(t, sb, []int32{2, 3, 4, 5, 6})

	// Replace all elements.
	for i := 7; i <= 10; i++ {
		sb.Add(createStats(int32(i)))
	}
	expectSize(t, sb, 5)
	expectElements(t, sb, []int32{6, 7, 8, 9, 10})
}


func TestStartMethod(t *testing.T){
	sb := memory.NewStatsBuffer(5)

	// Fill the buffer.
	for i := 1; i <= 5; i++ {
		expectSize(t, sb, i - 1)
		sb.Add(createStats(int32(i)))
	}
	expectSize(t, sb, 5)
	expectElements(t, sb, []int32{1, 2, 3, 4, 5})

	numStats := sb.Size()
	_, start_idx := sb.Start(numStats)
	if start_idx != 0 {
		t.Errorf("Expected start index 0, got %v. The buffer size is %v.", start_idx, numStats)
	}

}

func TestNextMethod(t *testing.T){
	sb := memory.NewStatsBuffer(5)

	// Fill the buffer.
	for i := 1; i <= 5; i++ {
		expectSize(t, sb, i - 1)
		sb.Add(createStats(int32(i)))
	}
	expectSize(t, sb, 5)
	expectElements(t, sb, []int32{1, 2, 3, 4, 5})

	numStats := sb.Size()
	_, start_idx := sb.Start(numStats)

	_, next_idx := sb.Next(start_idx)
	if next_idx != start_idx + 1 {
		t.Errorf("Expected next index %v, got %v. The buffer size is %v.", start_idx + 1, next_idx, numStats)
	}

	_, next_idx = sb.Next(sb.Size() - 1)
		if next_idx != 0 {
		t.Errorf("Expected next index 0, got %v. The buffer size is %v.", next_idx, numStats)
	}

}
