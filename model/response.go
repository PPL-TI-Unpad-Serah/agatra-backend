package model

type ErrorResponse struct {
	Error 		string `json:"error"`
}

type SuccessResponse struct {
	Message 	string 	
}

type CityResponse struct{
	Message 	string 	
	City		City
}

type CityArrayResponse struct{
	Message 	string 	
	Cities		[]City
}

type CenterResponse struct{
	Message string 	
	Center 	Center	
}

type CenterArrayResponse struct{
	Message 	string 	
	Centers 	[]Center	
}

type MachineResponse struct{
	Message 	string 	
	Machine 	Machine	
}

type MachineArrayResponse struct{
	Message 	string 	
	Machines 	[]Machine	
}

type LocationResponse struct{
	Message 	string 	
	Location 	Location	
}

type LocationArrayResponse struct{
	Message 	string 	
	Locations 	[]Location	
}

type LocationRangeResponse struct{
	Message 	string 	
	Locations 	[]Location_range
}

type VersionResponse struct{
	Message 	string 	
	Version 	Version	
}

type VersionArrayResponse struct{
	Message 	string 	
	Versions 	[]Version	
}

type TitleResponse struct{
	Message 	string 	
	Title 	Title	
}

type TitleArrayResponse struct{
	Message 	string 	
	Titles 	[]Title	
}

type UserResponse struct{
	Message 	string 	
	User 	User	
}

type UserArrayResponse struct{
	Message 	string 	
	Users 	[]User	
}

type SessionResponse struct{
	Message 	string 	
	Session 	Session	
}

type SessionArrayResponse struct{
	Message 	string 	
	Sessions 	[]Session	
}

type LoginResponse struct{
	ApiKey 		string
	User		User
}

type RegisterInput struct{
	Username			string
	Email				string
	Password			string
	Confirm_password	string
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}

func NewSuccessResponse(msg string) SuccessResponse {
	return SuccessResponse{
		Message: msg,
	}
}