package content

type CreateContentRequest struct {
	UserId      string `json:"userId" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
	Text        string `json:"text,omitempty"`
	MediaUrl    string `json:"mediaUrl" validate:"required"`
	Type        string `json:"type" validate:"required"`
	FolderId    string `json:"folderId" validate:"required"`
}

type UpdateContentRequest struct {
	Id          string  `json:"id" validate:"required"`
	DisplayName *string `json:"displayName,omitempty"`
	Text        *string `json:"text,omitempty"`
	MediaUrl    *string `json:"mediaUrl,omitempty"`
}

type Filter struct {
	UserId   string
	Search   string
	FolderId string
}

type ContentDto struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	DisplayName string `json:"display_name"`
	Text        string `json:"text"`
	MediaUrl    string `json:"media_url"`
	Type        string `json:"type"`
	FolderId    string `json:"folder_id"`
}

func MapContentToDto(model *Content) ContentDto {
	if model == nil {
		return ContentDto{}
	}

	return ContentDto{
		Id:          model.Id,
		UserId:      model.UserId,
		DisplayName: model.DisplayName,
		Text:        model.Text,
		MediaUrl:    model.MediaUrl,
		Type:        model.Type,
		FolderId:    model.FolderId,
	}
}

func MapContentSliceToDto(content []*Content) []*ContentDto {
	if content == nil {
		return make([]*ContentDto, 0)
	}

	result := make([]*ContentDto, len(content))
	for _, model := range content {
		result = append(result, &ContentDto{
			Id:          model.Id,
			UserId:      model.UserId,
			DisplayName: model.DisplayName,
			Text:        model.Text,
			MediaUrl:    model.MediaUrl,
			Type:        model.Type,
			FolderId:    model.FolderId,
		})
	}

	return result
}
