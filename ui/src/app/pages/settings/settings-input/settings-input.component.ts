import { NgFor, NgForOf, NgIf, NgSwitch, NgSwitchCase } from "@angular/common";
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, Output, booleanAttribute, forwardRef, inject } from "@angular/core";
import { ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR } from "@angular/forms";
import { TkdDatepickerComponent } from "src/app/components/datepicker";
import { TkdSelectComponent } from "src/app/components/select";
import { TkdSwitchDirective } from "src/app/components/switch";
import { TkdTimepickerComponent } from "src/app/components/timepicker";
import { FieldConfig } from "src/app/config.service";

@Component({
  selector: 'tkd-settings-input',
  templateUrl: './settings-input.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [
    NgSwitch,
    NgSwitchCase,
    NgIf,
    FormsModule,
    NgFor,
    NgForOf,
    TkdSelectComponent,
    TkdDatepickerComponent,
    TkdTimepickerComponent,
    TkdSwitchDirective,
  ],
  providers: [
    {provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => TkdSettingsInputComponent), multi: true}
  ]
})
export class TkdSettingsInputComponent<T> implements ControlValueAccessor {
  private readonly cdr = inject(ChangeDetectorRef);

  @Input()
  config!: FieldConfig;

  @Input()
  value?: T;

  @Output()
  valueChange = new EventEmitter<T>();

  @Input({transform: booleanAttribute})
  disabled = false;

  _onBlur = () => {};
  registerOnTouched(fn: any): void {
    this._onBlur = fn
  }

  private _onChange = (_: T) => {};
  registerOnChange(fn: any): void {
    this._onChange = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  writeValue(obj: T): void {
    this.value = obj;
    this.cdr.markForCheck();
  }

  updateValue(value: T) {
    this.value = value;
    this.valueChange.next(value);
    this._onChange(value);
  }
}
