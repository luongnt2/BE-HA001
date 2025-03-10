package locations

import (
	"BE-HA001/cmd/api/pkg"
	"BE-HA001/cmd/api/pkg/util"
	"BE-HA001/pkg/httputil"
	"errors"
	"log"
	"net/http"
)

func NewGetDistanceHandler() *GetDistanceHandler {
	return &GetDistanceHandler{}
}

type GetDistanceHandler struct {
}

func (h *GetDistanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := h.parseReq(r)
	err := h.validateReq(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ipLat, ipLng, cityLat, cityLng float64

	eg := pkg.NewErrGroupWithRecovery(r.Context())
	eg.Go(func() error {
		clientIP := r.Header.Get("Client-IP")
		ipLat, ipLng, err = util.GetLocationFromIP(r.Context(), clientIP)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		cityLat, cityLng, err = util.GetCityCoordinates(req)
		if err != nil {
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Printf("Error getting location distance: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
	}
	resp := util.Haversine(ipLat, ipLng, cityLat, cityLng)

	httputil.ResponseWrapSuccessJSON(w, resp)
}

func (h *GetDistanceHandler) parseReq(r *http.Request) string {
	city := r.URL.Query().Get("city")
	return city
}

func (h *GetDistanceHandler) validateReq(req string) error {
	if req == "" {
		return errors.New("city is required")
	}

	return nil
}
