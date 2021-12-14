import { Controller } from "@hotwired/stimulus";

export default class extends Controller {
  connect() {
    this.originalSrc = this.src;
    this.originalInnerHTML = this.innerHTML;
    this.restoreDisabled = false;
  }

  restore(event) {
    if (this.restoreDisabled) {
      return;
    }

    if (!this.frameSrcHasChanged) {
      return;
    }

    if (this.elementInFrame(event.target)) {
      return;
    }

    this.restoreContent();
  }

  disableRestore({ target: { value }, params: { unlessBlank } }) {
    if (unlessBlank && value == "") {
      this.restoreDisabled = false;
      return;
    }

    this.restoreDisabled = true;
  }

  elementInFrame(element) {
    return this.frame.contains(element);
  }

  restoreContent() {
    this.frame.innerHTML = this.originalInnerHTML;
  }

  get frame() {
    return this.element;
  }

  get frameSrcHasChanged() {
    return this.originalSrc != this.src;
  }

  get innerHTML() {
    return this.frame.innerHTML;
  }

  get src() {
    return this.frame.src;
  }
}
