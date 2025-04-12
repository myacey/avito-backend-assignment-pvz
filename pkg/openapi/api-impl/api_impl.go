package api_impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
	"github.com/myacey/avito-backend-assignment-pvz/pkg/openapi"
	"github.com/oapi-codegen/runtime/types"
)

// APIHandlerWithFuncs realizes openapi.ServerInterface
type APIHandlerWithFuncs struct {
	PostDummyLoginFunc                 func(c *gin.Context)
	PostLoginFunc                      func(c *gin.Context)
	PostProductsFunc                   func(c *gin.Context)
	GetPvzFunc                         func(c *gin.Context, params openapi.GetPvzParams)
	PostPvzFunc                        func(c *gin.Context)
	PostPvzPvzIdCloseLastReceptionFunc func(c *gin.Context, pvzId types.UUID)
	PostPvzPvzIdDeleteLastProductFunc  func(c *gin.Context, pvzId types.UUID)
	PostReceptionsFunc                 func(c *gin.Context)
	PostRegisterFunc                   func(c *gin.Context)
}

func (h *APIHandlerWithFuncs) PostDummyLogin(c *gin.Context) {
	h.PostDummyLoginFunc(c)
}

func (h *APIHandlerWithFuncs) PostLogin(c *gin.Context) {
	h.PostLoginFunc(c)
}

func (h *APIHandlerWithFuncs) PostProducts(c *gin.Context) {
	h.PostProductsFunc(c)
}

func (h *APIHandlerWithFuncs) GetPvz(c *gin.Context, params openapi.GetPvzParams) {
	h.GetPvzFunc(c, params)
}

func (h *APIHandlerWithFuncs) PostPvz(c *gin.Context) {
	h.PostPvzFunc(c)
}

func (h *APIHandlerWithFuncs) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId types.UUID) {
	h.PostPvzPvzIdCloseLastReceptionFunc(c, pvzId)
}

func (h *APIHandlerWithFuncs) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId types.UUID) {
	h.PostPvzPvzIdDeleteLastProductFunc(c, pvzId)
}

func (h *APIHandlerWithFuncs) PostReceptions(c *gin.Context) {
	h.PostReceptionsFunc(c)
}

func (h *APIHandlerWithFuncs) PostRegister(c *gin.Context) {
	h.PostRegisterFunc(c)
}

// NewAPIHandlerWithFuncs создаёт новый APIHandlerWithFuncs, назначая функции-обработчики.
// Здесь можно переиспользовать уже написанные функции из пакета handler.
func NewAPIHandlerWithFuncs(svc *service.Service) *APIHandlerWithFuncs {
	return &APIHandlerWithFuncs{
		PostDummyLoginFunc: func(c *gin.Context) {
			if err := handler.DummyLogin(c, &svc.UserService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostLoginFunc: func(c *gin.Context) {
			if err := handler.Login(c, &svc.UserService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostProductsFunc: func(c *gin.Context) {
			if err := handler.AddProductToReception(c, &svc.ReceptionService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		GetPvzFunc: func(c *gin.Context, params openapi.GetPvzParams) {
			// Здесь можно вызвать метод сервиса, возвращающий список ПВЗ с фильтрацией и пагинацией.
			// Для примера возвращаем пустой срез.
			c.JSON(http.StatusOK, []openapi.PVZ{})
		},
		PostPvzFunc: func(c *gin.Context) {
			if err := handler.CreatePvz(c, &svc.PvzService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostPvzPvzIdCloseLastReceptionFunc: func(c *gin.Context, pvzId types.UUID) {
			if err := handler.FinishReception(c, &svc.ReceptionService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostPvzPvzIdDeleteLastProductFunc: func(c *gin.Context, pvzId types.UUID) {
			if err := handler.DeleteLastProduct(c, &svc.ReceptionService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostReceptionsFunc: func(c *gin.Context) {
			if err := handler.CreateReception(c, &svc.ReceptionService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
		PostRegisterFunc: func(c *gin.Context) {
			if err := handler.Register(c, &svc.UserService); err != nil {
				c.JSON(http.StatusInternalServerError, openapi.Error{Message: err.Error()})
			}
		},
	}
}
