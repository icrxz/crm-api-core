package domain

type Billing struct {
	Key             string
	Option          string
	Type            BillingType
	Name            string
	IsSameFromOwner bool
}

type BillingType string

const (
	PIX BillingType = "PIX"
)
