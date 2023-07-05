import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit, inject } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { ConnectError } from '@bufbuild/connect';
import Cropper from 'cropperjs';
import { filter, take } from 'rxjs';
import { SELF_SERVICE } from 'src/app/clients';
import { ProfileService } from 'src/services/profile.service';

// Unfortunately there are no typings for blueimp-load-image
declare var loadImage: any;

export interface ImageCropperResult {
  imageData: Cropper.ImageData;
  cropData: Cropper.CropBoxData;
  blob?: Blob;
  dataUrl?: string;
}

@Component({
  selector: 'app-edit-avatar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
  ],
  templateUrl: './edit-avatar.component.html',
  styleUrls: ['./edit-avatar.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class EditAvatarComponent implements OnInit, OnDestroy {
  imageUrl: string = '';
  profileService = inject(ProfileService);
  selfService = inject(SELF_SERVICE);
  cdr = inject(ChangeDetectorRef);
  router = inject(Router);

  @Input() cropbox?: Cropper.CropBoxData;
  @Input() loadImageErrorText?: string;

  public isLoading: boolean = true;
  public cropper: Cropper | null = null;
  public loadError: any;

  constructor() { }

  ngOnInit() {
    this.profileService
      .profile
      .pipe(
        filter(profile => !!profile),
        take(1)
      )
      .subscribe(profile => this.imageUrl = profile!.user!.avatar || '');
  }

  ngOnDestroy() {
    if (this.cropper) {
      this.cropper.destroy();
      this.cropper = null;
    }
  }

  /**
   * Image loaded
   * @param ev
   */
  imageLoaded(ev: Event) {
    //
    // Unset load error state
    this.loadError = false;

    //
    // Setup image element
    const image = ev.target as HTMLImageElement;

    //
    // Image on ready event
    image.addEventListener("ready", () => {
      //
      // Unset loading state
      this.isLoading = false;

      //
      // Validate cropbox existance
      if (this.cropbox) {
        //
        // Set cropbox data
        this.cropper!.setCropBoxData(this.cropbox);
      }
    });

    //
    // Set cropperjs
    if (this.cropper) {
      this.cropper.destroy();
      this.cropper = null;
    }

    this.cropper = new Cropper(image, {
      aspectRatio: 1,
      viewMode: 1,
    });
  }

  /**
   * Image load error
   * @param event
   */
  imageLoadError(event: any) {
    //
    // Set load error state
    this.loadError = true;

    //
    // Unset loading state
    this.isLoading = false;
  }

  /**
   * Export canvas
   * @param base64
   */
  async save() {
    const canvas = this.cropper!.getCroppedCanvas();
    const url = canvas.toDataURL("image/png")

    try {
      await this.selfService.updateProfile({
        avatar: url,
        fieldMask: {
          paths: ["avatar"]
        }
      })

      await this.profileService.loadProfile();

      this.router.navigate(['/profile']);

    } catch(err) {
      this.loadError = ConnectError.from(err).rawMessage;
      this.cdr.markForCheck();
    }
  }

  loadImage(event: any) {
    let file: File = (event.target as any).files[0];
    console.log("loading new image")

    loadImage(
      file,
      (img: HTMLCanvasElement) => {
        this.imageUrl = img.toDataURL();
        this.cdr.markForCheck();
      },
      {
        maxWidth: 800,
        canvas: true
      }
    );
  }
}
