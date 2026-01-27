# Project Planning & Documentation Consolidation

**Timestamp**: 2026-01-27 19:30

## Summary

Consolidated multiple documentation files into a structured PLAN.md and updated AGENTS.md for AI agent guidelines. Established project roadmap for evolving LSCS Core from API Key Management to a full-featured web application.

## Problem

The project had scattered documentation across multiple files:
- `IMPROVEMENT_PLAN.md` - Detailed improvement plan
- `TODOs.md` - Comprehensive analysis with issues
- `GEMINI.md` - Outdated project overview
- `REFACTOR_SUMMARY.md` - Past refactoring notes

This made it difficult to track progress and understand current project state.

## Solution

1. Created comprehensive `PLAN.md` with 8 phases:
   - Phase 1: Security & Stability Fixes
   - Phase 2: Foundation Setup (goose, zerolog, swagger)
   - Phase 3: Authentication & Session Management
   - Phase 4: RBAC & Permissions
   - Phase 5: Member Management API
   - Phase 6: Image Upload (S3/Garage)
   - Phase 7: Next.js Frontend
   - Phase 8: Registration & Onboarding (future)

2. Updated `AGENTS.md` with coding guidelines for AI agents

3. Created `logs/` directory for tracking major changes

## Current State Assessment

### Tests
- **Status**: FAILING
- `auth/handler_test.go`: Panic - nil interface conversion
- `member/handler_test.go`: SQL mock column count mismatch
- Coverage: 74.1%

### Security Issues Identified
- CORS allows any origin (`https://*`, `http://*`)
- JWT tokens have no expiration
- `log.Fatalf` in database health check causes panic
- No input validation library

## Files Affected

- `PLAN.md` - Created (comprehensive project plan)
- `AGENTS.md` - Updated (AI agent guidelines)
- `logs/` - Created directory

## Files to Remove (Pending User Confirmation)

- `IMPROVEMENT_PLAN.md`
- `TODOs.md`
- `GEMINI.md`
- `REFACTOR_SUMMARY.md`

## Architecture Decisions

1. **Monorepo Structure**: Go API at root, Next.js in `web/` directory
2. **Auth Strategy**: Go API handles all auth (Google OAuth + sessions for web, JWT for API consumers)
3. **Sessions**: MySQL-based (no Redis for now)
4. **RBAC**: Separate `roles` and `member_roles` tables, not array field on members
5. **Image Upload**: Pre-signed URL approach (backend generates URL, frontend uploads directly to S3)

## Next Steps

1. Begin Phase 1: Security & Stability Fixes
2. Remove deprecated documentation files
3. Fix failing tests
4. Implement CORS, JWT expiration, input validation fixes
