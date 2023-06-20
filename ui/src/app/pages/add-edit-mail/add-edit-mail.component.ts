import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import { combineLatest } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

@Component({
  selector: 'app-add-edit-mail',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    ReactiveFormsModule,
  ],
  templateUrl: './add-edit-mail.component.html',
  styleUrls: ['./add-edit-mail.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AddEditMailComponent implements OnInit {
  isNew = true;

  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  activeRoute = inject(ActivatedRoute)
  destroyRef = inject(DestroyRef);
  router = inject(Router);
  cdr = inject(ChangeDetectorRef);

  id: string | null = null;
  saveAddrError: string | null = null;

  address = new FormControl('');
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
          const address = (profile?.emailAddresses || []).find(addr => addr.id === this.id);

          if (!address) {
            this.router.navigate(['../'])
          } else {
            this.address.setValue(address.address);
            this.primary.setValue(address.primary);
            this.verified = address.verified;

            this.cdr.markForCheck();
          }
        }
      })
  }

  async save() {
    try {
      await this.selfService.addEmailAddress({ email: this.address.value! })
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

    await this.selfService.deleteEmailAddress({ id: this.id! })
    await this.profileService.loadProfile();
    this.router.navigate(['../'])
  }

  async markAsPrimary() {
    if (this.isNew) {
      return
    }

    try {
      await this.selfService.markEmailAsPrimary({ id: this.id! })
      await this.profileService.loadProfile();
      this.router.navigate(['../'])
    } catch (err) {
      this.saveAddrError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }
}
