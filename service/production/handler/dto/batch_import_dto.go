package dto

type BatchImportRowResult struct {
	Row    int            `json:"row"`
	Status string         `json:"status"`
	Batch  *BatchResponse `json:"batch,omitempty"`
	Error  *string        `json:"error,omitempty"`
}

type BatchImportTotals struct {
	TotalRows int `json:"total_rows"`
	Created   int `json:"created"`
	Failed    int `json:"failed"`
}

type BatchImportResponse struct {
	Totals  BatchImportTotals      `json:"totals"`
	Results []BatchImportRowResult `json:"results"`
}
