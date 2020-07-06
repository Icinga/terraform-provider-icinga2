package iapi

import (
	"encoding/json"
	"fmt"
)

// GetNotification ...
func (server *Server) GetNotification(name string) ([]NotificationStruct, error) {

	var notifications []NotificationStruct
	results, err := server.NewAPIRequest("GET", "/objects/notifications/"+name, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &notifications); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return notifications, err
}

// CreateNotification ...
func (server *Server) CreateNotification(name, hostname, command, servicename string, interval int, users []string, vars map[string]string, templates []string) ([]NotificationStruct, error) {

	var newAttrs NotificationAttrs
	newAttrs.Command = command
	newAttrs.Users = users
	newAttrs.Servicename = servicename
	newAttrs.Interval = interval
	newAttrs.Vars = vars
	newAttrs.Templates = templates

	var newNotification NotificationStruct
	newNotification.Name = name
	newNotification.Type = "Notification"
	newNotification.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newNotification)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/notifications/"+name, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		notifications, err := server.GetNotification(name)
		return notifications, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// DeleteNotification ...
func (server *Server) DeleteNotification(name string) error {
	results, err := server.NewAPIRequest("DELETE", "/objects/notifications/"+name+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
