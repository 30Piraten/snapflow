### **8. Testing & Debugging (Detailed Guide)**

#### **8.1 Choosing a Testing Library for Go**
Go has a built-in testing framework (`testing` package), which is the most commonly used approach for writing tests in Go. However, for more advanced assertions, you can use third-party libraries like:

- **`testing` (Standard Library)** → Best for unit tests, table-driven tests, and benchmarking.
- **Testify (`github.com/stretchr/testify`)** → Provides better assertions and mocking capabilities.
- **GoMock (`github.com/golang/mock`)** → Useful for mocking AWS SDKs and other dependencies.

#### **8.2 Types of Tests**
Your project requires different levels of testing:

1. **Unit Tests**: Test individual functions to ensure they behave correctly.
2. **Integration Tests**: Ensure multiple components (API, S3, SQS, DynamoDB) work together as expected.
3. **End-to-End (E2E) Tests**: Simulate the entire user workflow from photo upload to print completion.
4. **Load Tests**: Check how the system handles multiple requests.
5. **AWS-Specific Debugging**: Ensure logs, SQS queues, and Lambda executions work correctly.

---

### **8.3 How to Write Tests for This Project**

#### **Unit Testing (`testing` & Testify)**

Unit tests validate functions like resizing an image or generating a presigned URL.

Example: Testing `ResizePhoto()`  
```go
package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizePhoto(t *testing.T) {
	// Mock input image
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	// Call the function
	resizedImg, err := ResizePhoto(buf.Bytes(), 500, 500)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resizedImg)
}
```
✔ **Why?**  
- Uses `Testify` for assertions (`assert.NoError` and `assert.NotNil`).
- Tests that the function correctly resizes the image.

---

Example: Testing `GeneratePresignedURL()`
```go
func TestGeneratePresignedURL(t *testing.T) {
	url, err := GeneratePresignedURL("test-bucket", "test-file.jpg")
	assert.NoError(t, err)
	assert.Contains(t, url, "https://")
}
```
✔ **Why?**  
- Ensures the presigned URL is generated without errors.
- Checks if the returned string contains `"https://"` (basic validation).

---

#### **Integration Testing (Database, S3, and SQS)**  
Integration tests ensure that your backend interacts correctly with AWS services.

Example: Testing S3 Upload  
```go
func TestUploadToS3(t *testing.T) {
	err := UploadToS3("test-bucket", "test-file.jpg", []byte("mock data"))
	assert.NoError(t, err)
}
```
✔ **Why?**  
- Mocks an S3 upload and ensures no errors occur.

**For actual AWS integration tests, you can use AWS SDK with environment variables to avoid hardcoding credentials.**

---

#### **Mocking AWS Services with GoMock**
To test AWS interactions without actually calling AWS, use **GoMock**.

1. **Install GoMock:**
   ```sh
   go get github.com/golang/mock/gomock
   go install github.com/golang/mock/mockgen
   ```

2. **Generate Mocks:**
   ```sh
   mockgen -source=s3_service.go -destination=mocks/s3_mock.go -package=mocks
   ```

3. **Example Test (Mocking S3 Upload)**
   ```go
   func TestUploadToS3Mock(t *testing.T) {
       ctrl := gomock.NewController(t)
       defer ctrl.Finish()

       mockS3 := mocks.NewMockS3API(ctrl)
       mockS3.EXPECT().PutObject(gomock.Any()).Return(nil, nil)

       err := UploadToS3(mockS3, "test-bucket", "test-file.jpg", []byte("mock data"))
       assert.NoError(t, err)
   }
   ```
✔ **Why?**  
- Uses `GoMock` to fake an S3 upload, allowing testing without real AWS interaction.

---

#### **End-to-End (E2E) Testing**
E2E tests verify the full workflow from **upload to print**.

1. **Install HTTP Testing Tool**
   ```sh
   go get github.com/go-resty/resty/v2
   ```
   
2. **Simulating a Full User Upload Flow**
   ```go
   func TestFullUploadFlow(t *testing.T) {
       client := resty.New()

       // Step 1: Upload Photo
       resp, err := client.R().
           SetFile("photo", "test.jpg").
           SetFormData(map[string]string{
               "name": "John Doe",
               "email": "john@example.com",
           }).
           Post("http://localhost:8080/upload")

       assert.NoError(t, err)
       assert.Equal(t, 200, resp.StatusCode())

       // Step 2: Check Print Status
       resp, err = client.R().
           Get("http://localhost:8080/status/123")

       assert.NoError(t, err)
       assert.Contains(t, resp.String(), "printed")
   }
   ```
✔ **Why?**  
- Simulates a real user request.
- Ensures the API handles uploads and updates statuses correctly.

---

#### **8.4 Load Testing**
To test performance, use **k6**.

1. **Install k6:**
   ```sh
   brew install k6  # (Mac)
   sudo apt install k6  # (Linux)
   ```

2. **Write Load Test Script (`test.js`)**
   ```js
   import http from "k6/http";
   import { check, sleep } from "k6";

   export default function () {
       let res = http.post("http://localhost:8080/upload", { photo: "test.jpg" });
       check(res, { "status is 200": (r) => r.status === 200 });
       sleep(1);
   }
   ```

3. **Run Load Test:**
   ```sh
   k6 run test.js
   ```
✔ **Why?**  
- Simulates multiple users uploading photos to test scalability.

---

### **8.5 AWS-Specific Debugging**
- **SQS Debugging**
  ```sh
  aws sqs receive-message --queue-url QUEUE_URL
  ```
  ✔ Checks if messages are in SQS.

- **Lambda Debugging**
  ```sh
  aws logs tail /aws/lambda/my-lambda-function
  ```
  ✔ Retrieves logs from AWS Lambda.

- **DynamoDB Debugging**
  ```sh
  aws dynamodb scan --table-name PHOTO_ORDERS
  ```
  ✔ Checks stored orders.

---

### **Final Thoughts**
- ✅ **Unit tests** ensure functions work correctly.
- ✅ **Integration tests** validate AWS interactions.
- ✅ **E2E tests** simulate a full user flow.
- ✅ **Load tests** check system performance.
- ✅ **AWS-specific debugging** ensures cloud resources function properly.

Would you like a **test plan** written in a structured format for the documentation? 🚀