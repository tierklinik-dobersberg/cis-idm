import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, inject } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { SELF_SERVICE } from 'src/app/clients';
import { SecurityCodeComponent } from 'src/app/shared/security-code/security-code.component';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-verify-phone',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    SecurityCodeComponent,
    RouterModule,
  ],
  templateUrl: './verify-phone.component.html',
  styleUrls: ['./verify-phone.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class VerifyPhoneComponent implements OnInit {
  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  route = inject(ActivatedRoute);
  router = inject(Router);
  cdr = inject(ChangeDetectorRef);

  error: string | null = null;

  securityCode: string = '';

  async ngOnInit() {
    await this.requestVerificationCode();
  }

  async requestVerificationCode() {
    try {
      await this.selfService.validatePhoneNumber({
        step: {
          case: 'id',
          value: this.route.snapshot.paramMap.get("id")  || '',
        }
      })
      this.error = null;
    }
    catch(err) {
      this.error = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }

  async verifyPhoneNumber() {
    try {
      await this.selfService.validatePhoneNumber({
        step: {
          case: 'code',
          value: this.securityCode,
        }
      })

      this.error = null;
      await this.profileService.loadProfile();
      await this.router.navigate(['/profile']);
    }
    catch(err) {
      this.error = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }
}
