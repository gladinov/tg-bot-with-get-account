package mapper

import (
	"bonds-report-service/internal/models/domain"
	"testing"

	factories "bonds-report-service/internal/models/domain/testing"

	"github.com/stretchr/testify/assert"
)

func TestMapOperationToOperationWithoutCustomTypes(t *testing.T) {
	op := factories.NewOperation()
	want := factories.NewOperationWithoutCustomTypes()

	got := MapOperationToOperationWithoutCustomTypes([]domain.Operation{op})

	assert.Equal(t, []domain.OperationWithoutCustomTypes{want}, got)
}
