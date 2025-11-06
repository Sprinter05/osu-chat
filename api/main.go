package api

const OSU_URL string = "https://osu.ppy.sh/api/v2"
const TIME_FORMAT string = "2006-01-02T15:04:05Z0700"

type OAuth struct {
	ClientId    int    `json:"client_id"`
	TokenSecret string `json:"token_secret"`
}
