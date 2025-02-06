### **Senior-Level Documentation Format for SnapFlow**

A well-structured senior-level documentation should be **comprehensive, structured, and technical** while also being easy to navigate. The goal is to ensure that developers, DevOps engineers, and even business stakeholders can understand the project's purpose, architecture, and implementation details.

---

## **Recommended Documentation Structure**

### **1. Introduction**
   - **Project Name**: SnapFlow
   - **Talk about the origin of Snapflow (Company X)**
        - briefly talk about Company X services
        - department and what role you play 
        - talk, briefly about why you started or created this project
   - **Purpose**: A streamlined photo processing and printing system.
   - **Target Audience**: Developers, DevOps engineers, and anyone looking to understand how the system is built and functions.
   - **Tech Stack**: 
     - **Backend**: Go (Fiber)
     - **Frontend**: HTML/CSS (for form submission)
     - **AWS Services**: S3, DynamoDB, SQS, Lambda
     - **Infrastructure**: Terraform (IaC)
   - **High-Level Overview**: Describe the end-to-end process, from user photo upload to printing completion.

---

### **2. System Architecture**
   - **Architecture Diagram (ASCII + Mermaid.js or Draw.io)**
   - **Component Breakdown**:
     - **Frontend**: HTML form for photo upload. (tell us what the form is doing)
     - **Backend API (Go/Fiber)**: Processes images and updates the database. (how is this done?)
     - **AWS Services**: S3 (storage), DynamoDB (metadata storage), SQS (print job queue), Lambda (print simulation). -> briefly highlight how they all align
   - **Flow of Data**:
     - **User → Backend API → S3 → DynamoDB → SQS → Lambda → DynamoDB Update → Customer Pick-Up.** -> give us a representation of this project

---

### **3. Installation and Setup**
   - **Prerequisites**:
     - Go installed (`>=1.18`) -> installation instructions: https://go.dev/doc/install
     - Terraform installed (`>=1.5`) -> installation instructions: https://developer.hashicorp.com/terraform/install
     - AWS CLI configured with required permissions.
   - **Environment Variables**:
     ```
     AWS_ACCESS_KEY_ID=
     AWS_SECRET_ACCESS_KEY=
     SQS_QUEUE_URL=
     DYNAMODB_TABLE_NAME=
     ```
   - **Step-by-Step Setup**:
     - Cloning the repository
        - `git clone git@github.com:30Piraten/snapflow.git`
     - Setting up the `.env` file
        - Use the predifined .env variables above
     - Running the Go backend
        - for the backend, you can have air installed and run `air`(see repo: https://github.com/air-verse/air) from the src/ dir
          or you can run `go run main.go || go run .` from the src/ dir
     - Deploying infrastructure using Terraform
        - to deploy the defined AWS services config with terraform
        run  the following command: 
            - `terraform init && terraform validate` 
            - `terraform plan && terraform apply`
---

### **4. Features and Workflow**
   - **User Upload Process**:
     - How the Go backend processes photo uploads.
     - Pre-signed URL generation for secure S3 uploads.
   - **Storage and Data Handling**:
     - How S3 and DynamoDB store processed data.
   - **Print Job Handling**:
     - How SQS queues and Lambda simulate printing.
   - **DynamoDB Status Updates**:
     - Transition states: `uploaded → processing → printed`
   - **Notifications (If SNS/SES were used)**:
     - Explain why they were omitted for this workflow.

---

### **5. API Reference (Go/Fiber Endpoints)**
   - **List of API Endpoints**:
     - `POST /upload`: Accepts file and user details.
     - `GET /status/{photo_id}`: Retrieves photo processing status.
   - **Request & Response Examples**:
     - JSON request bodies and expected responses.

---

### **6. Deployment and Infrastructure**
   - **Terraform Setup**:
     - How AWS resources are provisioned.
   - **How to Deploy**:
     - Steps for deploying to AWS (S3, Lambda, SQS, DynamoDB).
   - **Scaling Considerations**:
     - How SQS helps with asynchronous processing.

---

### **7. Error Handling & Logging**
   - **Common Errors and Fixes**:
     - DynamoDB update failures.
     - SQS delivery issues.
   - **Logging Strategy**:
     - AWS CloudWatch setup.
     - Go backend logging.

---

### **8. Testing & Validation**
   - **Unit Testing**:
     - Go test cases for API routes.
   - **Integration Testing**:
     - Simulating photo uploads.
   - **End-to-End Testing**:
     - Step-by-step testing of the full workflow.

---

### **9. Future Improvements & Roadmap**
   - **Potential Enhancements**:
     - Add auto-scaling to Lambda.
     - Improve frontend UI.
   - **Performance Optimizations**:
     - Reduce processing time.

---

### **10. Conclusion**
   - Summary of project success.
   - Final considerations before making it public.

--_