# SnapFlow - Sprint Print Documentation

## 1. Overview
**Project Name:** SnapFlow - Sprint Print

**Description:** A cloud-based photo printing workflow that allows users to upload, process, and print photos, with automated tracking and customer notifications. 

**Target Audience:** Developers, DevOps engineers, and cloud architects.  

**Key Technologies:** Go (Fiber), AWS (S3, DynamoDB, SQS, Lambda, SNS, SES), and Terraform.  

## 1.1. About Company X:
Company X operates through four major divisions:

1.  **Admin - Order and Payment:** The Admin Office is responsible for handling order placement and payment processing. When a customer or photographer arrives, they place an order by providing their details and selecting the desired photo format, which includes choosing the paper type and size—both of which are critical aspects of the order. Once the order is confirmed, a receipt is issued. This receipt and the photos (scanned or uploaded to the computer's hard drive) are then forwarded to the Editing Department. In summary, no editing or further processing occurs until payment and order details have been confirmed.

2.  **Editing:** The Editing Department (where I worked) is responsible for all photo enhancements, ranging from cropping and resizing to color correction and sharpness adjustments. No edits can be made without an order/payment receipt.
Upon receiving a receipt, the editor or retoucher scans or uploads the photos and organizes them into a timestamped folder structure for tracking. This ensures that all edits are documented, allowing for easy rollback in case of errors.
My department specializes in Sprint Print, a budget-friendly, express editing and printing service. Sprint Print prioritizes only two key editing operations:
    - Color correction (I removed the color correction logic, it was unnecessary at this point)
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
5. Although, Sprint print editors do not send notifications to customers, **I added SNS and SES to complete the workflow**.

**Workflow Simulation:**
1. An order is received and stored in DynamoDB (similar to saving details on a hard drive).
2. Editing is performed, and the file is sent for printing.
3. The printing request enters the SQS queue, ensuring orderly processing even if multiple requests arrive simultaneously.
4. Lambda processes the print job and updates DynamoDB with a "Printed" status once completed.
5. Lambda sends a notification to SNS which forwards the message to the subscribed email via SES
6. SES sends an email notification to the customer

Editors in the Sprint Print Department focus solely on editing and sending photos for printing—they do not handle customer notifications. The entire process is optimized for speed and efficiency.


## 2. System Architecture
### 2.1 Workflow Diagram
![snapflow2 arch](/arc/snapflow2.png)

---

### 2.2 Features and Workflow: 

#### 1. **User Upload Process**  
- **Photo Upload and Validation**  
  The frontend is a simple HTML form where users submit their details and upload photos for printing. When a user visits the site (`http://127.0.0.1:1234`), fills out the form, and submits their photos, the Go/Fiber backend processes the request.  
  - The backend validates the user’s details and checks whether the uploaded photos meet the required specifications:
    - Ensures the uploaded files are **JPEG or PNG**.
    - Validates the **MIME type** to prevent invalid formats.
    - Resizes the photos if necessary 
    - Confirms that the dimensions do not exceed 6000 X 6000.  

- **Pre-signed URL for Secure Uploads**  
  Once the validation and resizing are complete, the backend generates a **pre-signed URL** for each processed photo. The pre-signed URL allows the user to **directly upload** the photo to the **S3 bucket** without the backend handling large file transfers, reducing latency and costs. This also ensures that **AWS Lambda is not overloaded**, keeping the function focused on **processing and printing** rather than handling file uploads.  

#### 2. **Storage and Data Handling**  
- **Photo Storage in S3**  
  After validation and resizing, the photos are uploaded to an S3 bucket using pre-signed URLs. Each photo is stored in a structured format:  
  ```
  s3://snapflow-bucket/uploads/{customer_name}/{uniqueFilename}.jpg
  ```
  This ensures **each user’s photos are organized** in separate directories.  

- **User and Order Data in DynamoDB**  
  The Go backend updates the **DynamoDB table** with customer details and assigns an `"uploaded"` status to the order. The table tracks the progress of each photo throughout the workflow, transitioning from:  
  ```
  uploaded → processing → printed
  ```
  Once the print process is completed, the status is updated to `"printed"`, allowing the system to track job completion.  

#### 3. **Print Job Handling**  
- **SQS Message Queue for Print Requests**  
  After the photos are uploaded and customer metadata is stored, the backend triggers a **print request** by sending a message to an **AWS SQS queue** using the [`SendPrintRequest`](./src/config/sqs.go) function. This queue ensures that **each print job is processed in order**, preventing failures due to concurrent requests.  

- **Lambda Processing and Simulated Printing**  
  - The **AWS Lambda function** continuously polls the **SQS queue** for new messages.  
  - Once a print job request is received, Lambda **simulates the printing process** with a short delay to mimic a real-world print operation.  
  - The **print job metadata** is extracted, validated, and passed to the [`ProcessPrintJob`](./src/lambda/lambda.go) function, which updates the **DynamoDB status** to `"printed"`.  

#### 4. **DynamoDB Status Updates**  
Each print job undergoes a **state transition** in DynamoDB:  
```
uploaded → processing → printed
```
- `"uploaded"`: Initial status when photos are validated and uploaded.  
- `"processing"`: Status assigned when a print job is received by the Lambda function.  
- `"printed"`: Final status after the Lambda function completes the print simulation.  

#### 5. **Notifications: SNS and SES**  
- **Order Confirmation Notification (SES via SNS)**  
  - Once a uploads and details are validated and stored, the Go backend Calls the [`SendSNSNotification`](./src/config/sns.go) 
  function, which:  
       - Publishes a message to an **SNS topic** indicating that a new print order has been received.  
       - SNS **forwards** this message to an **SES email subscription**, triggering an email notification to the customer.  
       - The customer receives an **order confirmation email** stating that their photos are **queued for printing**.  

- **Print Completion Notification (SES via SNS)**  
  - Once the **printing process is completed**:
    1. The **AWS Lambda function** updates the **DynamoDB order status** from `"processing"` to `"printed"`.  
    2. The Lambda function **publishes a message** to an **SNS topic**, notifying that the print job is **completed**.  
    3. SNS forwards this message to **SES**, which:  
       - Sends a **final email notification** to the customer, informing them that their **photos are ready for pickup**.  

---

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
   - **SNS & SES:** Sends customer email notifications.

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
  - **Step-by-Step Setup**
    1. Cloning the repository:
        ```
        git clone git@github.com:30Piraten/snapflow.git
        ```
    2. Setting up the `.env` file:
        - Use the predifined .env variables above
    3. Running the Go backend:
        - For the backend, you can have [Air](https://github.com/air-verse/air) installed then run:
        
         ```
         air or 
         go run main.go
         ``` 
         
        from the terminal.
    4. Deploying infrastructure using Terraform:
        - To deploy the defined AWS services config with terraform
        run  the following command:

        ```
        terraform init && terraform validate
        terraform plan && terraform apply
        ```

## 4. Backend API
### 4.1 Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/submit-order` | Uploads photo and customer info |

### 4.2 Key Functions
- [`ProcessPrintJob()`](./src/lambda/lambda.go): Acts as a dummy printer and updates DynamoDB and sends SNS notification.
- [`HandleOrderSubmission()`](./src/routes/order.go): This is the main entry point for the order submission process. 
- [`GeneratePresignedURL()`](./src/url/presigned_url.go): Generates a pre-signed URL for the given order details.
- [`UploadToS3()`](./src/config/s3.go): Uploads resized photos using the pre-signed URL to an S3 bucket. 
- [`SendPrintRequest()`](./src/config/sqs.go): Sends a print job request from the backend to SQS.

## 5. 📌 Testing & Debugging
SnapFlow's testing includes:

- ✅ Terraform Validation & Drift Detection
- ✅ AWS Service-Specific Tests (S3, SQS, Lambda, DynamoDB)
- ✅ Integration Tests (End-to-End workflow)
- ✅ Security Tests (IAM role permissions)

For detailed test cases including AWS debugging, and execution steps see: **[docs](./docs/tests_plan.md)**. 
And for tests results checkout: **[tests](./docs/test_results.md)**

## 5.1. 📌 Workflow Video
Snapflow's video walkthrough: [Snapflow Demo](https://youtu.be/hdQmYGdg_WQ)

Here's a more polished and professional version of your **Future Enhancements** section:  

---

## **6. Future Enhancements** *(Planned for Next Iteration)*  

- **Customer Dashboard for Order Tracking** – Implement a frontend dashboard where customers can track their photo processing status in real time.  
- **SQS Retry Logic for Failures** – Introduce an automatic retry mechanism to handle message failures and improve system resilience.  
- **DynamoDB Update Rollback** – Implement a rollback mechanism to prevent cases where an image is printed but not correctly marked in the database.  
- **DynamoDB Streams Instead of SQS for Order Processing (Maybe)** – Replace SQS with DynamoDB Streams to trigger Lambda functions **only on new orders** (not on every update), reducing costs by leveraging DynamoDB’s built-in event system.  
- **ECS Fargate for Image Resizing** – Offload intensive image resizing from the Go/Fiber backend to **ECS Fargate**. The backend will enqueue jobs to **SQS**, and Fargate workers will handle them asynchronously, improving scalability. With this, the backend only has to do two things: generate presignedURL and uploadToS3.
- **Optimizing SNS & SES for Notifications** – Since Sprint Print customers wait on-site for their photos, an initial confirmation email is not needed. Instead, SES will send a single confirmation email **only after printing is complete**, streamlining the workflow.  

---

## **7. Conclusion**  
SnapFlow is a cloud-native Sprint Print system designed to streamline the photo printing process. It allows customers to upload photos, which are then processed and prepared for printing using a scalable and event-driven architecture.  

The system leverages **AWS services** for efficiency and automation:  
- **S3** for secure photo storage  
- **DynamoDB** for managing metadata  
- **SQS** to handle job queues  
- **Lambda** to process and simulate printing  
- **SNS & SES** for automated notifications  

SnapFlow was built to replicate the end-to-end workflow of my previous role, covering everything from **photo submission, processing, and optimization** to **final printing and customer notifications**. This project demonstrates expertise in **cloud infrastructure, automation, and scalable system design** while optimizing for **cost, performance, and reliability**. 