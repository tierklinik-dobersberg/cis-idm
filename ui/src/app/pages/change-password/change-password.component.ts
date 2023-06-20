import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
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

  current = '';
  new = '';
  newRepeat = '';

  async changePasswords() {
    await this.selfService.changePassword({
      newPassword: this.new,
      oldPassword: this.current,
    })

    this.router.navigate(["../"]);
  }
}
