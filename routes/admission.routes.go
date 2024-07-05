package routes

import (
	"net/http"
	"strconv"

	"easyauthapi/controllers"
	"easyauthapi/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func AdmissionRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	Admissions := router.Group("/admission", handlers...)
	{
		//==================================
		call := controllers.AdmissionController{}
		ctrl := call.NewController()
		//==================================

		Admissions.POST(
			"",
			validators.IsUser(),
			validators.CreateAdmissionValidator(),
			func(c *gin.Context) {
				ctrl.Create(c)
			},
		)

		//==================================

		//menggunakan Param dan bisa QueryArray
		Admissions.GET(
			"",
			validators.IsUser(),
			validators.GetAdmissionsValidator(),
			func(c *gin.Context) {
				modeStr := c.Query("m")
				if modeStr != "" {
					if mode, err := strconv.ParseBool(modeStr); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean value for 's'"})
						return
					} else if mode {
						ctrl.GetOne(c, "", "")
						return
					}
				}
				ctrl.Get(c, "", "")
			},
		)

		//hanya menggunakan QueryArray
		Admissions.GET(
			"/:optionparam/:optionsrcparam",
			validators.IsUser(),
			validators.GetAdmissionsValidator(),
			func(c *gin.Context) {
				modeStr := c.Query("m")
				optionparam := c.Param("optionparam")
				optionsrcparam := c.Param("optionsrcparam")
				if modeStr != "" {
					if mode, err := strconv.ParseBool(modeStr); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid boolean value for 's'"})
						return
					} else if mode {
						ctrl.GetOne(c, optionparam, optionsrcparam)
						return
					}
				}
				ctrl.Get(c, optionparam, optionsrcparam)
			},
		)

		//==================================
		/*
		   - lebih baik [-deleteOne-] dan [-updateOne-]
		   - karena [-updateMultiple-] dan [-deleteMultiple-] memiliki resiko
		   - jika ingin tetap [-Multiple-] maka di setup looping request pada [-Frontend-]
		*/
		//==================================

		Admissions.PUT(
			"",
			validators.IsUser(),
			validators.PathIdValidator(),
			validators.UpdateAdmissionValidator(),
			func(c *gin.Context) {
				ctrl.UpdateOne(c, "", "", false)
			},
		)

		Admissions.PUT(
			"/:optionparam/:optionsrcparam",
			validators.IsUser(),
			validators.PathIdValidator(),
			validators.UpdateAdmissionValidator(),
			func(c *gin.Context) {
				optionparam := c.Param("optionparam")
				optionsrcparam := c.Param("optionsrcparam")
				ctrl.UpdateOne(c, optionparam, optionsrcparam, true)
			},
		)

		//==================================

		Admissions.DELETE(
			"",
			validators.IsUser(),
			validators.PathIdValidator(),
			func(c *gin.Context) {
				ctrl.DeleteOne(c, "", "", false)
			},
		)

		Admissions.DELETE(
			"/:optionparam/:optionsrcparam",
			validators.IsUser(),
			validators.PathIdValidator(),
			func(c *gin.Context) {
				optionparam := c.Param("optionparam")
				optionsrcparam := c.Param("optionsrcparam")
				ctrl.DeleteOne(c, optionparam, optionsrcparam, true)
			},
		)

		//==================================
	}
}
