import { CommonModule, Location } from '@angular/common';
import { ChangeDetectionStrategy, Component, OnInit, inject } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { Profile } from '@tierklinik-dobersberg/apis';
import { L10nTranslateAsyncPipe } from 'angular-l10n';
import { Observable, repeat, take } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    L10nTranslateAsyncPipe,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './change-password.component.html',
  styleUrls: ['./change-password.component.css']
})
export class ChangePasswordComponent implements OnInit {
  private readonly selfService = inject(SELF_SERVICE);
  private readonly location = inject(Location);

  readonly profile: Observable<Profile | null> = inject(ProfileService).profile;

  changePasswordError: string | null = null;

  form = new FormGroup({
    current: new FormControl(''),
    new: new FormControl('', {
      validators: [
        Validators.required,
      ]
    }),
    newRepeat: new FormControl('', {
      validators: [
        Validators.required,
      ]
    }),
  }, {
    validators: (control: AbstractControl) => {
      const newValue = control.get('new')!;
      const repeatValue = control.get('newRepeat')!;

      if (newValue.dirty && repeatValue.dirty && newValue.value !== repeatValue.value) {
        newValue.setErrors({
          'password-mismatch': 'password-mismatch'
        })
        repeatValue.setErrors({
          'password-mismatch': 'password-mismatch'
        })
      } else  {
        newValue.setErrors(null)
        repeatValue.setErrors(null)
      }

      return null
    }
  })

  get current() { return this.form.get('current')! }
  get new() { return this.form.get('new')! }
  get newRepeat() { return this.form.get('newRepeat')! }

  ngOnInit(): void {
    this.profile
      .pipe(take(1))
      .subscribe(profile => {
        if (!profile) {
          // cannot happen due to route guard
          return;
        }

        if (profile.passwordAuthEnabled) {
          this.current.setValidators([
            Validators.required
          ])
        } else {
          this.current.setValidators([])
        }
      })
  }

  async changePasswords() {
    try {
      await this.selfService.changePassword({
        newPassword: this.new.value!,
        oldPassword: this.current.value!,
      })

      this.changePasswordError = '';
      this.location.back()
    } catch(err) {
      const connectErr = ConnectError.from(err);
      this.changePasswordError = connectErr.rawMessage;
    }
  }
}
