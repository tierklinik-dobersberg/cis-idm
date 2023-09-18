import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, inject } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-verify-mail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
  ],
  templateUrl: './verify-mail.component.html',
  styleUrls: ['./verify-mail.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class VerifyMailComponent {
  selfService = inject(SELF_SERVICE);
  route = inject(ActivatedRoute)
  router = inject(Router);
  profileService = inject(ProfileService);
  profile$ = this.profileService.profile;
  error: string | null = null;
  cdr = inject(ChangeDetectorRef);

  async ngOnInit() {
    try {
    const token = this.route.snapshot.queryParamMap.get("token")
    if (!token) {
      this.router.navigate(['/profile'])
      return
    }

    await this.selfService.validateEmail({
      kind: {
        case: 'token',
        value: token,
      }
    })

    await this.profileService.loadProfile();
    } catch(err) {
      this.error = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }

}
