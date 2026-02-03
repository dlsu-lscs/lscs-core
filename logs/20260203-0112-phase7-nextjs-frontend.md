# Phase 7: Next.js Frontend Setup

**Date**: 2026-02-03  
**Phase**: 7 - Next.js Frontend Setup

## Summary

Initialized Next.js 16 frontend with TypeScript, Tailwind CSS, and shadcn/ui components. Implemented authentication flow, dashboard, and profile management pages.

## Changes

### Project Setup

- Initialized Next.js 16 project in `web/` directory with App Router
- Configured TypeScript with path aliases (`@/*` → `./src/*`)
- Set up Tailwind CSS v4 with shadcn/ui-compatible CSS variables
- Configured pnpm as package manager

### Dependencies Installed

**UI Components** (Radix UI + shadcn/ui pattern):
- `@radix-ui/react-avatar` - Profile image display
- `@radix-ui/react-dialog` - Modal dialogs
- `@radix-ui/react-dropdown-menu` - User menu
- `@radix-ui/react-label` - Form labels
- `@radix-ui/react-select` - Select dropdowns
- `@radix-ui/react-separator` - Visual separator
- `@radix-ui/react-tabs` - Tab navigation
- `@radix-ui/react-toast` - Notifications
- `@radix-ui/react-slot` - Slot component for polymorphism

**State Management**:
- `@tanstack/react-query` + `@tanstack/react-query-devtools` - Server state
- `zustand` - Client state with persistence

**Utilities**:
- `class-variance-authority` - Component variant props
- `clsx` - Conditional class names
- `tailwind-merge` - Tailwind class merging
- `tailwindcss-animate` - Animation utilities
- `lucide-react` - Icon library

### Components Created (`web/src/components/ui/`)

- `button.tsx` - Button with variants (default, destructive, outline, secondary, ghost, link)
- `input.tsx` - Form input component
- `card.tsx` - Card components (Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter)

### Core Library (`web/src/lib/`)

- `utils.ts` - `cn()` utility for class name merging
- `api.ts` - API client with typed endpoints for auth, members, and upload
- `auth-store.ts` - Zustand store for authentication state

### Hooks (`web/src/hooks/`)

- `use-auth.ts` - Authentication hook with React Query integration

### Pages Created

**Login Page** (`/login`):
- Google Sign-In button
- Redirects to backend OAuth flow

**Dashboard** (`/dashboard`):
- Layout with navigation header
- Overview cards showing profile, contact, and academic info

**Profile Page** (`/profile`):
- Profile photo upload with pre-signed URL flow
- Basic information display (read-only)
- Editable contact and details section
- Image deletion support

### Authentication Flow

1. User clicks "Sign in with Google" on `/login`
2. Redirected to Go backend `/auth/google/login`
3. Google OAuth callback sets session cookie
4. Frontend uses session cookie for authenticated requests
5. `/auth/me` endpoint returns current user
6. React Query caches user data with `staleTime: 5min`

### Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Files Created

**Directories**:
- `web/src/components/ui/` - Reusable UI components
- `web/src/lib/` - Utilities and API client
- `web/src/hooks/` - React hooks
- `web/src/app/login/` - Login page
- `web/src/app/dashboard/` - Dashboard pages
- `web/src/app/profile/` - Profile page

**Key Files**:
- `web/src/components/ui/button.tsx`
- `web/src/components/ui/input.tsx`
- `web/src/components/ui/card.tsx`
- `web/src/lib/utils.ts`
- `web/src/lib/api.ts`
- `web/src/lib/auth-store.ts`
- `web/src/hooks/use-auth.ts`
- `web/src/components/providers.tsx`
- `web/src/app/login/page.tsx`
- `web/src/app/dashboard/layout.tsx`
- `web/src/app/dashboard/page.tsx`
- `web/src/app/profile/page.tsx`

## Updated Files

- `AGENTS.md` - Updated to include frontend guidelines and monorepo structure

## Next Steps (Phase 8)

Per PLAN.md, Phase 8 involves **Registration & Onboarding**:
- Registration requests table
- Self-registration flow
- Admin approval workflow
- Welcome flow after approval

However, there's also opportunity to:
- Add Members list page
- Add Member detail/edit pages with RBAC
- Add API Key dashboard for RND members
- Add MagicUI animations
