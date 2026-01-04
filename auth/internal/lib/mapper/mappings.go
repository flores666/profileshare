package mapper

import (
	"auth/internal/storage"
	"time"
)

type UserDto struct {
	Id              string    `json:"id"`
	Nickname        string    `json:"nickname"`
	Email           string    `json:"email"`
	CodeRequestedAt time.Time `json:"codeRequestedAt"`
	RoleId          string    `json:"roleId"`
	BannedBefore    time.Time `json:"bannedBefore"`
}

func MapUserToDto(user *storage.User) UserDto {
	if user == nil {
		return UserDto{}
	}

	return UserDto{
		Id:              user.Id,
		Nickname:        user.Nickname,
		Email:           user.Email,
		CodeRequestedAt: user.CodeRequestedAt,
		RoleId:          user.RoleId,
		BannedBefore:    user.BannedBefore,
	}
}

func MapUserSliceToDto(users []*storage.User) []*UserDto {
	if users == nil {
		return make([]*UserDto, 0)
	}

	result := make([]*UserDto, len(users))
	for _, item := range users {
		result = append(result, &UserDto{
			Id:              item.Id,
			Nickname:        item.Nickname,
			Email:           item.Email,
			CodeRequestedAt: item.CodeRequestedAt,
			RoleId:          item.RoleId,
			BannedBefore:    item.BannedBefore,
		})
	}

	return result
}
