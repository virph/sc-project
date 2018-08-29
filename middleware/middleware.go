package middleware

import (
	"fmt"
	"log"
	"net/http"
)

type Middleware interface {
	Check(http.HandlerFunc) error
	Next(Middleware)
}

type AuthMiddleware struct {
	nextMiddleware Middleware
}

func (m *AuthMiddleware) Check(h http.HandlerFunc) error {
	log.Println("authmiddleware check")
	if !false {
		log.Println("authmiddleware error")
		return fmt.Errorf("authmiddleware error")
	}
	if m.nextMiddleware != nil {
		return m.nextMiddleware.Check(h)
	}
	log.Println("authmiddleware jalan")
	return nil
}

func (m *AuthMiddleware) Next(n Middleware) {
	m.nextMiddleware = n
}

type CumiMiddleware struct {
	nextMiddleware Middleware
}

func (m *CumiMiddleware) Check(h http.HandlerFunc) error {
	log.Println("cumimiddleware check")
	if !false {
		log.Println("cumimiddleware error")
		return fmt.Errorf("cumimiddleware error")
	}
	if m.nextMiddleware != nil {
		return m.nextMiddleware.Check(h)
	}
	log.Println("cumimiddleware jalan")
	return nil
}

func (m *CumiMiddleware) Next(n Middleware) {
	m.nextMiddleware = n
}
