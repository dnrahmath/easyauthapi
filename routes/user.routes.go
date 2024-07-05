package routes

import (
	"easyauthapi/controllers"
	"easyauthapi/middlewares/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := router.Group("/users", handlers...)
	{
		//==================================
		call := controllers.UserController{}
		ctrl := call.NewController()
		//==================================

		users.POST(
			"",
			validators.IsUser(),
			validators.CreateUserValidatorByAdm(),
			func(c *gin.Context) {
				ctrl.Create(c)
			},
		)

		//==================================

		//menggunakan Param dan bisa QueryArray
		users.GET(
			"",
			validators.IsAdmin(),
			validators.GetUsersValidator(),
			func(c *gin.Context) {
				modeStr := c.Query("m")
				if modeStr != "" {
					if mode, err := strconv.ParseBool(modeStr); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean value for 's'"})
						return
					} else if mode {
						ctrl.GetOne(c)
						return
					}
				}
				ctrl.Get(c)
			},
		)

		//hanya menggunakan QueryArray
		users.GET(
			"/:optionparam/:optionsrcparam",
			validators.IsAdmin(),
			validators.GetUsersValidator(),
			func(c *gin.Context) {
				modeStr := c.Query("m")
				// optionparam := c.Param("optionparam")
				// optionsrcparam := c.Param("optionsrcparam")
				if modeStr != "" {
					if mode, err := strconv.ParseBool(modeStr); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean value for 's'"})
						return
					} else if mode {
						ctrl.GetOne(c)
						return
					}
				}
				ctrl.Get(c)
			},
		)

		//==================================
		/*
		   - lebih baik [-deleteOne-] dan [-updateOne-]
		   - karena [-updateMultiple-] dan [-deleteMultiple-] memiliki resiko
		   - jika ingin tetap [-Multiple-] maka di setup looping request pada [-Frontend-]
		*/
		//==================================

		users.PUT(
			"",
			validators.IsAdmin(),
			validators.PathIdValidator(),
			validators.UpdateUserValidatorByAdm(),
			func(c *gin.Context) {
				ctrl.UpdateOne(c)
			},
		)

		users.PUT(
			"/:optionparam/:optionsrcparam",
			validators.IsAdmin(),
			validators.PathIdValidator(),
			validators.UpdateUserValidatorByAdm(),
			func(c *gin.Context) {
				// optionparam := c.Param("optionparam")
				// optionsrcparam := c.Param("optionsrcparam")
				ctrl.UpdateOne(c)
			},
		)

		//==================================

		users.DELETE(
			"",
			validators.IsAdmin(),
			validators.PathIdValidator(),
			func(c *gin.Context) {
				ctrl.DeleteOne(c)
			},
		)

		users.DELETE(
			"/:optionparam/:optionsrcparam",
			validators.IsAdmin(),
			validators.PathIdValidator(),
			func(c *gin.Context) {
				// optionparam := c.Param("optionparam")
				// optionsrcparam := c.Param("optionsrcparam")
				ctrl.DeleteOne(c)
			},
		)

		//==================================
	}
}
