package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// UserController represents a user controller
type UserController struct {
	userService service.UserService
}

// NewUserController returns a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (us *UserController) RegisterUser(c *gin.Context) {
	var user model.User

	// Bind JSON to user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Printf("bound user: %+v", user)

	// Call service to register, WalletID is already a *string
	createdUser, err := us.userService.RegisterUser(
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		*user.WalletID, // safe to pass directly
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return created user
	c.JSON(201, createdUser)
}



func (us *UserController) AuthenticateUser(c *gin.Context) {
	var user model.User

	// Bind JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	log.Printf("Authenticating user with email: %s", user.Email)

	// Authenticate the user using service
	authenticatedUser, err := us.userService.AuthenticateUser(user.Email, user.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed: " + err.Error()})
		return
	}

	// Return authenticated user
	c.JSON(200, authenticatedUser)
}


func (us *UserController) GetUserByID(c *gin.Context) {
	// Get the user ID from the request parameters
	userParam := c.Param("id")
	userID, err := uuid.FromString(userParam)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the user
	user, err := us.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the user in the response
	c.JSON(200, user)
}

func (us *UserController) ListUsers(c *gin.Context) {
	// Call the service to get the user
	user, err := us.userService.ListUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the user in the response
	c.JSON(200, user)

}

func (us *UserController) UpdateUser(c *gin.Context) {
	var user model.User

	// get the user ID from the request parameters
	userParam := c.Param("id")
	userID, err := uuid.FromString(userParam)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Bind JSON data from the request body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user.ID = userID
	// Call the service to update the user
	if err := us.userService.UpdateUser(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Return the updated user in the response
	c.JSON(200, user)
}

func (us *UserController) DeleteUser(c *gin.Context) {
    userParam := c.Param("id")
    userID, err := uuid.FromString(userParam)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    err = us.userService.DeleteUser(userID)
    if err != nil {
        if err.Error() == "user not found" {
            c.JSON(404, gin.H{"error": "user not found"})
        } else {
            c.JSON(500, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(200, gin.H{"message": "user deleted successfully"})
}

