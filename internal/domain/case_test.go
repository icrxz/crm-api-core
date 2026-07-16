package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCaseFull(t *testing.T) {
	crmCase := Case{
		CaseID:  "case-1",
		OwnerID: "user-1",
		QueueID: "queue-1",
	}
	queue := Queue{QueueID: "queue-1", Name: "SP Mobile", Category: MobileQueueCategory}

	caseFull := NewCaseFull(crmCase, nil, nil, Product{}, Customer{}, Partner{}, Contractor{}, queue)

	assert.Equal(t, "case-1", caseFull.CaseID)
	assert.Equal(t, queue, caseFull.Queue)
}

func TestNewCaseFull_EmptyQueue(t *testing.T) {
	crmCase := Case{CaseID: "case-1"}

	caseFull := NewCaseFull(crmCase, nil, nil, Product{}, Customer{}, Partner{}, Contractor{}, Queue{})

	assert.Equal(t, Queue{}, caseFull.Queue)
}
