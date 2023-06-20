import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { SELF_SERVICE } from 'src/app/clients';

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
  selfService = inject(SELF_SERVICE);
  router = inject(Router);

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
