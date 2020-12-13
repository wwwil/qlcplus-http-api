package qlcplus

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Whether to print websocket messages, used for development.
var debug = false

// WebsocketConnectionHandler handles a websocket connection to a QLC+ instance.
type WebsocketConnectionHandler struct {
	// Address is the IP address or hostname, and port number, to access the QLC+ websocket.
	Address string
}

// GetWidgetsMap returns a map of widget IDs to names.
func (q *WebsocketConnectionHandler) GetWidgetsMap() (map[string]string, error) {
	request := "QLC+API|getWidgetsList"
	widgetsListString, err := q.makeRequest(request, "")
	if err != nil {
		return nil, err
	}
	widgetsList := strings.Split(widgetsListString, "|")
	widgetsMap := make(map[string]string)
	for i := 0; i < len(widgetsList) - 1; i = i + 2 {
		widgetsMap[widgetsList[i]] = widgetsList[i + 1]
	}
	return widgetsMap, nil
}

// GetWidgetIDByName gets the ID of a widget from its name. If multiple widgets have the same name the behaviour is
// undefined.
func (q *WebsocketConnectionHandler) GetWidgetIDByName(widgetNameToFind string) (string, error) {
	widgetsMap, err := q.GetWidgetsMap()
	if err != nil {
		return "", err
	}
	for widgetID, widgetName := range widgetsMap {
		if widgetName == widgetNameToFind {
			return widgetID, nil
		}
	}
	return "", nil
}

// GetWidgetStatusByID gets the status of a widget by ID.
func (q *WebsocketConnectionHandler) GetWidgetStatusByID(widgetID string) (string, error) {
	getWidgetStatus := "QLC+API|getWidgetStatus"
	return q.makeRequest(getWidgetStatus, widgetID)
}

// GetWidgetStatusByName gets the status of a widget by name.
func (q *WebsocketConnectionHandler) GetWidgetStatusByName(widgetName string) (string, error) {
	widgetID, err := q.GetWidgetIDByName(widgetName)
	if err != nil {
		return "", err
	}
	return q.GetWidgetStatusByID(widgetID)
}

// TODO: Add functions to getWidgetType.

// SetWidgetStatusByID sets the status of a widget by ID. How the widget behaves with the specified value depends on the
// widget type.
func (q *WebsocketConnectionHandler) SetWidgetStatusByID(widgetID, widgetValue string) (string, error) {
	return q.makeRequest(widgetID, widgetValue)
}

// SetWidgetStatusByName sets the status of a widget by name. How the widget behaves with the specified value depends on
// the widget type.
func (q *WebsocketConnectionHandler) SetWidgetStatusByName(widgetName, widgetValue string) (string, error) {
	widgetID, err := q.GetWidgetIDByName(widgetName)
	if err != nil {
		return "", err
	}
	return q.SetWidgetStatusByID(widgetID, widgetValue)
}

// makeRequest wraps the readWrite function to confirm the response is in the expected format and strip the prefix.
func (q *WebsocketConnectionHandler) makeRequest(request, value string) (string, error) {
	// Compose the message from the request and value if present.
	message := request
	if value != "" {
		message = fmt.Sprintf("%s|%s", request, value)
	}
	// Write the request message and read the response.
	// TODO: If the widget's current value is already the value to set then QLC+ will not respond. This causes the read
	//       request to hang until it hits the timeout, it does not affect other requests. This should be
	//       checked or handled in some way.
	response, err := q.writeRead(message)
	if err != nil {
		return "", err
	}
	// Verify the response is in the expected format.
	if strings.Count(response, request) != 1 {
		return "", fmt.Errorf("unexpected response from QLC+ websocket, sent message: \"%s\", got response \"%s\"", message, request)
	}
	// Return the response without the request string.
	return strings.TrimPrefix(response, request + "|"), nil
}

// writeRead creates a new QLC+ websocket connection and uses it to write the specified message then read the response.
func (q *WebsocketConnectionHandler) writeRead(message string) (string, error) {
	// Create the websocket URL.
	u := url.URL{Scheme: "ws", Host: q.Address, Path: "/qlcplusWS"}
	// Connect to the websocket.
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", err
	}
	// Set a 10 second timeout.
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))
	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer c.Close()
	// Write a message to the websocket.
	if debug {
		fmt.Printf("Writing websocket message: %s\n", message)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", err
	}
	// Read the response to the message.
	_, p, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Printf("Read websocket message: %s\n", string(p))
	}
	return string(p), nil
}
