package context_test

import "testing"

var testPasswordContextCases = []struct {
	name     string
	password string
}{
	{
		name:     "Data 1",
		password: "1234567890",
	},
	{
		name:     "Data 2",
		password: "qwertyuiop",
	},
	{
		name:     "Data 3",
		password: "1234567890qwertyuiop",
	},
}

func TestPasswordContext(t *testing.T) {
	for _, tc := range testPasswordContextCases {
		t.Run(tc.name, func(t *testing.T) {
			hashValue, err := dmCtx.Password().Hash(tc.password)
			if err != nil {
				t.Fatalf("Failed to hash the password: %v", err)
			}

			correct, err := dmCtx.Password().Verify(tc.password, hashValue)
			if err != nil {
				t.Fatalf("Failed to verify the password: %v", err)
			}
			if !correct {
				t.Fatalf("Failed to verify the hashed value which should be success: %v", err)
			}

			incorrect, err := dmCtx.Password().Verify(tc.password+"fake", hashValue)
			if err != nil {
				t.Fatalf("Failed to verify the password: %v", err)
			}
			if incorrect {
				t.Fatalf("Failed to verify the hashed value which should be fail: %v", err)
			}
		})
	}
}
