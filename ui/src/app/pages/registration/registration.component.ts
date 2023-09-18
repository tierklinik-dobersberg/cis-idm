import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, inject } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { browserSupportsWebAuthn, startRegistration } from '@simplewebauthn/browser';
import { firstValueFrom } from 'rxjs';
import { AUTH_SERVICE } from 'src/app/clients';
import { ConfigService } from 'src/app/config.service';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-registration',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule
  ],
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class RegistrationComponent {
  client = inject(AUTH_SERVICE);
  profile = inject(ProfileService);
  router = inject(Router);
  config = inject(ConfigService).config;
  currentRoute = inject(ActivatedRoute);
  cdr = inject(ChangeDetectorRef);
  http = inject(HttpClient);

  webauthnSupported = browserSupportsWebAuthn();
  usePassword = !this.webauthnSupported;
  username = '';
  password = '';
  passwordRepeat = '';
  email = '';
  token = '';
  errorMessage = '';

  ngOnInit() {
    const params = this.currentRoute.snapshot.queryParamMap;
    this.token = params.get("token") || "";
    this.email = params.get("mail") || "";
    this.username = params.get("name") || "";
  }

  async submit() {
    if (!this.usePassword) {
      try {
        const payload: any = await firstValueFrom(this.http.post(`/webauthn/registration/begin`, {
          username: this.username,
          token: this.token,
        }, { withCredentials: true}))

        const response = await startRegistration(payload.publicKey)

        await firstValueFrom(this.http.post(`/webauthn/registration/finish`, response, {
          withCredentials: true
        }))

      } catch(err) {
        this.errorMessage = JSON.stringify(err);
      }

      await this.router.navigate(['/login'])

      return
    }

    try {
      const result = await this.client.registerUser({
        registrationToken: this.token,
        password: this.password,
        username: this.username,
        email: this.email,
      })

      if (!result.accessToken) {
        throw new Error("unexpected response")
      }


      await this.profile.loadProfile();
      await this.router.navigate(['/profile'])

    } catch(err) {
      const connectErr = ConnectError.from(err);

      this.errorMessage = connectErr.rawMessage;
      this.cdr.markForCheck();
    }
  }
}
