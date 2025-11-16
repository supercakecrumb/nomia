# AGENT.md - AI Assistant Guide for Affirm Name Project

This guide is for AI assistants working with the human developer on the Affirm Name project. Read this carefully before starting work.

---

## Project Overview

**Affirm Name** is a name exploration app for trans/nonbinary users to discover gender-neutral names across different countries and time periods.

**Mission:** Help people choose names that affirm their identity while navigating real-world legal/social systems.

**Current Status:** Phase 1 MVP (Frontend) - COMPLETE ✅

---

## Project Structure

```
affirm-name-frontend/
├── AGENT.md                    # This file
├── README.md                   # Project overview + status tracker
├── ARCHITECTURE.md             # Full system architecture
├── architecture/               # Detailed design docs
│   ├── 00-overview.md
│   ├── 01-shared-contract.md
│   ├── 02-backend-carcass.md
│   ├── 03-frontend-carcass.md  # ← Read this for frontend work
│   ├── 04-integration-strategy.md
│   ├── 05-cross-cutting-concerns.md
│   └── 06-development-workflow.md
├── frontend/                   # React app (main work area)
│   ├── README.md              # Frontend-specific guide
│   ├── src/
│   │   ├── api/               # API client (mock & real modes)
│   │   ├── components/        # React components
│   │   ├── hooks/             # Custom hooks
│   │   ├── i18n/              # i18next config
│   │   ├── pages/             # Page components
│   │   ├── types/             # TypeScript types
│   │   └── utils/             # Helper functions
│   └── public/locales/        # Translations (en, ru)
├── spec-examples/             # Mock data fixtures
└── migrations/                # Database migrations (future)
```

---

## Critical Information

### 1. The Human Developer's Preferences

**Communication Style:**
- ❌ **NEVER** start with "Great", "Certainly", "Okay", "Sure"
- ✅ Be direct, technical, and to the point
- ✅ "I've updated the CSS" not "Great! I've updated the CSS"
- ✅ Explain decisions clearly without being conversational

**Code Style:**
- No excessive comments - code should be self-documenting
- TypeScript strict mode - everything properly typed
- React 19 patterns - functional components, hooks
- Clean, minimal implementations preferred

**What They Care About:**
- **User Experience** - Smooth, intuitive, accessible
- **Performance** - Backend is slow (7-8s), so optimize frontend
- **Visual Design** - Soft pastels, no harsh colors, smooth gradients
- **Production Quality** - No placeholders, no half-finished features

### 2. Key Technical Decisions

#### React Query (TanStack Query)
- Used for ALL data fetching
- Handles caching, request cancellation, loading states
- **Never** use native fetch without React Query

#### Debouncing Strategy
- **1000ms** debounce for filters (backend is slow!)
- Uses custom `useDebounce` hook
- Applied to: year range, gender balance, popularity filters

#### API Modes
- **Mock Mode**: Uses fixtures from `spec-examples/*.json`
- **Real Mode**: Connects to backend at `http://localhost:8080`
- Controlled by `VITE_API_MODE` in `.env.development`

#### URL State Synchronization
- All filter state synced to URL query params
- Users can bookmark/share filtered views
- Handled by `useFilters` hook

#### Internationalization (i18n)
- English and Russian supported
- Translation files in `public/locales/en|ru/*.json`
- **Always** add translations for both languages
- Use `useTranslation` hook, never hardcoded strings

---

## Known Issues & Workarounds

### 1. Backend Performance
**Problem:** API responses take 7-8 seconds per request

**Solutions Applied:**
- 1000ms debouncing on all filters
- Request cancellation with AbortSignal
- Loading banner to show progress
- Smart caching (30s stale time)
- Keep previous data visible during loading

### 2. Popularity Filters Disabled
**Location:** `frontend/src/components/filters/FilterBar.tsx`

**Why:** Backend too slow for good UX with complex filters

**Code:** Commented out with TODO markers

**Re-enable when:** Backend response time <500ms

### 3. Name Search Disabled
**Location:** `frontend/src/components/filters/FilterBar.tsx`

**Why:** Not yet implemented in backend

**Code:** Commented out with TODO markers

**Re-enable when:** Backend implements glob pattern matching

### 4. Tailwind CSS Version
**Using:** Tailwind CSS 3.4.1 (NOT v4)

**Why:** Initial v4 setup had rendering issues

