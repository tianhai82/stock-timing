package cloudtask

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AddCloudTasks add cloud task handlers to router
func AddCloudTasks(router *gin.RouterGroup) {
	router.POST("/analyze-stocks", analyzeStock)
}

func analyzeStock(c *gin.Context) {

}

// CreateTask creates a new task in your App Engine queue.
func CreateTask(projectID, locationID, queueID string, message gin.H) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#AppEngineHttpRequest
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/task/" + queueID,
				},
			},
		},
	}

	if message != nil {
		// Add a payload message if one is present.
		b, err := json.Marshal(message)
		if err == nil {
			req.Task.GetAppEngineHttpRequest().Body = b
		}
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}
