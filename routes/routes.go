package routes

import (

	"github.com/kemboy-svg/investment/controllers"
	"github.com/labstack/echo/v4"
)


func Routes(e *echo.Echo) {
	user := controllers.DepositController{}


	e.GET("/Deposits/", user.SyncDeposits)
	
	


}