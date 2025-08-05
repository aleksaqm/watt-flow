package util_test

import (
	"testing"
	"watt-flow/model"
	"watt-flow/util"
)

func TestGenerateMonthlyBillEmail(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		bill *model.Bill
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _ := util.GenerateMonthlyBillEmail(tt.bill)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GenerateMonthlyBillEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
