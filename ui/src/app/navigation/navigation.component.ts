import {
  ChangeDetectionStrategy,
  Component,
  ViewChild,
  inject,
} from '@angular/core';
import { Router } from '@angular/router';
import { ProfileService } from 'src/services/profile.service';
import { AUTH_SERVICE } from '../clients';
import { ConfigService } from '../config.service';
import { TkdSideNavComponent } from '../components/navigation';
import { LayoutService } from '@tierklinik-dobersberg/angular/layout';

export type NavMode = 'side' | 'over';

@Component({
  selector: 'app-nav',
  exportAs: 'appNav',
  templateUrl: './navigation.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavigationComponent {
  readonly config = inject(ConfigService);
  readonly layout = inject(LayoutService).withAutoUpdate();

  private readonly authService = inject(AUTH_SERVICE);
  private readonly profileService = inject(ProfileService);
  private readonly router = inject(Router);

  @ViewChild(TkdSideNavComponent, { static: true })
  sideNav!: TkdSideNavComponent;

  async logout() {
    try {
      await this.authService.logout({});

      localStorage.removeItem('access_token');

      // trigger a "reloading" of the profile.
      await this.profileService.loadProfile();

      this.router.navigate(['/login'], {
        queryParams: {
          logout: '1',
        },
      });
    } catch (err) {
      console.error(err);
    }
  }
}
