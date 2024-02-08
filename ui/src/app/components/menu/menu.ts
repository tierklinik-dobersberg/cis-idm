import { trigger, transition, style, animate } from '@angular/animations';
import { NgClass, NgIf, NgTemplateOutlet } from '@angular/common';
import { Input, AfterViewInit, Component, ContentChild, Directive, ElementRef, HostBinding, OnInit, Renderer2, inject, TemplateRef, HostListener, ViewChild } from "@angular/core";
import { RouterLinkActive } from "@angular/router";

@Directive({
  selector: '[tkd-menu]',
  standalone: true,
})
export class TkdMenuDirective {
  readonly elementRef = inject(ElementRef);

  @HostBinding('class')
  readonly classList = 'py-4 space-y-2 border-t border-b border-slate-200 dark:border-slate-600 dark:text-white [&>svg]:h-5 [&>svg]:w-5';
}

@Directive({
  selector: '[tkd-menu-item]',
  hostDirectives: [RouterLinkActive],
  standalone: true,
})
export class TkdMenuItemDirective implements OnInit {
  readonly elementRef = inject(ElementRef);

  private routerLink = inject(RouterLinkActive, {
    host: true,
    optional: true,
  })

  @HostBinding('class')
  readonly classList = 'flex flex-row gap-3 border-r-2 border-transparent justify-start items-center pr-3 py-2 text-base cursor-pointer hover:bg-slate-200 dark:hover:bg-slate-600 pl-12 [&>svg:first-child]:-ml-8 [&>svg]:h-5 [&>svg]:w-5 [&>ul[tkd-menu]]:w-full';

  ngOnInit(): void {
    if (this.routerLink) {
      this.routerLink.routerLinkActive = '!border-slate-600 !dark:border-slate-200 bg-slate-100 dark:bg-slate-600'
    }
  }
}

@Component({
  selector: 'li[tkd-sub-menu]',
  standalone: true,
  imports: [
    TkdMenuDirective,
    TkdMenuItemDirective,
    NgIf,
    NgTemplateOutlet,
    NgClass
  ],
  styles:[
    `
    :host {
      overflow: hidden;
    }
    `
  ],
  template: `
    <div tkd-menu-item class="z-[500] bg-white relative">
      <ng-container *ngIf="icon">
        <ng-container *NgTemplateOutlet="icon"></ng-container>
      </ng-container>

      <span class="flex-grow">{{ title }}</span>

      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5 transition duration-150 ease-in-out" [ngClass]="{'rotate-90': isOpen}">
        <path fill-rule="evenodd" d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
      </svg>
    </div>

    <div #menuContainer *ngIf="isOpen" [@animate] class="overflow-hidden border-l-2 border-dashed ml-6 border-gray-200 dark:border-gray-600">
      <ng-content></ng-content>
    </div>
  `,
  animations: [
    trigger('animate', [
      transition(':enter', [
        style({
          opacity: 0,
          transform: 'translateY(-100%)'
        }),
        animate('150ms ease-in-out', style({
          opacity: 1,
          transform: 'translateY(0%)'
        }))
      ]),
      transition(':leave', [
        style({
          opacity: 1,
          transform: 'translateY(0%)'
        }),
        animate('150ms ease-in-out', style({
          opacity: 0,
          transform: 'translateY(-100%)'
        }))
      ]),
    ])
  ]
})
export class TkdSubMenuComponent {
  private readonly renderer = inject(Renderer2);
  readonly elementRef = inject(ElementRef);

  @Input('tkd-sub-menu')
  title: string = '';

  @Input('tkdIcon')
  icon?: TemplateRef<any>;

  isOpen = false;

  @ViewChild('menuContainer', {read: ElementRef, static: false})
  menuContainer?: ElementRef<HTMLElement>;

  @HostListener('click', ['$event'])
  handleClick(event: MouseEvent) {
    let iter: HTMLElement | null = event.target as HTMLElement;
    while (iter !== this.elementRef.nativeElement && iter !== null) {
      if (iter === this.menuContainer?.nativeElement) {
        return
      }

      iter = iter.parentElement;
    }

    this.isOpen = !this.isOpen;
  }
}
