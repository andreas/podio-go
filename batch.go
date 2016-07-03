package podio

import "fmt"

type Batch struct {
	Id         int64  `json:"batch_id"`
	Name       string `json:"name"`
	Plugin     string `json:"plugin"`
	Status     string `json:"status"`
	Completed  int64  `json:"completed"`
	Skipped    int64  `json:"skipped"`
	Failed     int64  `json:"failed"`
	File       *File  `json:"file"`
	App        *App   `json:"app"`
	Space      *Space `json:"space"`
	CreatedOn  Time   `json:"created_on"`
	StartedOn  Time   `json:"started_on"`
	EndedOn    Time   `json:"ended_on"`
}

// https://developers.podio.com/doc/batch/get-batch-6144225
func (client *Client) GetBatch(batchId int64) (batch *Batch, err error) {
	path := fmt.Sprintf("/batch/%d", batchId)
	err = client.Request("GET", path, nil, nil, &batch)
	return
}
