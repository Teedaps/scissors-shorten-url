package routes

import(
    "github.com/Teedaps/shorten-url.git/database"
    "github.com/go-redis/redis/v8"
		"github.com/gofibre/fibre/v2"
)


func ResolveURL(c *fibre.Ctx) error{

	url := c.params("url")

	r := database.createClient(0)
  defer r.close()

	value, err := r.Get(database.Ctx, Url).Result()
	if err == redis.Nil{
		  return c.status(fibre.statusNotfound).Json(fibre.map{
				"error":"short not found in the database",
			})
	} else if err!=nil {
		  return c.status(fibre.statusInternalServerError).JSON(fibre.Map{
			    "error":"cannot connect to DB"
	       })
		}
  
    rInr := database.CreateClient(1)
		defer rInr.close()

		_= rInr.Incr(database.ctx, "counter")

		return c.Redirect(value, 301)

}
