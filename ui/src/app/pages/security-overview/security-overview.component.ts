import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, TrackByFunction, inject } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { startRegistration } from '@simplewebauthn/browser';
import { EnrollTOTPResponseStep1, RegisteredPasskey } from '@tierklinik-dobersberg/apis';
import { firstValueFrom, from, switchMap } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { TkdButtonDirective } from 'src/app/components/button';
import { SecurityCodeComponent } from 'src/app/shared/security-code/security-code.component';
import { ProfileService } from 'src/services/profile.service';

@Component({
  standalone: true,
  imports: [
    CommonModule,
    SecurityCodeComponent,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    TkdButtonDirective,
    RouterLink
  ],
  templateUrl: './security-overview.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SecurityOverviewComponent implements OnInit {
  cdr = inject(ChangeDetectorRef);
  profile$ = inject(ProfileService).profile;
  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  router = inject(Router);
  httpClient = inject(HttpClient);
  hasPublicKeyCreds = !!window.PublicKeyCredential;

  errMsg: string | null = null;

  enrollmentStep: EnrollTOTPResponseStep1 | null = null;
  code: string | null = null;
  recoveryCodes: string[] = [];
  disableTotpMode = false;
  passkeys: RegisteredPasskey[] = [];

  trackPassKey: TrackByFunction<RegisteredPasskey> = (_, key) => key.id;

  async ngOnInit() {
    await this.loadDevices();
  }

  async loadDevices() {
    try{
      const response = await this.selfService.getRegisteredPasskeys({})
      this.passkeys = response.passkeys;
      this.cdr.markForCheck();
    } catch(err) {
      console.error(err);
    }
  }

  async disableTotp() {
    if (!this.disableTotpMode) {
      this.disableTotpMode = true;
      this.cdr.markForCheck();

      return
    }

    try {
      await this.selfService.remove2FA({
        kind: {
          case: 'totpCode',
          value: this.code || '',
        }
      })

      await this.profileService.loadProfile();

      this.disableTotpMode = false;
      this.errMsg = null;
      this.code = null;
      this.cdr.markForCheck();

    } catch (err) {
      const connectErr = ConnectError.from(err);
      this.errMsg = connectErr.rawMessage;

      this.cdr.markForCheck();
    }
  }

  async removePasskey(id: string) {
    await this.selfService.removePasskey({id})
    await this.loadDevices();
  }

  async enableTotp() {
    try {

      if (this.enrollmentStep === null) {
        const res = await this.selfService.enroll2FA({
          kind: {
            case: "totpStep1",
            value: {}
          }
        })

        if (res.kind.case !== 'totpStep1') {
          throw new Error("unexpected server response")
        }

        this.enrollmentStep = res.kind.value;
      } else {
        const res = await this.selfService.enroll2FA({
          kind: {
            case: "totpStep2",
            value: {
              secret: this.enrollmentStep?.secret,
              secretHmac: this.enrollmentStep?.secretHmac,
              verifyCode: this.code || '',
            }
          }
        })

        if (res.kind.case !== 'totpStep2') {
          throw new Error('unexpected server response')
        }

        this.enrollmentStep = null;
        this.code = '';
        await this.profileService.loadProfile()
      }

      this.errMsg = null;
    } catch (err) {
      const connectErr = ConnectError.from(err);
      this.errMsg = connectErr.rawMessage;
    }

    this.cdr.markForCheck();
  }

  async generateRecoveryCodes() {
    try {
      let res = await this.selfService.generateRecoveryCodes({})
      this.recoveryCodes = res.recoveryCodes || [];

      await this.profileService.loadProfile();

    } catch (err) {
      const connectErr = ConnectError.from(err);
      this.errMsg = connectErr.rawMessage;
    }

    this.cdr.markForCheck();
  }

  async registerWebAuthN() {
    const profile = await firstValueFrom(this.profile$);
    this.httpClient.post<any>('/webauthn/registration/begin', {})
      .pipe(
        switchMap((credentialCreationOptions: any) => {
          return from(startRegistration(credentialCreationOptions.publicKey!))
        }),
        switchMap(credential => {
          return this.httpClient.post(`/webauthn/registration/finish`, credential)
        }),
      )
      .subscribe(response => {
        this.loadDevices();
      })
  }
}
