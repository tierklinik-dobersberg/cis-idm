import { Component, DestroyRef, ElementRef, OnInit, Renderer2, forwardRef, inject } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR } from "@angular/forms";
import { LayoutService } from "@tierklinik-dobersberg/angular/layout";
import { Timepicker } from 'tw-elements';

@Component({
  selector: 'tkd-timepicker',
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
    { provide: NG_VALUE_ACCESSOR, multi: true, useExisting: forwardRef(() => TkdTimepickerComponent) }
  ]
})
export class TkdTimepickerComponent implements OnInit, ControlValueAccessor {
  private readonly destroyRef = inject(DestroyRef);
  private readonly renderer = inject(Renderer2);
  private readonly element = inject(ElementRef);
  private readonly layout = inject(LayoutService).withAutoUpdate();

  private picker!: typeof Timepicker;
  value = '';

  ngOnInit() {
    this.picker = new Timepicker(this.element.nativeElement, {
      format24: true,
    });

    this.destroyRef.onDestroy(() => this.picker.dispose())

    this.layout
      .change
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        this.picker.close()
        this.picker.update({
          inline: this.layout.md,
        })
      })

    let cleanup = this.renderer.listen(this.element.nativeElement, 'blur', () => this._onBlur())
    this.destroyRef.onDestroy(cleanup);
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
}
