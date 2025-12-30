# Frontend Client

Vue 3 + TypeScript frontend for passwordless OTP authentication.

## Tech Stack

- Vue 3.5 (Composition API with `<script setup>`)
- TypeScript 5.9
- Vite 7.3
- Vue Router 4
- Tailwind CSS 4.1 + DaisyUI 5.5
- Firebase SDK 12.7 (Auth, Emulator support)
- pnpm

## Project Structure

```
client_demo/
├── src/
│   ├── components/
│   │   └── auth/              # Auth UI components (8 files)
│   │       ├── OTPRequestForm.vue
│   │       ├── OTPVerifyForm.vue
│   │       ├── SignupForm.vue
│   │       ├── FormInput.vue
│   │       ├── ErrorMessage.vue
│   │       └── SubmitButton.vue
│   ├── composables/
│   │   └── useAuthApi.ts      # Auth API calls (requestOTP, verifyOTP)
│   ├── config/
│   │   ├── api.ts             # API endpoint config
│   │   └── firebase.ts        # Firebase initialization
│   ├── pages/
│   │   ├── Login.vue          # Login/Signup page
│   │   └── Dashboard.vue      # Protected dashboard
│   ├── router/
│   │   └── index.ts           # Router with auth guards
│   └── main.ts                # App entry point
├── .env.example
├── package.json
├── vite.config.ts
└── tsconfig.json
```

## Quick Start

### 1. Configure Environment Variables

```bash
cp .env.example .env.local
```

Edit `.env.local`:

**For Emulator (Development):**

```bash
# Firebase config (dummy values OK with emulator)
VITE_FIREBASE_API_KEY=your_api_key_here
VITE_FIREBASE_AUTH_DOMAIN=demo-project.firebaseapp.com
VITE_FIREBASE_PROJECT_ID=demo-project
VITE_FIREBASE_STORAGE_BUCKET=demo-project.firebasestorage.app
VITE_FIREBASE_MESSAGING_SENDER_ID=123456789
VITE_FIREBASE_APP_ID=1:123456789:web:abcdef

# Emulator
VITE_USE_EMULATOR=true
VITE_AUTH_EMULATOR_URL=http://localhost:9099

# API
VITE_API_BASE_URL=http://localhost:8000
```

**For Production:**

```bash
# Use actual Firebase project credentials
VITE_FIREBASE_API_KEY=<your-production-api-key>
VITE_FIREBASE_AUTH_DOMAIN=<your-project>.firebaseapp.com
# ... (other Firebase config)

VITE_USE_EMULATOR=false
VITE_API_BASE_URL=https://api.yourdomain.com
```

### 2. Install Dependencies

```bash
pnpm install
```

### 3. Start Development Server

```bash
pnpm dev
```

**Client**: <http://localhost:5173>

## Development

### Commands

```bash
pnpm dev       # Start dev server (HMR enabled)
pnpm build     # Production build
pnpm preview   # Preview production build
```

### Build Output

Build artifacts are output to `dist/` directory.

## Features

### Authentication Flow

1. **Signup**: Email + password registration (Firebase standard auth)
2. **Login (OTP)**:
   - Enter email → Request OTP
   - Enter 6-digit OTP → Verify and auto-login
   - Redirect to dashboard

### Components

- **OTPRequestForm**: Email input for OTP request
- **OTPVerifyForm**: 6-digit OTP input
- **SignupForm**: Email + password registration
- **FormInput**: Reusable input field
- **ErrorMessage**: Error display component
- **SubmitButton**: Loading state button

### Router Guards

Authentication guard in `src/router/index.ts`:

- Redirect to `/login` if unauthenticated
- Redirect to `/dashboard` if already logged in (on login page)

### API Integration

Composable `useAuthApi.ts` provides:

- `requestOTP(email)`: Request OTP for email
- `verifyOTP(email, otp)`: Verify OTP and get custom token
- Loading/error state management

## Configuration Files

### `src/config/firebase.ts`

Firebase initialization with emulator support. Validates required environment variables.

### `src/config/api.ts`

API endpoint configuration (`VITE_API_BASE_URL`).

## Troubleshooting

**Environment variable errors:**

- Verify `.env.local` exists with all required variables
- Restart Vite server to reload env vars

**Firebase emulator connection:**

- Check `VITE_USE_EMULATOR=true`
- Verify emulator is running: `docker compose ps`

**OTP not received (dev):**

- Check server console for OTP code
- Dev mode prints OTP to server logs (no actual email sent)

**Build errors:**

```bash
rm -rf node_modules pnpm-lock.yaml
pnpm install
pnpm build
```

See [main README](../README.md) for more details.
