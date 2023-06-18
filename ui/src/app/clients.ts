import { InjectionToken } from "@angular/core";
import { Code, ConnectError, Interceptor, PromiseClient, Transport, createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { AuthService } from "@tkd/apis/gen/es/tkd/idm/v1/auth_service_connect.js";
import { SelfServiceService } from '@tkd/apis/gen/es/tkd/idm/v1/self_service_connect.js';

export const TRANSPORT = new InjectionToken<Transport>('TRANSPORT');
export const AUTH_SERVICE = new InjectionToken<AuthServiceClient>('AUTH_SERVICE');
export const SELF_SERVICE = new InjectionToken<SelfServiceClient>('SELF_SERVICE');

export type AuthServiceClient = PromiseClient<typeof AuthService>;
export type SelfServiceClient = PromiseClient<typeof SelfServiceService>;

const authHeader: Interceptor = (next) => async (req) => {
  const token = localStorage.getItem("access_token")
  if (!!token) {
    req.header.set("Authentication", `Bearer ${token}`)
  }

  return await next(req)
}

const retryRefreshToken: (transport: Transport) => Interceptor = (transport) => (next) => async (req) => {
  try {
    const result = await next(req)
    return result;

  } catch (err) {
    const connectErr = ConnectError.from(err);

    // don't retry the request if it was a Login.
    if (req.service.typeName === AuthService.typeName && req.method.name === 'Login') {
      throw err
    }

    if (connectErr.code === Code.Unauthenticated && (connectErr.details as any) === 'Expired') {
      const cli = createPromiseClient(AuthService, transport);

      console.log(`[DEBUG] call to ${req.service.typeName}/${req.method.name} not authenticated, trying to refresh token`)
      try {
        const token = await cli.refreshToken({})
        if (!!token.accessToken) {
          localStorage.setItem('access_token', token.accessToken.token);
        }
      } catch (refreshErr) {
        console.error("failed to refresh token", refreshErr)

        throw err;
      }

      // retry with a new access token.
      return await next(req);
    }

    throw err;
  }
}

export function transportFactory(): Transport {
  const retryTransport = createConnectTransport({baseUrl: 'http://localhost:8080', credentials: 'include'})

  return createConnectTransport({
    baseUrl: 'http://localhost:8080',
    credentials: 'include',
    interceptors: [
      retryRefreshToken(retryTransport),
      authHeader
    ],
  })
}

export function authServiceFactory(transport: Transport): AuthServiceClient {
  return createPromiseClient(AuthService, transport);
}

export function selfServiceFactory(transport: Transport): SelfServiceClient {
  return createPromiseClient(SelfServiceService, transport);
}

