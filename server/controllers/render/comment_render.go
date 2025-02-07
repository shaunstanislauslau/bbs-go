package render

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/markdown"
	"bbs-go/services"
	"html"
)

func BuildComments(comments []model.Comment) []model.CommentResponse {
	var ret []model.CommentResponse
	for _, comment := range comments {
		ret = append(ret, *BuildComment(comment))
	}
	return ret
}

func BuildComment(comment model.Comment) *model.CommentResponse {
	return _buildComment(&comment, true)
}

func _buildComment(comment *model.Comment, buildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &model.CommentResponse{
		CommentId:  comment.Id,
		User:       BuildUserInfoDefaultIfNull(comment.UserId),
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		QuoteId:    comment.QuoteId,
		Status:     comment.Status,
		CreateTime: comment.CreateTime,
	}

	if comment.Status == constants.StatusOk {
		if comment.ContentType == constants.ContentTypeMarkdown {
			content := markdown.ToHTML(comment.Content)
			ret.Content = handleHtmlContent(content)
		} else if comment.ContentType == constants.ContentTypeHtml {
			ret.Content = handleHtmlContent(comment.Content)
		} else {
			ret.Content = html.EscapeString(comment.Content)
		}
		ret.ImageList = buildImageList(comment.ImageList)
	} else {
		ret.Content = "内容已删除"
	}

	if buildQuote && comment.QuoteId > 0 {
		quote := _buildComment(services.CommentService.Get(comment.QuoteId), false)
		if quote != nil {
			ret.Quote = quote
		}
	}
	return ret
}
