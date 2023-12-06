import { CommonModule } from "@angular/common";
import { HttpClient, HttpClientModule } from "@angular/common/http";
import { ChangeDetectorRef, Component, DestroyRef, OnInit, TrackByFunction, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { ActivatedRoute, Router, RouterModule } from "@angular/router";
import { ConnectError } from "@bufbuild/connect";
import { browserSupportsWebAuthn, browserSupportsWebAuthnAutofill, startAuthentication } from '@simplewebauthn/browser';
import { AuthType, LoginResponse, RequiredMFAKind } from "@tkd/apis";
import { firstValueFrom, from, switchMap } from "rxjs";
import { AUTH_SERVICE } from "src/app/clients";
import { ConfigService } from "src/app/config.service";
import { SecurityCodeComponent } from "src/app/shared/security-code/security-code.component";
import { ProfileService } from "src/services/profile.service";

interface LoggedInUser {
  id: string;
  displayName: string;
  avatarUrl: string;
  username: string;
}

interface LoggedInUserHistory {
  users: LoggedInUser[];
}

type States = 'user-select' | 'username-input' | 'password-input' | 'totp-input' ;
const allStates: States[] = [
  'user-select',
  'username-input',
  'password-input',
  'totp-input'
]

function isValidState(s: any): s is States {
  return allStates.includes(s)
}

@Component({
  standalone: true,
  templateUrl: './login.component.html',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    CommonModule,
    RouterModule,
    SecurityCodeComponent,
    HttpClientModule,
  ]
})
export class LoginComponent implements OnInit {
  private readonly client = inject(AUTH_SERVICE);
  private readonly profile = inject(ProfileService);
  private readonly router = inject(Router);
  private readonly currentRoute = inject(ActivatedRoute)
  private readonly cdr = inject(ChangeDetectorRef)
  private readonly http = inject(HttpClient);
  private readonly destroyRef = inject(DestroyRef);

  readonly config = inject(ConfigService).config;

  set display(s: States) {
    this._display = s

    this.router.navigate(['.'], {
      queryParams: {
        s: this._display,
        redirect: this.currentRoute.snapshot.paramMap.get("redirect")
      }
    })
  }
  get display(): States {
    return this._display
  }
  private _display: States = 'username-input';

  username = '';
  password = '';
  code = '';
  loginErrorMessage = '';
  rememberMe = true;
  loggedInUsers: LoggedInUser[] = [];
  abortController = new AbortController();

  webauthnSupport = browserSupportsWebAuthn();
  autofillSupported = false;

  trackLoggedInUsers: TrackByFunction<LoggedInUser> = (_, user) => user.id;

  private state = '';
  selectedUser: LoggedInUser | null = null;

  async ngOnInit() {
    this.currentRoute.queryParamMap
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(params => {
        const s = params.get("s")
        if (isValidState(s)) {
          this._display = s

          this.cdr.markForCheck();
        }
      })

    this.autofillSupported = await browserSupportsWebAuthnAutofill();

    try {
      const state = localStorage.getItem("loggedInUsers");
      if (!!state) {
        const parsed: LoggedInUserHistory = JSON.parse(state);
        this.loggedInUsers = parsed.users || [];
      }

      if (this.loggedInUsers.length > 0) {
        this.display = 'user-select';
      }
    } catch (err) {
      this.display = 'username-input';
      console.error(err)
    }

    const justLoggedOut = this.currentRoute.snapshot.queryParamMap.has("logout");

    try {
      if (this.autofillSupported && !justLoggedOut) {
        const response: any = await firstValueFrom(this.http.post(`/webauthn/login/begin/`, {}, {
          withCredentials: true
        }));
        const loginValue = await startAuthentication(response.publicKey!, false)

        const loginResponse: any = await firstValueFrom(this.http.post(`/webauthn/login/finish`, loginValue, {
          withCredentials: true,
          params: {
            redirect: this.currentRoute.snapshot.queryParamMap.get("redirect") || ''
          }
        }))

        if (loginResponse.redirectTo) {
          window.location.href = loginResponse.redirectTo;

          return
        }

        await this.profile.loadProfile();

        const user = await firstValueFrom(this.profile.profile)
        if (this.rememberMe) {
          this.addLoggedInUser({
            id: user!.user!.id,
            displayName: user!.user!.displayName || user!.user!.username,
            username: user!.user!.username,
            avatarUrl: `/avatar/${user!.user!.id}`
          })
        }

        await this.router.navigate(['/profile'])
      } else if (this.autofillSupported) {
        this.autofillSupported = false;
      }

    } catch (err) {
      console.error(err);
    }
  }

