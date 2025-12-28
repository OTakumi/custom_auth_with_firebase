import { initializeApp } from "firebase/app";
import { getAuth, connectAuthEmulator } from "firebase/auth";

interface FirebaseEnv {
  VITE_FIREBASE_API_KEY: string;
  VITE_FIREBASE_AUTH_DOMAIN: string;
  VITE_FIREBASE_PROJECT_ID: string;
  VITE_FIREBASE_STORAGE_BUCKET: string;
  VITE_FIREBASE_MESSAGING_SENDER_ID: string;
  VITE_FIREBASE_APP_ID: string;
  VITE_FIREBASE_MEASUREMENT_ID?: string;
  VITE_USE_EMULATOR?: string;
  VITE_AUTH_EMULATOR_URL?: string;
}

const validateEnv = (): FirebaseEnv => {
  const requiredVars = [
    "VITE_FIREBASE_API_KEY",
    "VITE_FIREBASE_AUTH_DOMAIN",
    "VITE_FIREBASE_PROJECT_ID",
    "VITE_FIREBASE_STORAGE_BUCKET",
    "VITE_FIREBASE_MESSAGING_SENDER_ID",
    "VITE_FIREBASE_APP_ID",
  ] as const;

  const missing = requiredVars.filter((key) => !import.meta.env[key]);

  if (missing.length > 0) {
    const errorMsg = `Missing required Firebase environment variables: ${missing.join(", ")}. Please copy .env.example to .env.local and configure your Firebase credentials.`;
    console.error(errorMsg);
    throw new Error(errorMsg);
  }

  return import.meta.env as unknown as FirebaseEnv;
};

/**
 * Initialize Firebase application with environment-based configuration
 * Supports both production Firebase and local Emulator
 */
const initializeFirebase = () => {
  const env = validateEnv();

  const firebaseConfig = {
    apiKey: env.VITE_FIREBASE_API_KEY,
    authDomain: env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: env.VITE_FIREBASE_APP_ID,
    measurementId: env.VITE_FIREBASE_MEASUREMENT_ID,
  };

  const app = initializeApp(firebaseConfig);
  const auth = getAuth(app);

  // Connect to Firebase Emulator if enabled
  const useEmulator = env.VITE_USE_EMULATOR === "true";
  if (useEmulator) {
    const authEmulatorUrl =
      env.VITE_AUTH_EMULATOR_URL || "http://localhost:9099";
    connectAuthEmulator(auth, authEmulatorUrl, { disableWarnings: true });
    console.log(`🔧 Firebase Auth connected to Emulator: ${authEmulatorUrl}`);
  } else {
    console.log("🔥 Firebase Auth connected to production");
  }

  return { app, auth };
};

const { app, auth } = initializeFirebase();
export { app, auth };
