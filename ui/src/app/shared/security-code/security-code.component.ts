import { coerceNumberProperty } from '@angular/cdk/coercion';
import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, DestroyRef, ElementRef, Input, QueryList, ViewChildren, forwardRef, inject } from '@angular/core';
import { ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-security-code',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './security-code.component.html',
  styleUrls: ['./security-code.component.css'],
  providers: [
    { provide: NG_VALUE_ACCESSOR, multi: true, useExisting: forwardRef(() => SecurityCodeComponent) }
  ]
})
export class SecurityCodeComponent implements ControlValueAccessor {
  onDestroy = inject(DestroyRef).onDestroy;
  cdr = inject(ChangeDetectorRef);

  @ViewChildren('codeInput', {read: ElementRef})
  inputElements!: QueryList<ElementRef<HTMLInputElement>>;

  inputs = Array.from(Array(6).keys()).map(v => ``)
  isDisabled = false;

  @Input()
  set count(v: any) {
    const number = coerceNumberProperty(v);
    this.inputs = Array.from(Array(number).keys()).map(v => ``)
  }
  get count() { return this.inputs.length }

  onKeyDown(event: KeyboardEvent, index: number) {
    if (event.keyCode === 8 && index > 0) {
      this.inputElements.get(index)!.nativeElement.value = '';
      this.inputElements.get(index-1)?.nativeElement.focus();

      event.preventDefault();
    }
  }

  onInput(event: Event, index: number) {
    const [first,...rest] = (event.target as any).value as string;

    (event.target as any).value = first ?? '';

    const lastInputBox = index === this.inputElements.length-1

    const didInsertContent = first!==undefined
    if(didInsertContent && !lastInputBox) {
      // continue to input the rest of the string
      this.inputElements.get(index+1)!.nativeElement.focus()
      this.inputElements.get(index+1)!.nativeElement.value = rest.join('')
      this.inputElements.get(index+1)!.nativeElement.dispatchEvent(new Event('input'))
    }

    const code = this.inputElements.map(ref => ref.nativeElement.value).join('');
    this._onChange(code);
  }

  writeValue(obj: string): void {
    obj = obj ?? '';

    const rest = obj.split('');
    rest.forEach((ch, idx) => {
      this.inputElements.get(idx)!.nativeElement.value = ch;
    })
    this.cdr.markForCheck();
  }

  _onChange: (val: string) => void = () => {};
  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  _onTouch: () => void = () => {};
  registerOnTouched(fn: any): void {
      this._onTouch = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.isDisabled = isDisabled;
    this.cdr.markForCheck();
  }
}
