import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  DestroyRef,
  Input,
  TemplateRef,
  ViewChild,
  ViewContainerRef,
  booleanAttribute,
  inject,
} from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import {
  trigger,
  animate,
  style,
  transition,
  AnimationEvent,
} from '@angular/animations';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { TemplatePortal } from '@angular/cdk/portal';
import { takeUntil } from 'rxjs';
import { NgIf, NgTemplateOutlet } from '@angular/common';
import { coerceCssPixelValue } from '@angular/cdk/coercion';

export type NavMode = 'side' | 'over';

@Component({
  selector: 'tkd-nav',
  exportAs: 'tkdNav',
  standalone: true,
  templateUrl: './navigation.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [NgIf, NgTemplateOutlet],
  animations: [
    trigger('animate', [
      transition(':enter', [
        style({
          opacity: 0,
          transform: 'translateX(-100%)',
        }),
        animate(
          '0.15s ease-in-out',
          style({
            opacity: 1,
            transform: 'translateX(0)',
          })
        ),
      ]),
      transition(':leave', [
        style({
          opacity: 1,
          transform: 'translateX(0)',
        }),
        animate(
          '0.15s ease-in-out',
          style({
            opacity: 0,
            transform: 'translateX(-100%)',
          })
        ),
      ]),
    ]),
  ],
})
export class TkdSideNavComponent {
  readonly router = inject(Router);
  readonly overlay = inject(Overlay);
  readonly viewContainerRef = inject(ViewContainerRef);
  readonly destroyRef = inject(DestroyRef);
  readonly cdr = inject(ChangeDetectorRef);

  @ViewChild('navigationTemplate', { read: TemplateRef, static: true })
  readonly navTemplate!: TemplateRef<never>;

  private overlayRef: OverlayRef | null = null;

  @Input()
  mode: NavMode = 'side';

  @Input({ transform: coerceCssPixelValue })
  offsetTop = '0px';

  @Input({ transform: coerceCssPixelValue })
  offsetBottom = '0px';

  @Input({ transform: coerceCssPixelValue })
  offsetSide = '0px';

  @Input()
  position: 'left' | 'right' = 'left';

  open() {
    if (this.mode === 'side') {
      return;
    }

    if (this.overlayRef) {
      return;
    }

    let position = this.overlay.position().global().bottom(this.offsetBottom);

    switch (this.position) {
      case 'left':
        position = position.left(this.offsetSide);
        break;
      case 'right':
        position = position.right(this.offsetSide);
        break;
    }

    this.overlayRef = this.overlay.create({
      hasBackdrop: true,
      disposeOnNavigation: true,
      height: `calc(100vh - ${this.offsetTop} - ${this.offsetBottom})`,
      panelClass: 'z-[1035]',
      positionStrategy: position,
    });

    this.router.events
      .pipe(takeUntil(this.overlayRef.detachments()))
      .subscribe((evt) => {
        if (evt instanceof NavigationEnd) {
          this.close();
        }
      });

    this.overlayRef.backdropClick().subscribe(() => this.close());

    this.overlayRef.attach(
      new TemplatePortal(this.navTemplate, this.viewContainerRef)
    );
  }

  animationDone(event: AnimationEvent) {
    if (event.toState === 'void') {
      this.overlayRef?.dispose();
      this.overlayRef = null;
    }
  }

  close() {
    if (this.overlayRef) {
      this.overlayRef.detach();
    }
  }

  toggle() {
    if (this.overlayRef) {
      this.close();

      return;
    }

    this.open();
  }
}