  chooseDifferentAccount() {
    this.selectedUser = null;
    this.username = '';
    if (this.display === 'user-select') {
      this.display = 'username-input';
    } else {
      this.display = this.loggedInUsers.length > 0 ? 'user-select' : 'username-input';
    }
  }

  selectKnownAccount(user: LoggedInUser) {
    this.username = user.username;
    this.selectedUser = user;
    this.display = 'password-input';
  }

  removeAccount(user: LoggedInUser) {
    const newUsers = this.loggedInUsers.filter(u => u.id !== user.id);
    localStorage.setItem("loggedInUsers", JSON.stringify({
      users: newUsers,
    }))

    this.loggedInUsers = newUsers;
    if (this.loggedInUsers.length === 0) {
      this.display = 'username-input';
    }
  }

  private addLoggedInUser(usr: LoggedInUser) {
    const usersString: string = localStorage.getItem("loggedInUsers") || '{"users": []}';
    const users: LoggedInUserHistory = JSON.parse(usersString);

    if (users.users.some(u => u.id === usr.id)) {
      return;
    }

    users.users.push(usr)
    localStorage.setItem("loggedInUsers", JSON.stringify(users));
  }

  loginUsingWebAuthn() {
    this.http.post('/webauthn/login/begin/' + this.username, {}, {
      withCredentials: true,
    })
      .pipe(
        switchMap((data: any) => {
          return from(startAuthentication(data.publicKey!, this.autofillSupported));
        }),
        switchMap(creds => {
          return this.http.post<any>('/webauthn/login/finish', creds, {
            withCredentials: true,
            params: {
              redirect: this.currentRoute.snapshot.queryParamMap.get("redirect") || ''
            }
          })
        }),
      )
      .subscribe({
        next: async result => {
          if (!!result.redirectTo) {
            window.location.href = result.redirectTo;

            return
          }

          await this.profile.loadProfile();

          const user = await firstValueFrom(this.profile.profile)
          if (this.rememberMe) {
            this.addLoggedInUser({
              id: user!.user!.id,
              displayName: user!.user!.displayName || user!.user!.username,
              username: user!.user!.username,
              avatarUrl: `/avatar/${user!.user!.id}`
            })
          }

          await this.router.navigate(['/profile']);
        },
        error: err => {
          console.error(err);
          this.loginErrorMessage = JSON.stringify(err);
        }
      })
  }

  async submit() {
    if (this.display === 'username-input') {
      this.selectedUser = this.loggedInUsers.find(user => user.username === this.username) || null;
      this.display = 'password-input';
      return
    }

    try {
      let result: LoginResponse;
      if (this.display === 'password-input') {
        result = await this.client.login({
          authType: AuthType.PASSWORD,
          auth: {
            case: 'password',
            value: {
              password: this.password,
              username: this.username,
            }
          },
          requestedRedirect: this.currentRoute.snapshot.queryParamMap.get("redirect") || '',
        })
      } else if (this.display === 'totp-input') {
        result = await this.client.login({
          authType: AuthType.TOTP,
          auth: {
            case: 'totp',
            value: {
              code: this.code,
              state: this.state,
            }
          },
          requestedRedirect: this.currentRoute.snapshot.queryParamMap.get("redirect") || '',
        })
      } else {
        return;
      }

      if (result.response.case === "accessToken") {
        await this.profile.loadProfile();

        if (this.rememberMe) {
          const profile = await firstValueFrom(this.profile.profile)
          this.addLoggedInUser({
            id: profile!.user!.id,
            username: profile!.user!.username,
            displayName: profile?.user?.displayName || profile?.user?.username || '',
            avatarUrl: `/avatar/${profile!.user!.id}`
          })

        }

        if (!!result.redirectTo) {
          window.location.href = result.redirectTo;

          return
        }

        await this.router.navigate(['/profile'])
      } else {
        if (result.response.case === 'mfaRequired') {
          switch (result.response.value.kind) {
            case RequiredMFAKind.REQUIRED_MFA_KIND_TOTP:
              this.state = result.response.value.state;
              this.code = '';
              this.display = 'totp-input'
              break;

            default:
              throw new Error('unsupported mfa kind')
          }
        }
      }
    } catch (err) {
      const connectErr = ConnectError.from(err);

      this.loginErrorMessage = connectErr.rawMessage;
    }
  }
}
