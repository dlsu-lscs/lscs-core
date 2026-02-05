## 1. Backend API Client

- [x] 1.1 Add `requestKey` function to `web/src/lib/api.ts`
- [x] 1.2 Add `RequestKeyRequest` and `RequestKeyResponse` TypeScript interfaces
- [x] 1.3 Note: Endpoint uses Google Auth, not session cookies

## 2. Page Structure

- [x] 2.1 Create `web/src/app/request-key/` directory
- [x] 2.2 Create `page.tsx` with AuthenticatedLayout wrapper
- [x] 2.3 Add loading spinner state
- [x] 2.4 Add back navigation to `/api-keys`

## 3. Authentication Flow

- [x] 3.1 Check authentication state on page load
- [x] 3.2 Implement Google OAuth redirect for unauthenticated users
- [x] 3.3 Handle redirect param from OAuth callback
- [x] 3.4 Check RND AVP+ authorization

## 4. Form Implementation

- [x] 4.1 Add project name input field
- [x] 4.2 Add key type selector (Dev/Prod/Admin radio buttons)
- [x] 4.3 Add allowed origin input field
- [x] 4.4 Implement conditional required validation based on key type
- [x] 4.5 Add origin validation (localhost for dev, https for prod)

## 5. Form Submission & Display

- [x] 5.1 Connect form to `api.requestKey` function
- [x] 5.2 Handle loading state during submission
- [x] 5.3 Display success state with generated API key
- [x] 5.4 Add copy-to-clipboard functionality
- [x] 5.5 Handle and display error messages
