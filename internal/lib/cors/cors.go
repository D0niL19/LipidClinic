package cors

import "github.com/rs/cors"

func CorsSettings() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		ExposedHeaders:     []string{"Authorization"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
		Debug:              true,
	})

	return c
}
