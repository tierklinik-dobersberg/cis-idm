import { InjectionToken } from "@angular/core";
import { Code, ConnectError, Interceptor, PromiseClient, Transport, createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { AuthService, SelfServiceService } from "@tkd/apis";

export const TRANSPORT = new InjectionToken<Transport>('TRANSPORT');
export const AUTH_SERVICE = new InjectionToken<AuthServiceClient>('AUTH_SERVICE');
export const SELF_SERVICE = new InjectionToken<SelfServiceClient>('SELF_SERVICE');

export type AuthServiceClient = PromiseClient<typeof AuthService>;
export type SelfServiceClient = PromiseClient<typeof SelfServiceService>;

const retryRefreshToken: (transport: Transport) => Interceptor = (transport) => {
  let pendingRefresh: Promise<void> | null = null;

  return (next) => async (req) => {
    try {
      const result = await next(req)
      return result;

    } catch (err) {
      const connectErr = ConnectError.from(err);

      // don't retry the request if it was a Login.
      if (req.service.typeName === AuthService.typeName && req.method.name === 'Login') {
        throw err
      }

      if (connectErr.code === Code.Unauthenticated) {
        if (pendingRefresh === null) {
          let _resolve: any;
          let _reject: any;
          pendingRefresh = new Promise((resolve, reject) => {
            _resolve = resolve;
            _reject = reject;
          })

          pendingRefresh
            .catch(() => {})
            .then(() => pendingRefresh = null)

          const cli = createPromiseClient(AuthService, transport);

          console.log(`[DEBUG] call to ${req.service.typeName}/${req.method.name} not authenticated, trying to refresh token`)
          try {
            const token = await cli.refreshToken({})
            if (!!token.accessToken) {
              localStorage.setItem('access_token', token.accessToken.token);
            }

            _resolve();
          } catch (refreshErr) {
            console.error("failed to refresh token", refreshErr)

            _reject(err);

            throw err;
          }
        } else {
          // wait for the pending refresh to finish
          try {
            await pendingRefresh;
          } catch (_) {
            throw err;
          }
        }

        // retry with a new access token.
        return await next(req);
      }

      throw err;
    }
  }
}

export function transportFactory(): Transport {
  const retryTransport = createConnectTransport({baseUrl: 'http://localhost:8080', credentials: 'include'})

  return createConnectTransport({
    baseUrl: 'http://localhost:8080',
    credentials: 'include',
    interceptors: [
      retryRefreshToken(retryTransport),
    ],
  })
}

export function authServiceFactory(transport: Transport): AuthServiceClient {
  return createPromiseClient(AuthService, transport);
}

export function selfServiceFactory(transport: Transport): SelfServiceClient {
  return createPromiseClient(SelfServiceService, transport);
}

