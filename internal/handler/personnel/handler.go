package personnel

import (
	"example/internal/dto"
	"example/internal/mapper"
	"example/internal/storage"
	"example/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ds *storage.DataStorage
}

func New(ds *storage.DataStorage) *Handler {
	return &Handler{
		ds: ds,
	}
}

func (h *Handler) CreatePersonnel(c *gin.Context) {
	var personnel dto.PersonnelCreateRequest
	if err := c.ShouldBindJSON(&personnel); err != nil {
		c.JSON(400, dto.CommonError{
			Message: err.Error(),
			Code:    "invalid_request",
		})
		return
	}
	p, edus := mapper.ToModelFromRequest(&personnel)

	txErr := h.ds.Transaction(c.Request.Context(), func(subds *storage.DataStorage) error {
		var err error
		p, err = subds.CreatePersonnel(c.Request.Context(), p)
		if err != nil {
			return err
		}
		edus, err = subds.CreateEducations(c.Request.Context(), edus...)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		c.JSON(500, dto.CommonError{
			Message: txErr.Error(),
			Code:    "internal_error",
		})
		return
	}
	c.JSON(201, dto.PersonnelResponse{
		Person:     mapper.ToPersonnelDTO(p),
		Educations: mapper.ToEducationDTOs(edus),
	})
}

func (h *Handler) GetPersonnel(c *gin.Context) {
	id := c.Param("personnelId")
	if id == "" {
		c.JSON(400, dto.CommonError{
			Message: "personnelId can not be empty",
			Code:    "invalid_request",
		})
		return
	}
	p, err := h.ds.GetPersonnel(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.CommonError{
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}
	if p == nil {
		c.JSON(404, dto.CommonError{
			Message: "personnel not found",
			Code:    "not_found",
		})
		return
	}
	edus, err := h.ds.ListEducationsByPersonnelId(c.Request.Context(), p.Id)
	c.JSON(200, dto.PersonnelResponse{
		Person:     mapper.ToPersonnelDTO(p),
		Educations: mapper.ToEducationDTOs(edus),
	})
}

func (h *Handler) UpdatePersonnel(c *gin.Context) {
	var personnel dto.PersonnelCreateRequest
	if err := c.ShouldBindJSON(&personnel); err != nil {
		c.JSON(400, dto.CommonError{
			Message: err.Error(),
			Code:    "invalid_request",
		})
		return
	}
	id := c.Param("personnelId")
	if id == "" {
		c.JSON(400, dto.CommonError{
			Message: "personnelId can not be empty",
			Code:    "invalid_request",
		})
		return
	}
	currentPersonnel, err := h.ds.GetPersonnel(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.CommonError{
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}
	if currentPersonnel == nil {
		c.JSON(404, dto.CommonError{
			Message: "personnel not found",
			Code:    "not_found",
		})
		return
	}
	p, edus := mapper.ToModelFromRequest(&personnel, currentPersonnel)
	txErr := h.ds.Transaction(c.Request.Context(), func(subds *storage.DataStorage) error {
		var err error
		p, err = subds.CreatePersonnel(c.Request.Context(), p)
		if err != nil {
			return err
		}

		if err = subds.DeleteEducationsByPersonnelId(c.Request.Context(), p.Id); err != nil {
			return err
		}

		edus, err = subds.CreateEducations(c.Request.Context(), edus...)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		c.JSON(500, dto.CommonError{
			Message: txErr.Error(),
			Code:    "internal_error",
		})
		return
	}
	c.JSON(200, dto.PersonnelResponse{
		Person:     mapper.ToPersonnelDTO(p),
		Educations: mapper.ToEducationDTOs(edus),
	})
}

func (h *Handler) DeletePersonnel(c *gin.Context) {
	id := c.Param("personnelId")
	if id == "" {
		c.JSON(400, dto.CommonError{
			Message: "personnelId can not be empty",
			Code:    "invalid_request",
		})
		return
	}
	p, err := h.ds.GetPersonnel(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.CommonError{
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}
	if p == nil {
		c.JSON(404, dto.CommonError{
			Message: "personnel not found",
			Code:    "not_found",
		})
		return
	}
	txErr := h.ds.Transaction(c.Request.Context(), func(subds *storage.DataStorage) error {
		var err error
		if err = subds.DeleteEducationsByPersonnelId(c.Request.Context(), p.Id); err != nil {
			return err
		}

		if err = subds.DeletePersonnel(c.Request.Context(), p.Id); err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		c.JSON(500, dto.CommonError{
			Message: txErr.Error(),
			Code:    "internal_error",
		})
		return
	}
	c.JSON(204, map[string]interface{}{})
}

func (h *Handler) ListPersonnel(c *gin.Context) {
	var paging dto.PageRequest
	if err := c.ShouldBind(&paging); err != nil {
		c.JSON(400, dto.CommonError{
			Message: err.Error(),
			Code:    "invalid_request",
		})
		return
	}
	// TODO can validate/set default dto by using tag and gin binding but I don't like that
	if paging.Page <= 0 {
		paging.Page = 1
	}
	if paging.Size <= 0 {
		paging.Size = 100
	}
	ps, count, err := h.ds.ListPersonnels(c.Request.Context(), paging.Page, paging.Size)
	if err != nil {
		c.JSON(500, dto.CommonError{
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}
	resp := []*dto.PersonnelResponse{}
	// TODO can optimize this one
	for _, p := range ps {
		edus, err := h.ds.ListEducationsByPersonnelId(c.Request.Context(), p.Id)
		if err != nil {
			log.Errorw("error when query educations", "error", err)
			// Can return/break here if want. I use continue to don't add the personel into response list
			continue
		}
		resp = append(resp, &dto.PersonnelResponse{
			Person:     mapper.ToPersonnelDTO(p),
			Educations: mapper.ToEducationDTOs(edus),
		})
	}
	c.JSON(200, &dto.PagingDTO{
		Page:       paging.Page,
		Size:       paging.Size,
		Count:      len(resp),
		TotalCount: count,
		Data:       resp,
	})
}
