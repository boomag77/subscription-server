# API Documentation

## Endpoints

### 1. Health Check
- **URL**: `/ping`
- **Method**: `GET`
- **Description**: Returns a simple "pong!!! new SSH_KEY" response to verify the server is running.
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**: `pong!!! new SSH_KEY`

---

### 2. App Store Notifications
- **URL**: `/api/v1/notifications/apple/v2`
- **Method**: `POST`
- **Description**: Handles App Store Connect notifications (Server-to-Server).
- **Request**:
  - **Headers**: `Content-Type: application/json`
  - **Body**: JSON payload containing the signed notification data.
- **Response**:
  - **Status Code**: `200 OK` on success, `400 Bad Request` or `500 Internal Server Error` on failure.

---

### 3. Client Notifications (iOS)
- **URL**: `/api/v1/notifications/client/ios`
- **Method**: `POST`
- **Description**: Handles client notifications for iOS.
- **Request**:
  - **Headers**: `Content-Type: application/json`
  - **Body**: JSON payload.
- **Response**:
  - **Status Code**: `200 OK` on success, `405 Method Not Allowed` on invalid method.

---

### 4. Client Notifications (Android)
- **URL**: `/api/v1/notifications/client/android`
- **Method**: `POST`
- **Description**: Handles client notifications for Android.
- **Request**:
  - **Headers**: `Content-Type: application/json`
  - **Body**: JSON payload.
- **Response**:
  - **Status Code**: `200 OK` on success, `405 Method Not Allowed` on invalid method.

---

### 5. Client Request Status (iOS)
- **URL**: `/api/v1/requests/client/ios/status`
- **Method**: `GET`
- **Description**: Retrieves the status of a client request for iOS.
- **Request**:
  - **Query Parameters**:
    - `userToken` (required): The token identifying the user.
- **Response**:
  - **Status Code**: `200 OK` on success, `400 Bad Request` for missing parameters, `500 Internal Server Error` on failure.
  - **Body**:
    ```json
    {
      "expiresAt": "2025-08-29T12:00:00Z",
      "userToken": "user123",
      "productId": "com.example.product",
      "originalTransactionId": "1000000123456789",
      "isActive": true
    }
    ```

---

### 6. Client Request Status (Android)
- **URL**: `/api/v1/requests/client/android/status`
- **Method**: `GET`
- **Description**: Retrieves the status of a client request for Android.
- **Request**:
  - **Query Parameters**:
    - `userToken` (required): The token identifying the user.
- **Response**:
  - **Status Code**: `200 OK` on success, `400 Bad Request` for missing parameters, `500 Internal Server Error` on failure.
  - **Body**:
    ```json
    {
      "expiresAt": "2025-08-29T12:00:00Z",
      "userToken": "user123",
      "productId": "com.example.product",
      "originalTransactionId": "1000000123456789",
      "isActive": true
    }
    ```
