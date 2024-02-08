import { CommonModule } from '@angular/common';
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  DestroyRef,
  ElementRef,
  OnInit,
  ViewChild,
  inject,
} from '@angular/core';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { PartialMessage } from '@bufbuild/protobuf';
import { UpdateProfileRequest } from '@tierklinik-dobersberg/apis';
import { take } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ConfigService } from 'src/app/config.service';
import { ProfileService } from 'src/services/profile.service';
import { Datepicker } from 'tw-elements';
import { LayoutService } from '@tierklinik-dobersberg/angular/layout';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { TkdButtonDirective } from 'src/app/components/button';

@Component({
  selector: 'app-edit-profile',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, RouterModule, TkdButtonDirective],
  templateUrl: './edit-profile.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EditProfileComponent implements OnInit, AfterViewInit {
  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  cdr = inject(ChangeDetectorRef);
  router = inject(Router);
  config = inject(ConfigService).config;
  layout = inject(LayoutService);
  destroyRef = inject(DestroyRef);

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

  @ViewChild('datepicker', { read: ElementRef, static: true })
  datepickerEl!: ElementRef<HTMLInputElement>;

  ngAfterViewInit(): void {
    let dt = Datepicker.getOrCreateInstance(this.datepickerEl.nativeElement);

    this.layout.change
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        dt.close();
        dt.update({
          inline: this.layout.md,
          disableFuture: true,
          confirmDateOnSelect: true,
          startDay: 1,
          monthsFull: [
            'Jänner',
            'Feburar',
            'März',
            'April',
            'Mai',
            'Juni',
            'Juli',
            'August',
            'September',
            'Oktober',
            'November',
            'Dezember',
          ],
          monthsShort: [
            'Jän',
            'Feb',
            'Mär',
            'Apr',
            'Mai',
            'Jun',
            'Jul',
            'Aug',
            'Sep',
            'Okt',
            'Nov',
            'Dez',
          ],
          weekdaysFull: [
            'Sonntag',
            'Montag',
            'Dienstag',
            'Mittwoch',
            'Donnerstag',
            'Freitag',
            'Samstag',
          ],
          weekdaysNarrow: [
            'S',
            'M',
            'D',
            'M',
            'D',
            'F',
            'S',
          ],
          weekdaysShort: [
            'Son',
            'Mon',
            'Die',
            'Mit',
            'Don',
            'Fre',
            'Sam',
          ],
          title: 'Geburtstag auswählen'
        });
      });

    this.destroyRef.onDestroy(() => dt.dispose());
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

    if (this.username.dirty && this.config.features.allowUsernameChange) {
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
