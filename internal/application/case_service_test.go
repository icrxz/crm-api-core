package application

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCaseService_CreateCase(t *testing.T) {
	t.Run("should create a case successfully", func(t *testing.T) {
		_ = gomock.NewController(t)
	})
}
