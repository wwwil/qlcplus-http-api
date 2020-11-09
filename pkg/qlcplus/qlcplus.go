package qlcplus

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"strings"
)

// Connection is a connection to a QLC+ instance.
type Connection struct {
	// Address is the IP address or hostname, and port number, to access QLC+.
	Address string
}

// GetWidgetsMap returns a map of widget IDs to names.
func (q *Connection) GetWidgetsMap() (map[string]string, error) {
	message := "QLC+API|getWidgetsList"
	widgetsListString, err := q.writeRead(message)
	if err != nil {
		return nil, err
	}
	widgetsListString = strings.TrimPrefix(widgetsListString, message + "|")
	widgetsList := strings.Split(widgetsListString, "|")
	widgetsMap := make(map[string]string)
	for i := 0; i < len(widgetsList); i = i + 2 {
		widgetsMap[widgetsList[i]] = widgetsList[i + 1]
	}
	return widgetsMap, nil
}

// GetWidgetIDByName gets the ID of a widget from its name. If multiple widgets have the same name the behaviour is
// undefined.
func (q *Connection) GetWidgetIDByName(widgetNameToFind string) (string, error) {
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
func (q *Connection) GetWidgetStatusByID(widgetID string) (string, error) {
	return q.writeRead(fmt.Sprintf("QLC+API|getWidgetStatus|%s", widgetID))
}

// GetWidgetStatusByName gets the status of a widget by name.
func (q *Connection) GetWidgetStatusByName(widgetName string) (string, error) {
	widgetID, err := q.GetWidgetIDByName(widgetName)
	if err != nil {
		return "", err
	}
	return q.GetWidgetStatusByID(widgetID)
}

// SetWidgetStatusByID sets the status of a widget by ID. How the widget behaves with the specified value depends on the
// widget type.
func (q *Connection) SetWidgetStatusByID(widgetID, widgetValue string) (string, error) {
	return q.writeRead(fmt.Sprintf("%s|%s", widgetID, widgetValue))
}

// SetWidgetStatusByName sets the status of a widget by name. How the widget behaves with the specified value depends on
// the widget type.
func (q *Connection) SetWidgetStatusByName(widgetName, widgetValue string) (string, error) {
	widgetID, err := q.GetWidgetIDByName(widgetName)
	if err != nil {
		return "", err
	}
	return q.writeRead(fmt.Sprintf("%s|%s", widgetID, widgetValue))
}

// writeRead creates a new QLC+ websocket connection and uses it to write the specified message then read the response.
func (q *Connection) writeRead(message string) (string, error) {
	// Create the websocket URL.
	u := url.URL{Scheme: "ws", Host: q.Address, Path: "/qlcplusWS"}
	// Connect to the websocket.
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", err
	}
	defer c.Close()
	// Write a message to the websocket.
	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return "", err
	}
	// Read the response to the message.
	_, p, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(p), nil
}
