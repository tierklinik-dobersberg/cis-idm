import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { Address } from '@tierklinik-dobersberg/apis';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-address',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  templateUrl: './add-edit-address.component.html',
  styleUrls: ['./add-edit-address.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AddEditAddressComponent implements OnInit {
  isNew = true;

  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  activeRoute = inject(ActivatedRoute)
  destroyRef = inject(DestroyRef);
  router = inject(Router);
  cdr = inject(ChangeDetectorRef);

  saveAddrError: string | null = null;

  id: string | null = null;
  cityCode = new FormControl('');
  cityName = new FormControl('');
  countryName = new FormControl('');
  street = new FormControl('');
  extra = new FormControl('');

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
            this.router.navigate(['../'])
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

  async save() {
    const addr: Partial<Address> = {
      cityCode: this.cityCode.value!,
      cityName: this.cityName.value!,
      street: this.street.value!,
      extra: this.extra.value!,
    };

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

    await this.selfService.deleteAddress({ id: this.id! })
    await this.profileService.loadProfile();
    this.router.navigate(['../'])
  }
}

