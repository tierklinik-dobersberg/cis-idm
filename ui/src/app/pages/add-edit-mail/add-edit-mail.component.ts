import { CommonModule, Location } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { L10N_LOCALE, L10nTranslateAsyncPipe, L10nTranslatePipe } from 'angular-l10n';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { TkdBacklinkDirective } from 'src/app/components/backlink';
import { TkdButtonDirective } from 'src/app/components/button';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-mail',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TkdButtonDirective,
    TkdBacklinkDirective,
    L10nTranslateAsyncPipe
  ],
  templateUrl: './add-edit-mail.component.html',
  styleUrls: ['./add-edit-mail.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AddEditMailComponent implements OnInit {
  isNew = true;

  private readonly profileService = inject(ProfileService);
  private readonly selfService = inject(SELF_SERVICE);
  private readonly activeRoute = inject(ActivatedRoute)
  private readonly destroyRef = inject(DestroyRef);
  private readonly location = inject(Location);
  private readonly cdr = inject(ChangeDetectorRef);

  validationSent = false;

  id: string | null = null;
  saveAddrError: string | null = null;

  form = new FormGroup({
    address: new FormControl('', {
      validators: [
        Validators.required,
        Validators.email,
      ]
    }),
    primary: new FormControl(false),
    verified: new FormControl(false)
  })

  get address() { return this.form.get('address')! }
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

        if (!!this.id) {
          const address = (profile?.emailAddresses || []).find(addr => addr.id === this.id);

          if (!address) {
            this.location.back()
          } else {
            this.address.setValue(address.address);
            this.primary.setValue(address.primary);
            this.verified.setValue(address.verified);

            this.cdr.markForCheck();
          }
        }
      })
  }

  async save() {
    try {
      await this.selfService.addEmailAddress({ email: this.address.value! })
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

    await this.selfService.deleteEmailAddress({ id: this.id! })
    await this.profileService.loadProfile();
    this.location.back()
  }

  async validateEmail() {
    try {
      await this.selfService.validateEmail({
        kind: {
          case: 'emailId',
          value: this.id!,
        }
      })

      this.validationSent = true
    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }

  async markAsPrimary() {
    if (this.isNew) {
      return
    }

    try {
      await this.selfService.markEmailAsPrimary({ id: this.id! })
      await this.profileService.loadProfile();
      this.location.back()
    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }
}
