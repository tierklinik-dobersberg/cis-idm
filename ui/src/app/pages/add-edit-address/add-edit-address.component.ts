import { CommonModule, Location } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, RequiredValidator, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { Address } from '@tierklinik-dobersberg/apis';
import { L10N_LOCALE, L10nTranslateAsyncPipe, L10nTranslatePipe } from 'angular-l10n';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { TkdBacklinkDirective } from 'src/app/components/backlink';
import { TkdButtonDirective } from 'src/app/components/button';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-address',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TkdButtonDirective,
    TkdBacklinkDirective,
    L10nTranslateAsyncPipe
  ],
  templateUrl: './add-edit-address.component.html',
  styleUrls: ['./add-edit-address.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AddEditAddressComponent implements OnInit {
  isNew = true;

  private readonly profileService = inject(ProfileService);
  private readonly selfService = inject(SELF_SERVICE);
  private readonly activeRoute = inject(ActivatedRoute)
  private readonly destroyRef = inject(DestroyRef);
  private readonly location = inject(Location);
  private readonly cdr = inject(ChangeDetectorRef);

  saveAddrError: string | null = null;

  id: string | null = null;

  form = new FormGroup({
    cityCode: new FormControl('', {
      validators: Validators.required,
    }),

    cityName: new FormControl('', {
      validators: Validators.required
    }),

    street: new FormControl('', {
      validators: Validators.required
    }),

    extra: new FormControl('')
  })


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
          const address = (profile?.addresses || []).find(addr => addr.id === this.id);

          if (!address) {
            this.location.back()
          } else {
            this.cityCode.setValue(address.cityCode);
            this.cityName.setValue(address.cityName);

            this.street.setValue(address.street);
            this.extra.setValue(address.extra);

            this.cdr.markForCheck();
          }
        }
      })
  }

  get cityCode() { return this.form.get('cityCode')! }
  get cityName() { return this.form.get('cityName')! }
  get street() { return this.form.get('street')! }
  get extra() { return this.form.get('extra')! }

  async save() {
    const addr = this.form.value as Partial<Address>;

    try {

      if (this.isNew) {
        await this.selfService.addAddress({ ...addr })
      } else {
        const paths: string[] = [];
        if (this.cityCode.dirty) {
          paths.push("city_code");
        }

        if (this.cityName.dirty) {
          paths.push("city_name")
        }

        if (this.street.dirty) {
          paths.push("street")
        }

        if (this.extra.dirty) {
          paths.push("extra")
        }

        await this.selfService.updateAddress({
          id: this.id!,
          fieldMask: { paths },
          ...addr
        })
      }

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

    await this.selfService.deleteAddress({ id: this.id! })
    await this.profileService.loadProfile();
    this.location.back()
  }
}

