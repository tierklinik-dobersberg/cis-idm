import { NgIf, NgTemplateOutlet } from "@angular/common";
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ElementRef, HostListener, Input, Renderer2, TemplateRef, inject } from "@angular/core";

@Component({
  selector: 'tkd-img',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    NgIf,
    NgTemplateOutlet,
  ],
  template: `
  <img *ngIf="showImage && !!src; else: template" class="w-full h-full" [attr.src]="src" (error)="loadFallbackOnError()">
  <ng-template #template>
    <div class="flex flex-row items-center justify-center h-full w-full">
      <ng-container *ngIf="fallback">
        <ng-container *ngTemplateOutlet="$any(fallback)"></ng-container>
      </ng-container>
      <ng-container *ngIf="!fallback">
        <ng-content></ng-content>
      </ng-container>
    </div>
  </ng-template>
  `
})
export class TkdImageComponent {
  private element = inject(ElementRef);
  private cdr = inject(ChangeDetectorRef);

  @Input()
  src?: string = '';

  @Input()
  fallback?: string | TemplateRef<any>;

  showImage = true;

  loadFallbackOnError() {
    if (typeof this.fallback === 'string') {
      this.src = this.fallback;
      this.showImage = true;
    } else {
      this.showImage = false;
    }
  }
}
