package entity

import "time"

type (
	TenderID          string
	TenderName        string
	TenderDescription string
	TenderServiceType string
	TenderStatus      string
	OrganizationID    string
	TenderVersion     int32
)

const (
	TenderServiceTypeUndefined    TenderServiceType = ""
	TenderServiceTypeConstruction TenderServiceType = "Construction"
	TenderServiceTypeDelivery     TenderServiceType = "Delivery"
	TenderServiceTypeManufacture  TenderServiceType = "Manufacture"
)

func NewTenderServiceType(serviceType string) TenderServiceType {
	switch serviceType {
	case string(TenderServiceTypeConstruction):
		return TenderServiceTypeConstruction
	case string(TenderServiceTypeDelivery):
		return TenderServiceTypeDelivery
	case string(TenderServiceTypeManufacture):
		return TenderServiceTypeManufacture
	default:
		return TenderServiceTypeUndefined
	}
}

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

type Tender struct {
	ID             TenderID          `json:"id"`
	Name           TenderName        `json:"name"`
	Description    TenderDescription `json:"description"`
	ServiceType    TenderServiceType `json:"serviceType"`
	Status         TenderStatus      `json:"status"`
	OrganizationID OrganizationID    `json:"organizationId"`
	Version        TenderVersion     `json:"version"`
	CreatedAt      time.Time         `json:"createdAt"`
}

func (r *Tender) Apply(update *TenderUpdate) *Tender {
	if update.Name != "" {
		r.Name = TenderName(update.Name)
	}

	if update.Description != "" {
		r.Description = TenderDescription(update.Description)
	}

	if update.ServiceType != "" {
		r.ServiceType = TenderServiceType(update.ServiceType)
	}

	return r
}

type RequestTender struct {
	Tender   `json:",inline"`
	Username string `json:"creatorUsername"`
}

type TenderUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}
