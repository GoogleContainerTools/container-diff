/*
Copyright 2017 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package differs

import (
	"testing"
)

// TestLockUnlock runs some lock-unlock cycles to make sure that close,
// subsequent locking and unlocking doesn't cause any issues.
func TestLockUnlock(t *testing.T) {
	for i := 0; i < 20; i++ {
		if err := lock(); err != nil {
			t.Errorf("[RPMDiff test] lock() %d failed: %s", i, err)
		}
		if err := unlock(); err != nil {
			t.Errorf("[RPMDiff test] unlock() %d failed: %s", i, err)
		}
	}
}

// TestGoRoutineRace creates a race condition among goroutines.  Race conditions
// among processes cannot be tested here as we can't fork.
func TestGoRoutineRace(t *testing.T) {
	wait := make(chan int)
	couldLock := false
	if err := lock(); err != nil {
		t.Errorf("[RPMDiff test] lock() failed: %s", err)
	}

	go func() {
		wait <- 1
		if err := lock(); err != nil {
			t.Errorf("[RPMDiff test] lock() failed: %s", err)
		}
		couldLock = true
		if err := unlock(); err != nil {
			t.Errorf("[RPMDiff test] unlock() failed: %s", err)
		}
		wait <- 1
	}()

	<-wait
	if couldLock {
		t.Errorf("Other goroutine locked although lock wasn't released")
	}
	if err := unlock(); err != nil {
		t.Errorf("[RPMDiff test] unlock() failed: %s", err)
	}
	<-wait
	if !couldLock {
		t.Errorf("Other goroutine didn't lock although lock was released")
	}
}
