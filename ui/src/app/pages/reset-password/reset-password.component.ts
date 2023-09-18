import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { AUTH_SERVICE } from 'src/app/clients';
import { ConfigService } from 'src/app/config.service';

@Component({
  selector: 'app-reset-password',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ResetPasswordComponent implements OnInit {
  authService = inject(AUTH_SERVICE);
  config = inject(ConfigService).config;
  route = inject(ActivatedRoute)
  router = inject(Router)
  destroyRef = inject(DestroyRef);
  newPassword = '';
  newPasswordRepeat = '';

  display: 'request-reset' | 'reset' = 'request-reset';

  username = '';
  token = '';

  ngOnInit(): void {
    this.route
      .url
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(url => {
        if (!url.length) {
          return;
        }

        const last = url[url.length - 1];
        if (last.path === 'request-reset') {
          this.display = 'request-reset';
        } else {
          this.display = 'reset';
          this.token = this.route.snapshot.queryParamMap.get("token") || '';
        }
      })
  }

  async submit() {
    debugger;

    try {
      if (this.display === 'request-reset') {
        await this.authService.requestPasswordReset({
          kind: {
            case: 'email',
            value: this.username
          }
        })

        return
      }

      await this.authService.requestPasswordReset({
        kind: {
          case: 'passwordReset',
          value: {
            token: this.token,
            newPassword: this.newPassword
          }
        }
      })

      this.router.navigate(['/login']);

    } catch (err) {
      console.error(err)
    }
  }
}
