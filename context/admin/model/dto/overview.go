package admin_dto

type Overview struct {
	Albums       int64 `json:"albums"`
	Users        int64 `json:"users"`
	Cards        int64 `json:"cards"`
	Participants int64 `json:"participants"`
	Contacts     int64 `json:"contacts"`
	Collected    int64 `json:"collected"`
}