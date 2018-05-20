package main

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()
	negativeCases := []struct {
		name        string
		minRoutines int64
	}{
		{"negative", -1},
		{"zero", 0},
	}

	for _, tc := range negativeCases {
		_, resErr := New(tc.minRoutines)
		if resErr == nil {
			t.Errorf("Must return error if minRoutines == %d", tc.minRoutines)
		}
	}

	happyCases := []struct {
		expectedMinRs int64
	}{
		{1},
		{10},
		{100},
		{1000},
		{10000},
	}

	for _, tc := range happyCases {
		resPool, resErr := New(tc.expectedMinRs)
		if resErr != nil {
			t.Errorf("Function New must NOT return error if minRoutines == %d \n %v", tc.expectedMinRs, resErr)
		}
		if resPool.routines != tc.expectedMinRs {
			t.Errorf("Expected [%d] routines, but got [%d] routines", tc.expectedMinRs, resPool.routines)
		}
		resPool.Shutdown()
	}
}

func TestWorkerPool_Add(t *testing.T) {
	t.Parallel()
	cases := []struct {
		minRs      int64
		addNums    int64
		expectedRs int64
	}{
		{10, -1, 10},
		{10, 0, 10},
		{1, 100, 101},
		{100, 555, 655},
	}
	for _, tc := range cases {
		resPool, resErr := New(tc.minRs)
		if resErr != nil {
			t.Errorf("Function New must NOT return error \n %v", resErr)
		}
		resPool.Add(tc.addNums)
		if resPool.routines != tc.expectedRs {
			t.Errorf("Expected active routines [%d], but got [%d]", tc.expectedRs, resPool.routines)
		}
		resPool.Shutdown()
	}
}

func TestWorkerPool_Shutdown(t *testing.T) {
	t.Parallel()
	cases := []struct {
		minRs      int64
		expectedRs int64
	}{
		{999, 0},
		{100, 0},
		{565, 0},
	}
	for _, tc := range cases {
		resPool, resErr := New(tc.minRs)
		if resErr != nil {
			t.Errorf("Function New must NOT return error \n %v", resErr)
		}
		resPool.Shutdown()
		if resPool.routines != tc.expectedRs {
			t.Errorf("Expected active routines [%d], but got [%d]", tc.expectedRs, resPool.routines)
		}
		err := resPool.Shutdown()
		if err == nil {
			t.Error("Shutdown can't be executed twice")
		}
	}
}

type TestWorker struct {
	value int
	sum   int
}

func (tw *TestWorker) DoWork() {
	tw.sum += tw.value
}

func TestWorkerPool_Run(t *testing.T) {
	t.Parallel()
	cases := []struct {
		tw          TestWorker
		expectedSum int
		minRoutines int64
	}{
		{TestWorker{1, 0}, 10, 10},
		{TestWorker{2, 0}, 20, 100},
		{TestWorker{100, 0}, 1000, 1000},
		{TestWorker{-50, 0}, -500, 2},
		{TestWorker{0, 0}, 0, 5},
	}
	for _, tc := range cases {
		resPool, resErr := New(tc.minRoutines)
		if resErr != nil {
			t.Errorf("Function New must NOT return error \n %v", resErr)
		}
		for i := 0; i <= 9; i++ {
			resPool.Run(&tc.tw)
		}
		for i := 0; tc.expectedSum != tc.tw.sum; i++ {
			time.Sleep(time.Millisecond * 100)
			if i == 30 {
				t.Errorf("Expected sum [%d], but got [%d] \n", tc.expectedSum, tc.tw.sum)
			}
		}
		resPool.Shutdown()
	}
}