**Important:** Don't suggest upgrading to v4 without testing

---

## Common Patterns

### Adding a New Filter

1. **Update `useFilters` hook** (`frontend/src/hooks/useFilters.ts`):
   - Add to `FilterState` interface
   - Add state getter/setter
   - Add to URL sync effect
   - Add to `getApiParams()` method
   - Add to `resetFilters()` method

2. **Create filter component** (`frontend/src/components/filters/`):
   - Use Tailwind CSS for styling
   - Add translations to `public/locales/*/filters.json`
   - Handle both controlled value and onChange

3. **Add to FilterBar** (`frontend/src/components/filters/FilterBar.tsx`):
   - Import component
   - Add props to interface
   - Render in appropriate section

4. **Update API types** (`frontend/src/types/api.ts`):
   - Add parameter to `NamesFilterParams`

### Working with Gradients

**Important Pattern:** The gender balance bars use a **2x-wide shifting gradient**:

```typescript
style={{
  background: `linear-gradient(...)`,
  backgroundSize: '200% 100%',  // 2x wide
  backgroundPosition: `${value}% center`,  // Shifts with data
}}
```

**Why:** Avoids visible seams from overlapping colored divs

**Don't:** Create multiple overlay divs with different widths - causes borders

### React Query Pattern

```typescript
import { useQuery } from '@tanstack/react-query';

export const useSomeData = (params: Params) => {
  return useQuery({
    queryKey: ['someData', params],
    queryFn: ({ signal }) => fetchSomeData(params, signal),
    staleTime: 30000,  // 30 seconds
    gcTime: 300000,    // 5 minutes
    placeholderData: (prev) => prev,  // Keep old data visible
  });
};
```

**Always:**
- Include `signal` for request cancellation
- Use appropriate cache times
- Use `placeholderData` for smooth transitions

---

## Phase 2 Features (NOT YET IMPLEMENTED)

### Popularity Filters (Ready to Enable)
**Files:**
- `frontend/src/components/filters/PopularityFilterTrio.tsx` (exists)
- `frontend/src/components/filters/FilterBar.tsx` (commented out)
- `frontend/src/pages/NamesExplorerPage.tsx` (commented out)
- `frontend/src/hooks/useFilters.ts` (logic exists)

**Requirements:**
- Backend must respond in <500ms
- Test thoroughly before enabling

### Name Search (Not Implemented)
**Requirements:**
- Backend must implement glob pattern matching (*, ?)
- Add debouncing (1000ms)
- Case-insensitive matching

### Future Enhancements (frontend/README.md)
- Interactive trend charts on detail pages
- Dark mode support
- Export results (CSV/JSON)
- Loading skeletons
- Better mobile responsiveness
- Accessibility improvements

---

## Testing Checklist

Before claiming work is complete:

### Visual
- [ ] Check in both English and Russian
- [ ] Test at different screen sizes (desktop, tablet, mobile)
- [ ] Verify gradients are smooth (no borders/seams)
- [ ] Check color contrast and accessibility
- [ ] Test all interactive states (hover, focus, active)

### Functional
- [ ] All filters work correctly
- [ ] URL updates when filters change
- [ ] Filters persist on page reload
- [ ] Loading states show appropriately
- [ ] Error states handled gracefully
- [ ] Pagination works correctly

### Performance
- [ ] No console errors or warnings
- [ ] TypeScript compiles without errors
- [ ] Debouncing reduces API calls
- [ ] Previous data stays visible during loading
- [ ] Request cancellation works

### Code Quality
- [ ] Types properly defined
- [ ] No hardcoded strings (use i18n)
- [ ] Consistent code style
- [ ] No commented-out code (unless with TODO)
- [ ] Translations added for both languages

---

## Important Files to Know

### Core Logic
- `frontend/src/hooks/useFilters.ts` - **Most important** - manages all filter state
- `frontend/src/api/client.ts` - API client with mock/real modes
- `frontend/src/types/api.ts` - TypeScript types for all API responses

### UI Components
- `frontend/src/components/filters/FilterBar.tsx` - Main filter UI container
- `frontend/src/components/table/NamesTable.tsx` - Names display table
- `frontend/src/components/table/GenderBalanceBar.tsx` - Gradient visualization
- `frontend/src/pages/NamesExplorerPage.tsx` - Main explorer page

