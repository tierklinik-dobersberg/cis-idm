import { inject } from "@angular/core";
import { Router } from "@angular/router";
import { map, switchMap } from "rxjs";
import { ProfileService } from "src/services/profile.service";

export const authGuard = () => {
  const profileService = inject(ProfileService);
  const router = inject(Router);

  return profileService.ready
    .pipe(
      switchMap(() => profileService.profile),
      map(value => {
        return value === null ? router.navigate(['/login']) : true
      })
    )
}

export const notLoggedInGuard = () => {
  const profileService = inject(ProfileService);
  const router = inject(Router);

  return profileService.ready
    .pipe(
      switchMap(() => profileService.profile),
      map(value => {
        return value !== null ? router.navigate(['/profile']) : true
      })
    )
}
