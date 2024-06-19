package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFatal_TestingMode(t *testing.T) {
	TestingMode(true)

	defer func() {
		if p := recover(); p != nil {
			// capture panic
			assert.NotNil(t, p)
			assert.Equal(t, "[Panic error]", fmt.Sprintf("%v", p))
		}
	}()
	Fatal("Panic error")
}

func TestFatalf_TestingMode(t *testing.T) {
	TestingMode(true)

	defer func() {
		if p := recover(); p != nil {
			// capture panic
			assert.NotNil(t, p)
			assert.Equal(t, "Panic error: [cannot handle]", fmt.Sprintf("%v", p))
		}
	}()
	Fatalf("Panic error: %v", "cannot handle")
}

func TestFatalw_TestingMode(t *testing.T) {
	TestingMode(true)

	defer func() {
		if p := recover(); p != nil {
			// capture panic
			assert.NotNil(t, p)
			assert.Equal(t, "Panic error", fmt.Sprintf("%v", p))
		}
	}()
	Fatalw("Panic error", "key", "value")
}
