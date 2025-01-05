package dto

type ResourceDetailCustomerDTO struct {
	ID      string            `json:"id"`
	AppName string            `json:"appName"`
	Version string            `json:"version"`
	Build   string            `json:"build"`
	Message string            `json:"message"`
	Data    DetailCustomerDTO `json:"data"`
}

type DetailCustomerDTO struct {
	ID                   uint                 `json:"id" example:"1"`
	Firstname            string               `json:"firstname" example:"Jhon"`
	Lastname             string               `json:"lastname" example:"Doe"`
	Email                string               `json:"email" example:"your@mail.com"`
	Phone                *string              `json:"phone" example:"081111116542"`
	VerifiedPhoneNumber  *string              `json:"verifiedPhoneNumber" example:"081111116542"`
	DateOfBirth          string               `json:"dateOfBirth" example:"1991-04-13"`
	IdentityNumber       *string              `json:"identityNumber" example:"3274031122110001"`
	MemberStatus         string               `json:"memberStatus" example:"Crew"`
	MemberExpiredDate    string               `json:"memberExpiredDate" example:"1991-04-13"`
	MemberPoint          *uint                `json:"memberPoint" example:"100"`
	ExpiryPoint          string               `json:"expiryPoint" example:"1991-04-13"`
	TotalCoupons         *uint                `json:"totalCoupons" example:"10"`
	Gender               string               `json:"gender" example:"male"`
	SubscribeNewsletter  bool                 `json:"isSubscribeNewsletter" example:"true"`
	IsActive             bool                 `json:"isActive" example:"true"`
	IsVerified           bool                 `json:"isVerified" example:"true"`
	CapillaryFraudStatus CapillaryFraudStatus `json:"capillaryFraudStatus"`
}

type CapillaryFraudStatus struct {
	StatusLabel         *string `json:"statusLabel" example:"Confirmed"`
	StatusLabelReason   *string `json:"statusLabelReason" example:"FRAUD"`
	StatusLastUpdatedOn *string `json:"statusLastUpdatedOn" example:"2024-10-28T13:19:47+05:30"`
}
