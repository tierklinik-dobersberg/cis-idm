import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
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
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ResetPasswordComponent implements OnInit {
  private readonly authService = inject(AUTH_SERVICE);
  private readonly route = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly cdr = inject(ChangeDetectorRef)
  private readonly destroyRef = inject(DestroyRef);

  readonly config = inject(ConfigService).config;

  newPassword = '';
  newPasswordRepeat = '';

  display: 'request-reset' | 'reset' | 'sent' = 'request-reset';

  username = '';
  token = '';
  resetError: string | null = null;

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
    this.resetError = null;
    this.cdr.markForCheck();

    try {
      if (this.display === 'request-reset') {
        await this.authService.requestPasswordReset({
          kind: {
            case: 'email',
            value: this.username
          }
        })

        this.display = 'sent'
        this.cdr.markForCheck();

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
      const cerr = ConnectError.from(err)
      this.resetError = cerr.rawMessage
      this.cdr.markForCheck()
    }
  }
}
