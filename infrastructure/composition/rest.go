package composition

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aboglioli/big-brother/composition"
	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/infrastructure/events"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// StartREST starts API server
func StartREST() {
	// Dendencies resolution
	eventManager, err := events.GetManager()
	if err != nil {
		log.Fatal(err)
	}

	compositionRepository, rawErr := composition.NewRepository()
	if rawErr != nil {
		log.Fatal(err)
	}

	compositionService := composition.NewService(compositionRepository, eventManager)

	// Create context
	rest := &RESTContext{
		compositionService: compositionService,
	}

	// Start Gin server
	c := config.Get()
	server := gin.Default()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(cors.New(corsConfig))

	// Static server
	server.Use(static.Serve("/", static.LocalFile("www", false)))

	// Routes
	server.GET("/v1/composition/:compositionId", rest.GetByID)
	server.POST("/v1/composition", rest.Post)
	server.PUT("/v1/composition/:compositionId", rest.Put)
	server.DELETE("/v1/composition/:compositionId", rest.Delete)

	server.Run(fmt.Sprintf(":%d", c.Port))
}

type RESTContext struct {
	compositionService composition.Service
}

// GetByID finds a Composition by ID
/**
* @api {get} /v1/comosition/:compositionId GetByID
* @apiName Find by ID
* @apiGroup Composition
*
* @apiParam {String} compositionId Composition ID
*
* @apiDescription Gets a composition by ID
*
* @apiSuccessExample {json} Response
* HTTP/1.1 200 OK
* {
*   "composition": {
*     "id": "9dc9c429b9aa2a3c82801007",
*     "name": "Comp 7",
*     "cost": 475.75,
*     "unit": {
*       "quantity": 3,
*       "unit": "u"
*     },
*     "stock": {
*       "quantity": 1,
*       "unit": "u"
*     },
*     "dependencies": [
*       {
*         "of": "9dc9c429b9aa2a3c82801005",
*         "quantity": {
*           "quantity": 2,
*           "unit": "u"
*         },
*         "subvalue": 82
*       },
*       {
*         "of": "9dc9c429b9aa2a3c82801006",
*         "quantity": {
*           "quantity": 1.5,
*           "unit": "u"
*         },
*         "subvalue": 393.75
*       }
*     ],
*     "autoupdateCost": true,
*     "enabled": true,
*     "validated": true,
*     "usesUpdatedSinceLastChange": true,
*     "createdAt": "2019-11-11T22:15:59.301Z",
*     "updatedAt": "2019-11-15T01:35:19.024Z"
*  }
* }
 */
func (r *RESTContext) GetByID(c *gin.Context) {
	compID := c.Param("compositionId")

	comp, err := r.compositionService.GetByID(compID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"composition": comp,
	})
}

// Post creates a new Composition
/**
* @api {post} /v1/composition Create
* @apiName Post
* @apiGroup Composition
*
* @apiParam {String} [name=""] Name
* @apiParam {String} [cost=0] Initial cost
* @apiParam {Quantity} unit Composition base unit
* @apiParam {Quantity} [stock] Stock quantity. Same units as "unit".
* @apiParam {[]Dependency} [dependencies] Dependencies: foreign key "of", and "quantity".
* @apiParam {Boolean} [autoupdateCost=true] Auto update cost based on dependencies.
*
* @apiDescription Creates a new Composition. "id" is optional but it can be
* specified. In case "id" was not specified, a new ObjectID would be assigned.
* The only required field in body is "unit". "cost" can be set to any value
* greater than 0 (default: 0). If "dependencies" are added, "cost" will be
* calculated automatically based on these, unless "autoupdateCost" is set to
* "false".
*
* @apiExample {json} Body
* {
*   "composition": {
*     "id": "9dc9c429b9aa2a3c82801007",
*     "name": "Comp 7",
*     "cost": 0,
*     "unit": {
*       "quantity": 3,
*       "unit": "u"
*     },
*     "stock": {
*       "quantity": 1,
*       "unit": "u"
*     },
*     "dependencies": [
*       {
*         "of": "9dc9c429b9aa2a3c82801005",
*         "quantity": {
*           "quantity": 2,
*           "unit": "u"
*         },
*         "subvalue": 82
*       },
*       {
*         "of": "9dc9c429b9aa2a3c82801006",
*         "quantity": {
*           "quantity": 1.5,
*           "unit": "u"
*         },
*         "subvalue": 393.75
*       }
*     ],
*     "autoupdateCost": true,
*   }
* }
*
* @apiSuccessExample {json} Response
* HTTP/1.1 200 OK
*
* {
*   "composition": {
*     "id": "9dc9c429b9aa2a3c82801007",
*     "name": "Comp 7",
*     "cost": 475.75,
*     "unit": {
*       "quantity": 3,
*       "unit": "u"
*     },
*     "stock": {
*       "quantity": 1,
*       "unit": "u"
*     },
*     "dependencies": [
*       {
*         "of": "9dc9c429b9aa2a3c82801005",
*         "quantity": {
*           "quantity": 2,
*           "unit": "u"
*         },
*         "subvalue": 82
*       },
*       {
*         "of": "9dc9c429b9aa2a3c82801006",
*         "quantity": {
*           "quantity": 1.5,
*           "unit": "u"
*         },
*         "subvalue": 393.75
*       }
*     ],
*     "autoupdateCost": true,
*     "enabled": true,
*     "validated": false,
*     "usesUpdatedSinceLastChange": true,
*     "createdAt": "2019-11-11T22:15:59.301Z",
*     "updatedAt": "2019-11-15T01:35:19.024Z"
*   },
*	"status": "CREATED"
* }
 */
