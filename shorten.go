package routes


import(
	"os"
"time"
"strconv"
"github.com/Teedaps/shorten-url.git/database"
"github.com/Teedaps/shorten-url.git/helpers"
"github.com/go-redis/redis/v8"
"github.com/gofibre/fibre/v2"
"github.com/asaskevich/govalidator"
"github.com/google/uuid"

)

type request struct {
    URL             string             'json:"url"'
		Customshort     string             'json:"short"'
		Expiry          time.Duration      'json:"expiry"'
}

type response struct{
    URL               string            'json:"url"'
		Customshort       string            'json:"short"'
		Expiry            time.Duration     'json:"expiry"'
		XRateRemaining    int               'json:"rate_limit"'
		XRatelimitRest    time.Duration     'json:"rate_limit_reset"'
	}

	func shortenURL(c  *fibre.ctx) error {

		body := new(request)

		if err := c.Bodyparser(&body); err!= nil {
			 return c.status(fibre.statusBadRequest).JSON(fiber.map{"error":"cannot parse JSON"})
	}

	//implement rate limiting

  r2 := database.CreateClient(1)
	defer r2.close()
	val, err := r2.Get(database.Ctx, c.Ip()).Result
	if err == redis.Nil{
		_= r2.set(database.ctx, c.Ip(),os.Getenv("API_QUOTA"), 30*60*time.second).Err()
	} else{
		   val, _ = r2.Get(database.ctx, c.IP()).Result()
			 valInt, _ := strconv.Atoi(val)
			 if valInt <m 0 {
				   limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
					 return c.status(fiber.StatusServiceUnavailable).JSON(fibre.Map{
						   "error": "Rate limit exceeded",
							 "rate_limit_rest": limit / time.Nanosecond / time.Minute,
					 })
			 }
	}


	 //check if the input if an actual URL 

   if !govalidator.IsURL(body.URL){
       return c.status(fibre.statusBadRequest).JSON(fibre.map{"error":"Invalid URL"})
	}

	 //check for domain error

    if !helpers.RemoveDomainError(body.URL){
		    return c.status(fibre.statusServiceunavailable).JSON(fibre.map{"error":"you can't hack the system (:")})
	}

   //enforce https, SSL

   body.URL = helpers.EnforceHTTP(body.URL)

   var id string

	 if body.CustomShort == ""{
		   id = uuid.New().string()[:6]
	 } else {
		   id = body.CustomShort
	 }

	 r := database.CreateClient(0)
   defer r.close()

	 val, _ = r.Get(database.Ctx, id).Result()
  if val != ""{
		  return c.status(fibre.statusForbidden).JSON(fibre.Map{
				  "error": "URL custom short is already in use",
			})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}


	error = r.set(database.Ctx, id, body.URL,body.Expiry*360*time.second).Err()

  if err !=nil {
		  return c.status(fibre.statusInternalServerError).JSON(fibre.Map{
				 "err":"Unable to connect to server",
			})
	}

  resp := response{
		  URL:            body.URL,
			CustomShort:    "",
			Expiry:         body.Expiry,
			XRateRemaining: 10,
			XRateLimitReset:30,
	}

  r2.Decr(database.Ctx, c.IP())

  val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	response.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id 

	return c.Status(fibre.StatusOK).JSON(resp)
}