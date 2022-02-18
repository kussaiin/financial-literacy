package routes

import 	(
	"github.com/gin-gonic/gin"
	"github.com/user/financial-literacy/controllers"
)

func HandleRequests(incomingRequest *gin.Engine)  {
	incomingRequest.GET("/finance/transactions", controllers.GetTransactions())
	incomingRequest.GET("/finance/transactions/:transaction_id", controllers.GetTransaction())
	incomingRequest.POST("/finance/transactions", controllers.CreateTransaction())
	incomingRequest.PUT("/finance/transactions/:transaction_id", controllers.UpdateTransaction())
	incomingRequest.DELETE("/finance/transactions/:transaction_id", controllers.DeleteTransaction())

	incomingRequest.GET("/finance/categories", controllers.GetCategories())
	incomingRequest.GET("/finance/categories/:category_id", controllers.GetCategory())
	incomingRequest.POST("/finance/categories", controllers.CreateCategory())
	incomingRequest.PUT("/finance/categories/:category_id", controllers.UpdateCategory())
	incomingRequest.DELETE("/finance/categories/:category_id", controllers.DeleteCategory())
}
