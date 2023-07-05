import { CommonModule } from "@angular/common";
import { Component, inject } from "@angular/core";
import { Router, RouterModule } from "@angular/router";
import { Profile } from '@tkd/apis/gen/es/tkd/idm/v1/user_pb';
import { Observable } from 'rxjs';
import { AUTH_SERVICE } from "src/app/clients";
import { ConfigService } from "src/app/config.service";
import { ProfileService } from 'src/services/profile.service';

@Component({
  standalone: true,
  templateUrl: './profile.component.html',
  styles: [
    `:host {
      @apply flex flex-col gap-8;
    }`
  ],
  imports: [
    CommonModule,
    RouterModule
  ]
})
export class ProfileComponent {
  authService = inject(AUTH_SERVICE);
  router = inject(Router)
  profileService = inject(ProfileService)
  config = inject(ConfigService).config;
  profile: Observable<Profile | null> = inject(ProfileService).profile;

  async logout() {
    try {
      await this.authService.logout({})

      localStorage.removeItem("access_token")

      // trigger a "reloading" of the profile.
      await this.profileService.loadProfile();

      this.router.navigate(["/login"], {
        queryParams: {
          'logout': '1'
        }
      })
    } catch(err) {
      console.error(err);
    }
  }
}
