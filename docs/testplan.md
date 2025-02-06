# SnapFlow Test Plan

## 1. Introduction
This document outlines the testing strategy for **SnapFlow**, covering **unit, integration, end-to-end (E2E), load, and AWS-specific debugging tests**. The goal is to ensure system reliability, performance, and correctness.

---

## 2. Testing Scope
### 2.1 Components Under Test
1. **Backend API (Go/Fiber)**
   - Photo upload handling
   - Image resizing
   - Presigned URL generation
   - DynamoDB updates
   - SQS job submission

2. **AWS Services**
   - **Amazon S3**: Stores processed images.
   - **Amazon DynamoDB**: Tracks user details and order status.
   - **Amazon SQS**: Holds print job requests.
   - **AWS Lambda**: Simulates printing & updates the status.

3. **User Flow**
   - Photo submission
   - Order processing
   - Print status updates

---

## 3. Testing Approach
### 3.1 Unit Testing
- **Purpose**: Validate individual functions in the Go backend.
- **Tooling**: Go `testing` package, Testify.
- **Scope**:
  - `ResizePhoto()`
  - `GeneratePresignedURL()`
  - `UploadToS3()`
  - `SendToSQS()`
  - `UpdateDynamoDB()`
- **Example**:  
  ```go
  func TestResizePhoto(t *testing.T) {
      img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
      buf := new(bytes.Buffer)
      jpeg.Encode(buf, img, nil)

      resizedImg, err := ResizePhoto(buf.Bytes(), 500, 500)

      assert.NoError(t, err)
      assert.NotNil(t, resizedImg)
  }
  ```

---

### 3.2 Integration Testing
- **Purpose**: Validate interactions between backend and AWS services.
- **Tooling**: GoMock, LocalStack (AWS mocking).
- **Scope**:
  - Can the backend upload a photo to S3?
  - Does SQS receive print job messages?
  - Can Lambda process messages and update DynamoDB?
- **Example**:  
  ```go
  func TestUploadToS3(t *testing.T) {
      err := UploadToS3("test-bucket", "test-file.jpg", []byte("mock data"))
      assert.NoError(t, err)
  }
  ```

---

### 3.3 End-to-End (E2E) Testing
- **Purpose**: Simulate the entire user journey.
- **Tooling**: Resty (HTTP API testing).
- **Scope**:
  - Upload photo → Check order status → Verify print completion.
- **Example**:  
  ```go
  func TestFullUploadFlow(t *testing.T) {
      client := resty.New()

      // Step 1: Upload Photo
      resp, err := client.R().
          SetFile("photo", "test.jpg").
          SetFormData(map[string]string{"name": "John"}).
          Post("http://localhost:8080/upload")

      assert.NoError(t, err)
      assert.Equal(t, 200, resp.StatusCode())

      // Step 2: Check Print Status
      resp, err = client.R().Get("http://localhost:8080/status/123")
      assert.Contains(t, resp.String(), "printed")
  }
  ```

---

### 3.4 Load Testing
- **Purpose**: Measure system performance under load.
- **Tooling**: `k6`
- **Scope**:
  - 100 concurrent users uploading images.
  - Peak load test.
- **Example Script (`test.js`)**:
  ```js
  import http from "k6/http";
  import { check, sleep } from "k6";

  export default function () {
      let res = http.post("http://localhost:8080/upload", { photo: "test.jpg" });
      check(res, { "status is 200": (r) => r.status === 200 });
      sleep(1);
  }
  ```
- **Run Test**:  
  ```sh
  k6 run test.js
  ```

---

### 3.5 AWS-Specific Debugging
- **SQS Debugging**  
  ```sh
  aws sqs receive-message --queue-url QUEUE_URL
  ```
  ✅ Checks if messages are in the queue.

- **Lambda Logs**  
  ```sh
  aws logs tail /aws/lambda/my-lambda-function
  ```
  ✅ Retrieves logs for debugging.

- **DynamoDB Verification**  
  ```sh
  aws dynamodb scan --table-name PHOTO_ORDERS
  ```
  ✅ Confirms order status updates.

---

## 4. Pass/Fail Criteria
| Test Type     | Criteria |
|--------------|----------|
| **Unit Tests** | All functions return expected results. |
| **Integration Tests** | API interacts correctly with AWS services. |
| **E2E Tests** | Users can complete the photo upload-to-print workflow. |
| **Load Tests** | API handles 100+ concurrent uploads without failures. |
| **AWS Debugging** | All AWS logs confirm expected behavior. |

---

## 5. Test Execution Plan
| Phase | Description | Owner |
|-------|------------|-------|
| **Unit Testing** | Test core functions. | Backend Dev |
| **Integration Testing** | Ensure AWS interactions work. | DevOps |
| **E2E Testing** | Simulate full user flow. | QA |
| **Load Testing** | Evaluate performance under stress. | QA |
| **AWS Debugging** | Validate cloud behavior. | DevOps |

---

## 6. Conclusion
This test plan ensures **SnapFlow** operates correctly across all components, preventing regressions before making it public.

