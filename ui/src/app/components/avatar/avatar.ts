import { TkdImageComponent } from '../image';
import { SecurityCodeComponent } from './../../shared/security-code/security-code.component';
import { ChangeDetectionStrategy, Component, Input } from "@angular/core";

@Component({
  selector: 'tkd-avatar',
  styles: [
    `
    :host {
      @apply inline-block rounded-full shadow;
    }
    `
  ],
  template: `
    <tkd-img class="w-full h-full rounded-full overflow-hidden block" [src]="src">
      <!-- Fallback Avatar -->
      <div class="bg-primary-100 text-primary-700 flex-grow self-stretch flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" fill="currentColor" stroke="currentColor" class="w-1/2 h-1/2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
        </svg>
      </div>
    </tkd-img>
  `,
  standalone: true,
  imports: [
    TkdImageComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TkdAvatarComponent {
  @Input()
  src?: string = '';
}
