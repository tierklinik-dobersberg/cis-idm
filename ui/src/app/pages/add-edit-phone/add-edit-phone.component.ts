import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-phone',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  templateUrl: './add-edit-phone.component.html',
  styleUrls: ['./add-edit-phone.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AddEditPhoneComponent {
  isNew = true;

  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  activeRoute = inject(ActivatedRoute)
  destroyRef = inject(DestroyRef);
  router = inject(Router);
  cdr = inject(ChangeDetectorRef);

  id: string | null = null;
  saveAddrError: string | null = null;

  number = new FormControl('');
  primary = new FormControl(false);
  verified = false;

  ngOnInit(): void {
    combineLatest([
      this.activeRoute.paramMap,
      this.profileService.profile,
    ])
      .pipe(
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe(([params, profile]) => {
        this.id = params.get("id");
        this.isNew = !this.id;
        this.saveAddrError = null;
        this.verified = false;

        if (!!this.id) {
          const phoneNumber = (profile?.phoneNumbers || []).find(addr => addr.id === this.id);

          if (!phoneNumber) {
            this.router.navigate(['../'])
          } else {
            this.number.setValue(phoneNumber.number);
            this.primary.setValue(phoneNumber.primary);
            this.verified = phoneNumber.verified;

            this.cdr.markForCheck();
          }
        }
      })
  }

  async save() {
    try {
      await this.selfService.addPhoneNumber({ number: this.number.value! })
      await this.profileService.loadProfile();
      this.router.navigate(['../'])

    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }

  async deleteAddress() {
    if (this.isNew) {
      return
    }

    await this.selfService.deletePhoneNumber({ id: this.id! })
    await this.profileService.loadProfile();
    this.router.navigate(['../'])
  }

  async markAsPrimary() {
    if (this.isNew) {
      return
    }

    try {
      await this.selfService.markPhoneNumberAsPrimary({ id: this.id! })
      await this.profileService.loadProfile();
      this.router.navigate(['../'])
    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }
}
