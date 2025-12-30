import { ref } from "vue";
import { getApiUrl, API_ENDPOINTS } from "../config/api";

/**
 * OTP Request Response
 */
interface OTPRequestResponse {
  message: string;
  otp?: string; // OTP code (development only)
}

/**
 * OTP Verify Response
 */
interface OTPVerifyResponse {
  token: string; // Firebase Custom Token
}

/**
 * API Error Response
 */
interface ApiErrorResponse {
  error: string;
}

/**
 * Composable for authentication API calls
 */
export const useAuthApi = () => {
  const loading = ref(false);
  const error = ref<string>("");

  /**
   * Request OTP code
   * @param email - User email address
   * @returns OTP request response
   */
  const requestOTP = async (
    email: string,
  ): Promise<OTPRequestResponse | null> => {
    loading.value = true;
    error.value = "";

    try {
      const response = await fetch(getApiUrl(API_ENDPOINTS.AUTH.OTP_REQUEST), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      });

      if (!response.ok) {
        const errorData: ApiErrorResponse = await response.json();
        error.value = errorData.error || "OTPリクエストに失敗しました";
        return null;
      }

      const data: OTPRequestResponse = await response.json();
      return data;
    } catch (err) {
      error.value = "ネットワークエラーが発生しました";
      console.error("OTP request error:", err);
      return null;
    } finally {
      loading.value = false;
    }
  };

  /**
   * Verify OTP code and get custom token
   * @param email - User email address
   * @param otp - OTP code
   * @returns Custom token
   */
  const verifyOTP = async (
    email: string,
    otp: string,
  ): Promise<string | null> => {
    loading.value = true;
    error.value = "";

    try {
      const response = await fetch(getApiUrl(API_ENDPOINTS.AUTH.OTP_VERIFY), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, otp }),
      });

      if (!response.ok) {
        const errorData: ApiErrorResponse = await response.json();
        error.value = errorData.error || "OTP検証に失敗しました";
        return null;
      }

      const data: OTPVerifyResponse = await response.json();
      return data.token;
    } catch (err) {
      error.value = "ネットワークエラーが発生しました";
      console.error("OTP verify error:", err);
      return null;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    requestOTP,
    verifyOTP,
  };
};
