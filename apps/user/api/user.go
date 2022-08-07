package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"

	"codeup.aliyun.com/baber/go/keyauth/apps/user"
)

func (h *handler) CreateUser(r *restful.Request, w *restful.Response) {
	req := user.NewCreateUserRequest()

	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.CreateUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, set)
}

func (h *handler) QueryUser(r *restful.Request, w *restful.Response) {
	req := user.NewQueryUserRequestFromHTTP(r.Request)
	set, err := h.service.QueryUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (h *handler) DescribeUser(r *restful.Request, w *restful.Response) {
	req := user.NewDescribeUserRequest(r.PathParameter("id"))
	ins, err := h.service.DescribeUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, ins)
}

func (h *handler) UpdateUser(r *restful.Request, w *restful.Response) {
	req := user.NewPutUserRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Data); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.UpdateUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (h *handler) PatchUser(r *restful.Request, w *restful.Response) {
	req := user.NewPatchUserRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Data); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.UpdateUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (h *handler) DeleteUser(r *restful.Request, w *restful.Response) {
	req := user.NewDeleteUserRequestWithID(r.PathParameter("id"))
	set, err := h.service.DeleteUser(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
