package queries

type CustomResponse struct {
	Status  int
	Message string
}

var (
	ErrDefault             = CustomResponse{Status: 500, Message: "Default Error"}
	ErrInternalServerError = CustomResponse{Status: 500, Message: "Internal Server Error"}
	ErrCreateUser          = CustomResponse{Status: 400, Message: "Create User Error"}
	ErrGetUser             = CustomResponse{Status: 404, Message: "User Id Not Found"}
	ErrLoginUser           = CustomResponse{Status: 404, Message: "User Login Failed"}
	SuccessCreateUser      = CustomResponse{Status: 200, Message: "Create User Success"}
	SuccessDeleteUser      = CustomResponse{Status: 200, Message: "Delete User Success"}
	SuccessGetUser         = CustomResponse{Status: 200, Message: "Get User Name Success"}
	SuccessLoginUser       = CustomResponse{Status: 200, Message: "Login User Success"}
)
