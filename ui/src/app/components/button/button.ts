import { ChangeDetectorRef, Directive, HostBinding, Input, OnChanges, SimpleChanges, inject } from '@angular/core';

export type TkdButtonType = 'primary' | 'secondary' | 'tertiary' | 'light' | 'dark';
export type TkdButtonSize = 'default' | 'small' | 'large';

const commonClasses = 'inline-block rounded font-medium uppercase leading-normal transition duration-150 ease-in-out focus:outline-none focus:ring-0 text-xs';

const buttonMap: {[key in TkdButtonType]: string} = {
  'primary': 'bg-primary  text-white shadow-[0_4px_9px_-4px_#3b71ca]  hover:bg-primary-600 hover:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.3),0_4px_18px_0_rgba(59,113,202,0.2)] focus:bg-primary-600 focus:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.3),0_4px_18px_0_rgba(59,113,202,0.2)]  active:bg-primary-700 active:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.3),0_4px_18px_0_rgba(59,113,202,0.2)] dark:shadow-[0_4px_9px_-4px_rgba(59,113,202,0.5)] dark:hover:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.2),0_4px_18px_0_rgba(59,113,202,0.1)] dark:focus:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.2),0_4px_18px_0_rgba(59,113,202,0.1)] dark:active:shadow-[0_8px_9px_-4px_rgba(59,113,202,0.2),0_4px_18px_0_rgba(59,113,202,0.1)]',
  'secondary': 'bg-primary-100  text-primary-700  hover:bg-primary-accent-100 focus:bg-primary-accent-100  active:bg-primary-accent-200',
  'tertiary': 'text-primary hover:text-primary-600 focus:text-primary-600  active:text-primary-700',
  'light': 'inline-block rounded bg-neutral-50 px-6 pb-2 pt-2.5 text-xs font-medium uppercase leading-normal text-neutral-800 shadow-[0_4px_9px_-4px_#cbcbcb] transition duration-150 ease-in-out hover:bg-neutral-100 hover:shadow-[0_8px_9px_-4px_rgba(203,203,203,0.3),0_4px_18px_0_rgba(203,203,203,0.2)] focus:bg-neutral-100 focus:shadow-[0_8px_9px_-4px_rgba(203,203,203,0.3),0_4px_18px_0_rgba(203,203,203,0.2)] focus:outline-none focus:ring-0 active:bg-neutral-200 active:shadow-[0_8px_9px_-4px_rgba(203,203,203,0.3),0_4px_18px_0_rgba(203,203,203,0.2)] dark:shadow-[0_4px_9px_-4px_rgba(251,251,251,0.3)] dark:hover:shadow-[0_8px_9px_-4px_rgba(251,251,251,0.1),0_4px_18px_0_rgba(251,251,251,0.05)] dark:focus:shadow-[0_8px_9px_-4px_rgba(251,251,251,0.1),0_4px_18px_0_rgba(251,251,251,0.05)] dark:active:shadow-[0_8px_9px_-4px_rgba(251,251,251,0.1),0_4px_18px_0_rgba(251,251,251,0.05)]',
  'dark': 'inline-block rounded bg-neutral-800 px-6 pb-2 pt-2.5 text-xs font-medium uppercase leading-normal text-neutral-50 shadow-[0_4px_9px_-4px_rgba(51,45,45,0.7)] transition duration-150 ease-in-out hover:bg-neutral-800 hover:shadow-[0_8px_9px_-4px_rgba(51,45,45,0.2),0_4px_18px_0_rgba(51,45,45,0.1)] focus:bg-neutral-800 focus:shadow-[0_8px_9px_-4px_rgba(51,45,45,0.2),0_4px_18px_0_rgba(51,45,45,0.1)] focus:outline-none focus:ring-0 active:bg-neutral-900 active:shadow-[0_8px_9px_-4px_rgba(51,45,45,0.2),0_4px_18px_0_rgba(51,45,45,0.1)] dark:bg-neutral-900 dark:shadow-[0_4px_9px_-4px_#030202] dark:hover:bg-neutral-900 dark:hover:shadow-[0_8px_9px_-4px_rgba(3,2,2,0.3),0_4px_18px_0_rgba(3,2,2,0.2)] dark:focus:bg-neutral-900 dark:focus:shadow-[0_8px_9px_-4px_rgba(3,2,2,0.3),0_4px_18px_0_rgba(3,2,2,0.2)] dark:active:bg-neutral-900 dark:active:shadow-[0_8px_9px_-4px_rgba(3,2,2,0.3),0_4px_18px_0_rgba(3,2,2,0.2)]'
}

const sizeMap: {[key in TkdButtonSize]: string} = {
  'default': 'px-6 pb-2 pt-2.5',
  'small': 'px-4 pb-[5px] pt-[6px]',
  'large': 'px-7 pb-2.5 pt-3',
}

@Directive({
  selector: '[tkd-button]',
  standalone: true,
})
export class TkdButtonDirective implements OnChanges {
  private readonly cdr = inject(ChangeDetectorRef);

  @Input('tkd-button')
  tkdButtonType: TkdButtonType | '' = 'secondary';

  @Input('tkdSize')
  tkdSize: TkdButtonSize = 'default';

  @HostBinding('class')
  classList: string = '';

  ngOnChanges(changes: SimpleChanges): void {
    if ('tkdButtonType' in changes || 'tkdSize' in changes) {
      this._updateStyle()
    }
  }

  private _updateStyle() {
    let classes = [
      commonClasses
    ];

    classes.push(buttonMap[this.tkdButtonType|| 'secondary'])
    classes.push(sizeMap[this.tkdSize]);

    this.classList = classes.join(' ');

    this.cdr.markForCheck();
  }
}
