import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  inject
} from '@angular/core';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { PartialMessage } from '@bufbuild/protobuf';
import { UpdateProfileRequest } from '@tierklinik-dobersberg/apis';
import { take } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { TkdButtonDirective } from 'src/app/components/button';
import { TkdDatepickerComponent } from 'src/app/components/datepicker';
import { ConfigService } from 'src/app/config.service';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-edit-profile',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    TkdButtonDirective,
    TkdDatepickerComponent,
  ],
  templateUrl: './edit-profile.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EditProfileComponent implements OnInit {
  private readonly profileService = inject(ProfileService);
  private readonly selfService = inject(SELF_SERVICE);
  private readonly cdr = inject(ChangeDetectorRef);
  private readonly router = inject(Router);
  config = inject(ConfigService).config;

  firstName = new FormControl('');
  lastName = new FormControl('');
  displayName = new FormControl('');
  birthday = new FormControl('');
  username = new FormControl('');

  ngOnInit(): void {
    this.profileService.profile.pipe(take(1)).subscribe((profile) => {
      this.firstName.setValue(profile?.user?.firstName || '');
      this.lastName.setValue(profile?.user?.lastName || '');
      this.displayName.setValue(profile?.user?.displayName || '');
      this.birthday.setValue(profile?.user?.birthday || '');
      this.username.setValue(profile?.user?.username || '');

      this.cdr.markForCheck();
    });
  }

  async saveProfile() {
    let user: PartialMessage<UpdateProfileRequest> = {};
    let fieldSet: string[] = [];

    if (this.firstName.dirty) {
      user.firstName = this.firstName.value!;
      fieldSet.push('first_name');
    }

    if (this.lastName.dirty) {
      user.lastName = this.lastName.value!;
      fieldSet.push('last_name');
    }

    if (this.displayName.dirty) {
      user.displayName = this.displayName.value!;
      fieldSet.push('display_name');
    }

    if (this.birthday.dirty) {
      user.birthday = this.birthday.value!;
      fieldSet.push('birthday');
    }

    if (this.username.dirty && this.config.userNameChange) {
      user.username = this.username.value!;
      fieldSet.push('username');
    }

    if (fieldSet.length === 0) {
      this.router.navigate(['../']);
      return;
    }

    await this.selfService.updateProfile({
      ...user,
      fieldMask: {
        paths: fieldSet,
      },
    });

    await this.profileService.loadProfile();
    await this.router.navigate(['../']);
  }
}
