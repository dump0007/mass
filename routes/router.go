package routes

import(
	"mass/controllers"
	"github.com/gin-gonic/gin"
)


func Routes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("sign-in",controllers.SignIn())
	incomingRoutes.POST("upload-excel",controllers.SendMsg())

	incomingRoutes.POST("estimate-gas",controllers.EstimateGas())
	incomingRoutes.POST("send-msg",controllers.SendMsg())
	incomingRoutes.POST("execute-txn",controllers.ExecuteTxns())
	incomingRoutes.POST("get-orders",controllers.GetOrders())

	incomingRoutes.POST("upload-nft-address",controllers.CreateNftOrder())

	// incomingRoutes.POST("setup-excel",controllers.SetupExcel())
}

