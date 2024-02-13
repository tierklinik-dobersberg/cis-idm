import { Pipe, PipeTransform } from "@angular/core";
import { FieldConfig } from "src/app/config.service";
import type { Struct } from "@bufbuild/protobuf";

const simpleTypes = ['string', 'number', 'bool', 'date', 'time'];

@Pipe({
  name: 'simpleFields',
  pure: true,
  standalone: true
})
export class SimpleFieldsPipe implements PipeTransform {
  transform(list: FieldConfig[] | null): FieldConfig[] {
    if (list === null) {
      return []
    }

    return list.filter(c => simpleTypes.includes(c.type) && c.visibility !== 'private')
  }
}

@Pipe({
  name: 'complexFields',
  pure: true,
  standalone: true
})
export class ComplexFieldsPipe implements PipeTransform {
  transform(list: FieldConfig[] | null): FieldConfig[] {
    if (list === null) {
      return []
    }

    return list.filter(c => !simpleTypes.includes(c.type) && c.visibility !== 'private' && c.type !== 'any')
  }
}

@Pipe({
  name: 'fieldValue',
  pure: true,
  standalone: true
})
export class FieldValuePipe implements PipeTransform {
  transform(data: Struct | null | undefined, field: FieldConfig):any {
    if (!data) {
      return null;
    }

    return data.fields[field.name]?.kind?.value || null;
  }
}

@Pipe({
  name: 'fieldPath',
  pure: true,
  standalone: true
})
export class FieldPathPipe implements PipeTransform {
  transform(value: (FieldConfig | string)[]) {
    return value.map(f => {
      if (typeof f === 'object') {
        return f.name
      }

      return f
    }) .join('.')
  }
}
