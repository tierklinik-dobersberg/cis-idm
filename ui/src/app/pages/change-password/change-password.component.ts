import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { Profile } from '@tierklinik-dobersberg/apis';
import { Observable } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
  ],
  templateUrl: './change-password.component.html',
  styleUrls: ['./change-password.component.css']
})
export class ChangePasswordComponent {
  private readonly selfService = inject(SELF_SERVICE);
  private readonly router = inject(Router);

  readonly profile: Observable<Profile | null> = inject(ProfileService).profile;

  changePasswordError: string | null = null;

  current = '';
  new = '';
  newRepeat = '';

  async changePasswords() {
    try {
      await this.selfService.changePassword({
        newPassword: this.new,
        oldPassword: this.current,
      })

      this.changePasswordError = '';
      this.router.navigate(["../"]);
    } catch(err) {
      const connectErr = ConnectError.from(err);
      this.changePasswordError = connectErr.rawMessage;
    }
  }
}
