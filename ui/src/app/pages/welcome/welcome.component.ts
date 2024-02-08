import { AsyncPipe, NgClass, NgFor, NgForOf, NgIf } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  DestroyRef,
  OnInit,
  inject,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { filter } from 'rxjs';
import { DisplayNamePipe } from '@tierklinik-dobersberg/angular/pipes';
import { TkdButtonDirective } from 'src/app/components/button';
import { ProfileService } from 'src/services/profile.service';
import { TkdAvatarComponent } from 'src/app/components/avatar';
import { trigger, transition, style, animate } from '@angular/animations';
import { Router, RouterModule } from '@angular/router';
import { EMail, PhoneNumber } from '@tierklinik-dobersberg/apis';
import { SELF_SERVICE } from 'src/app/clients';

@Component({
  standalone: true,
  templateUrl: './welcome.component.html',
  imports: [
    NgIf,
    NgClass,
    NgFor,
    NgForOf,
    AsyncPipe,
    TkdButtonDirective,
    DisplayNamePipe,
    TkdAvatarComponent,
    RouterModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  animations: [
    trigger('avatar', [
      transition(':enter', [
        style({
          opacity: 0,
          transform: 'scale(0)',
        }),
        animate(
          '150ms ease-in-out',
          style({
            opacity: 1,
            transform: 'scale(1)',
          })
        ),
      ]),
    ]),
    trigger('slideIn', [
      transition(':enter', [
        style({
          transform: 'translateX(-100%)',
        }),
        animate(
          '150ms 150ms ease-in-out',
          style({
            transform: 'translateX(0)',
          })
        ),
      ]),
    ]),
  ],
})
export class WelcomePageComponent implements OnInit {
  readonly profile = inject(ProfileService).profile;
  private readonly selfService = inject(SELF_SERVICE);
  private readonly profileService = inject(ProfileService);
  private readonly destroyRef = inject(DestroyRef);
  private readonly cdr = inject(ChangeDetectorRef);

  avatarMissing = false;
  profileIncomplete = false;
  addressMissing = false;
  phoneMissing = false;
  emailMissing = false;

  primaryMailMissing = false;
  primaryPhoneMissing = false;

  get everythingDone() {
    return !this.avatarMissing &&
      !this.profileIncomplete &&
      !this.addressMissing &&
      !this.phoneMissing &&
      !this.emailMissing &&
      !this.primaryMailMissing &&
      !this.primaryPhoneMissing;
  }

  ngOnInit(): void {
    this.profile
      .pipe(
        filter((profile) => !!profile),
        takeUntilDestroyed(this.destroyRef)
      )
      .subscribe((profile) => {
        const user = profile?.user;
        if (!user) {
          // this shouldn't happen
          return;
        }

        this.avatarMissing = !user.avatar;
        this.profileIncomplete =
          !user.displayName ||
          !user.firstName ||
          !user.lastName ||
          !user.birthday;

        this.phoneMissing = !profile.phoneNumbers?.length;
        this.emailMissing = !profile.emailAddresses?.length;
        this.primaryMailMissing = !user.primaryMail;
        this.primaryPhoneMissing = !user.primaryPhoneNumber;
        this.addressMissing = !profile.addresses?.length;

        this.cdr.markForCheck();
      });
  }

  async markAsPrimaryMail(email: EMail) {
    await this.selfService.markEmailAsPrimary({
      id: email.id
    })

    await this.profileService.loadProfile();
  }

  async markAsPrimaryPhone(phone: PhoneNumber) {
    await this.selfService.markPhoneNumberAsPrimary({
      id: phone.id,
    })

    await this.profileService.loadProfile();
  }
}
