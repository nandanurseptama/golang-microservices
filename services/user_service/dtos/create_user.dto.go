package dtos

import "encoding/json"

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func (dto CreateUserDto) ToMap() (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	bytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &inInterface)
	if err != nil {
		return nil, err
	}
	return inInterface, nil
}
