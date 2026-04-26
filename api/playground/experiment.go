package playground

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type key string

const (
	requestId key = "requestId"
	userId    key = "userId"
)

// worker
func getUser(ctx context.Context) (map[string]string, error) {
	id, ok := ctx.Value(requestId).(string)
	uId, ok2 := ctx.Value(userId).(string)
	if !ok {
		id = "unknown"
	}
	if !ok2 {
		uId = "unknown"
	}
	select {
	case <-time.After(1 * time.Second):
		response := map[string]string{
			"status":  "Job done",
			"req_id":  id,
			"user_id": uId,
		}
		return response, nil
	case <-ctx.Done():
		fmt.Println("Closed work!")
		return nil, ctx.Err()
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), requestId, "1209456789v123456")
	ctx = context.WithValue(ctx, userId, "098765432123467890ifdfghjkl0987654ecvhu65rdcg4")
	result, err := getUser(ctx)

	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
