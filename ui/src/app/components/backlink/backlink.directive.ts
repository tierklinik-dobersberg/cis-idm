import { Directive } from '@angular/core';
import { Location } from "@angular/common";
import { HostListener, inject } from "@angular/core";

@Directive({
  selector: '[tkd-backlink]',
  standalone: true,
})
export class TkdBacklinkDirective {
  private readonly location = inject(Location);

  @HostListener('click')
  goBack() {
    this.location.back();
  }
}
