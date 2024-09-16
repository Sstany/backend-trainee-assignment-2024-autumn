package service

import "errors"

var (
	ErrUserNotExists   = errors.New("user not exists")
	ErrNotEnoughRights = errors.New("user does not have enough rights")

	ErrTenderNotFound      = errors.New("tender not found")
	ErrWrongInputFormat    = errors.New("wrong format or parameters")
	ErrTenderOrBidNotFound = errors.New("tender or bid not found")
	ErrBidNotFound         = errors.New("bid not found")
)
