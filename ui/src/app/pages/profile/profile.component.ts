import { CommonModule } from "@angular/common";
import { Component, inject } from "@angular/core";
import { Router } from "@angular/router";
import { Profile } from '@tkd/apis/gen/es/tkd/idm/v1/user_pb';
import { Observable } from 'rxjs';
import { AUTH_SERVICE } from "src/app/clients";
import { ProfileService } from 'src/services/profile.service';

@Component({
  standalone: true,
  templateUrl: './profile.component.html',
  imports: [
    CommonModule
  ]
})
export class ProfileComponent {
  authService = inject(AUTH_SERVICE);
  router = inject(Router)
  profileService = inject(ProfileService)
  profile: Observable<Profile | null> = inject(ProfileService).profile;

  async logout() {
    try {
      await this.authService.logout({})

      localStorage.removeItem("access_token")

      // trigger a "reloading" of the profile.
      await this.profileService.loadProfile();

      this.router.navigate(["/login"])
    } catch(err) {
      console.error(err);
    }
  }
}
