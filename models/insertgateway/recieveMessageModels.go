package insertgateway

type Request struct{
	Key string
	Value string
	Topic string
}

type OkResponse struct{
	message string
}

type ErrorResponse struct{
	message string
}