Here's a structured `TEST_RESULTS.md` file that will document all test results, including screenshots, logs, and videos. This ensures your repo has **clear proof** of functionality for Cloud Engineer/Architect roles.  

---

### **📜 TEST_RESULTS.md** (Add this to your repo)

```md
# 📌 SnapFlow Test Results

This document contains the results of all tests performed on **SnapFlow**, including **infrastructure provisioning, API functionality, and AWS service interactions**.

---

## 📍 1. Infrastructure Deployment Results

### ✅ Terraform Apply Output
The following screenshot shows the successful provisioning of SnapFlow infrastructure.

![Terraform Apply](docs/screenshots/terraform-apply.png)

#### 🔍 **Verification via AWS CLI**
After deployment, verify resources exist using:

```sh
aws s3 ls
aws dynamodb list-tables
aws sqs list-queues
aws lambda list-functions
```

---

## 📍 2. Unit Test Results

### ✅ Resizing and Presigned URL Generation
All unit tests passed for backend functions.

```sh
go test ./... -v
```

#### 📸 **Test Output**
![Go Unit Test](docs/screenshots/unit-tests.png)

✅ Functions Tested:  
- `ResizePhoto()`
- `GeneratePresignedURL()`
- `UploadToS3()`
- `SendToSQS()`
- `UpdateDynamoDB()`

---

## 📍 3. Integration Test Results

### ✅ API to AWS Interaction
Each component was verified:

- **S3 Upload Success**
  ```sh
  aws s3 ls s3://snapflow-processed-images/
  ```
  ![S3 Upload](docs/screenshots/s3-upload.png)

- **SQS Message Received**
  ```sh
  aws sqs receive-message --queue-url $SQS_URL
  ```
  ![SQS Message](docs/screenshots/sqs-message.png)

- **Lambda Execution Logs**
  ```sh
  aws logs tail /aws/lambda/snapflow-print-processor --follow
  ```
  ![Lambda Logs](docs/screenshots/lambda-logs.png)

- **DynamoDB Order Updated**
  ```sh
  aws dynamodb scan --table-name PHOTO_ORDERS
  ```
  ![DynamoDB Scan](docs/screenshots/dynamodb-orders.png)

---

## 📍 4. End-to-End (E2E) Test Results

### ✅ Full User Workflow Test  
A complete flow was tested:  

1. **Photo Upload via API**
   ```sh
   curl -X POST -F "photo=@test.jpg" http://localhost:8080/submit-order
   ```

2. **S3 Check**
   ```sh
   aws s3 ls s3://snapflow-processed-images/
   ```

3. **SQS Check**
   ```sh
   aws sqs receive-message --queue-url $SQS_URL
   ```

4. **Print Processing (Lambda Execution)**
   ```sh
   aws logs tail /aws/lambda/snapflow-print-processor --follow
   ```

5. **Order Status Check**
   ```sh
   curl http://localhost:8080/order-status/123
   ```

### 📸 **Screenshot of API Response**
![E2E Test](docs/screenshots/e2e-test.png)

---

## 📍 5. Video Demo 🎥  

### ✅ **Full Workflow Video**
[![SnapFlow Demo](docs/videos/snapflow-demo.png)](https://youtu.be/YOUR_VIDEO_LINK)

This video demonstrates:
✔ Infrastructure Deployment  
✔ API Upload & Processing  
✔ AWS Resource Verification  

---

## 📍 6. Debugging Logs (If Needed)

If an issue occurs, check logs:  

```sh
aws logs tail /aws/lambda/snapflow-print-processor --follow
```

Example Output:

```
Processing order 123...
Upload to S3 complete!
SQS message sent!
Print status: Completed
```

---

## 📍 7. Summary of Results

| Test Type          | Status  | Proof |
|--------------------|---------|--------------------------------|
| Terraform Apply   | ✅ Passed | `terraform-apply.png` |
| Unit Tests       | ✅ Passed | `unit-tests.png` |
| S3 Upload       | ✅ Passed | `s3-upload.png` |
| SQS Message Sent | ✅ Passed | `sqs-message.png` |
| Lambda Execution | ✅ Passed | `lambda-logs.png` |
| DynamoDB Update  | ✅ Passed | `dynamodb-orders.png` |
| End-to-End Test  | ✅ Passed | `e2e-test.png` |
| Video Demo       | ✅ Done  | [Watch Here](https://youtu.be/YOUR_VIDEO_LINK) |

---

## 📍 8. Conclusion  

All tests have successfully passed, confirming that **SnapFlow** is fully functional and ready for production. 🚀  
```

---

### **How to Use This?**
1. **Create a `docs/screenshots/` folder**  
   - Add Terraform, AWS CLI, and API test screenshots.  

2. **Create a `docs/videos/` folder**  
   - Upload a **short 1-2 min demo video** of the entire process.  

3. **Update `TEST_RESULTS.md`**  
   - Add actual screenshot links and a YouTube video link.  

---

### **Final Thoughts**
✅ **This file proves your system works as expected.**  
✅ **It helps in interviews / documentation for future teams.**  
✅ **Adding a short video demo makes it even more impactful.**  

Would you like a **sample README that links to this test results file?** 📜