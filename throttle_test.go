package throttle

import (
	"testing"
	"time"
)

func TestBucketTake(t *testing.T) {
	b := NewBucket(100, 10)
	n, wait := b.Take(5)
	if n != 5 {
		t.Errorf("Tokens got %d, expected 5", n)
	}
	if wait != 0 {
		t.Errorf("Wait must be 0")
	}

	n, wait = b.Take(10)
	if n != 5 {
		t.Errorf("Tokens got %d, expected 5", n)
	}
	if wait < 40*time.Millisecond || 50*time.Millisecond < wait {
		t.Errorf("Invalid wait time: %v", wait)
	}

	time.Sleep(10 * time.Millisecond)
	n, wait = b.Take(5)
	if n != 1 {
		t.Errorf("Tokens got %d, expected 1", n)
	}
	if wait < 30*time.Millisecond || 40*time.Millisecond < wait {
		t.Errorf("Invalid wait time: %v", wait)
	}

	time.Sleep(wait)
	n, wait = b.Take(4)
	if n != 4 {
		t.Errorf("Tokens got %d, expected 4", n)
	}
	if wait != 0 {
		t.Errorf("Wait must be 0")
	}
}

func TestBucketTakeLarge(t *testing.T) {
	b := NewBucket(100, 10)
	n, wait := b.Take(30)
	if n != 10 {
		t.Errorf("Tokens got %d, expected 10", n)
	}
	if wait != 100*time.Millisecond {
		t.Errorf("Incorrect wait time: %v", wait)
	}
}

func TestBucketCapacity(t *testing.T) {
	b := NewBucket(10000, 10)
	time.Sleep(10 * time.Millisecond)
	n, _ := b.Take(10)
	if n != 10 {
		t.Errorf("Tokens got %d, expected 10", n)
	}
	n, _ = b.Take(10)
	if n != 0 {
		t.Errorf("Tokens got %d, expected 0", n)
	}
	time.Sleep(10 * time.Millisecond)
	n, _ = b.Take(100)
	if n != 10 {
		t.Errorf("Tokens got %d, expected 10", n)
	}
}

func TestBucketTakeExactly(t *testing.T) {
	b := NewBucket(1000, 50)
	begin := time.Now()
	b.TakeExactly(100)
	took := time.Now().Sub(begin)
	if took < 50*time.Millisecond || 60*time.Millisecond < took {
		t.Errorf("Invalid wait time: %v", took)
	}
}
