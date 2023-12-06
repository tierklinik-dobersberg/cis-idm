import { InjectionToken } from "@angular/core";
import { ActivatedRoute, Router } from "@angular/router";
import { Code, ConnectError, Interceptor, PromiseClient, Transport, createPromiseClient } from "@bufbuild/connect";
import { createConnectTransport } from "@bufbuild/connect-web";
import { AuthService, SelfServiceService } from "@tierklinik-dobersberg/apis";

export const TRANSPORT = new InjectionToken<Transport>('TRANSPORT');
export const AUTH_SERVICE = new InjectionToken<AuthServiceClient>('AUTH_SERVICE');
export const SELF_SERVICE = new InjectionToken<SelfServiceClient>('SELF_SERVICE');

export type AuthServiceClient = PromiseClient<typeof AuthService>;
export type SelfServiceClient = PromiseClient<typeof SelfServiceService>;

const retryRefreshToken: (transport: Transport, activatedRoute: ActivatedRoute, router: Router) => Interceptor = (transport, activatedRoute, router) => {
  let pendingRefresh: Promise<void> | null = null;

  return (next) => async (req) => {
    try {
      const result = await next(req)
      return result;

    } catch (err) {
      const connectErr = ConnectError.from(err);

      // don't retry the request if it was a Login or RefreshToken.
      if (req.service.typeName === AuthService.typeName && (req.method.name === 'Login' || req.method.name == 'RefreshToken')) {
        console.log("skipping retry as requested service is " + `${req.service.typeName}/${req.method.name}`)

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
            let redirect = activatedRoute.snapshot.queryParamMap.get("redirect");
            if (!redirect && router.getCurrentNavigation() !== null) {
              redirect = router.getCurrentNavigation()!.extractedUrl.queryParamMap.get("redirect")
            }

            const res = await cli.refreshToken({
              requestedRedirect: redirect || '',
            })

            if (res.redirectTo !== '') {
              window.location.href = res.redirectTo;
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

export function transportFactory(route: ActivatedRoute, router: Router): Transport {
  const retryTransport = createConnectTransport({baseUrl: '', credentials: 'include'})

  return createConnectTransport({
    baseUrl: '',
    credentials: 'include',
    jsonOptions: {
      ignoreUnknownFields: true
    },
    interceptors: [
      retryRefreshToken(retryTransport, route, router),
    ],
  })
}

export function authServiceFactory(transport: Transport): AuthServiceClient {
  return createPromiseClient(AuthService, transport);
}

export function selfServiceFactory(transport: Transport): SelfServiceClient {
  return createPromiseClient(SelfServiceService, transport);
}