### Configuration
- `frontend/.env.development` - Environment variables (API mode, URL)
- `frontend/vite.config.ts` - Vite configuration (aliases, etc.)
- `frontend/tailwind.config.js` - Custom color palette
- `frontend/src/i18n/config.ts` - i18next setup

---

## Color Palette

**Primary (Purple):** `#8B5CF6` - Main brand color  
**Secondary (Teal):** `#14B8A6` - Accent color  
**Accent (Coral):** `#F472B6` - Highlights

**Gender Colors (Soft Pastels):**
- Female: Pink `#FFB3D9` → `#FFE5F0`
- Neutral: Lavender `#E6D9FF` → `#F5F0FF`
- Male: Blue `#99D5FF` → `#E5F2FF`

**Important:** All gender colors are **soft pastels**, never harsh/bright

---

## Do's and Don'ts

### ✅ DO
- Read architecture docs before making changes
- Test in both mock and real API modes
- Add translations for both EN and RU
- Use React Query for data fetching
- Follow existing patterns and conventions
- Ask clarifying questions if uncertain
- Explain technical decisions clearly

### ❌ DON'T
- Use conversational filler ("Great!", "Sure!")
- Hardcode strings (use i18n)
- Skip TypeScript types (strict mode)
- Create multiple colored overlay divs (causes seams)
- Suggest upgrading Tailwind to v4
- Add comments stating the obvious
- Leave console.log statements
- Use `any` type in TypeScript
- Forget to update both language translations

---

## Development Workflow

1. **Read relevant architecture docs**
2. **Understand the existing patterns**
3. **Make changes following conventions**
4. **Test in browser** (http://localhost:5173)
5. **Verify TypeScript compiles**
6. **Check both languages**
7. **Report completion with clear summary**

### Running the App

```bash
cd frontend
npm install
npm run dev  # Starts on http://localhost:5173
```

**With Real API:**
```bash
VITE_API_MODE=real npm run dev
```

### Important Commands

```bash
npm run dev      # Development server
npm run build    # Production build
npm run preview  # Preview production build
npm run lint     # Lint check
npm run type-check  # TypeScript check
```

---

## Communication Templates

### When Starting Work
```
I'll [action] by [approach]. This will [outcome].
```

### When Asking Questions
```
I need clarification on [specific thing]. Should I [option A] or [option B]?
```

### When Completing Work
```
I've [completed action]. Changes:
1. [specific change]
2. [specific change]

Result: [outcome]
```

### When Encountering Issues
```
I found [issue]. The problem is [cause]. 
I can fix it by [solution]. Should I proceed?
```

---

## Key Insights from Phase 1

### What Worked Well
- 1000ms debouncing effectively reduced API load
- React Query caching dramatically improved perceived performance
- 2x-wide shifting gradient eliminated border seams
- URL state sync made the app shareable/bookmarkable
- i18next made bilingual support straightforward

### Lessons Learned
- Backend performance dictates frontend strategy
- Multiple gradient overlays create visible seams
- Soft pastels > bright colors for gender data
- Accessibility matters from day one
- Clear documentation enables faster iteration

---

## Getting Help

1. **Architecture Questions:** Read `architecture/*.md` files
2. **Frontend Specifics:** Check `frontend/README.md`
3. **API Details:** See `architecture/01-shared-contract.md`
4. **Design Patterns:** Look at existing components
5. **Unknown Territory:** Ask the human developer

---

## Success Metrics

Your work is successful when:
- ✅ Features work smoothly in the browser
- ✅ Code follows existing patterns
- ✅ TypeScript compiles without errors
- ✅ No console warnings or errors
- ✅ Both languages display correctly
- ✅ Performance is maintained or improved
- ✅ The human developer is satisfied

---

## Final Notes

This is a **real project** serving **real users** with **real needs**. The trans and nonbinary community deserves tools that are:
- **Respectful** - Treating gender as a spectrum, not a binary
- **Accessible** - Everyone can use it, regardless of ability
- **Reliable** - It works well and consistently
- **Beautiful** - Soft, welcoming, affirming design

Your work directly impacts people making one of the most important decisions of their lives. Take it seriously, do it well.

---

**Last Updated:** November 2025  
**Phase:** 1 (MVP Frontend) - Complete  
**Next Phase:** Backend optimization → Re-enable popularity/search filters

Good luck! 🎨✨