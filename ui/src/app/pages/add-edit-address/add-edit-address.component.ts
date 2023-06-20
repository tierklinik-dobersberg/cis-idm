import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ProfileService } from 'src/services/profile.service';
import { SELF_SERVICE } from 'src/app/clients';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { combineLatest } from 'rxjs';
import { AddAddressRequest } from '@tkd/apis/gen/es/tkd/idm/v1/self_service_pb.js';
import { Address } from '@tkd/apis/gen/es/tkd/idm/v1/user_pb';

@Component({
  selector: 'app-add-edit-address',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
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
        cityCode:  this.cityCode.value!,
        cityName: this.cityName.value!,
        street: this.street.value!,
        extra: this.extra.value!,
      };

    if (this.isNew) {
      await this.selfService.addAddress({...addr})
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
        fieldMask: {paths},
        ...addr
      })
    }

    await this.profileService.loadProfile();
    this.router.navigate(['../'])
  }
}

