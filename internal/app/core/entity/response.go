package entity

import "strconv"

type ResponseError struct {
	Reason string `json:"reason"`
}

type RequestLimitOffset struct {
	Limit  int
	Offset int
}

func ParseRequestLimitOffset(limit, offset string) *RequestLimitOffset {
	lo := new(RequestLimitOffset)

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil
		}
		lo.Limit = l
	}

	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return nil
		}
		lo.Offset = o
	}

	return lo
}
