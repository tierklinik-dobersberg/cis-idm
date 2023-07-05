import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnInit, inject } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { AUTH_SERVICE } from 'src/app/clients';

@Component({
  selector: 'app-refresh',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
  ],
  templateUrl: './refresh.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RefreshComponent implements OnInit {
  authService = inject(AUTH_SERVICE);
  activeRoute = inject(ActivatedRoute);
  router = inject(Router);

  async ngOnInit() {
      const redirect = this.activeRoute.snapshot.queryParamMap.get("redirect")
      if (!redirect) {
        this.router.navigate(['/profile'])

        return
      }

      try {
        const res = await this.authService.refreshToken({
          requestedRedirect: redirect,
        })

        if (!!res.redirectTo) {
          window.location.href = res.redirectTo;

          return
        }
      } catch(err) {
        console.error(err);

        this.router.navigate(['/profile'])
      }
  }
}

