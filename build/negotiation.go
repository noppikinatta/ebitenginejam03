package build

type Negotiation struct {
	Size           PointF
	Vendors        []*Vendor
	Managers       []*Manager
	Money          int
	ApprovedEquips []*Equip
}
