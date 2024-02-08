import { NgModule } from "@angular/core";
import { TkdMenuDirective, TkdMenuItemDirective, TkdSubMenuComponent } from "./menu";

@NgModule({
  imports: [
    TkdMenuDirective,
    TkdMenuItemDirective,
    TkdSubMenuComponent,
  ],
  exports: [
    TkdMenuDirective,
    TkdMenuItemDirective,
    TkdSubMenuComponent,
  ]
})
export class TkdMenuModule {}
