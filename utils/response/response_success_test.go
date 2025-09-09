package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse_Structure(t *testing.T) {
	// Test the SuccessResponse struct
	response := SuccessResponse{
		Status:  200,
		Message: "Test message",
		Data:    []string{"test", "data"},
	}

	// Test that all fields are properly set
	assert.Equal(t, 200, response.Status)
	assert.Equal(t, "Test message", response.Message)
	assert.Equal(t, []string{"test", "data"}, response.Data)
}

func TestSuccessResponse_JSONSerialization(t *testing.T) {
	// Test JSON marshaling
	response := SuccessResponse{
		Status:  200,
		Message: "Success",
		Data:    map[string]interface{}{"key": "value"},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Verify the data is the same
	assert.Equal(t, response.Status, unmarshaled.Status)
	assert.Equal(t, response.Message, unmarshaled.Message)

	// For interface{} comparison, we need to convert to expected type
	expectedData := map[string]interface{}{"key": "value"}
	actualData := unmarshaled.Data.(map[string]interface{})
	assert.Equal(t, expectedData, actualData)
}

func TestSuccessResponse_EmptyData(t *testing.T) {
	// Test with empty data slice
	response := SuccessResponse{
		Status:  200,
		Message: "Success with empty data",
		Data:    []interface{}{},
	}

	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, response.Status, unmarshaled.Status)
	assert.Equal(t, response.Message, unmarshaled.Message)
	assert.Empty(t, unmarshaled.Data)
}

func TestSuccessResponse_NilData(t *testing.T) {
	// Test with nil data
	response := SuccessResponse{
		Status:  200,
		Message: "Success with nil data",
		Data:    nil,
	}

	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	var unmarshaled SuccessResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, response.Status, unmarshaled.Status)
	assert.Equal(t, response.Message, unmarshaled.Message)
	assert.Nil(t, unmarshaled.Data)
}
