# SnapFlow - Sprint Print Documentation

## 1. Overview
**Project Name:** SnapFlow - Sprint Print (Sprint Print is a budget-friendly, express editing and printing service.)

**Description:** A cloud-based photo printing workflow that allows users to upload, process, and print photos, with automated tracking and customer notifications. 

**Target Audience:** Developers, DevOps engineers, and cloud architects.  

**Key Technologies:** Go (Fiber), AWS (S3, DynamoDB, SQS, Lambda, SNS, SES), and Terraform.  

## 1.1. About Company X:
There are four major divisons in Company X.

1.  **Admin - Order and Payment:** The Admin Office is responsible for handling order placement and payment processing. When a customer or photographer arrives, they place an order by providing their details and selecting the desired photo format, which includes choosing the paper type and size—both of which are critical aspects of the order. Once the order is confirmed, a receipt is issued. This receipt and the photos (scanned or uploaded to the computer's hard drive) are then forwarded to the Editing Department. In summary, no editing or further processing occurs until payment and order details have been confirmed.

2.  **Editing:** The Editing Department (where I worked) is responsible for all photo enhancements, ranging from cropping and resizing to color correction and sharpness adjustments. No edits can be made without an order/payment receipt.
Upon receiving a receipt, the editor or retoucher scans or uploads the photos and organizes them into a timestamped folder structure for tracking. This ensures that all edits are documented, allowing for easy rollback in case of errors.
My department specializes in Sprint Print, a budget-friendly, express editing and printing service. Sprint Print prioritizes only two key editing operations:
    - Color correction
    - Resizing

    Photos designated for Sprint Print are marked accordingly before being sent to the Printing Department.
    
3.  **Printing:** Once the Editing Department sends the edited photos for printing, the Sprint Print workflow is initiated (only for marked photos). The key difference with Sprint Print is that no notifications are sent to customers. Instead, customers wait to collect their photos, often processed within minutes.
Sprint Print is intended for small orders (typically fewer than five photos), allowing for a low-cost, high-speed workflow with minimal overhead.
   
4.  **Delivery | collection:** The Delivery Department manages the collection and delivery of printed photos, including framed orders. For Sprint Print, customers collect their photos directly from this department immediately after printing is completed.

## 1.2. Purpose of the Project:
The goal was to simulate the Sprint Print workflow in the Editing Department using Amazon Web Services (AWS). The system is designed as follows:

1. Amazon S3 stores photo files (blob storage).
2. Amazon DynamoDB stores order details and metadata.
3. AWS Lambda simulates the printing process (and could be integrated with a physical printer).
4. Amazon SQS queues print requests, decoupling the workflow for scalability.
5. I added SNS and SES to make the workflow more robust. (SNS handles notification service, while SES delivers the email). The SNS operation is handled by Lambda.

**Workflow Simulation:**
1. An order is received and stored in DynamoDB (similar to saving details on a hard drive).
2. Editing is performed, and the file is sent for printing.
3. The printing request enters the SQS queue, ensuring orderly processing even if multiple requests arrive simultaneously.
4. Lambda processes the print job and updates DynamoDB with a "Printed" status once completed.

Editors in the Sprint Print Department focus solely on editing and sending photos for printing—they do not handle customer notifications. The entire process is optimized for speed and efficiency.


## 2. System Architecture
### 2.1 Workflow Diagram
![snapflow2 arch](/arc/snapflow2.png)

### **2.2 Features and Workflow**
   - **User Upload Process**:
     - How the Go backend processes photo uploads.
     - Pre-signed URL generation for secure S3 uploads.
   - **Storage and Data Handling**:
     - How S3 and DynamoDB store processed data.
   - **Print Job Handling**:
     - How SQS queues and Lambda simulate printing.
   - **DynamoDB Status Updates**:
     - Transition states: `uploaded → processing → printed`
   - **Notifications SNS/SES**:
     - Explain how they are used.

### 2.3 Component Breakdown
1. **User Uploads:**
   - Basic HTML form for photo and user input. 
2. **Go/Fiber Backend:**
   - Processes photo, generates a signed URL, uploads to S3, updates DynamoDB, and sends print request to SQS.
3. **AWS Services:**
   - **S3:** Stores resized photos.
   - **DynamoDB:** Tracks customer data and photo status.
   - **SQS:** Holds print job requests.
   - **Lambda:** Processes print jobs, updates status, sends notifications.
   - **SNS & SES:** Sends customer email notifications. (SES has not been included)

## 3. Installation & Setup
### 3.1 Prerequisites
- [Go](https://go.dev/doc/install) installed
- [Air](https://github.com/air-verse/air) installed
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) configured
- [Terraform](https://developer.hashicorp.com/terraform/) installed
- Required AWS permissions

### 3.2 Setting Up the Environment
- Clone repository
- Configure `.env` file with:
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY`
  - `S3_BUCKET_NAME`
  - `DYNAMODB_TABLE_NAME`
  - `SQS_QUEUE_URL`
  - `SNS_TOPIC_ARN`

### 3.3 Deploying AWS Infrastructure
  - **Step-by-Step Setup**:
     - Cloning the repository
        - `git clone git@github.com:30Piraten/snapflow.git`
     - Setting up the `.env` file
        - Use the predifined .env variables above
     - Running the Go backend
        - for the backend, you can have air installed and run `air`(see repo: ) from the src/ dir
          or you can run `go run main.go || go run .` from the src/ dir
     - Deploying infrastructure using Terraform
        - to deploy the defined AWS services config with terraform
        run  the following command: 
            - `terraform init && terraform validate` 
            - `terraform plan && terraform apply`


## 4. Backend API
### 4.1 Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/upload` | Uploads photo and customer info |
| GET | `/status/:photo_id` | Retrieves print status |

### 4.2 Key Functions -> tell us what the key functions do 
- `ResizePhoto()`
- `GeneratePresignedURL()`
- `UploadToS3()`
- `UpdateDynamoDB()`
- `SendToSQS()`

### 5. Error Handling -> add to future enhancements 
- Retry logic for SQS failures
- DynamoDB update rollback

## 6. Testing & Debugging
### 6.1 Unit Tests
- Test `ResizePhoto()` function
- Test `GeneratePresignedURL()`

### 6.2 Integration Tests
- Simulate photo upload
- Manually check S3 and DynamoDB entries

### 6.3 AWS-Specific Debugging
- Checking SQS logs
- Debugging Lambda via CloudWatch

## 7. Deployment & CI/CD -> tied to next project
- Use GitHub Actions for automated deployment.
- Run Terraform in CI/CD for AWS infrastructure updates.

## 8. Future Enhancements
- Implement real printing service.
- Add customer dashboard for tracking.

## 9. Conclusion
- Recap of features and workflow.
- How to contribute to the project.
:: This is recap of my 

