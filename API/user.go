package API

import (
	service "agatra/Service"
	"agatra/middleware"
	"agatra/model"
	"net/http"
	"strconv"
	"time"

	// "errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserAPI interface {
	Login(u *gin.Context)
	Register(u *gin.Context)
	Logout(u *gin.Context)
	AddUser(u *gin.Context)
	UpdateUser(u *gin.Context)
	DeleteUser(u *gin.Context)
	GetUserByID(u *gin.Context)
	GetUserList(u *gin.Context)
	GetPrivileged(u *gin.Context)
	Profile(u *gin.Context)
}

type userAPI struct {
	userService service.UserService
	sessionService service.SessionService
}

func NewUserAPI(userService service.UserService, sessionService service.SessionService) *userAPI{
	return &userAPI{userService, sessionService}
}

func (ua *userAPI) Login(u *gin.Context) {
	var user model.User_login
	if err := u.BindJSON(&user); err != nil {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}
	if user.Username == "" || user.Password == "" {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("login data is empty"))
		return
	}
	dbUser, _ := ua.userService.GetByName(user.Username)
	if dbUser.Name == "" || dbUser.ID == 0 {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("user not found"))
		return
	}
	if !middleware.CheckPasswordHash(user.Password, dbUser.Password) {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("wrong email or password"))
		return
	}
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claims{
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(model.JwtKey)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.NewErrorResponse("error signing claims"))
		return
	}
	session := model.Session{
		Token:  tokenString,
		Email:  dbUser.Email,
		Expiry: expirationTime,
	}
	_, err = ua.sessionService.SessionAvailEmail(session.Email)
	if err != nil {
		err = ua.sessionService.AddSessions(session)
	} else {
		err = ua.sessionService.UpdateSessions(session)
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return model.JwtKey, nil
	})

	if err != nil {
		u.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	if !token.Valid {
		u.JSON(http.StatusUnauthorized, model.NewErrorResponse("invalid token"))
		return
	}

	u.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"data": gin.H{
			"apiKey": tokenString,
			"user":gin.H{
				"id":dbUser.ID,
				"username":dbUser.Name,
				"email":dbUser.Email,
				"role":dbUser.Role,
			},
		},
	})
}

func (ua *userAPI) Register(u *gin.Context) {
	var user model.RegisterInput
	if err := u.BindJSON(&user); err != nil {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}
	if user.Email == "" || user.Password == "" || user.Username == "" {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}else if user.Password != user.Confirm_password{
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("password and confirm password doesn't match"))
		return
	}
	_, exists := ua.userService.GetByEmail(user.Email)
	if exists {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("email already exists"))
		return
	}

	hashedPw, err := middleware.HashPassword(user.Password)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.NewErrorResponse("Skill Issue at Hashing"))
	}

	var result model.User = model.User{
		Name:  user.Username,
		Email: user.Email,
		Password: hashedPw,
		Role: "member",
	}
	err = ua.userService.Store(&result)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.NewErrorResponse("Error Storing Data"))
		return
	}
	u.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (ua *userAPI) Logout(u *gin.Context) {
	u.JSON(http.StatusOK, model.NewSuccessResponse("logout success"))
}

func (ua *userAPI) AddUser(u *gin.Context) {
	var newUser model.User
	if err := u.ShouldBindJSON(&newUser); err != nil {
		u.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" {
		u.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	err := ua.userService.Store(&newUser)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	u.JSON(http.StatusOK, model.SuccessResponse{Message: "add User success"})
}

func (ua *userAPI) UpdateUser(u *gin.Context) {
	var newUser model.User
	if err := u.ShouldBindJSON(&newUser); err != nil {
		u.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	UserID, err := strconv.Atoi(u.Param("id"))
	if err != nil {
		u.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid User ID"})
		return
	}
	err = ua.userService.Update(UserID, newUser)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	u.JSON(http.StatusOK, model.SuccessResponse{Message: "User update success"})
}

func (ua *userAPI) DeleteUser(u *gin.Context) {
	UserID, err := strconv.Atoi(u.Param("id"))
	if err != nil {
		u.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid User ID"})
		return
	}
	err = ua.userService.Delete(UserID)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	u.JSON(http.StatusOK, model.SuccessResponse{Message: "User delete success"})
}

func (ua *userAPI) GetUserByID(u *gin.Context) {
	UserID, err := strconv.Atoi(u.Param("id"))
	if err != nil {
		u.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid User ID"})
		return
	}

	User, err := ua.userService.GetByID(UserID)
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.UserResponse
	result.User = *User 
	result.Message = "User with ID " + strconv.Itoa(UserID) + " Found"

	u.JSON(http.StatusOK, result)
}

func (ua *userAPI) GetUserList(u *gin.Context) {
	name := u.Query("name")
	if name != ""{
		User, err := ua.userService.SearchName(name)
		if err != nil {
			u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var result model.UserArrayResponse
		result.Users = User 
		result.Message = "Getting All Users From Search Result Success"

		u.JSON(http.StatusOK, result)
	}else{
		User, err := ua.userService.GetList()
		if err != nil {
			u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var result model.UserArrayResponse
		result.Users = User 
		result.Message = "Getting All Users Success"

		u.JSON(http.StatusOK, result)
	}
	
}

func (ua *userAPI) GetPrivileged(u *gin.Context) {
	User, err := ua.userService.GetPrivileged()
	if err != nil {
		u.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.UserArrayResponse
	result.Users = User 
	result.Message = "Getting All Privileged Users Success"

	u.JSON(http.StatusOK, result)
}

func (ua *userAPI) Profile(u *gin.Context) {
	email := u.Keys["email"].(string)

	compare, boo := ua.userService.GetByEmail(email)
	if !boo {
		u.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Trouble finding user"})
	}

	u.JSON(http.StatusOK, gin.H{
		"message": "Get User Profile Success",
		"data": gin.H{
			"id":compare.ID,
			"username":compare.Name,
			"email":compare.Email,
			"role":compare.Role,
		},
	})
}

