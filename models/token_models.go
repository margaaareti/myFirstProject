package models

import uuid2 "github.com/gofrs/uuid"

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   uuid2.UUID
	RefreshUuid  uuid2.UUID
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUUID  string
	RefreshUUID string
	UserId      uint64
}
