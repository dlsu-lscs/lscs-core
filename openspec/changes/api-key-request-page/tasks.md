## 1. Backend API Client

- [ ] 1.1 Add `requestKey` function to `web/src/lib/api.ts`
- [ ] 1.2 Add `RequestKeyRequest` and `RequestKeyResponse` TypeScript interfaces
- [ ] 1.3 Note: Endpoint uses Google Auth, not session cookies

## 2. Page Structure

- [ ] 2.1 Create `web/src/app/request-key/` directory
- [ ] 2.2 Create `page.tsx` with AuthenticatedLayout wrapper
- [ ] 2.3 Add loading spinner state
- [ ] 2.4 Add back navigation to `/api-keys`

## 3. Authentication Flow

- [ ] 3.1 Check authentication state on page load
- [ ] 3.2 Implement Google OAuth redirect for unauthenticated users
- [ ] 3.3 Handle redirect param from OAuth callback
- [ ] 3.4 Check RND AVP+ authorization

## 4. Form Implementation

- [ ] 4.1 Add project name input field
- [ ] 4.2 Add key type selector (Dev/Prod/Admin radio buttons)
- [ ] 4.3 Add allowed origin input field
- [ ] 4.4 Implement conditional required validation based on key type
- [ ] 4.5 Add origin validation (localhost for dev, https for prod)

## 5. Form Submission & Display

- [ ] 5.1 Connect form to `api.requestKey` function
- [ ] 5.2 Handle loading state during submission
- [ ] 5.3 Display success state with generated API key
- [ ] 5.4 Add copy-to-clipboard functionality
- [ ] 5.5 Handle and display error messages