func (r *RESTContext) Post(c *gin.Context) {
	var body composition.CreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	comp, err := r.compositionService.Create(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "CREATED",
		"composition": comp,
	})
}

// Put updates a new Composition
/**
* @api {put} /v1/composition/:compositionId Update
* @apiName Put
* @apiGroup Composition
*
* @apiParam {String} compositionId Composition ID
*
* @apiParam {String} [name] Name
* @apiParam {String} [cost] Initial cost
* @apiParam {Quantity} [unit] Composition base unit. Cannot be changed.
* @apiParam {Quantity} [stock] Stock quantity. Same units as "unit".
* @apiParam {[]Dependency} [dependencies] Dependencies: foreign key "of", and "quantity".
* @apiParam {Boolean} [autoupdateCost] Auto update cost based on dependencies.
*
* @apiDescription Updates an existing Composition based on its ID. "cost" can
* be set to any value greater than 0 (default: 0). If "dependencies" are added,
* "cost" will be calculated automatically based on these, unless
* "autoupdateCost" is set to "false". All fields are optional. "unit" unit
* cannot be changed of type.
*
* @apiExample {json} Body
* {
*   "composition": {
*     "name": "Comp 8",
*     "cost": 125,
*     "unit": {
*       "quantity": 2,
*       "unit": "u"
*     },
*     "stock": {
*       "quantity": 0,
*       "unit": "u"
*     },
*     "dependencies": [
*       {
*         "of": "9dc9c429b9aa2a3c82801005",
*         "quantity": {
*           "quantity": 2,
*           "unit": "u"
*         },
*         "subvalue": 82
*       },
*     ],
*     "autoupdateCost": true,
*   }
* }
*
* @apiSuccessExample {json} Response
* HTTP/1.1 200 OK
*
* {
*   "composition": {
*     "id": "9dc9c429b9aa2a3c82801007",
*     "name": "Comp 8",
*     "cost": 475.75,
*     "unit": {
*       "quantity": 3,
*       "unit": "u"
*     },
*     "stock": {
*       "quantity": 1,
*       "unit": "u"
*     },
*     "dependencies": [
*       {
*         "of": "9dc9c429b9aa2a3c82801005",
*         "quantity": {
*           "quantity": 2,
*           "unit": "u"
*         },
*         "subvalue": 82
*       },
*       {
*         "of": "9dc9c429b9aa2a3c82801006",
*         "quantity": {
*           "quantity": 1.5,
*           "unit": "u"
*         },
*         "subvalue": 393.75
*       }
*     ],
*     "autoupdateCost": true,
*     "enabled": true,
*     "validated": true,
*     "usesUpdatedSinceLastChange": false,
*     "createdAt": "2019-11-11T22:15:59.301Z",
*     "updatedAt": "2019-11-15T01:35:19.024Z"
*   },
*	"status": "UPDATED"
* }
 */
func (r *RESTContext) Put(c *gin.Context) {
	compID := c.Param("compositionId")
	var body composition.UpdateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	comp, err := r.compositionService.Update(compID, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "UPDATED",
		"composition": comp,
	})
}

// Delete deletes an existing Composition
/**
* @api {delete} /v1/composition/:compositionId Delete
* @apiName Delete
* @apiGroup Composition
*
* @apiParam {String} compositionId Composition ID
*
* @apiDescription Deletes an existing Composition by ID. To be deleted it
* doesn't have to be used as dependency in another Composition. Soft deleting
* (set "enabled" to false).
*
* @apiSuccessExample {json} Response
* HTTP/1.1 200 OK
*
* {
*	"status": "DELETED"
* }
 */
func (r *RESTContext) Delete(c *gin.Context) {
	compID := c.Param("compositionId")

	if err := r.compositionService.Delete(compID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "DELETED",
	})
}
