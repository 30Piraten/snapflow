### **3. Test S3 Upload with Postman**

To test the presigned S3 URL using Postman, follow these steps:

---

#### **Step 1: Get the Presigned URL**
1. Use Postman to send a `POST` request to your API endpoint (e.g., `/api/upload`).
2. Include the required payload in the request body (e.g., email and photo name).
3. The response should return a `presigned_url` along with metadata like `photo_id`.

---

#### **Step 2: Upload the Photo Using Postman**
1. Open a new tab in Postman.
2. Set the **method** to `PUT`.
3. Paste the `presigned_url` into the URL field.
4. Go to the **Headers** tab and ensure you include:
   - `Content-Type`: Set this to the MIME type of the file (e.g., `image/jpeg`, `image/png`).
5. Go to the **Body** tab:
   - Select **binary**.
   - Click "Select File" and choose the photo you want to upload.
6. Click **Send**.

---

#### **Step 3: Verify the Upload**
1. After sending the request, check the response:
   - A `200 OK` status code confirms a successful upload.
2. Go to the AWS S3 console and confirm the photo is stored in the correct path:
   ```
   s3://your-originals-bucket/customer_email/photo_id
   ```

---

### **6. Full Documentation for Day 1**

Here’s the compiled documentation for everything completed on **Day 1**:

---

#### **Frontend and API Workflow**
1. **Frontend Upload**:
   - A user provides their email and selects a photo to upload.
   - The form submits a `POST` request to the `/api/upload` endpoint.

2. **API Workflow**:
   - **Endpoint**: `/api/upload`
     - **Method**: `POST`
     - **Parameters**:
       - `email`: Customer's email (string).
       - `photo_name`: Name of the uploaded photo (string).
     - **Response**:
       - `presigned_url`: URL to upload the photo to S3 (string).
       - `photo_id`: Unique identifier for the photo (string).
       - `upload_timestamp`: Timestamp of the upload (string).

---

#### **S3 Bucket Configuration**
1. **Bucket Name**: `your-originals-bucket`
2. **Folder Structure**:
   - Photos are stored in folders organized by customer email:
     ```
     s3://your-originals-bucket/customer_email/photo_id
     ```
3. **Lifecycle Rule**:
   - Photos are deleted after 2 days.

4. **Access**:
   - Presigned URLs provide secure upload access.
   - Public access is blocked for the bucket.

---

#### **DynamoDB Table Configuration**
1. **Table Name**: `CustomerPhotos`
2. **Schema**:
   - **Partition Key**: `customer_email` (string).
   - **Sort Key**: `photo_id` (string).
   - **Attributes**:
     - `upload_timestamp` (number): When the photo was uploaded.
     - `photo_status` (string): Current status (e.g., `uploaded`, `processed`).
     - `processed_s3_location` (string): S3 path of the processed photo.

3. **Settings**:
   - **Capacity Mode**: On-demand.
   - **Backup**: Point-in-Time Recovery enabled.

---

#### **Testing Workflow**
1. **Frontend Test**:
   - Submit an email and upload a photo.
   - Verify that the request is processed and a presigned URL is returned.

2. **S3 Test**:
   - Upload a photo using the presigned URL (via Postman or the frontend).
   - Confirm the photo is stored in the correct folder in S3.

3. **DynamoDB Test**:
   - Verify that metadata is stored in DynamoDB after a photo upload:
     - Partition Key: `customer_email`
     - Sort Key: `photo_id`
     - Attributes:
       - `upload_timestamp`
       - `photo_status`: `uploaded`
       - `processed_s3_location`: Empty at this stage.

---

#### **Documentation Notes**
1. **Key Information**:
   - S3 bucket name and structure.
   - DynamoDB table schema and attributes.
   - API endpoint and parameters.
2. **Troubleshooting**:
   - **Issue**: Presigned URL upload fails.
     - Check IAM permissions for S3.
   - **Issue**: Metadata not stored in DynamoDB.
     - Review API logs for errors.

---