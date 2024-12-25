package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Define server settings.

  // Converts the string environment variable into an integer.

  readTimeoutSecondsCount, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
  if err != nil || readTimeoutSecondsCount == 0 {
    readTimeoutSecondsCount = 30 
  }
	// Return Fiber configuration.
	return fiber.Config{
		EnableTrustedProxyCheck: true,
   
    // Converts the integer readTimeoutSecondsCount into a time.Duration value (in seconds)
		ReadTimeout:             time.Second * time.Duration(readTimeoutSecondsCount),
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
	}
}

