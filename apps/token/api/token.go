package api

import (
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
)

func (h *handler) IssueToken(r *restful.Request, w *restful.Response) {
	req := token.NewIssueTokenRequest()

	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.IssueToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, set)
}

func (h *handler) ValidateToken(r *restful.Request, w *restful.Response) {
	/**
	access_token从哪里获取
	1. url上传
	2. Custom Header传
	3. Authorization header上传
	*/
	auth := r.HeaderParameter("Authorization")
	req := token.NewValidateTokenRequest(strings.Trim(auth, "Bearer "))
	ins, err := h.service.ValidateToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, ins)
}

func (h *handler) RevokeToken(r *restful.Request, w *restful.Response) {
	req := token.NewRevokeTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	set, err := h.service.RevokeToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
