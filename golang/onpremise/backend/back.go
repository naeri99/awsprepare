package backend

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "time"
    "operation/login"
    "operation/storage"
    "operation/master"
    "operation/worker"
    "operation/controller"
)

func generateCluster(){
    clientset:= login.GetClient()
    firstprocess := storage.Deploying(clientset)
    secondprocess := master.DeployingMaster(clientset,firstprocess)
    thirdprocess := worker.DeployingWorker(clientset,secondprocess )
    fourprocess := controller.DeployingController(clientset,thirdprocess )
    // Wait for the process to complete
    <- fourprocess
    close(fourprocess)
}

func deleteCluster(){
    clientset:= login.GetClient()
    firstdelete := controller.DeletingController(clientset)
    seconddelete := worker.DeletingWorker(clientset, firstdelete )
    thirddelete := master.DeletingMaster(clientset,seconddelete )
    fourdelete := storage.DeletingStorage(clientset, thirddelete)
    <-fourdelete
    close(fourdelete)
}


type RequestBody struct {
    Message string `json:"message"`
}

func StartRouter() {
    // Create a new Gin router
    router := gin.Default()

    // Apply the CORS middleware to the router
    router.Use(cors.New(cors.Config{
        AllowAllOrigins:  true, // Allows all origins
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    // Define a GET route
    router.GET("/", func(c *gin.Context) {
        c.String(200, "Hello, this is a GET request with CORS enabled for all origins!\n")
    })

    router.GET("/create", func(c *gin.Context) {
	generateCluster()
        c.String(200, "create\n")
    })

    router.GET("/delete", func(c *gin.Context) {
	deleteCluster()
        c.String(200, "delete\n")
    })


    // Define a POST route
    router.POST("/", func(c *gin.Context) {
        var requestBody RequestBody

        // Bind the JSON body to the struct
        if err := c.BindJSON(&requestBody); err != nil {
            c.String(400, "Invalid request body")
            return
        }

        c.String(200, "Received POST request with message: %s", requestBody.Message)
    })

    // Start the server on port 8080
    router.Run(":8070")
}

