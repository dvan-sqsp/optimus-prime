import { Directive, ElementRef, EventEmitter, HostListener, Input, Output } from '@angular/core';

@Directive({
  standalone: true,
  selector: '[clickOutside]'
})
export class ClickOutsideDirective {
  @Input() clickOutsideEnabled = true;
  @Output() clickOutside = new EventEmitter<void>();

  constructor(private elementRef: ElementRef) {}

  @HostListener('document:click', ['$event.target'])
  public onClick(target: any): void {
    const clickedInside = this.elementRef.nativeElement.contains(target);
    if (this.clickOutsideEnabled && !clickedInside) {
      this.clickOutside.emit();
    }
  }
}
