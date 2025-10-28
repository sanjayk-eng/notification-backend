package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sanjay/api/query"
	"sanjay/api/service"
	"sanjay/api/util"
	"sanjay/config/connection"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connections = make(map[string]*websocket.Conn)

type HandlerMeth struct {
	TwillowHandler *service.Twillow
	RedishHandler  *service.RedisCline
	Query          *query.QueryMeth
}

func NewHandlerAttach(twillowHandler *service.Twillow, redishHandler *service.RedisCline, queryHandler *query.QueryMeth) *HandlerMeth {
	return &HandlerMeth{
		TwillowHandler: twillowHandler,
		RedishHandler:  redishHandler,
		Query:          queryHandler,
	}
}

func (h *HandlerMeth) SendSMS() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := util.OTPRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if err := util.Validate.Struct(req); err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "Invalid  Input Data")
			return
		}
		otp := service.GenerateOTP()
		fullPhone := fmt.Sprintf("+%s%s", req.CountryCode, req.Phone)
		_, err := h.TwillowHandler.SendMessage(fullPhone, otp)
		if err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, fmt.Sprintf(" Failed to send OTP: %v", err))
			return
		}
		if err := h.RedishHandler.StoreOTP(context.Background(), fullPhone, otp); err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, "Invalid  Input Data")
			return
		}
		util.JsonRespond(c, http.StatusOK, fmt.Sprintf(" OTP sent successfully! otp: %s", otp))
	}
}

func (h *HandlerMeth) VerifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := util.User{}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "Invalid JSON format")
			return
		}
		if err := util.Validate.Struct(req); err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "Invalid  Input Data")
			return
		}
		fullPhone := fmt.Sprintf("+%s%s", req.CountryCode, req.Phone)
		val, err := h.RedishHandler.GetOTP(context.Background(), fullPhone)
		if err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "too late")
			return
		}
		if val != req.OTP {
			util.JsonRespond(c, http.StatusBadRequest, "invalid capture")
			return
		}
		if err := h.RedishHandler.DeleteOTP(context.Background(), fullPhone); err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, "failed to remove otp")
			return
		}
		accessToken, err := service.GenerateAccessToken(req.Phone)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := service.GenerateRefreshToken(req.Phone)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true)
		c.Header("Access-Control-Expose-Headers", "token")
		c.Header("token", accessToken)
		h.Query.StoredPhoneWithCountryCode(req.CountryCode, req.Phone)
		util.JsonRespond(c, http.StatusOK, fmt.Sprintf("success"))
	}
}

func (h *HandlerMeth) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneVal, exists := c.Get("phone")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		phone, ok := phoneVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone format"})
			return
		}
		user, err := h.Query.GetDataPhone(phone)
		if err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, "failed to get user")
		}

		// Here you can fetch profile from DB or a store using `phone`
		profileData := map[string]interface{}{
			"fullName":     user.FullName,
			"profilePhoto": user.ProfilePhoto,
			"countryCode":  user.CountryCode,
			"phone":        user.Phone,
		}

		c.JSON(http.StatusOK, profileData)
	}
}
func (h *HandlerMeth) UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneVal, exists := c.Get("phone")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		phone, ok := phoneVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone format"})
			return
		}
		user := util.AllDetails{}
		if err := c.ShouldBindJSON(&user); err != nil {
			util.JsonRespond(c, http.StatusBadRequest, "failed to bind data")
			return
		}
		fmt.Println("user", user)
		if err := h.Query.UpdateUserProfile(phone, user.ProfilePhoto, user.FullName); err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, "failed to update user")
			return
		}
		profileData := map[string]interface{}{
			"fullName":     user.FullName,
			"profilePhoto": user.ProfilePhoto,
			"countryCode":  user.CountryCode,
			"phone":        phone,
		}
		c.JSON(http.StatusOK, profileData)
	}
}
func (h *HandlerMeth) LiveConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		activeUsers := []util.AllDetails{}

		for phone := range connections {
			user, err := h.Query.GetDataPhone(phone) // fetch full details from DB
			if err != nil {
				continue // skip if user not found
			}
			activeUsers = append(activeUsers, *user)
		}

		util.JsonRespond(c, 200, gin.H{
			"count": len(activeUsers),
			"users": activeUsers, // full details
		})
	}
}

func (h *HandlerMeth) WsConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneVal, exists := c.Get("phone")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		phone, ok := phoneVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone format"})
			return
		}

		conn, err := connection.GetWsClient(c)
		if err != nil {
			util.JsonRespond(c, http.StatusInternalServerError, "Failed to create WebSocket connection")
			return
		}

		connections[phone] = conn
		ctx := context.Background()
		// Subscribe to the global chat channel
		pubsub := h.RedishHandler.SubscribeChannel(ctx)
		defer pubsub.Close()

		// Goroutine to broadcast incoming messages from Redis to this client
		go func() {
			for msg := range pubsub.Channel() {
				var wsMsg util.WSMessage
				if err := json.Unmarshal([]byte(msg.Payload), &wsMsg); err != nil {
					log.Println("PubSub unmarshal error:", err)
					continue
				}

				if err := conn.WriteJSON(wsMsg); err != nil {
					log.Println("Client disconnected:", err)
					delete(connections, phone)
					break
				}
			}
		}()

		// Fetch user metadata
		user, err := h.Query.GetDataPhone(phone)
		if err != nil {
			user = &util.AllDetails{
				Phone:        phone,
				FullName:     "Unknown",
				ProfilePhoto: "",
			}
		}

		// Broadcast "join" notification
		// joinMsg := util.WSMessage{
		// 	Phone:    phone,
		// 	FullName: user.FullName,
		// 	Message:  fmt.Sprintf("%s has joined the chat", user.FullName),
		// 	Type:     "join",
		// }
		// h.RedishHandler.PushMessage(ctx, joinMsg)
		// h.RedishHandler.PublishMessage(ctx, joinMsg)

		// Send last 100 messages to the newly connected client
		history, _ := h.RedishHandler.GetMessages(ctx)
		for _, msg := range history {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Failed to send history message:", err)
			}
		}

		// Listen for incoming messages from this client
		for {
			var incoming struct {
				Message string `json:"message"`
			}
			if err := conn.ReadJSON(&incoming); err != nil {
				log.Println("Read error:", err)
				delete(connections, phone)
				break

				// Broadcast "leave" notification
				// leaveMsg := util.WSMessage{
				// 	Phone:    phone,
				// 	FullName: user.FullName,
				// 	Message:  fmt.Sprintf("%s has left the chat", user.FullName),
				// 	Type:     "leave",
				// }
				// h.RedishHandler.PushMessage(ctx, leaveMsg)
				// h.RedishHandler.PublishMessage(ctx, leaveMsg)
				// break
			}

			// Create chat message
			wsMsg := util.WSMessage{
				Phone:        phone,
				FullName:     user.FullName,
				ProfilePhoto: user.ProfilePhoto,
				Message:      incoming.Message,
				Type:         "message",
			}

			// Save and broadcast
			h.RedishHandler.PushMessage(ctx, wsMsg)
			h.RedishHandler.PublishMessage(ctx, wsMsg)
		}
	}
}
