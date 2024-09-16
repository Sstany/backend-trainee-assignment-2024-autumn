package entity

import "time"

type (
	BidId          string
	BidName        string
	BidDescription string
	BidStatus      string
	BidAuthorType  string
	BidAuthorId    string
	BidVersion     int32

	BidReviewIdString          string
	BidReviewDescriptionstring string
)

type Bid struct {
	ID          BidId          `json:"id"`
	Name        BidName        `json:"name"`
	Description BidDescription `json:"description"`
	Status      BidStatus      `json:"organization" bindig:"status"`
	TenderID    TenderID       `json:"tenderId" `
	AuthorType  BidAuthorType  `json:"authorType"`
	AuthorID    BidAuthorId    `json:"AuthorId"`
	Version     BidVersion     `json:"version"`
	CreatedAt   time.Time
}

type BidReview struct {
	ID          BidReviewIdString          `json:"id"`
	Description BidReviewDescriptionstring `json:"description"`
	CreatedAt   string                     `json:"createdAt"`
}

type BidReviewDecision string

const (
	BidReviewApproved BidReviewDecision = "Approved"
	BidReviewRejected BidReviewDecision = "Rejected"
)

const (
	BidAuthorOrganization BidAuthorType = "Organization"
	BidAuthorUser         BidAuthorType = "User"
)
