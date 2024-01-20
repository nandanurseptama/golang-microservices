package entities

import "encoding/json"

type UserEntity struct {
	Id        string  `json:"id,omitempty"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

func UserEntityFromMap(d map[string]interface{}) (*UserEntity, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	user := UserEntity{}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserEntity) GetUpdatedAt() string {
	if u.UpdatedAt == nil {
		return ""
	}
	return *u.UpdatedAt
}
