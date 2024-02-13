import { Component, DestroyRef, ElementRef, EventEmitter, Input, OnChanges, OnInit, Output, Renderer2, SimpleChanges, booleanAttribute, forwardRef, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR } from "@angular/forms";
import { LayoutService } from "@tierklinik-dobersberg/angular/layout";
import { Datepicker } from 'tw-elements';

@Component({
  selector: 'tkd-datepicker',
  exportAs: 'tkdDatepicker',
  standalone: true,
  imports: [
    FormsModule,
  ],
  template: `
    <input type="text" [(ngModel)]="value" (ngModelChange)="handleChange()" class="tkd-input w-full">
  `,
  styles: [
    `:host {
      @apply block relative;
    }`
  ],
  providers: [
    { provide: NG_VALUE_ACCESSOR, multi: true, useExisting: forwardRef(() => TkdDatepickerComponent) }
  ]
})
export class TkdDatepickerComponent implements OnInit, OnChanges, ControlValueAccessor {
  private readonly destroyRef = inject(DestroyRef);
  private readonly renderer = inject(Renderer2);
  private readonly element = inject(ElementRef);
  private readonly layout = inject(LayoutService).withAutoUpdate();

  private picker!: typeof Datepicker;
  value = '';

  @Input({transform: booleanAttribute})
  disableFuture = false;

  @Input()
  title = '';

  @Output()
  openChange = new EventEmitter<boolean>();

  private isOpen = false;

  ngOnChanges(changes: SimpleChanges) {
    if (this.picker) {
      this.picker.update({disableFuture: this.disableFuture, title: this.title});
    }
  }

  ngOnInit() {
    this.picker = new Datepicker(this.element.nativeElement, {
      disableFuture: this.disableFuture,
      format: "yyyy-mm-dd",
      confirmDateOnSelect: true,
      startDay: 1,
      monthsFull: [
        'Jänner',
        'Feburar',
        'März',
        'April',
        'Mai',
        'Juni',
        'Juli',
        'August',
        'September',
        'Oktober',
        'November',
        'Dezember',
      ],
      monthsShort: [
        'Jän',
        'Feb',
        'Mär',
        'Apr',
        'Mai',
        'Jun',
        'Jul',
        'Aug',
        'Sep',
        'Okt',
        'Nov',
        'Dez',
      ],
      weekdaysFull: [
        'Sonntag',
        'Montag',
        'Dienstag',
        'Mittwoch',
        'Donnerstag',
        'Freitag',
        'Samstag',
      ],
      weekdaysNarrow: ['S', 'M', 'D', 'M', 'D', 'F', 'S'],
      weekdaysShort: ['Son', 'Mon', 'Die', 'Mit', 'Don', 'Fre', 'Sam'],
      title: this.title,
    });

    this.destroyRef.onDestroy(() => this.picker.dispose())

    this.layout
      .change
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        if (this.isOpen) {
          this.picker.close()
        }

        this.picker.update({
          inline: this.layout.md,
        })
      })

    let cleanup = this.renderer.listen(this.element.nativeElement, 'blur', () => this._onBlur())
    this.destroyRef.onDestroy(cleanup);

    cleanup = this.renderer.listen(this.element.nativeElement, 'open.te.datepicker', () => {
      this.isOpen = true;
      this.openChange.next(true)
    })
    this.destroyRef.onDestroy(cleanup)

    cleanup = this.renderer.listen(this.element.nativeElement, 'close.te.datepicker', () => {
      this.isOpen = false;
      this.openChange.next(false)
    })
    this.destroyRef.onDestroy(cleanup)
  }

  handleChange() {
    this._onChange(this.value);
  }

  setDisabledState(isDisabled: boolean): void {
    if (isDisabled) {
      this.renderer.setAttribute(this.element.nativeElement, 'disabled', '')
    } else {
      this.renderer.removeAttribute(this.element.nativeElement, 'disabled')
    }
  }

  _onChange = (_: any) => { }
  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  _onBlur = () => { }
  registerOnTouched(fn: any): void {
    this._onBlur = fn
  }

  writeValue(obj: string): void {
    this.value = obj;
  }

  open() {
    this.picker?.open()
  }

  close() {
    this.picker?.close();
  }
}
