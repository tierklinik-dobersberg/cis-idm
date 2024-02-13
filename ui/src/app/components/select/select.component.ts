import { DestroyRef, Directive, ElementRef, HostBinding, OnInit, Renderer2, forwardRef, inject } from "@angular/core";
import { ControlValueAccessor, NG_VALUE_ACCESSOR, NgModel } from "@angular/forms";
import { Select } from 'tw-elements';

@Directive({
  selector: 'select[tkd-select]',
  standalone: true,
  providers: [
    { provide: NG_VALUE_ACCESSOR, multi: true, useExisting: forwardRef(() => TkdSelectComponent) }
  ]
})
export class TkdSelectComponent implements OnInit, ControlValueAccessor {
  private readonly destroyRef = inject(DestroyRef);
  private readonly renderer = inject(Renderer2);
  private readonly element = inject(ElementRef);

  private select!: typeof Select;
  private value = '';

  ngOnInit() {
    this.select = new Select(this.element.nativeElement, {}, {
      selectInput: 'tkd-input cursor-pointer'
    })
    this.destroyRef.onDestroy(() => this.select.dispose())

    let isFirst = true;
    let cleanup = this.renderer.listen(this.element.nativeElement, 'valueChange.te.select', (event: any) => {
      if (!event.value) {
        return
      }

      if (isFirst) {
        isFirst = false;
        return
      }

      if (event.value !== this.value) {
        this._onChange(event.value);
        this.value = event.value;
      }
    })

    this.destroyRef.onDestroy(cleanup);

    cleanup = this.renderer.listen(this.element.nativeElement, 'blur', () => this._onBlur())
    this.destroyRef.onDestroy(cleanup);
  }

  setDisabledState(isDisabled: boolean): void {
    if (isDisabled) {
      this.renderer.setAttribute(this.element.nativeElement, 'disabled', '')
    } else {
      this.renderer.removeAttribute(this.element.nativeElement, 'disabled')
    }
  }

  _onChange = (_: any) =>  {}
  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  _onBlur = () => {}
  registerOnTouched(fn: any): void {
    this._onBlur = fn
  }

  writeValue(obj: string): void {
    this.value = obj;
    this.select.setValue(obj)
  }
}
