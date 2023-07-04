package main

import(
	"fmt"
	"log"
	"os"
  "github.com/gofibre/fibre/v2"
	"github.com/gofibre/fibre/v2/middleware/logger"
	"github.com/Teedaps/shorten-url.git/routes"
	"github.com/joho/godotenv"
)
      
func setupRoutes(app #fibre.App){
	 app.Get("/url", routes.ResolveURL)
	 app.post("/api/v1", routes.shortenURL)
}

func main(){
    error := godotenv.Load()

    if error!= nil {
			   fmt.println(err)
		}

    app : = fibre.New()

    app.use(logger.New())

    setupRoutes(app)

	log.fatal(app.Listen(os.Getenv("APP_PORT")))

	}