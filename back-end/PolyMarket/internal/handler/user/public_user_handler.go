package user

import (
	"net/http"

	"X402AiPolyMarket/PolyMarket/internal/logic/user"
	"X402AiPolyMarket/PolyMarket/internal/svc"
	"X402AiPolyMarket/PolyMarket/internal/utils"

	"github.com/gorilla/mux"
)

func GetPublicUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		if address == "" {
			utils.ParamError(w, "Missing address parameter")
			return
		}

		l := user.NewPublicUserLogic(r.Context(), svcCtx)
		resp, err := l.GetPublicUser(address)
		if err != nil {
			if customErr, ok := utils.IsCustomError(err); ok {
				utils.Error(w, customErr.Code, customErr.Msg)
			} else {
				utils.ServerError(w, err.Error())
			}
			return
		}

		utils.Success(w, resp)
	}
}

