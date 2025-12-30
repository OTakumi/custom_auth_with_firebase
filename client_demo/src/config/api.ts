/**
 * API Configuration
 * Centralized configuration for backend API requests
 */

interface ApiConfig {
  baseUrl: string;
}

const validateApiEnv = (): ApiConfig => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL;

  if (!baseUrl) {
    const errorMsg =
      "Missing VITE_API_BASE_URL environment variable. Please check your .env.local file.";
    console.error(errorMsg);
    throw new Error(errorMsg);
  }

  return { baseUrl };
};

const apiConfig = validateApiEnv();

/**
 * Get the full API endpoint URL
 * @param path - API endpoint path (e.g., "/auth/otp")
 * @returns Full URL
 */
export const getApiUrl = (path: string): string => {
  return `${apiConfig.baseUrl}${path}`;
};

/**
 * API base URL
 */
export const API_BASE_URL = apiConfig.baseUrl;

/**
 * API endpoints
 */
export const API_ENDPOINTS = {
  AUTH: {
    OTP_REQUEST: "/auth/otp",
    OTP_VERIFY: "/auth/verify",
    HEALTH: "/health",
  },
} as const;
