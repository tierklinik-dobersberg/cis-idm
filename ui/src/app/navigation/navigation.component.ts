import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  OnInit,
  ViewChild,
  inject,
} from '@angular/core';
import { ConfigService } from '../config.service';
import { AUTH_SERVICE } from '../clients';
import { ProfileService } from 'src/services/profile.service';
import { Router } from '@angular/router';
import { Sidenav, initTE } from 'tw-elements';

@Component({
  selector: 'app-nav',
  templateUrl: './navigation.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavigationComponent implements OnInit {
  readonly config = inject(ConfigService);
  readonly authService = inject(AUTH_SERVICE);
  readonly profileService = inject(ProfileService);
  readonly router = inject(Router);

  @ViewChild('sidenav', { read: ElementRef, static: true })
  sidenavElement!: ElementRef<HTMLElement>;

  ngOnInit() {
    initTE({ Sidenav });
  }

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
