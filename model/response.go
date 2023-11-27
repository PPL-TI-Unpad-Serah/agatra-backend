package model

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type CityResponse struct {
	Message string `json:"message"`
	City    City   `json:"data"`
}

type CityArrayResponse struct {
	Message string `json:"message"`
	Cities  []City `json:"data"`
}

type CenterResponse struct {
	Message string `json:"message"`
	Center  Center `json:"data"`
}

type CenterArrayResponse struct {
	Message string   `json:"message"`
	Centers []Center `json:"data"`
}

type MachineResponse struct {
	Message string  `json:"message"`
	Machine Machine `json:"data"`
}

type MachineArrayResponse struct {
	Message  string    `json:"message"`
	Machines []Machine `json:"data"`
}

type LocationResponse struct {
	Message  string   `json:"message"`
	Location Location `json:"data"`
}

type LocationArrayResponse struct {
	Message   string     `json:"message"`
	Locations []Location `json:"data"`
}

type LocationRangeResponse struct {
	Message   string           `json:"message"`
	Locations []Location_range `json:"data"`
}

type VersionResponse struct {
	Message string  `json:"message"`
	Version Version `json:"data"`
}

type VersionArrayResponse struct {
	Message  string    `json:"message"`
	Versions []Version `json:"data"`
}

type TitleResponse struct {
	Message string `json:"message"`
	Title   Title  `json:"data"`
}

type TitleArrayResponse struct {
	Message string  `json:"message"`
	Titles  []Title `json:"data"`
}

type UserResponse struct {
	Message string `json:"message"`
	User    User   `json:"data"`
}

type UserArrayResponse struct {
	Message string `json:"message"`
	Users   []User `json:"data"`
}

type SessionResponse struct {
	Message string  `json:"message"`
	Session Session `json:"data"`
}

type SessionArrayResponse struct {
	Message  string    `json:"message"`
	Sessions []Session `json:"data"`
}

type LoginResponse struct {
	ApiKey string    	`json:"apiKey"`
	User   User		    `json:"user"`
}

type RegisterInput struct {
	Username         string
	Email            string
	Password         string
	Confirm_password string
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