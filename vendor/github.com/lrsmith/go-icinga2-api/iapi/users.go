package iapi

import (
	"encoding/json"
	"fmt"
)

// GetUser ...
func (server *Server) GetUser(name string) ([]UserStruct, error) {

	var users []UserStruct
	results, err := server.NewAPIRequest("GET", "/objects/users/"+name, nil)
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

	if unmarshalErr := json.Unmarshal(jsonStr, &users); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return users, err
}

// CreateUser ...
func (server *Server) CreateUser(name, email string) ([]UserStruct, error) {

	var newAttrs UserAttrs
	newAttrs.Email = email

	var newUser UserStruct
	newUser.Name = name
	newUser.Type = "User"
	newUser.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newUser)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the users.
	results, err := server.NewAPIRequest("PUT", "/objects/users/"+name, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		users, err := server.GetUser(name)
		return users, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// DeleteUser ...
func (server *Server) DeleteUser(name string) error {
	results, err := server.NewAPIRequest("DELETE", "/objects/users/"+name+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
