package service

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V string `json:"v"`
}

type LowercaseRequest struct {
	S string `json:"s"`
}

type LowercaseResponse struct {
	V string `json:"v"`
}

type CountRequest struct {
	S string `json:"s"`
}

type CountResponse struct {
	V int `json:"v"`
}
