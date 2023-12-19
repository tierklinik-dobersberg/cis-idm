import { inject } from "@angular/core";
import { ActivatedRouteSnapshot, Router } from "@angular/router";
import { map, switchMap } from "rxjs";
import { ProfileService } from "src/services/profile.service";

export const authGuard = (route: ActivatedRouteSnapshot) => {
  const profileService = inject(ProfileService);
  const router = inject(Router);

  return profileService.ready
    .pipe(
      switchMap(() => profileService.profile),
      map(value => {
        console.log(route.routeConfig?.path)

        return value === null ? router.navigate(['/login'], {
          queryParamsHandling: 'merge',
        }) : true
      })
    )
}

export const notLoggedInGuard = (route: ActivatedRouteSnapshot) => {
  const profileService = inject(ProfileService);
  const router = inject(Router);

  return profileService.ready
    .pipe(
      switchMap(() => profileService.profile),
      map(value => {
        // if there requested route is /login and a force=yyy query parameter
        // is set than let the user open the login page instead of redirecting
        // to /profile.
        if (route.routeConfig?.path === 'login') {
          if (route.queryParamMap.get("force")) {
            return true;
          }
        }

        return value !== null ? router.navigate(['/profile']) : true
      })
    )
}
