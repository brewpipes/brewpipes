/**
 * Authentication and authorization types.
 */

/**
 * JWT token pair returned from login/refresh endpoints.
 */
export interface AuthTokens {
  access_token: string
  refresh_token: string
}
