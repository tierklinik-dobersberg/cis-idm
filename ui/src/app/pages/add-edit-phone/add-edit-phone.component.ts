import { CommonModule, Location } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { L10nTranslateAsyncPipe } from 'angular-l10n';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { TkdBacklinkDirective } from 'src/app/components/backlink';
import { TkdButtonDirective } from 'src/app/components/button';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-phone',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    TkdButtonDirective,
    TkdBacklinkDirective,
    L10nTranslateAsyncPipe
  ],
  templateUrl: './add-edit-phone.component.html',
  styleUrls: ['./add-edit-phone.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AddEditPhoneComponent {
  isNew = true;

  private readonly profileService = inject(ProfileService);
  private readonly selfService = inject(SELF_SERVICE);
  private readonly activeRoute = inject(ActivatedRoute)
  private readonly destroyRef = inject(DestroyRef);
  private readonly location = inject(Location);
  private readonly cdr = inject(ChangeDetectorRef);

  id: string | null = null;
  saveAddrError: string | null = null;

  form = new FormGroup({
    number: new FormControl('', {
      validators: [
        Validators.required,
        Validators.pattern(/^[0-9 \+-]+$/)
      ]
    }),
    primary: new FormControl(false),
    verified: new FormControl(false),
  })

  get number() { return this.form.get('number')! }
  get primary() { return this.form.get('primary')! }
  get verified() { return this.form.get('verified')! }

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

        this.form.reset();

        if (!!this.id) {
          const phoneNumber = (profile?.phoneNumbers || []).find(addr => addr.id === this.id);

          if (!phoneNumber) {
            this.location.back()
          } else {
            this.number.setValue(phoneNumber.number);
            this.primary.setValue(phoneNumber.primary);
            this.verified.setValue(phoneNumber.verified)

            this.cdr.markForCheck();
          }
        }
      })
  }

  async save() {
    try {
      await this.selfService.addPhoneNumber({ number: this.number.value! })
      await this.profileService.loadProfile();
      this.location.back()

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
    this.location.back()
  }

  async markAsPrimary() {
    if (this.isNew) {
      return
    }

    try {
      await this.selfService.markPhoneNumberAsPrimary({ id: this.id! })
      await this.profileService.loadProfile();
      this.location.back()
    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }
}
