import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ProfileService } from 'src/services/profile.service';
import { SELF_SERVICE } from 'src/app/clients';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { combineLatest } from 'rxjs';

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
        const addressID = params.get("id");
        if (!!addressID) {
          const address = (profile?.addresses || []).find(addr => addr.id === addressID);

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
}

