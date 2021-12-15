package dto

type ResponseMetaList struct {
	NumOfResults int64  `json:"num_of_results,omitempty"`
	NumOfPages   int64  `json:"num_of_pages,omitempty"`
	CurrentPage  int32  `json:"current_page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	NodeID       string `json:"node_id,omitempty"`
} // @Name ResponseMetaList

type ResponseMeta struct {
}
